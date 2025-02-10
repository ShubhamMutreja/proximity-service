package api

import (
	"proximityService/database"
	"proximityService/models"
)

type service struct {
	DataStore database.ProximityDBService
}

func NewService(dbService database.ProximityDBService) Service {
	return &service{
		DataStore: dbService,
	}
}

type Service interface {
	ListAllBusiness(req models.NearbySearchRequest) models.BusinessResponse
	BulkCreateBusiness(reqs []models.Business) models.BusinessResponse
	CreateBusiness(req models.Business) models.BusinessResponse
	UpdateBusiness(req models.Business) models.BusinessResponse
	DeleteBusiness(req models.Business) models.BusinessResponse
	GetBusiness(req models.Business) models.BusinessResponse
	GetNearbyBusinesses(req models.NearbySearchRequest) models.BusinessResponse
}

func (s *service) ListAllBusiness(req models.NearbySearchRequest) models.BusinessResponse {
	data := s.DataStore.GetAllBusinessesFromDB(req)

	var resp models.BusinessResponse
	resp.Action = "List of All Businesses for Debug"
	resp.Businesses = data
	return resp
}

func (s *service) BulkCreateBusiness(reqs []models.Business) models.BusinessResponse {
	var resp models.BusinessResponse

	resp.Action = "The Following Businesses Has Been Successfully Registered"
	for _, req := range reqs {
		data := s.DataStore.PublishNewBusinessToDB(req)
		resp.Businesses = append(resp.Businesses, data)
	}
	return resp
}

func (s *service) CreateBusiness(req models.Business) models.BusinessResponse {
	data := s.DataStore.PublishNewBusinessToDB(req)

	var resp models.BusinessResponse
	resp.Action = "The Following Business Has Been Successfully Registered"
	resp.Businesses = append(resp.Businesses, data)
	return resp
}

func (s *service) UpdateBusiness(req models.Business) models.BusinessResponse {
	data := s.DataStore.UpdateBusinessInDB(req)

	var resp models.BusinessResponse
	resp.Action = "The Following Business Has Been Successfully Updated"
	resp.Businesses = append(resp.Businesses, data)
	return resp
}

func (s *service) DeleteBusiness(req models.Business) models.BusinessResponse {
	data := s.DataStore.DeleteBusinessFromDB(req)

	var resp models.BusinessResponse
	resp.Action = "The Following Business Has Been Successfully Deleted"
	resp.Businesses = append(resp.Businesses, data)
	return resp
}

func (s *service) GetBusiness(req models.Business) models.BusinessResponse {
	data := s.DataStore.GetBusinessFromDB(req)

	var resp models.BusinessResponse
	resp.Action = "The Following Business Was Requested"
	resp.Businesses = append(resp.Businesses, data)
	return resp
}

func (s *service) GetNearbyBusinesses(req models.NearbySearchRequest) models.BusinessResponse {
	data := s.DataStore.GetNearbyBusinessesFromQuadTree(req)

	var resp models.BusinessResponse
	resp.Action = "The Following Business Was Requested By The User"
	resp.Businesses = append(resp.Businesses, data...)
	return resp
}
