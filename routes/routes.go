package routes

import (
	"net/http"

	"github.com/emelianrus/jenkins-release-notes-parser/handlers"
	"github.com/emelianrus/jenkins-release-notes-parser/storage/redisStorage"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(redis *redisStorage.RedisStorage) *gin.Engine {
	router := gin.Default()

	handler := handlers.ProjectService{
		Redis: redis,
	}

	// router.GET("/books", handlers.GetBooks)
	// router.GET("/books/:isbn", handlers.GetBookByISBN)
	// // router.DELETE("/books/:isbn", handlers.DeleteBookByISBN)
	// // router.PUT("/books/:isbn", handlers.UpdateBookByISBN)
	// router.POST("/books", handlers.PostBook)

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{http.MethodGet, http.MethodPatch, http.MethodPost, http.MethodHead, http.MethodDelete, http.MethodOptions},
		AllowHeaders:     []string{"Content-Type", "X-XSRF-TOKEN", "Accept", "Origin", "X-Requested-With", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// =================== routes ===================

	// =================== helpers ===================
	router.GET("/", handlers.RedirectToRoot)
	router.GET("/ping", handlers.Ping)
	router.GET("/redis/status", handler.RedisStatus)

	router.POST("/watcher-list", handler.EditWatcherList)
	router.GET("/watcher-list", handler.GetWatcherList)

	router.GET("/potential-updates", handler.GetPotentialUpdates)

	// GET all projects
	router.GET("/projects", handler.GetAllProjects)
	// https://api.github.com/repos/OWNER/REPO/releases
	router.GET("/project/:owner/:repo/releases", handler.GetProjectReleaseNotes)

	router.GET("/api/stats", handler.GetApiStats)

	return router
}
