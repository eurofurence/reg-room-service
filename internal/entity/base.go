package entity

import (
	"time"

	"gorm.io/gorm"
)

type Base struct {
	ID        string `gorm:"primaryKey; type:varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;NOT NULL"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
