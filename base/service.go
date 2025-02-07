package base

import "proximityService/database"

type service struct {
	dataStoreSvc *database.ProximityDBService
}

type Service interface {
	ListAllBusiness()
	CreateBusiness()
	UpdateBusiness()
	DeleteBusiness()
	GetBusiness()
	GetNearbyBusinesses()
}

func NewService(dbSvc database.ProximityDBService) Service {
	return &service{
		dataStoreSvc: &dbSvc,
	}
}

func (s *service) ListAllBusiness() {

}
func (s *service) CreateBusiness() {

}
func (s *service) UpdateBusiness() {

}
func (s *service) DeleteBusiness() {

}
func (s *service) GetBusiness() {

}

func (s *service) GetNearbyBusinesses() {

}
