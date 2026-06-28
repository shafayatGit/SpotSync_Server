# 🚗 SpotSync API

Smart Parking & EV Charging Reservation System

## Features
- JWT authentication
- bcrypt password hashing
- Driver and Admin roles
- Parking zone management
- EV charging reservation
- Reservation cancellation
- Dynamic availability calculation
- Concurrency-safe booking system using database transaction and row locking

## Tech Stack

- Go 1.22+
- Echo
- GORM
- PostgreSQL
- JWT
- Validator

## Architecture

SpotSync follows Clean Architecture:

Client
↓
Handler Layer
↓
Service Layer
↓
Repository Layer
↓
PostgreSQL Database

### Layers

DTO
- Request and response structures

Handler
- HTTP handling
- Validation
- JWT extraction

Service
- Business logic
- Authentication
- Reservation rules

Repository
- Database operations

Models
- Database tables

## Project Structure

spotsync-api/

cmd/
- main.go

models/
- user.go
- parking_zone.go
- reservation.go

dto/
- auth.dto.go
- zone.dto.go
- reservation.dto.go

handler/
- auth.handler.go
- zone.handler.go
- reservation.handler.go

service/
- auth.service.go
- zone.service.go
- reservation.service.go

repository/
- user.repository.go
- zone.repository.go
- reservation.repository.go

middleware/
- jwt.middleware.go

## Database Design

Users:
- id
- name
- email
- password
- role
- created_at
- updated_at

Parking Zones:
- id
- name
- type
- total_capacity
- price_per_hour
- created_at
- updated_at

Reservations:
- id
- user_id
- zone_id
- license_plate
- status
- created_at
- updated_at

## API Endpoints

Authentication:

POST /api/v1/auth/register

POST /api/v1/auth/login


Parking Zones:

POST /api/v1/zones
(Admin only)

GET /api/v1/zones

GET /api/v1/zones/:id


Reservations:

POST /api/v1/reservations

GET /api/v1/reservations/my-reservations

DELETE /api/v1/reservations/:id

GET /api/v1/reservations
(Admin only)


## Security

- JWT middleware
- Role based authorization
- Password encryption
- Input validation
- Protected routes


## Environment Variables

Create .env

PORT=8000

DATABASE_URL=

JWT_SECRET=


## Run Locally

Install dependencies:

go mod tidy

Run:

go run cmd/main.go


## Hot Reload

Install Air:

go install github.com/air-verse/air@latest

Run:

air


## Deployment

Backend:
- Render
- Railway
- Fly.io

Database:
- NeonDB
- Supabase
- Aiven


## Response Format

Success:

{
 "success": true,
 "message": "Operation successful",
 "data": {}
}


Error:

{
 "success": false,
 "message": "Error message"
}


## Author

Md. Shafayat Hossain Patowary
