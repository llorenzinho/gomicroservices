package server

import (
	"fmt"
	"users/config"
	"users/controllers"

	"github.com/gin-gonic/gin"
	ampq "github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
)

func getServices() (*config.Config, *gorm.DB, *ampq.Connection) {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
	db, err := config.CreateDatabase(&cfg.Database)
	if err != nil {
		panic(err)
	}
	config.AutoMigrate(db)

	rabbit, err := config.CreateRabbitMQClient(&cfg.RabbitMQConfig)
	if err != nil {
		panic(fmt.Errorf("failed to connect to RabbitMQ: %w", err))
	}

	return cfg, db, rabbit
}

func Run() {
	// Load configuration
	cfg, db, rabbit := getServices()
	defer func() {
		if err := rabbit.Close(); err != nil {
			fmt.Printf("Failed to close RabbitMQ connection: %v\n", err)
		}
	}()

	server := gin.Default()
	server.GET("/", func(c *gin.Context) {
		name := c.Query("name")
		if name == "" {
			name = "World"
		}
		c.JSON(200, gin.H{
			"message": fmt.Sprintf("Hello, %s!", name),
		})
	})

	userController := controllers.NewUserController(db, rabbit, cfg)
	jwtController := controllers.NewJwtController(db, cfg)

	userController.Bind(server)
	jwtController.Bind(server)

	if err := server.Run(fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)); err != nil {
		panic(err)
	}
}
