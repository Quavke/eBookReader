package config

import (
	"ebookr/pkg/controllers"
	"ebookr/pkg/models"
	"ebookr/pkg/repositories"
	"ebookr/pkg/routers"
	"ebookr/pkg/services"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)
type Config struct {
		Server struct {
			ServerPort int      `mapstructure:"SERVER_PORT"`
		}											`mapstructure:"server"`
		DB struct {
			DBHost     string   `mapstructure:"DB_HOST"`
			DBPort     int      `mapstructure:"DB_PORT"`
			DBName     string   `mapstructure:"DB_NAME"`
			DBUser     string   `mapstructure:"DB_USER"`
			DBPassword string   `mapstructure:"DB_PASSWORD"`
			TimeZone   string   `mapstructure:"TIME_ZONE"`
			SSLMode    string   `mapstructure:"SSLMode"`
			DSN        string
		}   									  `mapstructure:"db"`
		// JWT struct {
		// 	secretKey 	[]byte `mapstructure:"SECRET_KEY"`
		// } `mapstructure:"jwt"`
}
type App struct {
	router *gin.Engine
	cfg    *Config
}

func NewConfig() (*Config, error) {
	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config: %w", err)
		}
	}

	v.AutomaticEnv()

	var cfg Config

	if err := v.Unmarshal(&cfg); err != nil { // , viper.DecodeHook(hook)
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}


func NewApp(cfg *Config) *App {
	router := gin.Default()
	fmt.Printf("DB password: %v", cfg.DB.DBPassword)
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		cfg.DB.DBHost, cfg.DB.DBUser, cfg.DB.DBPassword, cfg.DB.DBName, cfg.DB.DBPort, cfg.DB.SSLMode, cfg.DB.TimeZone,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
  Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		errorMsg := fmt.Sprintf("Cannot create database. Err: %v", err.Error())
		log.Fatal(errorMsg)
	}
	db.AutoMigrate(&models.Author{}, &models.Book{}, &models.User{})
	
	bookRepo := repositories.NewGormBookRepo(db)
	bookService := services.NewBookService(bookRepo)
	bookController := controllers.NewBookController(bookService)

	// userRepo := repositories.NewGormUserRepo(db)
	// userService := services.NewUserService(userRepo)
	// userController := controllers.NewUserController(userService)

	routers.RegisterBookRoutes(router, bookController)
	// routers.RegisterUserRoutes(router, userController, middlewares.AuthMiddleware(cfg.JWT.secretKey))
	
	return &App{
		router: router,
		cfg:    cfg,
	}
}

func (a *App) Run() error {
	return a.router.Run(fmt.Sprintf(":%d", a.cfg.Server.ServerPort))
}