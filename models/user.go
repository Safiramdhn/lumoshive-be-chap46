package models

type User struct {
	ID       int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name     string `json:"name" gorm:"type:varchar(100);not null" binding:"required"`                    // Name must not be empty
	Email    string `json:"email" gorm:"type:varchar(100);uniqueIndex;not null" binding:"required,email"` // Email must not be empty and must be valid
	Password string `json:"-" gorm:"type:varchar(255);not null" binding:"required,min=8"`                 // Password must not be empty and have at least 8 characters
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`    // Email must not be empty and must be valid
	Password string `json:"password" binding:"required,min=8"` // Password must not be empty and have at least 8 characters
}

type LoginResponse struct {
	Token string `json:"token"`
}

func UserSeed() []User {
	return []User{
		{Name: "John Doe", Email: "john@example.com", Password: "password123"},
	}
}
