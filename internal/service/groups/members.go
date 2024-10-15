package groupservice

import (
	"context"
	"errors"
	"fmt"
	aulogging "github.com/StephanHCB/go-autumn-logging"
	"github.com/eurofurence/reg-room-service/internal/application/common"
	"github.com/eurofurence/reg-room-service/internal/entity"
	"github.com/eurofurence/reg-room-service/internal/repository/config"
	"github.com/eurofurence/reg-room-service/internal/repository/downstreams"
	"github.com/eurofurence/reg-room-service/internal/repository/downstreams/attendeeservice"
	"github.com/eurofurence/reg-room-service/internal/repository/downstreams/mailservice"
	"github.com/eurofurence/reg-room-service/internal/service/rbac"
	"gorm.io/gorm"
	"math/rand"
	"net/url"
)

// AddGroupMemberParams is the request type for the AddMemberToGroup operation.
//
// See OpenAPI spec for more details.
type AddGroupMemberParams struct {
	// GroupID is the ID of the group where a user should be added
	GroupID string
	// BadgeNumber is the registration number of a user
	BadgeNumber int64
	// Nickname is the nickname of a registered user that should receive
	// an invitation Email.
	Nickname string
	// Code is the invite code that can be used to join a group.
	Code string
	// Force is an admin only flag that allows to bypass the
	// validations.
	Force bool
}

