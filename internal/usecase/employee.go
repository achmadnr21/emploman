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

type EmployeeUsecase struct {
	empRepo  domain.EmployeeInterface
	roleRepo domain.RoleInterface
	unitRepo domain.UnitInterface
	s3Repo   domain.S3Interface
}

func NewEmployeeUsecase(empRepo domain.EmployeeInterface, roleRepo domain.RoleInterface, unitRepo domain.UnitInterface, s3Repo domain.S3Interface) *EmployeeUsecase {
	return &EmployeeUsecase{
		empRepo:  empRepo,
		roleRepo: roleRepo,
		unitRepo: unitRepo,
		s3Repo:   s3Repo,
	}
}

func (eu *EmployeeUsecase) authorize(proposerId string, requireAdd bool) (*domain.Employee, *domain.Role, error) {
	user, err := eu.empRepo.FindByID(proposerId)
	if err != nil {
		return nil, nil, &utils.NotFoundError{Message: "user not found"}
	}
	role, err := eu.roleRepo.FindByID(user.RoleID)
	if err != nil {
		return nil, nil, &utils.NotFoundError{Message: "user role not found"}
	}
	if requireAdd && !role.CanAddEmployee {
		return nil, nil, &utils.UnauthorizedError{Message: "not authorized"}
	}
	if !requireAdd && !(role.CanAddEmployee || role.CanAssignEmployeeInternal || role.CanAssignEmployeeGlobal) {
		return nil, nil, &utils.UnauthorizedError{Message: "not authorized"}
	}
	return user, role, nil
}

func (eu *EmployeeUsecase) GetAll(proposerId string) ([]domain.Employee, error) {
	// cek proposer
	if _, _, err := eu.authorize(proposerId, false); err != nil {
		return nil, err
	}
	// get all employee
	employees, err := eu.empRepo.FindAll()
	if err != nil {
		return nil, &utils.InternalServerError{Message: "failed to get employees"}
	}
	// return employees
	return employees, nil
}

func (eu *EmployeeUsecase) GetByNIP(proposerId string, nip string) (*domain.Employee, error) {
	// cek proposer
	if _, _, err := eu.authorize(proposerId, false); err != nil {
		return nil, err
	}
	// get employee by nip
	employee, err := eu.empRepo.FindByNIP(nip)
	if err != nil {
		return nil, &utils.NotFoundError{Message: "employee not found"}
	}
	// return employee
	return employee, nil
}
func (eu *EmployeeUsecase) GetByUnit(proposerId string, unitId int) ([]domain.Employee, error) {
	// cek proposer
	if _, _, err := eu.authorize(proposerId, false); err != nil {
		return nil, err
	}
	// check wether unit exists
	unit, err := eu.unitRepo.FindByID(unitId)
	if err != nil {
		return nil, &utils.NotFoundError{Message: "unit not found"}
	}

	// get employee by unit
	employees, err := eu.empRepo.FindByUnit(unit.ID)
	if err != nil {
		return nil, &utils.NotFoundError{Message: "employee not found"}
	}
	// return employee
	return employees, nil
}

func (eu *EmployeeUsecase) Search(proposerId string, input string) ([]domain.Employee, error) {
	// cek proposer
	if _, _, err := eu.authorize(proposerId, false); err != nil {
		return nil, err
	}
	// get employee by input
	employees, err := eu.empRepo.Search(input)
	if err != nil {
		return nil, &utils.NotFoundError{Message: "employee not found"}
	}
	// return employee
	return employees, nil
}

