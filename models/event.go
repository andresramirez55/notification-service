package models

import (
	"time"

	"gorm.io/gorm"
)

type Event struct {
	ID                uint           `json:"id" gorm:"primaryKey"`
	Title             string         `json:"title" gorm:"not null"`
	Description       string         `json:"description"`
	Date              time.Time      `json:"date" gorm:"not null"`
	Time              string         `json:"time"` // Format: "HH:MM"
	Location          string         `json:"location"`
	Email             string         `json:"email" gorm:"not null"`
	Phone             string         `json:"phone" gorm:"not null"`
	ReminderDay       bool           `json:"reminder_day" gorm:"default:true"`        // Reminder on the same day
	ReminderDayBefore bool           `json:"reminder_day_before" gorm:"default:true"` // Reminder one day before
	IsAllDay          bool           `json:"is_all_day" gorm:"default:false"`
	Color             string         `json:"color" gorm:"default:'#007AFF'"`
	Priority          string         `json:"priority" gorm:"default:'medium'"`
	Category          string         `json:"category"`
	NotifyFamily      bool           `json:"notify_family" gorm:"default:false"`
	NotifyPapa        bool           `json:"notify_papa" gorm:"default:false"`
	NotifyMama        bool           `json:"notify_mama" gorm:"default:false"`
	ChildTag          string         `json:"child_tag"`
	SelectedChildren  string         `json:"selected_children" gorm:"type:text"`
	FamilyMembers     string         `json:"family_members" gorm:"type:text"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}
