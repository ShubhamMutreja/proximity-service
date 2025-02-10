package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"proximityService/models"
	"slices"

	_ "github.com/lib/pq"

	quadtreeservice "proximityService/quadTreeService"
)

const (
	CREATETABLE    = "CREATE TABLE IF NOT EXISTS business (id SERIAL PRIMARY KEY, name TEXT, longitude FLOAT, latitude FLOAT, phone VARCHAR, city VARCHAR, state VARCHAR, zipcode INT)"
	INSERTQUERY    = "INSERT INTO business (name, latitude, longitude, phone, city, state, zipcode) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id"
	UPDATEQUERY    = "UPDATE business SET name = $1, latitude = $2, longitude = $3, phone = $4, city = $5, state = $6, zipcode = $7 WHERE id = $8"
	SELECTALLQUERY = "SELECT * FROM business"
	SELECTONEQUERY = "SELECT * FROM business WHERE id = $1"
	DELETEONEQUERY = "DELETE FROM business WHERE id = $1"
)

type proximityDBService struct {
	QT quadtreeservice.QuadTree
	DB *sql.DB
}

type ProximityDBService interface {
	GetAllBusinessesFromDB(req models.NearbySearchRequest) []models.Business
	GetBusinessFromDB(req models.Business) models.Business
	PublishNewBusinessToDB(req models.Business) models.Business
	UpdateBusinessInDB(req models.Business) models.Business
	DeleteBusinessFromDB(req models.Business) models.Business
	GetNearbyBusinessesFromQuadTree(req models.NearbySearchRequest) []models.Business
}

func InitDataBase() *sql.DB {
	connStr := fmt.Sprintf("postgres://postgres:%s@%s/%s?sslmode=disable", os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_DB"))
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	//create the table if it doesn't exist
	_, err = db.Exec(CREATETABLE)

	if err != nil {
		log.Fatal(err)
	}
	log.Print("database running")
	return db
}

// QuadTree Opeations
// get nearby businesses
func InitQuadTree() quadtreeservice.QuadTree {
	//latitude range is from -90 to 90 and longitude range is from -180 to 180
	startLatitude := -90.0
	startLongitude := -180.0
	width := 180.0
	height := 360.0
	maxSizeAtLeafNodes := 3
	return *quadtreeservice.NewQuadTree(startLatitude, startLongitude, width, height, maxSizeAtLeafNodes)
}

func NewDBService(qT quadtreeservice.QuadTree, dB *sql.DB) ProximityDBService {
	return &proximityDBService{
		QT: qT,
		DB: dB,
	}
}

// Database operations
// get all users
func (ds *proximityDBService) GetAllBusinessesFromDB(req models.NearbySearchRequest) []models.Business {
	rows, err := ds.DB.Query(SELECTALLQUERY)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	businesses := []models.Business{}
	for rows.Next() {
		var u models.Business
		if err := rows.Scan(&u.ID, &u.Name, &u.Location.Longitude, &u.Location.Latitude, &u.Phone, &u.City, &u.State, &u.ZipCode); err != nil {
			log.Println(err)
		}
		dist := quadtreeservice.GetDistance(req.UserLocation, u.Location)
		u.Dist = &dist
		if dist <= req.Radius {
			businesses = append(businesses, u)
		}
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
	}
	slices.SortFunc(businesses, func(a, b models.Business) int {
		if *a.Dist-*b.Dist < 0 {
			return -1
		} else if *a.Dist-*b.Dist > 0 {
			return 1
		} else {
			return 0
		}
	})
	//db.QT.UpdateQuadTree(businesses)
	return businesses
}

// get user by id
func (ds *proximityDBService) GetBusinessFromDB(req models.Business) models.Business {
	var u models.Business
	err := ds.DB.QueryRow(SELECTONEQUERY, req.ID).Scan(&u.ID, &u.Name, &u.Location.Longitude, &u.Location.Latitude, &u.Phone, &u.City, &u.State, &u.ZipCode)
	if err != nil {
		log.Println(err)
	}
	return u
}

// create user
func (ds *proximityDBService) PublishNewBusinessToDB(req models.Business) models.Business {
	err := ds.DB.QueryRow(INSERTQUERY, req.Name, req.Location.Latitude, req.Location.Longitude, req.Phone, req.City, req.State, req.ZipCode).Scan(&req.ID)
	if err != nil {
		log.Println(err)
	}
	ds.QT.UpdateQuadTree([]models.Business{req})
	return req
}

// update user
func (ds *proximityDBService) UpdateBusinessInDB(req models.Business) models.Business {
	err := ds.DB.QueryRow(SELECTONEQUERY, req.ID).Scan(&req.ID, &req.Name, &req.Location.Longitude, &req.Location.Latitude, &req.Phone, &req.City, &req.State, &req.ZipCode)
	ds.QT.DeleteFromQuadTree([]models.Business{req})

	_, err = ds.DB.Exec(UPDATEQUERY, &req.ID, &req.Name, &req.Location.Latitude, &req.Location.Longitude, &req.Phone, &req.City, &req.State, &req.ZipCode)
	if err != nil {
		log.Println(err)
	}
	ds.QT.UpdateQuadTree([]models.Business{req})
	return req
}

// delete user
func (ds *proximityDBService) DeleteBusinessFromDB(req models.Business) models.Business {
	err := ds.DB.QueryRow(SELECTONEQUERY, req.ID).Scan(&req.ID, &req.Name, &req.Location.Longitude, &req.Location.Latitude, &req.Phone, &req.City, &req.State, &req.ZipCode)
	if err != nil {
		return models.Business{}
	} else {
		_, err := ds.DB.Exec(DELETEONEQUERY, req.ID)
		if err != nil {
			return models.Business{}
		}
	}
	ds.QT.DeleteFromQuadTree([]models.Business{req})
	return req
}

func (ds *proximityDBService) GetNearbyBusinessesFromQuadTree(req models.NearbySearchRequest) []models.Business {
	partialResponse := ds.QT.GetNearbyEntities(req)
	var finalResponse []models.Business
	for _, businessData := range partialResponse {
		iD := businessData.BusinessID
		singleBusiness := ds.GetBusinessFromDB(models.Business{ID: iD})
		singleBusiness.Dist = &businessData.Dist
		finalResponse = append(finalResponse, singleBusiness)
	}
	slices.SortFunc(finalResponse, func(a, b models.Business) int {
		if *a.Dist-*b.Dist < 0 {
			return -1
		} else if *a.Dist-*b.Dist > 0 {
			return 1
		} else {
			return 0
		}
	})
	return finalResponse
}
