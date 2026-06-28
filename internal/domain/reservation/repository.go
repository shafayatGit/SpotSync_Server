package reservation

import (
	"errors"
	"spotsync/internal/domain/parking_zone"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrZoneFull            = errors.New("parking zone is full")
	ErrZoneNotFound        = errors.New("parking zone not found")
	ErrReservationNotFound = errors.New("reservation not found")
)

type ReservationRepository interface {
	CreateWithLock(userID, zoneID uint, licensePlate string) (*Reservation, error)
	FindMyReservations(userID uint) ([]Reservation, error)
	FindAll() ([]Reservation, error)
	FindByID(id uint) (*Reservation, error)
	UpdateStatus(id uint, status string) error
}

type reservationRepository struct {
	db *gorm.DB
}

func NewReservationRepository(db *gorm.DB) ReservationRepository {
	return &reservationRepository{db: db}
}

func (r *reservationRepository) CreateWithLock(userID, zoneID uint, licensePlate string) (*Reservation, error) {
	var res *Reservation

	err := r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Lock the parking zone row
		var zone parking_zone.ParkingZone
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&zone, zoneID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrZoneNotFound
			}
			return err
		}

		// 2. Count active reservations in this zone
		var activeCount int64
		if err := tx.Model(&Reservation{}).
			Where("zone_id = ? AND status = ?", zoneID, "active").
			Count(&activeCount).Error; err != nil {
			return err
		}

		// 3. Verify capacity
		if activeCount >= int64(zone.TotalCapacity) {
			return ErrZoneFull
		}

		// 4. Create reservation
		res = &Reservation{
			UserID:       userID,
			ZoneID:       zoneID,
			LicensePlate: licensePlate,
			Status:       "active",
		}
		if err := tx.Create(res).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *reservationRepository) FindMyReservations(userID uint) ([]Reservation, error) {
	var reservations []Reservation
	if err := r.db.Preload("Zone").Where("user_id = ?", userID).Find(&reservations).Error; err != nil {
		return nil, err
	}
	return reservations, nil
}

func (r *reservationRepository) FindAll() ([]Reservation, error) {
	var reservations []Reservation
	if err := r.db.Preload("User").Preload("Zone").Find(&reservations).Error; err != nil {
		return nil, err
	}
	return reservations, nil
}

func (r *reservationRepository) FindByID(id uint) (*Reservation, error) {
	var res Reservation
	if err := r.db.First(&res, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrReservationNotFound
		}
		return nil, err
	}
	return &res, nil
}

func (r *reservationRepository) UpdateStatus(id uint, status string) error {
	result := r.db.Model(&Reservation{}).Where("id = ?", id).Update("status", status)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrReservationNotFound
	}
	return nil
}
