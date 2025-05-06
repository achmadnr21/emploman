package usecase

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/achmadnr21/emploman/internal/domain"
	"github.com/achmadnr21/emploman/internal/utils"
	"github.com/disintegration/imaging"
)

type MeUsecase struct {
	empRepo  domain.EmployeeInterface
	roleRepo domain.RoleInterface
	unitRepo domain.UnitInterface
	s3Repo   domain.S3Interface
}

func NewMeUsecase(empRepo domain.EmployeeInterface, roleRepo domain.RoleInterface, unitRepo domain.UnitInterface, s3Repo domain.S3Interface) *MeUsecase {
	return &MeUsecase{
		empRepo:  empRepo,
		roleRepo: roleRepo,
		unitRepo: unitRepo,
		s3Repo:   s3Repo,
	}
}

func (eu *MeUsecase) GetMe(proposerId string) (*domain.Employee, error) {
	// cek proposer
	proposer, err := eu.empRepo.FindByID(proposerId)
	if err != nil {
		return nil, &utils.NotFoundError{Message: "user not found"}
	}
	proposer.Password = "" // clear password for security
	// return employee
	return proposer, nil
}

func (eu *MeUsecase) UpdateMe(proposerId string, employee *domain.Employee) (*domain.Employee, error) {
	// cek proposer
	proposer, err := eu.empRepo.FindByID(proposerId)
	if err != nil || proposer == nil {
		return nil, &utils.NotFoundError{Message: "user not found"}
	}
	if employee.RoleID != "" {
		// unauthorized in this endpoint
		return nil, &utils.UnauthorizedError{Message: "user not authorized"}
	}
	// UpdateMe hanya membolehkan update data diri, yaitu:
	// Phone Number, Address, Religion
	// full name checking
	if employee.FullName != "" || len(employee.FullName) > 3 || utils.IsAlpha(employee.FullName) {
		return nil, &utils.UnauthorizedError{Message: "user not authorized to change full name"}
	}
	// place of birth checking
	if employee.PlaceOfBirth != "" || len(employee.PlaceOfBirth) > 3 || utils.IsAlpha(employee.PlaceOfBirth) {
		return nil, &utils.UnauthorizedError{Message: "user not authorized to change place of birth"}
	}
	// date of birth checking jika is Zero atau tidak ada isinya
	if !employee.DateOfBirth.IsZero() {
		return nil, &utils.UnauthorizedError{Message: "user not authorized to change date of birth"}
	}
	// Gender checking
	if employee.Gender != "" && len(employee.Gender) == 1 {
		return nil, &utils.UnauthorizedError{Message: "user not authorized to change gender"}
	}
	// phone number checking
	if employee.PhoneNumber != "" && len(employee.PhoneNumber) > 6 && utils.IsNumeric(employee.PhoneNumber) {
		proposer.PhoneNumber = employee.PhoneNumber
	}
	// address checking
	if employee.Address != "" && len(employee.Address) > 6 {
		proposer.Address = employee.Address
	}
	// NPWP checking
	if employee.NPWP != "" && len(employee.NPWP) == 16 {
		proposer.NPWP = employee.NPWP
	}
	// grade id checking
	if employee.GradeID > 0 {
		return nil, &utils.UnauthorizedError{Message: "user not authorized to change grade"}
	}
	// religion id checking
	if employee.ReligionID != "" && len(employee.ReligionID) == 3 {
		proposer.ReligionID = employee.ReligionID
	}
	// echelon id checking
	if employee.EchelonID > 0 {
		return nil, &utils.UnauthorizedError{Message: "user not authorized to change echelon"}
	}

	newEmp, err := eu.empRepo.Update(proposer)
	if err != nil {
		return nil, &utils.InternalServerError{Message: "failed to update employee data"}
	}
	newEmp.Password = "" // clear password for security
	return newEmp, nil
}

func (eu *MeUsecase) UploadPPMe(proposerId string, file *multipart.FileHeader) (string, error) {
	// cek proposer
	proposer, err := eu.empRepo.FindByID(proposerId)
	if err != nil {
		return "", &utils.NotFoundError{Message: "user not found"}
	}

	// Buka file dari form
	src, err := file.Open()
	if err != nil {
		return "", &utils.BadRequestError{Message: "failed to open uploaded file"}
	}
	defer src.Close()

	// Ekstensi file
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		return "", &utils.BadRequestError{Message: "only .jpg, .jpeg, or .png files are allowed"}
	}

	// Decode gambar
	img, _, err := image.Decode(src)
	if err != nil {
		return "", &utils.BadRequestError{Message: "invalid image format"}
	}

	// Resize ke 200x200
	resizedImg := imaging.Resize(img, 200, 200, imaging.Lanczos)

	// Encode ke JPG
	var buf bytes.Buffer
	err = jpeg.Encode(&buf, resizedImg, &jpeg.Options{Quality: 90})
	if err != nil {
		return "", &utils.InternalServerError{Message: "failed to encode image"}
	}

	// Upload ke S3
	s3FolderPath := "pictureprofile"
	filename := fmt.Sprintf("%s/pp_%s.jpg", s3FolderPath, proposer.ID)
	url, err := eu.s3Repo.UploadFile(filename, buf.Bytes(), "image/jpeg")
	if err != nil {
		return "", &utils.InternalServerError{Message: "failed to upload file to S3"}
	}

	// Update PhotoURL
	proposer.PhotoURL = url
	newEmp, err := eu.empRepo.Update(proposer)
	if err != nil {
		return "", &utils.InternalServerError{Message: "failed to update employee"}
	}

	return newEmp.PhotoURL, nil
}
