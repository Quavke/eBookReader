package config

import (
	"ebookr/pkg/controllers"
	"ebookr/pkg/models"
	"ebookr/pkg/repositories"
	"ebookr/pkg/routers"
	"ebookr/pkg/services"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)
type Config struct {
		Server struct {
			Port int      `mapstructure:"PORT"`
		}											`mapstructure:"server"`
		DB struct {
			Host     string   `mapstructure:"HOST"`
			Port     int      `mapstructure:"PORT"`
			Name     string   `mapstructure:"NAME"`
			User     string   `mapstructure:"USER"`
			Password string   `mapstructure:"PASSWORD"`
			TimeZone   string   `mapstructure:"TIME_ZONE"`
			SSLMode    string   `mapstructure:"SSL_Mode"`
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


func NewApp(cfg *Config) (*App, error) {
	router := gin.Default()
	v1 := router.Group("/api/v1")
	// fmt.Printf("DB password: %v", cfg.DB.Password)
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		cfg.DB.Host, cfg.DB.User, cfg.DB.Password, cfg.DB.Name, cfg.DB.Port, cfg.DB.SSLMode, cfg.DB.TimeZone,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
  Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	sqlDB, err := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(1 * time.Hour)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	db.AutoMigrate(&models.Author{}, &models.Book{}, &models.User{})
	
	bookRepo := repositories.NewGormBookRepo(db)
	bookService := services.NewBookService(bookRepo)
	bookController := controllers.NewBookController(bookService)

	authorRepo := repositories.NewGormAuthorRepo(db)
	authorService := services.NewAuthorService(authorRepo)
	authorController := controllers.NewAuthorController(authorService)

	userRepo := repositories.NewGormUserRepo(db)
	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	routers.RegisterBookRoutes(v1, bookController)
	routers.RegisterAuthorRoutes(v1, authorController)
	routers.RegisterUserRoutes(v1, userController) // , middlewares.AuthMiddleware()
	return &App{
		router: router,
		cfg:    cfg,
	}, nil
}

func (a *App) Run() error {
	return a.router.Run(fmt.Sprintf(":%d", a.cfg.Server.Port))
}