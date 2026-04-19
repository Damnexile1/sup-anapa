package models

import "time"

type Instructor struct {
	ID          int         `db:"id"`
	Name        string      `db:"name"`
	Photo       string      `db:"photo"`
	Description string      `db:"description"`
	Phone       string      `db:"phone"`
	WalkTypes   []*WalkType `db:"-"`
	CreatedAt   time.Time   `db:"created_at"`
	UpdatedAt   time.Time   `db:"updated_at"`
}

type WalkType struct {
	ID           int       `db:"id"`
	InstructorID int       `db:"instructor_id"`
	Name         string    `db:"name"`
	Price        int       `db:"price"`
	MaxPeople    int       `db:"max_people"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

type Slot struct {
	ID            int        `db:"id"`
	Date          time.Time  `db:"date"`
	StartTime     string     `db:"start_time"`
	EndTime       string     `db:"end_time"`
	Price         int        `db:"price"`
	MaxPeople     int        `db:"max_people"`
	InstructorID  int        `db:"instructor_id"`
	WalkTypeID    int        `db:"walk_type_id"`
	WalkTypeName  string     `db:"-"`
	Status        string     `db:"status"`
	HoldExpiresAt *time.Time `db:"hold_expires_at"`
	CreatedAt     time.Time  `db:"created_at"`
	UpdatedAt     time.Time  `db:"updated_at"`
}

type Booking struct {
	ID          int       `db:"id"`
	SlotID      int       `db:"slot_id"`
	ClientName  string    `db:"client_name"`
	ClientPhone string    `db:"client_phone"`
	ClientEmail string    `db:"client_email"`
	PeopleCount int       `db:"people_count"`
	Status      string    `db:"status"` // pending, confirmed, cancelled
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type Admin struct {
	ID           int       `db:"id"`
	Username     string    `db:"username"`
	PasswordHash string    `db:"password_hash"`
	CreatedAt    time.Time `db:"created_at"`
}

type WeatherCache struct {
	ID          int       `db:"id"`
	Date        time.Time `db:"date"`
	AirTemp     float64   `db:"air_temp"`
	WaterTemp   float64   `db:"water_temp"`
	WindSpeed   float64   `db:"wind_speed"`
	CloudCover  int       `db:"cloud_cover"`
	Description string    `db:"description"`
	CachedAt    time.Time `db:"cached_at"`
}
