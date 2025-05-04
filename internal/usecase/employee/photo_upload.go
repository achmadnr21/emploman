package usecase_employee

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/achmadnr21/emploman/internal/utils"
	"github.com/disintegration/imaging"
)

func (eu *EmployeeUsecase) UploadPP(proposerId string, nip string, file *multipart.FileHeader) (string, error) {
	// cek proposer
	proposer, err := eu.empRepo.FindByID(proposerId)
	if err != nil {
		return "", &utils.NotFoundError{Message: "user not found"}
	}

	// cek role
	pr, err := eu.roleRepo.FindByID(proposer.RoleID)
	if err != nil {
		return "", &utils.NotFoundError{Message: "user role not found"}
	}

	// dapatkan employee dengan nip
	employee, err := eu.empRepo.FindByNIP(nip)
	if err != nil {
		return "", &utils.NotFoundError{Message: "employee not found"}
	}
	if employee.ID != proposer.ID && !(pr.CanAddEmployee) {
		return "", &utils.UnauthorizedError{Message: "user not authorized"}
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
	filename := fmt.Sprintf("%s/pp_%s.jpg", s3FolderPath, employee.ID)
	url, err := eu.s3Repo.UploadFile(filename, buf.Bytes(), "image/jpeg")
	if err != nil {
		return "", &utils.InternalServerError{Message: "failed to upload file to S3"}
	}

	// Update PhotoURL
	employee.PhotoURL = url
	newEmp, err := eu.empRepo.Update(employee)
	if err != nil {
		return "", &utils.InternalServerError{Message: "failed to update employee"}
	}

	return newEmp.PhotoURL, nil
}
