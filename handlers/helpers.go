package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func RedirectToRoot(c *gin.Context) {
	c.Redirect(http.StatusMovedPermanently, "/servers")
}

func (s *ProjectService) RedisStatus(c *gin.Context) {
	err := s.Redis.Status()

	if err == nil {
		c.JSON(http.StatusOK, "")
	} else {
		c.JSON(http.StatusOK, err.Error())
	}
}
