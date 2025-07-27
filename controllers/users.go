package controllers

import (
	"errors"
	"time"
	"users/config"
	"users/models"
	"users/utils"

	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	ampq "github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
)

type UserController struct {
	DB       *gorm.DB
	RabbitMQ *ampq.Connection
	Config   *config.Config
}

func NewUserController(db *gorm.DB, rabbitMQ *ampq.Connection, cfg *config.Config) *UserController {
	return &UserController{
		DB:       db,
		RabbitMQ: rabbitMQ,
		Config:   cfg,
	}
}

func (uc *UserController) publishUserCreatedEvent(user *models.User) error {
	ch, err := uc.RabbitMQ.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	queue, err := ch.QueueDeclare(
		uc.Config.RabbitMQConfig.UserCreateQueue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	body, err := json.Marshal(user)
	if err != nil {
		return err
	}

	err = ch.Publish(
		"",         // exchange
		queue.Name, // routing key
		false,      // mandatory
		false,      // immediate
		ampq.Publishing{
			DeliveryMode: ampq.Persistent,
			ContentType:  "application/json",
			Body:         []byte(body),
			Timestamp:    time.Now(),
			Type:         "user.created",
		},
	)
	return err
}

func (uc *UserController) publishUserUpdatedEvent(user *models.User) error {
	ch, err := uc.RabbitMQ.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	queue, err := ch.QueueDeclare(
		uc.Config.RabbitMQConfig.UserUpdateQueue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	body, err := json.Marshal(user)
	if err != nil {
		return err
	}

	err = ch.Publish(
		"",         // exchange
		queue.Name, // routing key
		false,      // mandatory
		false,      // immediate
		ampq.Publishing{
			DeliveryMode: ampq.Persistent,
			ContentType:  "application/json",
			Body:         []byte(body),
			Timestamp:    time.Now(),
			Type:         "user.updated",
		},
	)
	return err
}

func (uc *UserController) CreateUser(c *gin.Context) {
	var body models.CreateUserRequest
	if err := c.ShouldBindBodyWithJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	encrypted, err := utils.EncryptString(body.Password)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to encrypt password"})
		return
	}

	user := models.User{
		Username: body.Username,
		Email:    body.Email,
		Password: encrypted,
	}
	tx := uc.DB.Create(&user)
	if tx.Error != nil {
		var pgErr *pgconn.PgError
		if errors.As(tx.Error, &pgErr) {
			if pgErr.Code == "23505" { // Unique violation
				c.JSON(409, gin.H{"error": "Username or email already exists"})
				return
			}
		}
		c.JSON(500, gin.H{"error": tx.Error.Error()})
		return
	}
	c.JSON(201, user)
	uc.publishUserCreatedEvent(&user)
}

func (uc *UserController) UpdateUser(c *gin.Context) {
	claim, code, err := utils.ValidateJwtHelper(c, uc.Config.JWT)
	if err != nil {
		c.JSON(code, gin.H{"error": err.Error()})
		return
	}
	var body models.UpdateUserRequest
	if err := c.ShouldBindBodyWithJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	tx := uc.DB.Model(&models.User{Id: claim.UserID}).Updates(models.User{
		Username: body.Username,
		Email:    body.Email,
	})
	if tx.Error != nil {
		var pgErr *pgconn.PgError
		if errors.As(tx.Error, &pgErr) {
			if pgErr.Code == "23505" { // Unique violation
				c.JSON(409, gin.H{"error": "Username or email already exists"})
				return
			}
		}
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			c.JSON(404, gin.H{"error": "User not found"})
			return
		}
		// Handle other errors
		c.JSON(500, gin.H{"error": tx.Error.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "User updated successfully"})
	uc.publishUserUpdatedEvent(&models.User{
		Id:       claim.UserID,
		Username: body.Username,
		Email:    body.Email,
	})
}

func (uc *UserController) Me(c *gin.Context) {
	claim, code, err := utils.ValidateJwtHelper(c, uc.Config.JWT)
	if err != nil {
		c.JSON(code, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	tx := uc.DB.First(&user, claim.UserID)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			c.JSON(404, gin.H{"error": "User not found"})
			return
		}
		c.JSON(500, gin.H{"error": tx.Error.Error()})
		return
	}
	c.JSON(200, user)
}

func (uc *UserController) createUserRouter(server *gin.Engine) *gin.RouterGroup {
	api := server.Group("/api")
	users := api.Group("/users")
	return users
}

func (uc *UserController) Bind(server *gin.Engine) {
	router := uc.createUserRouter(server)
	router.POST("", uc.CreateUser)
	router.PUT("", uc.UpdateUser)
	router.GET("/me", uc.Me)
}
