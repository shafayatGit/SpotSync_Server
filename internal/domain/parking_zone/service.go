package parking_zone

import (
	"spotsync/internal/domain/parking_zone/dto"
)

type ParkingZoneService interface {
	CreateZone(req *dto.CreateParkingZoneRequest) (*dto.ParkingZoneResponse, error)
	GetAllZones() ([]dto.ParkingZoneResponse, error)
	GetZoneByID(id uint) (*dto.ParkingZoneResponse, error)
}

type parkingZoneService struct {
	repo ParkingZoneRepository
}

func NewParkingZoneService(repo ParkingZoneRepository) ParkingZoneService {
	return &parkingZoneService{repo: repo}
}

func (s *parkingZoneService) CreateZone(req *dto.CreateParkingZoneRequest) (*dto.ParkingZoneResponse, error) {
	zone := &ParkingZone{
		Name:          req.Name,
		Type:          req.Type,
		TotalCapacity: req.TotalCapacity,
		PricePerHour:  req.PricePerHour,
	}

	if err := s.repo.Create(zone); err != nil {
		return nil, err
	}

	return &dto.ParkingZoneResponse{
		ID:            zone.ID,
		Name:          zone.Name,
		Type:          zone.Type,
		TotalCapacity: zone.TotalCapacity,
		PricePerHour:  zone.PricePerHour,
		CreatedAt:     zone.CreatedAt,
		UpdatedAt:     zone.UpdatedAt,
	}, nil
}

func (s *parkingZoneService) GetAllZones() ([]dto.ParkingZoneResponse, error) {
	zones, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	counts, err := s.repo.GetActiveReservationsCounts()
	if err != nil {
		counts = make(map[uint]int)
	}

	var res []dto.ParkingZoneResponse
	for _, zone := range zones {
		activeCount := counts[zone.ID]
		available := zone.TotalCapacity - activeCount
		if available < 0 {
			available = 0
		}

		res = append(res, dto.ParkingZoneResponse{
			ID:             zone.ID,
			Name:           zone.Name,
			Type:           zone.Type,
			TotalCapacity:  zone.TotalCapacity,
			AvailableSpots: &available,
			PricePerHour:   zone.PricePerHour,
			CreatedAt:      zone.CreatedAt,
		})
	}
	return res, nil
}

func (s *parkingZoneService) GetZoneByID(id uint) (*dto.ParkingZoneResponse, error) {
	zone, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	activeCount, err := s.repo.GetActiveReservationsCount(id)
	if err != nil {
		activeCount = 0
	}

	available := zone.TotalCapacity - activeCount
	if available < 0 {
		available = 0
	}

	return &dto.ParkingZoneResponse{
		ID:             zone.ID,
		Name:           zone.Name,
		Type:           zone.Type,
		TotalCapacity:  zone.TotalCapacity,
		AvailableSpots: &available,
		PricePerHour:   zone.PricePerHour,
		CreatedAt:      zone.CreatedAt,
	}, nil
}
