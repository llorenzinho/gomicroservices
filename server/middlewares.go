package server

import (
	"users/config"

	"github.com/gin-gonic/gin"
	ampq "github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
)

func Database(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(config.DATABASE_KEY, db)
		c.Next()
	}
}

func RabbitMQ(conn *ampq.Connection) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(config.RABBITMQ_KEY, conn)
		c.Next()
	}
}

func Config(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(config.CONFIG_KEY, cfg)
		c.Next()
	}
}
