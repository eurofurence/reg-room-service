package entity

type Room struct {
	Base

	// Name is the name of the room
	Name string `gorm:"type:varchar(80) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;NOT NULL;uniqueIndex:room_room_name_uidx"`

	// Flags is a comma-separated list of flags, with both leading and trailing comma. The allowed flags are configuration dependent
	Flags string `gorm:"type:varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci"`

	// Comments are optional, not processed in any way
	Comments string `gorm:"type:varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci" testdiff:"ignore"`

	// Size is the size of the room
	Size uint
}

type RoomMember struct {
	Member

	// TODO references to get integrity check!

	// RoomID references the room to which the attendee belongs
	RoomID string `gorm:"type:varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;NOT NULL;index:room_room_member_roomid"`
}
