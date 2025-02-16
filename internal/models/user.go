package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email    string         `gorm:"not null; unique; index"`
	FullName string         `gorm:"not null"`
	Password string         `gorm:"not null"`
	Journals []JournalEntry `gorm:"foreignKey:UserID"`
}

// --------------------------
// Dtos
// --------------------------
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	FullName string `json:"full_name" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}
