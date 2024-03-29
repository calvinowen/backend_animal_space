package models

import (
	"time"

	"gorm.io/gorm"
)

type Animal struct {
	gorm.Model
	Name        string
	Image       string
	Gender      string
	Type        string
	Age         int
	Month       int
	Description string
	Quantity    int
	Status      string
	Price       int
	CategoryID  int
	Category    Category `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ShelterID   int
	Shelter     Shelter `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type GetAllAnimal struct {
	Id              int
	Name            string
	Gender          string
	Type            string
	Age             int
	Month           int
	Description     string
	Quantity        int
	Price           int
	Status          string
	Image           string
	Category_id     int
	Category_name   string
	Shelter_id      int
	Shelter_name    string
	Shelter_Phone   string
	Shelter_Address string
	Created_at      time.Time
}
