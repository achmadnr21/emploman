package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	ServicePort   int
	DbHost        string
	DbPort        int
	DbName        string
	DbUser        string
	DbPassword    string
	DbSsl         string
	JwtSecret     string
	RefreshSecret string
	S3endpoint    string
	S3accesskey   string
	S3secretkey   string
	S3bucket      string
	S3region      string
	S3pathstyle   bool
}

func (c *Config) LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	c.ServicePort, _ = strconv.Atoi(os.Getenv("SERVICE_PORT"))
	c.DbHost = os.Getenv("DB_HOST")
	c.DbPort, _ = strconv.Atoi(os.Getenv("DB_PORT"))
	c.DbName = os.Getenv("DB_NAME")
	c.DbUser = os.Getenv("DB_USER")
	c.DbPassword = os.Getenv("DB_PASSWORD")
	c.DbSsl = os.Getenv("DB_SSL")
	c.JwtSecret = os.Getenv("JWT_SECRET")
	c.RefreshSecret = os.Getenv("REFRESH_SECRET")

	c.S3endpoint = os.Getenv("S3_ENDPOINT")
	c.S3accesskey = os.Getenv("S3_ACCESS_KEY_ID")
	c.S3secretkey = os.Getenv("S3_SECRET_ACCESS_KEY")
	c.S3bucket = os.Getenv("S3_BUCKET")
	c.S3region = os.Getenv("S3_REGION")
	c.S3pathstyle, err = strconv.ParseBool(os.Getenv("S3_USE_PATH_STYLE"))
	if err != nil {
		c.S3pathstyle = false
		log.Println("[Error]:", err)
	}
	if c.ServicePort == 0 {
		c.ServicePort = 8080
	}
}
