package api

import (
	"proximityService/database"
	"proximityService/models"
)

func ListAllBusiness(req models.NearbySearchRequest) models.BusinessResponse {
	svc := database.GetDataStore()
	data := svc.GetAllBusinessesFromDB(req)

	var resp models.BusinessResponse
	resp.Action = "List of Businesses for Debug"
	resp.Businesses = data
	return resp
}

func BulkCreateBusiness(reqs []models.Business) models.BusinessResponse {
	svc := database.GetDataStore()
	var resp models.BusinessResponse
	resp.Action = "The Following Businesses Has Been Succesfully Registered"
	for _, req := range reqs {
		data := svc.PublishNewBusinessToDB(req)
		resp.Businesses = append(resp.Businesses, data)
	}
	return resp
}

func CreateBusiness(req models.Business) models.BusinessResponse {
	svc := database.GetDataStore()
	data := svc.PublishNewBusinessToDB(req)

	var resp models.BusinessResponse
	resp.Action = "The Following Business Has Been Succesfully Registered"
	resp.Businesses = append(resp.Businesses, data)
	return resp
}

func UpdateBusiness(req models.Business) models.BusinessResponse {
	svc := database.GetDataStore()
	data := svc.UpdateBusinessInDB(req)

	var resp models.BusinessResponse
	resp.Action = "The Following Business Has Been Succesfully Updated"
	resp.Businesses = append(resp.Businesses, data)
	return resp
}

func DeleteBusiness(req models.Business) models.BusinessResponse {
	svc := database.GetDataStore()
	data := svc.DeleteBusinessFromDB(req)

	var resp models.BusinessResponse
	resp.Action = "The Following Business Has Been Succesfully Deleted"
	resp.Businesses = append(resp.Businesses, data)
	return resp
}

func GetBusiness(req models.Business) models.BusinessResponse {
	svc := database.GetDataStore()
	data := svc.GetBusinessFromDB(req)

	var resp models.BusinessResponse
	resp.Action = "The Following Business Was Requested"
	resp.Businesses = append(resp.Businesses, data)
	return resp
}

func GetNearbyBusinesses(req models.NearbySearchRequest) models.BusinessResponse {
	svc := database.GetDataStore()
	data := svc.GetNearbyBusinessesFromQuadTree(req)

	var resp models.BusinessResponse
	resp.Action = "The Following Business Was Requested By The User"
	resp.Businesses = append(resp.Businesses, data...)
	return resp
}
