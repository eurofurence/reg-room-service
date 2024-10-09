package entity

import (
	"gorm.io/gorm"
	"time"
)

// Group is a group of attendees that wish to be assigned to a Room together.
type Group struct {
	Base

	// Name is the name of the group
	Name string `gorm:"type:varchar(80) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;NOT NULL;uniqueIndex:room_group_name_uidx"`

	// Flags is a comma-separated list of flags, with both leading and trailing comma. The allowed flags are configuration dependent
	Flags string `gorm:"type:varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci"`

	// Comments are optional, not processed in any way
	Comments string `gorm:"type:varchar(4096) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci" testdiff:"ignore"`

	// MaximumSize defaults to a value from service configuration, but we store it here so admins can increase it manually for some groups
	MaximumSize int64

	// Owner is the badge number (attendee ID) of the attendee owning the group. Ownership can be passed to another attendee.
	Owner int64
}

// GroupMember associates attendees to a group, either as a member or as an invited member.
type GroupMember struct {
	Member

	// GroupID references the group to which the member belongs (or has been invited)
	//
	// Note: foreign key constraint added programmatically in MysqlRepository.Migrate()
	GroupID string `gorm:"type:varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;NOT NULL;index:room_group_member_grpid"`

	// IsInvite is true if the member has been invited, or false if the member has already joined
	IsInvite bool
}

type GroupBan struct {
	// ID contains the badge number of the attendee (an attendee can only either be in a
	ID int64 `gorm:"primaryKey;autoIncrement:false"`
	// GroupID references the group from which the member has been banned
	//
	// Note: foreign key constraint added programmatically in MysqlRepository.Migrate()
	GroupID   string `gorm:"primaryKey;type:varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// Flags is a comma-separated list of flags, with both leading and trailing comma. The allowed flags are configuration dependent
	Flags string `gorm:"type:varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci"`

	// Comments are optional, not processed in any way
	Comments string `gorm:"type:varchar(4096) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci" testdiff:"ignore"`
}
