package config

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Quavke/eBookReader/pkg/controllers"
	"github.com/Quavke/eBookReader/pkg/middlewares"
	"github.com/Quavke/eBookReader/pkg/models"
	"github.com/Quavke/eBookReader/pkg/repositories"
	"github.com/Quavke/eBookReader/pkg/routers"
	"github.com/Quavke/eBookReader/pkg/services"
	"github.com/redis/go-redis/v9"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)
type Config struct {
		IsProd bool 					`mapstructure:"IS_PROD"`
		Server struct {
			Port int      			`mapstructure:"PORT"`
		}											`mapstructure:"server"`
		DB struct {
			Host     	 string   `mapstructure:"HOST"`
			Port     	 int      `mapstructure:"PORT"`
			Name     	 string   `mapstructure:"NAME"`
			User     	 string   `mapstructure:"USER"`
			TimeZone   string   `mapstructure:"TIME_ZONE"`
			SSLMode    string   `mapstructure:"SSL_MODE"`
			DSN        string
		}   									`mapstructure:"db"`
		Redis struct {
			Host string 				`mapstructure:"HOST"`
			Port int 						`mapstructure:"PORT"`
		}											`mapstructure:"redis"`		
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
	if cfg.DB.Host == "" || cfg.DB.User == "" || cfg.DB.Name == "" {
    return nil,fmt.Errorf("db host/user/name are required")
  }
	if gin.Mode() == gin.ReleaseMode {
		cfg.IsProd = true
	}
	return &cfg, nil
}


func NewApp(cfg *Config) (*App, error) {
	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Set("isProd", cfg.IsProd)
		c.Next()
	})
	err := godotenv.Load()
	if err != nil{
		return nil, err
	}
	v1 := router.Group("/api/v1")
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		cfg.DB.Host, cfg.DB.User, os.Getenv("DB_PASSWORD"), cfg.DB.Name, cfg.DB.Port, cfg.DB.SSLMode, cfg.DB.TimeZone,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
  Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}
	sqlDB, err := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(1 * time.Hour)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}
	db.AutoMigrate(&models.Author{}, &models.Book{}, &models.UserDB{})

	context := context.Background()

	addr := fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port)
	if cfg.Redis.Host == "" || cfg.Redis.Port == 0 {
		return nil, fmt.Errorf("redis host/port are required")
	}

	client := redis.NewClient(&redis.Options{
        Addr:	  addr,
        Password: os.Getenv("REDIS_PASSWORD"),
        DB:		  0,
        Protocol: 2,
				ReadTimeout: 5 * time.Second,
				WriteTimeout: 5 * time.Second,
  })
	
	bookRepo := repositories.NewGormBookRepo(db)
	bookService := services.NewBookService(bookRepo, context, client)
	bookController := controllers.NewBookController(bookService)
	
	authorRepo := repositories.NewGormAuthorRepo(db)
	authorService := services.NewAuthorService(authorRepo, context, client)
	authorController := controllers.NewAuthorController(authorService)

	userRepo := repositories.NewGormUserRepo(db)
	userService := services.NewUserService(userRepo, context, client)
	userController := controllers.NewUserController(userService)

	AuthMiddleware := middlewares.AuthMiddleware(userRepo)
	BooksMiddleware := middlewares.BooksMiddleware(userRepo)

	routers.RegisterBookRoutes(v1, bookController, AuthMiddleware, BooksMiddleware)
	routers.RegisterAuthorRoutes(v1, authorController, AuthMiddleware)
	routers.RegisterUserRoutes(v1, userController, AuthMiddleware)
	return &App{
		router: router,
		cfg:    cfg,
	}, nil
}

func (a *App) Run() error {
	return a.router.Run(fmt.Sprintf(":%d", a.cfg.Server.Port))
}