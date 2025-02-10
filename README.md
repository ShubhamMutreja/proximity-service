# Proximity Service (Quad Tree Based Location Indexing) (WIP)

Proximity Service is a high-performance location indexing system built using **Go**, **PostgreSQL**, and **Docker**. It efficiently stores, updates, and retrieves business locations using a **Quad Tree** data structure.

## 📌 Features

- ✅ **CRUD for Businesses** – Add, update, delete, and retrieve business locations  
- ✅ **Quad Tree Indexing** – Efficiently store and query spatial data
- ✅ **Nearby Search** – Fetch businesses closest to a user's location  
- ✅ **PostgreSQL Integration** – Store business data in a relational database  
- ✅ **Docker & Docker Compose** – Easily deploy with PostgreSQL  
- ✅ **Extensible Architecture** – UI integration planned for future updates  

## 🚀 Getting Started

### 📌 Prerequisites

Make sure you have the following installed:

- **Go** (≥1.22)  
- **Docker & Docker Compose**
- **PostgreSQL**

### 🛠️ Installation

```sh
git clone https://github.com/your-username/proximity-service.git
cd proximity-service
docker compose up -d --build .
```
This will:

- Start the Go API service
- Start the PostgreSQL database
- Create necessary tables automatically

To stop the service, run:
```sh
docker compose down
```

## 🔗 API Endpoints
- 📍 Business Management (CRUD)

| Method | Endpoint    |                 Description                  |
| :---:   | :---: |:--------------------------------------------:|
| GET | /businesses   |             List all businesses              |
| GET | /business/   |   Get a business by ID (?id=<business_id>)   |
| POST | /business/create   |            Create a new business             |
| PUT | /business/update   |         Update an existing business          |
| DELETE | /business/delete   | Delete a business by ID  (?id=<business_id>) |

### Request structure for POST, PUT
```json
{
    "name": "Biryani By handi",
    "location": {
        "longitude": 37.7749,
        "latitude": -122.4194
    },
    "phone": "+1-415-555-1234",
    "city": "San Francisco",
    "state": "CA",
    "zip_code": "94103"
}
```

- 📌 Nearby Search
  
| Method | Endpoint    | Description   |
| :---:  | :---: | :---: |
| GET | /search/nearby   | Get businesses near a user location (?latitude=<userLatitude>&id=<userLongitude&radius=<radius> ) |


## 🏗️ Technologies Used
- Golang – Backend API
- PostgreSQL – Database for storing businesses
- Docker & Docker Compose – Containerized deployment
- Quad Tree Service – Efficient spatial search

## 🎯 Future Enhancements
- ✅ UI for business registration and search
- ✅ Advanced caching mechanisms
- ✅ Realtime updating QuadTree for Movable business like food trucks or even Cabs/Taxis
