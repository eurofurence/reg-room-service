package entity

import (
	"time"
)

type Member struct {
	// ID contains the badge number of the attendee (an attendee can only either be in a
	// group or invited, and can only ever be in one room at the same time).
	ID        int64 `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	// intentionally not supplying DeletedAt -- don't want soft delete

	// Nickname caches the nickname of the attendee
	Nickname string `gorm:"type:varchar(80) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;NOT NULL"`

	// AvatarURL caches the url to obtain the avatar for this attendee, points to an image such as a png or jpg
	AvatarURL string `gorm:"type:varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci"`

	// Flags is a comma-separated list of flags such as "has_key", with a leading and trailing comma
	Flags string `gorm:"type:varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci"`
}