func (g *groupService) AddMemberToGroup(ctx context.Context, req *AddGroupMemberParams) (string, error) {
	adminPerm, loggedInAttendee, err := g.groupMembershipAuthCheck(ctx)
	if err != nil {
		return "", err
	}

	requestedAttendee, err := g.validateRequestedAttendee(ctx, req.BadgeNumber)
	if err != nil {
		return "", err
	}

	grp, gm, err := g.groupMembershipExisting(ctx, req.GroupID, req.BadgeNumber) // gm may be nil if not exists
	if err != nil {
		return "", err
	}

	informOwnerTemplate := ""
	informMemberTemplate := ""
	inviteCode := ""

	banned, err := g.DB.HasGroupBan(ctx, req.GroupID, req.BadgeNumber)
	if err != nil {
		return "", errGroupRead(ctx, err.Error())
	}

	if gm == nil {
		// no existing membership entry

		gm = g.DB.NewEmptyGroupMembership(ctx, req.GroupID, requestedAttendee.ID, requestedAttendee.Nickname)

		if adminPerm && req.Force {
			// admin mode, directly add and even allow cross-user additions

			if banned {
				aulogging.Infof(ctx, "group ban removed through force add - group %s badge %d by %s", req.GroupID, req.BadgeNumber, common.GetSubject(ctx))
				err := g.DB.RemoveGroupBan(ctx, req.GroupID, requestedAttendee.ID)
				if err != nil {
					return "", err
				}
			}

			gm.IsInvite = false
			gm.InvitationCode = ""
			gm.Comments = "forced join by admin " + common.GetSubject(ctx)

			err := g.DB.AddGroupMembership(ctx, gm)
			if err != nil {
				return "", errGroupWrite(ctx, err.Error())
			}

			informOwnerTemplate = "group-member-joined"
		} else if req.BadgeNumber == loggedInAttendee.ID {
			// trying to add self - invite coming from the joining attendee

			if banned {
				aulogging.Warnf(ctx, "user tried to circumvent ban - group %s badge %d by %s", req.GroupID, req.BadgeNumber, common.GetSubject(ctx))
				return "", common.NewForbidden(ctx, common.AuthForbidden, common.Details("you cannot join this group - please stop trying"))
			}

			gm.IsInvite = true
			gm.InvitationCode = "" // join request by owner, so no code
			gm.Comments = "self join request by " + common.GetSubject(ctx)

			err := g.DB.AddGroupMembership(ctx, gm)
			if err != nil {
				return "", errGroupWrite(ctx, err.Error())
			}

			informOwnerTemplate = "group-member-request"
		} else if grp.Owner == loggedInAttendee.ID {
			// owner trying to invite another attendee - check nickname matches

			if req.Nickname != requestedAttendee.Nickname {
				return "", common.NewBadRequest(ctx, common.GroupInviteMismatch, common.Details("nickname did not match - you need to know the nickname to be able to invite this attendee"))
			}

			if banned {
				aulogging.Infof(ctx, "group ban removed through owner add - group %s badge %d by %s", req.GroupID, req.BadgeNumber, common.GetSubject(ctx))
				err := g.DB.RemoveGroupBan(ctx, req.GroupID, req.BadgeNumber)
				if err != nil {
					return "", err
				}
			}

			gm.IsInvite = true
			gm.InvitationCode = rollInvitationCode()
			gm.Comments = "invite by owner " + common.GetSubject(ctx)

			inviteCode = fmt.Sprintf("?code=%s", gm.InvitationCode)

			err = g.DB.AddGroupMembership(ctx, gm)
			if err != nil {
				return "", errGroupWrite(ctx, err.Error())
			}

			informMemberTemplate = "group-invited" // you have been invited and here's your link
		} else {
			return "", common.NewForbidden(ctx, common.AuthForbidden, common.Details("only the group owner or an admin can invite other people into a group"))
		}
	} else {
		// existing membership (possibly invitation)

		if grp.ID != gm.GroupID {
			return "", common.NewConflict(ctx, common.GroupMemberConflict, common.Details("this attendee is already invited to another group or in another group"))
		}

		if !gm.IsInvite {
			return "", common.NewConflict(ctx, common.GroupMemberDuplicate, common.Details("this attendee is already a member of this group"))
		}

		if adminPerm && req.Force {
			// admin mode, directly add and allow cross-user additions

			if banned {
				// this is a rare timing edge case, normally an invitation or application with an active ban cannot happen
				aulogging.Infof(ctx, "group ban override and remove through force add - group %s badge %d by %s", req.GroupID, req.BadgeNumber, common.GetSubject(ctx))
				err := g.DB.RemoveGroupBan(ctx, req.GroupID, req.BadgeNumber)
				if err != nil {
					return "", err
				}
			}

			gm.IsInvite = false
			gm.InvitationCode = ""

			err = g.DB.UpdateGroupMembership(ctx, gm)
			if err != nil {
				return "", errGroupWrite(ctx, err.Error())
			}

			informOwnerTemplate = "group-member-joined"
		} else if req.BadgeNumber == loggedInAttendee.ID {
			// self accept after invite

			if req.Code != gm.InvitationCode {
				aulogging.Infof(ctx, "invited user failed to join due to invitation code mismatch - group %s badge %d by %s", req.GroupID, req.BadgeNumber, common.GetSubject(ctx))
				return "", common.NewForbidden(ctx, common.AuthForbidden, common.Details("you must provide the invitation code you were sent in order to join"))
			}

			gm.IsInvite = false

			err = g.DB.UpdateGroupMembership(ctx, gm)
			if err != nil {
				return "", errGroupWrite(ctx, err.Error())
			}

			informOwnerTemplate = "group-member-joined"
		} else if grp.Owner == loggedInAttendee.ID {
			// owner accept after apply

			if banned {
				// this is a rare timing edge case, normally an application with an active ban cannot happen
				aulogging.Infof(ctx, "group ban removed through owner add - group %s badge %d by %s", req.GroupID, req.BadgeNumber, common.GetSubject(ctx))
				err := g.DB.RemoveGroupBan(ctx, req.GroupID, req.BadgeNumber)
				if err != nil {
					return "", err
				}
			}

			gm.IsInvite = false
			// keep invitation code, multiple clicks should be idempotent

			err = g.DB.UpdateGroupMembership(ctx, gm)
			if err != nil {
				return "", errGroupWrite(ctx, err.Error())
			}

			informMemberTemplate = "group-application-accepted" // you have been added to the group
		} else {
			return "", common.NewForbidden(ctx, common.AuthForbidden, common.Details("only the group owner or an admin can accept invitations from others into a group"))
		}
	}

	_ = g.sendInfoMails(ctx, informOwnerTemplate, informMemberTemplate, grp, req.BadgeNumber, inviteCode)
	// can still see results in regsys, so do not fail at this point

	return inviteCode, nil
}

// RemoveGroupMemberParams is the request type for the RemoveMemberFromGroup operation.
//
// See OpenAPI spec for more details.
type RemoveGroupMemberParams struct {
	// GroupID is the ID of the group where a user should be added
	GroupID string
	// BadgeNumber is the registration number of a user
	BadgeNumber int64
	// AutoDeny future invitations (effectively creates or removes a ban)
	AutoDeny bool
}

