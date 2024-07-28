package router

import (
	"messaggio/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type messageService interface {
	NewMessage(text string) error
	GetStatistic() (models.Statistic, error)
}
type router struct {
	ginRouter      *gin.Engine
	messageService messageService
}

func New(messageService messageService) router {
	router := router{
		ginRouter:      gin.New(),
		messageService: messageService,
	}
	router.ginRouter.POST("/messages", router.newMessage())
	router.ginRouter.GET("/messages/statistic", router.getStatistics())

	return router
}

func (r *router) StartServer() {
	r.ginRouter.Run(":8080")
}

func (r *router) newMessage() func(c *gin.Context) {
	return func(c *gin.Context) {
		var data struct {
			Text string `json:"text"`
		}

		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, "Bad request")
			return
		}
		err := r.messageService.NewMessage(data.Text)
		if err != nil {
			c.JSON(http.StatusInternalServerError, "Internal Server Error")
		}
		c.Status(http.StatusAccepted)
	}
}

func (r *router) getStatistics() func(c *gin.Context) {
	return func(c *gin.Context) {
		stat, err := r.messageService.GetStatistic()
		if err != nil {
			logrus.Error("Get statistic error", err)
		}
		c.JSON(http.StatusOK, stat)
	}

}
