package repositories

import (
	"notification-service/models"
	"time"

	"gorm.io/gorm"
)

type EventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) *EventRepository {
	return &EventRepository{db: db}
}

// GetEventsByDateRange obtiene eventos en un rango de fechas
func (r *EventRepository) GetEventsByDateRange(startDate, endDate time.Time) ([]*models.Event, error) {
	var events []*models.Event

	err := r.db.Where("date >= ? AND date < ?", startDate, endDate).Find(&events).Error
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (r *EventRepository) GetEventsForTomorrow() ([]*models.Event, error) {
	tomorrow := time.Now().AddDate(0, 0, 1)
	start := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, tomorrow.Location())
	end := start.Add(24 * time.Hour)

	var events []*models.Event
	err := r.db.Where("date >= ? AND date < ? AND reminder_day_before = ?", start, end, true).Find(&events).Error
	return events, err
}

func (r *EventRepository) GetEventsForToday() ([]*models.Event, error) {
	today := time.Now()
	start := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
	end := start.Add(24 * time.Hour)

	var events []*models.Event
	err := r.db.Where("date >= ? AND date < ? AND reminder_day = ?", start, end, true).Find(&events).Error
	return events, err
}
