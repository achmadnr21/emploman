package main

import (
	"database/sql"
	"fmt"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"

	"github.com/achmadnr21/emploman/config"
	"github.com/achmadnr21/emploman/infrastructure/objectstorage"
	pgsql "github.com/achmadnr21/emploman/infrastructure/rdbms"
	"github.com/achmadnr21/emploman/internal/utils"
	gin_api "github.com/achmadnr21/emploman/service"

	"github.com/achmadnr21/emploman/internal/middleware"

	"github.com/achmadnr21/emploman/internal/repository"

	"github.com/achmadnr21/emploman/internal/usecase"
	emp "github.com/achmadnr21/emploman/internal/usecase/employee"

	"github.com/achmadnr21/emploman/internal/handler"
)

func main() {

	fmt.Println("========> Welcome to Emploman API")
	sc, err := service_init()
	if err != nil {
		fmt.Println("[Error] initializing service : ", err)
		panic("Service initialization failed")
	}
	var db *sql.DB = pgsql.GetPG()
	if db == nil {
		fmt.Println("[Error] Database connection is nil")
		panic("Database connection is nil")
	}
	defer db.Close()

	var s3client *s3.S3 = objectstorage.GetS3()
	if s3client == nil {
		fmt.Println("[Error] S3 client is nil")
		panic("S3 client is nil")
	}

	// Initialize the REST API

	var api gin_api.RESTapi
	var apiV *gin.RouterGroup = api.Init(gin.ReleaseMode)
	if apiV == nil {
		fmt.Println("[Error] API initialization failed")
		panic("API initialization failed")
	} else {
		fmt.Println("[Info] API initialized successfully")
	}

	// ========================= Dependency Injection =========================
	// Repository initialization
	roleRepo := repository.NewRoleRepository(db)
	// religionRepo := repository.NewReligionRepository(db)
	// gradeRepo := repository.NewGradeRepository(db)
	// echelonRepo := repository.NewEchelonRepository(db)
	employeeRepo := repository.NewEmployeeRepository(db)
	unitRepo := repository.NewUnitRepository(db)
	positionRepo := repository.NewPositionRepository(db)
	// employeeAssignmentRepo := repository.NewEmployeeAssignmentRepository(db)
	s3Repo := repository.NewS3Repository(s3client, sc.S3bucket)
	printRepo := repository.NewPrintRepository(db)

	// Usecase initialization
	authUsecase := usecase.NewAuthUsecase(employeeRepo, roleRepo)
	empUsecase := emp.NewEmployeeUsecase(employeeRepo, roleRepo, unitRepo, s3Repo)
	meUsecase := usecase.NewMeUsecase(employeeRepo, roleRepo, unitRepo, s3Repo)
	printUsecase := usecase.NewPrintUsecase(printRepo, employeeRepo, roleRepo, unitRepo)
	unitUsecase := usecase.NewUnitUsecase(unitRepo, roleRepo)
	positionUsecase := usecase.NewPositionUsecase(positionRepo, roleRepo)

	// Handler initialization
	authHandler := handler.NewAuthHandler(authUsecase)
	empHandler := handler.NewEmployeeHandler(empUsecase)
	meHandler := handler.NewMeHandler(meUsecase)
	printHandler := handler.NewPrintHandler(printUsecase)
	unitHandler := handler.NewUnitHandler(unitUsecase)
	positionHandler := handler.NewPositionHandler(positionUsecase)

	// ========================= API Routing =========================

	// 0. Create Simple open Ping endpoint
	apiV.GET("/ping", HandlePing)

	// 1. Auth
	auth := apiV.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.RefreshToken)
	}

	// 2. Employee Management (Admin)
	employee := apiV.Group("/employee")
	employee.Use(middleware.JWTAuthMiddleware)
	{
		// Basic CRUD
		employee.GET("", empHandler.GetAll)              // GET /employees
		employee.POST("", empHandler.Add)                // POST /employees
		employee.GET("/:nip", empHandler.GetByNIP)       // GET /employees/:nip
		employee.PUT("/:nip", empHandler.UpdateEmployee) // PUT /employees/:nip

		// Upload profile picture
		employee.POST("/:nip/profile-picture", empHandler.UploadPP) // POST /employees/:nip/profile-picture

		// Filtering
		employee.GET("/unit/:unit_id", empHandler.GetByUnit) // GET /employees/unit/:unit_id
		employee.GET("/search", empHandler.Search)           // GET /employees/search

		// Promotion
		employee.PUT("/:nip/promote", empHandler.Promote) // PUT /employees/:nip/promote

	}

	// 3. Me Management
	me := apiV.Group("/me")
	me.Use(middleware.JWTAuthMiddleware)
	{
		me.GET("", meHandler.GetMe)
		me.PUT("", meHandler.UpdateMe)
		me.POST("/profile-picture", meHandler.UploadPPMe)
	}

	// 4. Unit Management
	unit := apiV.Group("/unit")
	unit.Use(middleware.JWTAuthMiddleware)
	{
		unit.GET("", unitHandler.GetAllUnit) // GET /units
		unit.POST("", unitHandler.AddUnit)   // POST /units
		unit.GET("/:id", unitHandler.GetUnitByID)
		unit.PUT("/:id", unitHandler.UpdateUnit)
		unit.DELETE("/:id", unitHandler.DeleteUnit)
		unit.GET("/search", unitHandler.SearchUnit) // GET /units/search
	}
	// 5. position Management
	position := apiV.Group("/position")
	position.Use(middleware.JWTAuthMiddleware)
	{
		position.GET("", positionHandler.GetAllPosition) // GET /positions
		position.POST("", positionHandler.AddPosition)   // POST /positions
		position.GET("/:id", positionHandler.GetPositionByID)
		position.PUT("/:id", positionHandler.UpdatePosition)
		position.DELETE("/:id", positionHandler.DeletePosition)
		position.GET("/search", positionHandler.SearchPosition) // GET /positions/search
	}
	// 6. Religion Management
	// religion := apiV.Group("/religion")
	// religion.Use(middleware.JWTAuthMiddleware)
	// {
	// 	religion.GET("", religionHandler.GetAllReligion) // GET /religions
	// 	religion.POST("", religionHandler.AddReligion)   // POST /religions
	// 	religion.GET("/:id", religionHandler.GetReligionByID)
	// 	religion.PUT("/:id", religionHandler.UpdateReligion)
	// 	religion.DELETE("/:id", religionHandler.DeleteReligion)
	// 	religion.GET("/search", religionHandler.SearchReligion) // GET /religions/search
	// }
	// 7. Grade Management
	// grade := apiV.Group("/grade")
	// grade.Use(middleware.JWTAuthMiddleware)
	// {
	// 	grade.GET("", gradeHandler.GetAllGrade) // GET /grades
	// 	grade.POST("", gradeHandler.AddGrade)   // POST /grades
	// 	grade.GET("/:id", gradeHandler.GetGradeByID)
	// 	grade.PUT("/:id", gradeHandler.UpdateGrade)
	// 	grade.DELETE("/:id", gradeHandler.DeleteGrade)
	// 	grade.GET("/search", gradeHandler.SearchGrade) // GET /grades/search
	// }
	// 8. Echelon Management
	// echelon := apiV.Group("/echelon")
	// echelon.Use(middleware.JWTAuthMiddleware)
	// {
	// 	echelon.GET("", echelonHandler.GetAllEchelon) // GET /echelons
	// 	echelon.POST("", echelonHandler.AddEchelon)   // POST /echelons
	// 	echelon.GET("/:id", echelonHandler.GetEchelonByID)
	// 	echelon.PUT("/:id", echelonHandler.UpdateEchelon)
	// 	echelon.DELETE("/:id", echelonHandler.DeleteEchelon)
	// 	echelon.GET("/search", echelonHandler.SearchEchelon) // GET /echelons/search
	// }
	// // 9. Role Management
	// role := apiV.Group("/role")
	// role.Use(middleware.JWTAuthMiddleware)
	// {
	// 	role.GET("", authHandler.GetAllRole) // GET /roles
	// 	role.POST("", authHandler.AddRole)   // POST /roles
	// 	role.GET("/:id", authHandler.GetRoleByID)
	// 	role.PUT("/:id", authHandler.UpdateRole)
	// 	role.DELETE("/:id", authHandler.DeleteRole)
	// 	role.GET("/search", authHandler.SearchRole) // GET /roles/search
	// }

	// 100. Print Management
	print := apiV.Group("/print")
	print.Use(middleware.JWTAuthMiddleware)
	{
		printEmp := print.Group("/employee")
		{
			printEmp.GET("/:nip", printHandler.PrintByNIP)
			printEmp.GET("/unit/:unit_id", printHandler.PrintByUnitID)
			printEmp.GET("/all", printHandler.PrintAll)
		}
	}

	// ========================== Start HTTP API =========================
	service_config := fmt.Sprintf(":%d", sc.ServicePort)
	fmt.Printf("\nService running on port %s \n", service_config)
	api.Router.Run(service_config)

}

func HandlePing(c *gin.Context) {
	c.JSON(200, utils.ResponseSuccess("Pong", &struct {
		Developer string `json:"developer"`
	}{
		Developer: "Achmad Nashruddin Riskynanda",
	}))
}

func service_init() (config.Config, error) {
	var envload config.Config
	envload.LoadConfig()
	// database configuration
	err := pgsql.InitPG(envload.DbHost, int32(envload.DbPort), envload.DbUser, envload.DbPassword, envload.DbName, envload.DbSsl)
	if err != nil {
		return config.Config{}, fmt.Errorf("[Error] initializing PostgreSQL configuration : %v", err)
	}

	// s3 configuration
	err = objectstorage.InitS3(envload.S3accesskey, envload.S3secretkey, envload.S3endpoint, envload.S3bucket, envload.S3region)
	if err != nil {
		return config.Config{}, fmt.Errorf("[Error] initializing S3 configuration : %v", err)
	}

	// jwt configuration
	utils.JwtInit(envload.JwtSecret, envload.RefreshSecret)
	return envload, nil
}
