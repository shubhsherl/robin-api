package models

import (
	"github.com/RobinHoodArmyHQ/robin-api/pkg/nanoid"
	"time"
)

type CheckIn struct {
	ID                uint64        `json:"-" gorm:"primaryKey"`
	CheckinID         nanoid.NanoID `json:"checkin_id,omitempty" gorm:"checkin_id"`
	UserID            nanoid.NanoID `json:"user_id,omitempty"`
	EventID           nanoid.NanoID `json:"event_id,omitempty"`
	PhotoIDs          []int64       `json:"photo_ids,omitempty" gorm:"serializer:json"`
	Description       string        `json:"description,omitempty"`
	NoOfPeopleServed  uint64        `json:"no_of_people_served,omitempty"`
	NoOfStudentTaught uint64        `json:"no_of_student_taught,omitempty"`
	CreatedAt         time.Time     `json:"created_at,omitempty"`
	UpdatedAt         time.Time     `json:"updated_at,omitempty"`
}

func (CheckIn) TableName() string {
	return "checkins"
}
