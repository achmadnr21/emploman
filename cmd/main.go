package main

import (
	"database/sql"
	"fmt"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"

	"github.com/achmadnr21/emploman/config"
	"github.com/achmadnr21/emploman/infrastructure/objectstorage"
	pgsql "github.com/achmadnr21/emploman/infrastructure/rdbms"
	"github.com/achmadnr21/emploman/internal/handler"
	"github.com/achmadnr21/emploman/internal/middleware"
	"github.com/achmadnr21/emploman/internal/repository"
	"github.com/achmadnr21/emploman/internal/usecase"
	"github.com/achmadnr21/emploman/internal/utils"
	gin_api "github.com/achmadnr21/emploman/service"
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
	// positionRepo := repository.NewPositionRepository(db)
	// employeeAssignmentRepo := repository.NewEmployeeAssignmentRepository(db)
	s3Repo := repository.NewS3Repository(s3client, sc.S3bucket)

	// Usecase initialization
	authUsecase := usecase.NewAuthUsecase(employeeRepo, roleRepo)
	empUsecase := usecase.NewEmployeeUsecase(employeeRepo, roleRepo, unitRepo, s3Repo)
	// Handler initialization
	authHandler := handler.NewAuthHandler(authUsecase)
	empHandler := handler.NewEmployeeHandler(empUsecase)

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
		// CRUD
		employee.GET("", empHandler.GetAll)
		employee.POST("", empHandler.Add)
		employee.PUT("/:nip", empHandler.UpdateEmployee)
		employee.POST("/uploadpp/:nip", empHandler.UploadPP)
		employee.GET("/:nip", empHandler.GetByNIP)
		employee.GET("/unit/:unit_id", empHandler.GetByUnit)
		employee.GET("/search", empHandler.Search)
		employee.PUT("/uprole/:nip", empHandler.Promote)
		employee.PUT("/downrole/:nip", empHandler.Demote)
	}

	me := apiV.Group("/me")
	me.Use(middleware.JWTAuthMiddleware)
	{
		me.GET("", empHandler.GetMe)
		me.PUT("", empHandler.UpdateMe)
		me.POST("/uploadpp", empHandler.UploadPPMe)
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
