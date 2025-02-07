package database

import (
	"database/sql"
	"log"
	"proximityService/models"
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
	QT interface{}
	DB *sql.DB
}

var dataStore proximityDBService

type ProximityDBService interface {
	GetAllBusinessesFromDB() []models.Business
	GetBusinessFromDB(req models.Business) models.Business
	PublishNewBusinessToDB(req models.Business) models.Business
	UpdateBusinessInDB(req models.Business) models.Business
	DeleteBusinessFromDB(req models.Business) models.Business
	GetNearbyBusinessesFromQuadTree(req models.NearbySearchRequest) []models.Business
}

func GetDataStore() proximityDBService {
	return dataStore
}

var (
	UNAMEDB string = "postgres"
	PASSDB  string = "postgres123"
	HOSTDB  string = "postgres"
	DBNAME  string = "businessdata"
)

func InitQuadTree() interface{} {
	return nil
}
func InitDataBase() *sql.DB {
	//connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", UNAMEDB, PASSDB, HOSTDB, DBNAME)
	connStr := "postgres://postgres:postgres123@localhost:5432/db?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	//defer db.Close()

	//create the table if it doesn't exist
	_, err = db.Exec(CREATETABLE)

	if err != nil {
		log.Fatal(err)
	}
	log.Print("database running")
	return db
}
func NewDBService(qT interface{}, dB *sql.DB) ProximityDBService {
	dataStore = proximityDBService{
		QT: qT,
		DB: dB,
	}
	return &dataStore
}

// Database operations
// get all users
func (db *proximityDBService) GetAllBusinessesFromDB() []models.Business {
	rows, err := db.DB.Query(SELECTALLQUERY)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	businesses := []models.Business{}
	for rows.Next() {
		var u models.Business
		if err := rows.Scan(&u.ID, &u.Name, &u.Location.Latitude, &u.Location.Longitude, &u.Phone, &u.City, &u.State, &u.ZipCode); err != nil {
			log.Println(err)
		}
		businesses = append(businesses, u)
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
	}
	return businesses
}

// get user by id
func (db *proximityDBService) GetBusinessFromDB(req models.Business) models.Business {
	var u models.Business
	err := db.DB.QueryRow(SELECTONEQUERY, req.ID).Scan(&u.ID, &u.Name, &u.Location.Latitude, &u.Location.Longitude, &u.Phone, &u.City, &u.State, &u.ZipCode)
	if err != nil {
		log.Println(err)
	}
	return u
}

// create user
func (db *proximityDBService) PublishNewBusinessToDB(req models.Business) models.Business {
	err := db.DB.QueryRow(INSERTQUERY, req.Name, req.Location.Latitude, req.Location.Longitude, req.Phone, req.City, req.State, req.ZipCode).Scan(&req.ID)
	if err != nil {
		log.Println(err)
	}
	return req
}

// update user
func (db *proximityDBService) UpdateBusinessInDB(req models.Business) models.Business {
	_, err := db.DB.Exec(UPDATEQUERY, &req.ID, &req.Name, &req.Location.Latitude, &req.Location.Longitude, &req.Phone, &req.City, &req.State, &req.ZipCode)
	if err != nil {
		log.Println(err)
	}
	return req
}

// delete user
func (db *proximityDBService) DeleteBusinessFromDB(req models.Business) models.Business {
	err := db.DB.QueryRow(SELECTONEQUERY, req.ID).Scan(&req.ID, &req.Name, &req.Location.Latitude, &req.Location.Longitude, &req.Phone, &req.City, &req.State, &req.ZipCode)
	if err != nil {
		return models.Business{}
	} else {
		_, err := db.DB.Exec(DELETEONEQUERY, req.ID)
		if err != nil {
			return models.Business{}
		}
	}
	return req
}

// QuadTree Opeations
// get nearby businesses
func (qt *proximityDBService) GetNearbyBusinessesFromQuadTree(req models.NearbySearchRequest) []models.Business {
	return []models.Business{}
}
