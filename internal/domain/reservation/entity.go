package reservation

import (
	"time"
	"spotsync/internal/domain/parking_zone"
	"spotsync/internal/domain/user"
)

type Reservation struct {
	ID           uint                     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID       uint                     `gorm:"not null" json:"user_id"`
	User         user.User                `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"user,omitempty"`
	ZoneID       uint                     `gorm:"not null" json:"zone_id"`
	Zone         parking_zone.ParkingZone `gorm:"foreignKey:ZoneID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"zone,omitempty"`
	LicensePlate string                   `gorm:"type:varchar(15);not null" json:"license_plate"`
	Status       string                   `gorm:"type:varchar(20);default:'active';not null" json:"status"`
	CreatedAt    time.Time                `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time                `gorm:"autoUpdateTime" json:"updated_at"`
}