func (eu *EmployeeUsecase) Add(proposerId string, employee *domain.Employee) (*domain.Employee, error) {
	// check employee.RoleID should be empty
	if employee.RoleID != "" {
		return nil, &utils.BadRequestError{Message: "Invalid Payload"}
	}
	employee.RoleID = "USR"
	if _, _, err := eu.authorize(proposerId, true); err != nil {
		return nil, err
	}
	// validate employee input
	if err := validateEmployeeInput(employee); err != nil {
		return nil, err
	}
	// check wether employee already exists
	existingEmployee, err := eu.empRepo.FindByNIP(employee.NIP)
	if err == nil && existingEmployee == nil {
		fmt.Printf("existing : %v\nerr: %v\n", existingEmployee, err)
		return nil, &utils.InternalServerError{Message: "failed to check employee"}
	}
	if existingEmployee != nil {
		return nil, &utils.ConflictError{Message: "employee already exists"}
	}

	// now, hash the password
	hashedPassword, err := utils.HashPassword(employee.Password)
	if err != nil {
		return nil, &utils.InternalServerError{Message: "failed to hash password"}
	}
	employee.Password = hashedPassword
	employee.PhotoURL = "https://s3.nevaobjects.id/emploman/pictureprofile/defaultprofile.jpg"
	// save employee
	newEmployee, err := eu.empRepo.Save(employee)
	newEmployee.Password = "" // clear password for security
	if err != nil {
		return nil, &utils.InternalServerError{Message: "failed to save employee"}
	}
	return newEmployee, nil
}
func (eu *EmployeeUsecase) UploadPP(proposerId string, nip string, file *multipart.FileHeader) (string, error) {
	// cek proposer
	proposer, err := eu.empRepo.FindByID(proposerId)
	if err != nil {
		return "", &utils.NotFoundError{Message: "user not found"}
	}

	// cek role
	pr, err := eu.roleRepo.FindByID(proposer.RoleID)
	if err != nil {
		fmt.Println("error role", err)
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
func (eu *EmployeeUsecase) UpdateEmployee(proposerId string, nip string, employee *domain.Employee) (*domain.Employee, error) {
	// cek proposer
	proposer, err := eu.empRepo.FindByID(proposerId)
	if err != nil || proposer == nil {
		return nil, &utils.NotFoundError{Message: "user not found"}
	}
	// cek role
	pr, err := eu.roleRepo.FindByID(proposer.RoleID)
	if err != nil {
		fmt.Println("error role", err)
		return nil, &utils.NotFoundError{Message: "user role not found"}
	}
	// check if user is authorized to update employee or is self update are allowed
	if !(pr.CanAddEmployee || pr.CanAssignEmployeeInternal || pr.CanAssignEmployeeGlobal) {
		return nil, &utils.UnauthorizedError{Message: "user not authorized"}
	}
	// check wether employee already exists
	existingEmployee, err := eu.empRepo.FindByNIP(nip)
	if err != nil || existingEmployee == nil {
		return nil, &utils.InternalServerError{Message: "failed to check employee or employee not found"}
	}

	if employee.RoleID != "" {
		// unauthorized in this endpoint
		return nil, &utils.UnauthorizedError{Message: "user not authorized"}
	}

	// full name checking
	if employee.FullName != "" && len(employee.FullName) > 3 && utils.IsAlpha(employee.FullName) {
		existingEmployee.FullName = employee.FullName
	}
	// place of birth checking
	if employee.PlaceOfBirth != "" && len(employee.PlaceOfBirth) > 3 && utils.IsAlpha(employee.PlaceOfBirth) {
		existingEmployee.PlaceOfBirth = employee.PlaceOfBirth
	}
	// date of birth checking jika is Zero atau tidak ada isinya
	if !employee.DateOfBirth.IsZero() {
		existingEmployee.DateOfBirth = employee.DateOfBirth
	}
	// Gender checking
	if employee.Gender != "" && len(employee.Gender) == 1 {
		existingEmployee.Gender = employee.Gender
	}
	// phone number checking
	if employee.PhoneNumber != "" && len(employee.PhoneNumber) > 6 && utils.IsNumeric(employee.PhoneNumber) {
		existingEmployee.PhoneNumber = employee.PhoneNumber
	}
	// address checking
	if employee.Address != "" && len(employee.Address) > 6 {
		existingEmployee.Address = employee.Address
	}
	// NPWP checking
	if employee.NPWP != nil && len(*employee.NPWP) == 16 {
		existingEmployee.NPWP = employee.NPWP
	}
	// grade id checking
	if employee.GradeID > 0 {
		existingEmployee.GradeID = employee.GradeID
	}
	// religion id checking
	if employee.ReligionID != "" && len(employee.ReligionID) == 3 {
		existingEmployee.ReligionID = employee.ReligionID
	}
	// echelon id checking
	if employee.EchelonID > 0 {
		existingEmployee.EchelonID = employee.EchelonID
	}

	// print the existing employee
	fmt.Println("updated existing employee : ", existingEmployee)

	newEmp, err := eu.empRepo.Update(existingEmployee)
	if err != nil {
		return nil, &utils.InternalServerError{Message: "failed to update employee"}
	}
	newEmp.Password = "" // clear password for security
	return newEmp, nil
}

func (eu *EmployeeUsecase) Promote(proposerId string, nip string, roleID string) (*domain.Employee, error) {
	// get proposer
	proposer, err := eu.empRepo.FindByID(proposerId)
	proposerRole := proposer.RoleID
	if err != nil {
		return nil, &utils.NotFoundError{Message: "user not found"}
	}
	// get employee
	employee, err := eu.empRepo.FindByNIP(nip)
	employeeRole := employee.RoleID
	if err != nil {
		return nil, &utils.NotFoundError{Message: "employee not found"}
	}
	if employee.RoleID == roleID {
		return nil, &utils.BadRequestError{Message: "employee already has this role"}
	}
	isValid := eu.hasValidPath(proposerRole, employeeRole, roleID)
	if !isValid {
		return nil, &utils.UnauthorizedError{Message: "user not authorized"}
	}
	employee.RoleID = roleID
	// save employee
	newEmployee, err := eu.empRepo.Update(employee)
	if err != nil {
		return nil, &utils.BadRequestError{Message: "role not found"}
	}
	newEmployee.Password = "" // clear password for security
	return newEmployee, nil
}

// ==================================================================== PRIVATE FUNCTIONS ====================================================================
func (eu *EmployeeUsecase) hasValidPath(proposerRole string, employeeRole string, roleID string) bool {
	// check using uc
	promoteList, err := eu.roleRepo.FindPromoteRole(proposerRole)
	fmt.Println("promote list : ", promoteList)
	if err != nil {
		return false
	}

	return validateRolePromotion(employeeRole, roleID, promoteList)
}

// ==================================================================== UTILITIES ====================================================================

func validateRolePromotion(currRole string, targetRole string, roleList []domain.RolePromotion) bool {
	var currExists bool = false
	var targetExists bool = false
	for _, role := range roleList {
		// curr exists on the list
		if currRole == role.FromRoleID || currRole == role.ToRoleID {
			currExists = true
		}
		if targetRole == role.FromRoleID || targetRole == role.ToRoleID {
			targetExists = true
		}

	}
	return currExists && targetExists
}
func validateEmployeeInput(e *domain.Employee) error {

	if e.RoleID == "" || len(e.RoleID) != 3 {
		return &utils.BadRequestError{Message: "role_id is required or invalid"}
	}
	if e.NIP == "" || len(e.NIP) != 18 || !utils.IsNumeric(e.NIP) {
		return &utils.BadRequestError{Message: "nip is required or invalid"}
	}
	if e.Password == "" || len(e.Password) < 8 {
		return &utils.BadRequestError{Message: "password is required or invalid"}
	}
	if e.FullName == "" || len(e.FullName) < 3 || !utils.IsAlpha(e.FullName) {
		return &utils.BadRequestError{Message: "full_name is required or invalid"}
	}
	if e.PlaceOfBirth == "" || len(e.PlaceOfBirth) < 3 || !utils.IsAlpha(e.PlaceOfBirth) {
		return &utils.BadRequestError{Message: "place_of_birth is required or invalid"}
	}
	if e.DateOfBirth.IsZero() {
		return &utils.BadRequestError{Message: "date_of_birth is required or invalid"}
	}
	if e.Gender == "" || len(e.Gender) != 1 {
		return &utils.BadRequestError{Message: "gender is required or invalid"}
	}
	if e.PhoneNumber == "" || len(e.PhoneNumber) < 6 || !utils.IsNumeric(e.PhoneNumber) {
		return &utils.BadRequestError{Message: "phone_number is required or invalid"}
	}
	if e.Address == "" || len(e.Address) < 6 {
		return &utils.BadRequestError{Message: "address is required or invalid"}
	}
	if e.GradeID <= 0 {
		return &utils.BadRequestError{Message: "grade_id is required or invalid"}
	}
	if e.ReligionID == "" || len(e.ReligionID) != 3 {
		return &utils.BadRequestError{Message: "religion_id is required or invalid"}
	}
	if e.EchelonID <= 0 {
		return &utils.BadRequestError{Message: "echelon_id is required or invalid"}
	}
	return nil
}
