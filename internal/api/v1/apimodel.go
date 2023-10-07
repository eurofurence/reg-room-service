package v1

type Group struct {
	// The internal primary key of the group, in the form of a UUID. Only set when reading groups, completely ignored when you send a group to us.
	ID string `yaml:"id" json:"id"`
	// The name of the group, must be unique, but otherwise just used for display purposes
	Name string `yaml:"name" json:"name"`
	// A list of flags as declared in configuration. Flags are used to store yes/no-style information about the group such as \"wheelchair\", etc.
	Flags []string `yaml:"flags,omitempty" json:"flags,omitempty"`
	// Optional comments the owner wishes to make regarding the group. Not processed in any way.
	Comments *string `yaml:"comments,omitempty" json:"comments,omitempty"`
	// if set higher than 0 (the default), will limit the number of people that can join the group. Note that there is also a configuration item that globally limits the size of groups, e.g. to the maximum room size.
	MaximumSize *int32 `yaml:"maximum_size,omitempty" json:"maximum_size,omitempty"`
	// the badge number of the group owner. Must be a member of the group. If you are not an admin, you can only create groups with yourself as owner.
	Owner int32 `yaml:"owner" json:"owner"`
	// the current group members. READ ONLY, provided for ease of use of the API, but completely ignored in all write requests. Please use the relevant subresource API endpoints to manipulate group membership.
	Members []Member `yaml:"members,omitempty" json:"members,omitempty"`
	// the current outstanding invites for this group. READ ONLY, provided for ease of use of the API, but completely ignored in all write requests. Please use the relevant subresource API endpoints to send/revoke invites.
	Invites []Member `yaml:"invites,omitempty" json:"invites,omitempty"`
}

type GroupList struct {
	Groups []Group `yaml:"groups" json:"groups"`
}

type Member struct {
	// badge number (id in the attendee service).
	ID int32 `yaml:"id" json:"id"`
	// The nickname of the attendee, proxied from that attendee service.
	Nickname string `yaml:"nickname" json:"nickname"`
	// A url to obtain the avatar for this attendee, points to an image such as a png or jpg. May require the same authentication this API expects.
	Avatar *string `yaml:"avatar,omitempty" json:"avatar,omitempty"`
	// Set to true if this person has been given a key to the room, for groups this can only be set if already assigned a room.
	HasKey bool `yaml:"hasKey" json:"hasKey"`
}

type Room struct {
	// The internal primary key of the room, in the form of a UUID. Only set when reading rooms, completely ignored when you send a room to us.
	ID *string `yaml:"id,omitempty" json:"id,omitempty"`
	// The name of the room, must be unique, but otherwise just used for display purposes
	Name string `yaml:"name" json:"name"`
	// A comma separated list of flags as declared in configuration. Flags are used to store yes/no-style information about the room.
	Flags *string `yaml:"flags,omitempty" json:"flags,omitempty"`
	// Optional comment. Not processed in any way.
	Comments *string `yaml:"comments,omitempty" json:"comments,omitempty"`
	// the maximum room size, usually the number of sleeping spots/beds in the room.
	Size int32 `yaml:"size" json:"size"`
	// the assigned room members. READ ONLY, provided for ease of use of the API, but completely ignored in all write requests. Please use the relevant subresource API endpoints to manipulate individual or group assignments.
	Members []Member `yaml:"members,omitempty" json:"members,omitempty"`
}

type RoomList struct {
	Rooms []Room `yaml:"rooms" json:"rooms"`
}

// Empty defines a type which is used for empty responses.
type Empty struct {
}
