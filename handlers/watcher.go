package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (s *ProjectService) GetWatcherList(c *gin.Context) {
	watcherList, err := s.Redis.GetWatcherData()
	if err != nil {
		logrus.Errorln("can not get watcher list")
		logrus.Errorln(err)
	}
	c.JSON(http.StatusOK, watcherList)
}
func (s *ProjectService) EditWatcherList(c *gin.Context) {
	logrus.Infoln("EditWatcherList route reached")

	var body map[string]string
	if err := c.BindJSON(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	logrus.Infof("Received request body: %+v\n", body)

	err := s.Redis.SetWatcherData(body)
	if err != nil {
		logrus.Errorln("can not set watcher list to DB")
		logrus.Errorln(err)
	}

	c.String(http.StatusOK, "EditWatcherList")
}
