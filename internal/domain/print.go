package domain

import "time"

type PrintEmployee struct {
	NIP          string    `json:"nip"`
	FullName     string    `json:"full_name"`
	PlaceOfBirth string    `json:"place_of_birth"`
	Address      string    `json:"address"`
	DateOfBirth  time.Time `json:"date_of_birth"`
	Gender       string    `json:"gender"`
	Grade        string    `json:"golongan"`
	Echelon      string    `json:"echelon"`
	Jabatan      string    `json:"jabatan"`
	TempatTugas  string    `json:"tempat_tugas"`
	Religion     string    `json:"religion"`
	Unit         string    `json:"unit_kerja"`
	PhoneNumber  string    `json:"phone_number"`
	NPWP         string    `json:"npwp"`
	PhotoURL     string    `json:"photo_url"`
}
type PrintInterface interface {
	PrintAll() ([]PrintEmployee, error)
	PrintByNIP(id string) (*PrintEmployee, error)
	PrintByUnit(unitID int) ([]PrintEmployee, error)
}
