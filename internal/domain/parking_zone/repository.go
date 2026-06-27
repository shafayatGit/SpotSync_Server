package parking_zone

import (
	"gorm.io/gorm"
)

type ParkingZoneRepository interface {
	Create(zone *ParkingZone) error
	FindAll() ([]ParkingZone, error)
	FindByID(id uint) (*ParkingZone, error)
	GetActiveReservationsCount(zoneID uint) (int, error)
	GetActiveReservationsCounts() (map[uint]int, error)
}

type parkingZoneRepository struct {
	db *gorm.DB
}

func NewParkingZoneRepository(db *gorm.DB) ParkingZoneRepository {
	return &parkingZoneRepository{db: db}
}

func (r *parkingZoneRepository) Create(zone *ParkingZone) error {
	return r.db.Create(zone).Error
}

func (r *parkingZoneRepository) FindAll() ([]ParkingZone, error) {
	var zones []ParkingZone
	if err := r.db.Find(&zones).Error; err != nil {
		return nil, err
	}
	return zones, nil
}

func (r *parkingZoneRepository) FindByID(id uint) (*ParkingZone, error) {
	var zone ParkingZone
	if err := r.db.First(&zone, id).Error; err != nil {
		return nil, err
	}
	return &zone, nil
}

func (r *parkingZoneRepository) GetActiveReservationsCount(zoneID uint) (int, error) {
	var count int64
	if err := r.db.Table("reservations").
		Where("zone_id = ? AND status = ?", zoneID, "active").
		Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

func (r *parkingZoneRepository) GetActiveReservationsCounts() (map[uint]int, error) {
	var results []struct {
		ZoneID uint `gorm:"column:zone_id"`
		Count  int  `gorm:"column:count"`
	}
	if err := r.db.Table("reservations").
		Select("zone_id, COUNT(*) as count").
		Where("status = ?", "active").
		Group("zone_id").
		Scan(&results).Error; err != nil {
		return nil, err
	}
	counts := make(map[uint]int)
	for _, res := range results {
		counts[res.ZoneID] = res.Count
	}
	return counts, nil
}
