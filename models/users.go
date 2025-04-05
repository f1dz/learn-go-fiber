package models

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" gorm:"unique" validate:"required,email"`
	Password string `json:"-" validate:"required,password"`
}
