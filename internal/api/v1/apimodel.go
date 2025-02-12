package modelsv1

import "time"

var _ = time.Now

type Error struct {
	// The time at which the error occurred, formatted as ISO datetime according to spec.
	Timestamp string `json:"timestamp"`
	// An internal trace id assigned to the error. Used to find logs associated with errors across our services. Display to the user as something to communicate to us with inquiries about the error.
	Requestid string `json:"requestid"`
	// A keyed description of the error. We do not write human readable text here because the user interface will be multi language.  At this time, there are these values: - auth.unauthorized (token missing completely or invalid) - auth.forbidden (permissions missing)
	Message string `json:"message"`
	// Optional additional details about the error. If available, will usually contain English language technobabble.
	Details map[string][]string `json:"details,omitempty"`
}

type Group struct {
	// The internal primary key of the group, in the form of a UUID. Only set when reading groups, completely ignored when you send a group to us.
	ID string `yaml:"id" json:"id"`
	// The name of the group, must be unique, but otherwise just used for display purposes
	Name string `yaml:"name" json:"name"`
	// A list of flags as declared in configuration. Flags are used to store yes/no-style information about the group such as \"wheelchair\", etc.
	Flags []string `yaml:"flags" json:"flags"`
	// Optional comments the owner wishes to make regarding the group. Not processed in any way.
	Comments *string `yaml:"comments,omitempty" json:"comments,omitempty"`
	// if set higher than 0 (the default), will limit the number of people that can join the group. Note that there is also a configuration item that globally limits the size of groups, e.g. to the maximum room size.
	MaximumSize int64 `yaml:"maximum_size" json:"maximum_size"`
	// the badge number of the group owner. Must be a member of the group. If you are not an admin, you can only create groups with yourself as owner.
	Owner int64 `yaml:"owner" json:"owner"`
	// the current group members. READ ONLY, provided for ease of use of the API, but completely ignored in all write requests. Please use the relevant subresource API endpoints to manipulate group membership.
	Members []Member `yaml:"members,omitempty" json:"members,omitempty"`
	// the current outstanding invites for this group. READ ONLY, provided for ease of use of the API, but completely ignored in all write requests. Please use the relevant subresource API endpoints to send/revoke invites.
	Invites []Member `yaml:"invites,omitempty" json:"invites,omitempty"`
}

type GroupCreate struct {
	// The name of the group, must be unique, but otherwise just used for display purposes
	Name string `yaml:"name" json:"name"`
	// A list of flags as declared in configuration. Flags are used to store yes/no-style information about the group such as \"wheelchair\", etc.
	Flags []string `yaml:"flags,omitempty" json:"flags,omitempty"`
	// Optional comments the owner wishes to make regarding the group. Not processed in any way.
	Comments *string `yaml:"comments,omitempty" json:"comments,omitempty"`
	// if set higher than 0 (the default), will limit the number of people that can join the group. Note that there is also a configuration item that globally limits the size of groups, e.g. to the maximum room size.
	MaximumSize int64 `yaml:"maximum_size" json:"maximum_size"`
	// the badge number of the group owner. If you are not an admin, you can only create groups with yourself as owner. Defaults to yourself.
	Owner int64 `yaml:"owner" json:"owner"`
}

type GroupList struct {
	Groups []*Group `yaml:"groups" json:"groups"`
}

type Member struct {
	// badge number (id in the attendee service).
	ID int64 `yaml:"id" json:"id"`
	// The nickname of the attendee, proxied from that attendee service.
	Nickname string `yaml:"nickname" json:"nickname"`
	// A url to obtain the avatar for this attendee, points to an image such as a png or jpg. May require the same authentication this API expects.
	Avatar *string `yaml:"avatar,omitempty" json:"avatar,omitempty"`
	// A list of membership flags as declared in configuration. Flags are used to store yes/no-style information.
	Flags []string `yaml:"flags,omitempty" json:"flags,omitempty"`
}

type Room struct {
	// The internal primary key of the room, in the form of a UUID. Only set when reading rooms, completely ignored when you send a room to us.
	ID string `yaml:"id" json:"id"`
	// The name of the room, must be unique, but otherwise just used for display purposes
	Name string `yaml:"name" json:"name"`
	// A list of flags as declared in configuration. Flags are used to store yes/no-style information about the room.
	Flags []string `yaml:"flags" json:"flags"`
	// Optional comment. Not processed in any way.
	Comments *string `yaml:"comments,omitempty" json:"comments,omitempty"`
	// the maximum room size, usually the number of sleeping spots/beds in the room.
	Size int64 `yaml:"size" json:"size"`
	// the assigned room occupants. READ ONLY, provided for ease of use of the API, but completely ignored in all write requests. Please use the relevant subresource API endpoints to manipulate assignments.
	Occupants []Member `yaml:"occupants,omitempty" json:"occupants,omitempty"`
}

type RoomCreate struct {
	// The name of the room, must be unique, but otherwise just used for display purposes
	Name string `yaml:"name" json:"name"`
	// A list of flags as declared in configuration. Flags are used to store yes/no-style information about the room.
	Flags []string `yaml:"flags" json:"flags"`
	// Optional comment. Not processed in any way.
	Comments *string `yaml:"comments,omitempty" json:"comments,omitempty"`
	// the room size, usually the number of sleeping spots/beds in the room.
	Size int64 `yaml:"size" json:"size"`
}

type RoomList struct {
	Rooms []*Room `yaml:"rooms" json:"rooms"`
}

// Countdown contains information about the time until the secret is revealed, which is needed for the registration.
type Countdown struct {
	// CurrentTimeIsoDateTime is the current time on the server.
	CurrentTimeIsoDateTime string `json:"currentTime"`
	// TargetTimeIsoDateTime is the time at which the countdown ends (may depend on authorization, e.g. staff may register earlier than normal users).
	TargetTimeIsoDateTime string `json:"targetTime"`
	// CountdownSeconds is the number of seconds until the countdown ends
	// (may depend on authorization, e.g. staff may register earlier than normal users). Stays at 0 if the countdown is over.
	CountdownSeconds int64 `json:"countdown"`
	// Secret is the secret code word you'll need to give the hotel
	// (may depend on authorization, e.g. staff gets a different code word that allows earlier room booking).
	// Will be missing before your countdown has reached 0.
	Secret string `json:"secret,omitempty"`
}

// Empty defines a type which is used for empty responses.
type Empty struct {
}