func (g *groupService) RemoveMemberFromGroup(ctx context.Context, req *RemoveGroupMemberParams) error {
	adminPerm, loggedInAttendee, err := g.groupMembershipAuthCheck(ctx)
	if err != nil {
		return err
	}

	grp, gm, err := g.groupMembershipExisting(ctx, req.GroupID, req.BadgeNumber) // gm may be nil if not exists
	if err != nil {
		return err
	}

	if gm == nil {
		return common.NewNotFound(ctx, common.GroupMemberNotFound, common.Details("this attendee is not in any group"))
	}

	if grp.ID != gm.GroupID {
		return common.NewConflict(ctx, common.GroupMemberConflict, common.Details("this attendee is invited to a different group or in a different group"))
	}

	adjustBan := false
	informOwnerTemplate := ""
	informMemberTemplate := ""
	if adminPerm {
		adjustBan = true
		if gm.IsInvite {
			aulogging.Infof(ctx, "admin void group invitation - group %s badge %d by admin %s", req.GroupID, req.BadgeNumber, common.GetSubject(ctx))
		} else {
			aulogging.Infof(ctx, "admin kick from group - group %s badge %d by admin %s", req.GroupID, req.BadgeNumber, common.GetSubject(ctx))
			informOwnerTemplate = "group-member-removed" // member was removed
			informMemberTemplate = "group-member-kicked" // you have been removed
		}
	} else if grp.Owner == loggedInAttendee.ID {
		adjustBan = true
		if gm.IsInvite {
			aulogging.Infof(ctx, "declined join request - group %s badge %d by owner %s", req.GroupID, req.BadgeNumber, common.GetSubject(ctx))
			informMemberTemplate = "group-request-declined" // your request to join was declined
		} else {
			aulogging.Infof(ctx, "kick from group - group %s badge %d by owner %s", req.GroupID, req.BadgeNumber, common.GetSubject(ctx))
			informMemberTemplate = "group-member-kicked"
		}
	} else if gm.ID == loggedInAttendee.ID {
		if gm.IsInvite {
			aulogging.Infof(ctx, "declined group invitation - group %s badge %d by self", req.GroupID, req.BadgeNumber)
			informOwnerTemplate = "group-request-declined" // your request to join was declined
		} else {
			aulogging.Infof(ctx, "left group - group %s badge %d by self", req.GroupID, req.BadgeNumber)
			informOwnerTemplate = "group-member-left" // member left
		}
	} else {
		return common.NewForbidden(ctx, common.AuthForbidden, common.Details("only the group owner or an admin can remove other people from a group"))
	}

	banned, err := g.DB.HasGroupBan(ctx, req.GroupID, req.BadgeNumber)
	if err != nil {
		return errGroupRead(ctx, err.Error())
	}

	err = g.DB.DeleteGroupMembership(ctx, req.BadgeNumber)
	if err != nil {
		return errGroupWrite(ctx, err.Error())
	}

	_ = g.sendInfoMails(ctx, informOwnerTemplate, informMemberTemplate, grp, req.BadgeNumber, "")
	// can still see results in regsys, so do not fail at this point

	if adjustBan {
		if req.AutoDeny && !banned {
			aulogging.Infof(ctx, "group ban added - group %s badge %d by %s", req.GroupID, req.BadgeNumber, common.GetSubject(ctx))
			comment := fmt.Sprintf("group ban added by %s", common.GetSubject(ctx))
			err := g.DB.AddGroupBan(ctx, req.GroupID, req.BadgeNumber, comment)
			if err != nil {
				return err
			}
		} else if !req.AutoDeny && banned {
			aulogging.Infof(ctx, "group ban removed - group %s badge %d by %s", req.GroupID, req.BadgeNumber, common.GetSubject(ctx))
			err := g.DB.RemoveGroupBan(ctx, req.GroupID, req.BadgeNumber)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// mails

func (g *groupService) sendInfoMails(ctx context.Context, informOwnerTemplate string, informMemberTemplate string, grp *entity.Group, memberID int64, inviteCode string) error {
	if informOwnerTemplate == "" && informMemberTemplate == "" {
		return nil
	}

	owner, err := g.AttSrv.GetAttendee(ctx, grp.Owner)
	if err != nil {
		aulogging.WarnErrf(ctx, err, "failed to obtain attendee info for group owner %d: %s", grp.Owner, err.Error())
		return err
	}

	member, err := g.AttSrv.GetAttendee(ctx, memberID)
	if err != nil {
		aulogging.WarnErrf(ctx, err, "failed to obtain attendee info for group member %d: %s", memberID, err.Error())
		return err
	}

	if informOwnerTemplate != "" {
		mailRequest := mailservice.MailSendDto{
			CommonID: informOwnerTemplate,
			Lang:     owner.RegistrationLanguage,
			To:       []string{owner.Email},
			Variables: map[string]string{
				"nickname":            owner.Nickname,
				"groupname":           grp.Name,
				"object_badge_number": fmt.Sprintf("%d", member.ID),
				"object_nickname":     member.Nickname,
			},
		}

		err := g.MailSrv.SendEmail(ctx, mailRequest)
		if err != nil {
			aulogging.WarnErrf(ctx, err, "failed to send email to group owner %d about %s: %s", grp.Owner, informOwnerTemplate, err.Error())
			return err
		}
	}

	if informMemberTemplate != "" {
		conf, err := config.GetApplicationConfig()
		if err != nil {
			aulogging.WarnErrf(ctx, err, "bug - application config not loaded - failed to send email to group member %d about %s: %s", memberID, informMemberTemplate, err.Error())
			return err
		}

		requestURL, ok := ctx.Value(common.CtxKeyRequestURL{}).(*url.URL)
		if !ok {
			aulogging.WarnErrf(ctx, err, "bug - request URL not in context - failed to send email to group member %d about %s: %s", memberID, informMemberTemplate, err.Error())
			return errors.New("could not retrieve base URL from context - this is an implementation error")
		}

		mailRequest := mailservice.MailSendDto{
			CommonID: informMemberTemplate,
			Lang:     member.RegistrationLanguage,
			To:       []string{member.Email},
			Variables: map[string]string{
				"nickname":  member.Nickname,
				"groupname": grp.Name,
			},
		}
		if inviteCode != "" {
			// TODO this is probably not quite correct
			mailRequest.Variables["url"] = conf.Service.JoinLinkBaseURL + requestURL.Path + inviteCode
		}

		err = g.MailSrv.SendEmail(ctx, mailRequest)
		if err != nil {
			aulogging.WarnErrf(ctx, err, "failed to send email to group member %d about %s: %s", memberID, informMemberTemplate, err.Error())
			return err
		}
	}

	return nil
}

// internals

func (g *groupService) validateRequestedAttendee(ctx context.Context, badgeNo int64) (attendeeservice.Attendee, error) {
	if badgeNo <= 0 {
		return attendeeservice.Attendee{}, common.NewBadRequest(ctx, common.GroupDataInvalid, common.Details("attendee badge number must be positive integer"))
	}

	attendee, err := g.AttSrv.GetAttendee(ctx, badgeNo)
	if err != nil {
		if errors.Is(err, downstreams.ErrDownStreamNotFound) {
			return attendeeservice.Attendee{}, common.NewNotFound(ctx, common.NoSuchAttendee, common.Details("no such attendee"))
		} else {
			return attendeeservice.Attendee{}, common.NewBadGateway(ctx, common.DownstreamAttSrv, common.Details("failed to look up invited attendee - internal error, see logs for details"))
		}
	}

	if err := g.checkAttending(ctx, badgeNo); err != nil {
		return attendeeservice.Attendee{}, err
	}

	return attendee, nil
}

func (g *groupService) groupMembershipAuthCheck(ctx context.Context) (bool, attendeeservice.Attendee, error) {
	validator, err := rbac.NewValidator(ctx)
	if err != nil {
		aulogging.ErrorErrf(ctx, err, "Could not retrieve RBAC validator from context. [error]: %v", err)
		return false, attendeeservice.Attendee{}, errCouldNotGetValidator(ctx)
	}

	if validator.IsAdmin() || validator.IsAPITokenCall() {
		attendee, _ := g.loggedInUserValidRegistrationBadgeNo(ctx)

		// admin requests are allowed through even if the admin does not have a valid registration
		return true, attendee, nil
	}

	attendee, err := g.loggedInUserValidRegistrationBadgeNo(ctx)
	if err != nil {
		return false, attendeeservice.Attendee{}, err
	}
	return false, attendee, nil
}

func (g *groupService) groupMembershipExisting(ctx context.Context, groupID string, badgeNo int64) (*entity.Group, *entity.GroupMember, error) {
	grp, err := g.DB.GetGroupByID(ctx, groupID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, common.NewNotFound(ctx, common.GroupIDNotFound, common.Details("this group does not exist"))
		} else {
			return nil, nil, errGroupRead(ctx, err.Error())
		}
	}

	gm, err := g.DB.GetGroupMembershipByAttendeeID(ctx, badgeNo)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// maybe ok, does not have a membership record
			return grp, nil, nil
		} else {
			return grp, nil, errInternal(ctx, err.Error())
		}
	}

	return grp, gm, nil
}

func rollInvitationCode() string {
	return randomHumanReadableString(8) // 40 bits - about one in a trillion
}

func randomHumanReadableString(length int) string {
	text := ""
	charset := []rune("ABCDEFGHJKLMNPQRSTUVWXYZ23456789") // 32 characters = 5 bits

	for i := 1; i < length-1; i++ {
		idx := rand.Intn(len(charset))
		text += fmt.Sprintf("%c", charset[idx])
	}

	return text
}
