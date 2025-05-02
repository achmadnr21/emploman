package domain

import (
	"time"
)

type Employee struct {
	ID           string    `json:"id" db:"id"`
	RoleID       string    `json:"role_id" db:"role_id"`
	NIP          string    `json:"nip" db:"nip"`
	Password     string    `json:"password" db:"password"`
	FullName     string    `json:"full_name" db:"full_name"`
	PlaceOfBirth string    `json:"place_of_birth" db:"place_of_birth"`
	DateOfBirth  time.Time `json:"date_of_birth" db:"date_of_birth"`
	Gender       string    `json:"gender" db:"gender"`
	PhoneNumber  string    `json:"phone_number" db:"phone_number"`
	PhotoURL     string    `json:"photo_url" db:"photo_url"`
	Address      string    `json:"address" db:"address"`
	NPWP         *string   `json:"npwp" db:"npwp"`
	GradeID      int       `json:"grade_id" db:"grade_id"`
	ReligionID   string    `json:"religion_id" db:"religion_id"`
	EchelonID    int       `json:"echelon_id" db:"echelon_id"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	ModifiedAt   time.Time `json:"modified_at" db:"modified_at"`
}

type EmployeeInterface interface {
	FindAll() ([]Employee, error)
	FindByID(id string) (*Employee, error)
	Save(employee *Employee) (*Employee, error)
	Update(employee *Employee) (*Employee, error)
	UploadProfileImage(id string, fileName string) (string, error)
	Delete(id string) error
	FindByNIP(nip string) (*Employee, error)
	FindByEmail(email string) (*Employee, error)
}
