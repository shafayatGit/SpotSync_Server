package dto

import "time"

type ReservationZoneResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type ReservationUserResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type MyReservationResponse struct {
	ID           uint                    `json:"id"`
	LicensePlate string                  `json:"license_plate"`
	Status       string                  `json:"status"`
	Zone         ReservationZoneResponse `json:"zone"`
	CreatedAt    time.Time               `json:"created_at"`
}

type AdminReservationResponse struct {
	ID           uint                    `json:"id"`
	UserID       uint                    `json:"user_id"`
	User         ReservationUserResponse `json:"user"`
	ZoneID       uint                    `json:"zone_id"`
	Zone         ReservationZoneResponse `json:"zone"`
	LicensePlate string                  `json:"license_plate"`
	Status       string                  `json:"status"`
	CreatedAt    time.Time               `json:"created_at"`
	UpdatedAt    time.Time               `json:"updated_at"`
}

type CreateReservationResponse struct {
	ID           uint      `json:"id"`
	UserID       uint      `json:"user_id"`
	ZoneID       uint      `json:"zone_id"`
	LicensePlate string    `json:"license_plate"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
