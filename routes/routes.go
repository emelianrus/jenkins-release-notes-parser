package routes

import (
	"net/http"

	"github.com/emelianrus/jenkins-release-notes-parser/db"
	"github.com/emelianrus/jenkins-release-notes-parser/handlers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(db *db.Redis) *gin.Engine {
	router := gin.Default()

	handler := handlers.ProjectService{
		Redis: db,
	}

	// router.GET("/books", handlers.GetBooks)
	// router.GET("/books/:isbn", handlers.GetBookByISBN)
	// // router.DELETE("/books/:isbn", handlers.DeleteBookByISBN)
	// // router.PUT("/books/:isbn", handlers.UpdateBookByISBN)
	// router.POST("/books", handlers.PostBook)

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{http.MethodGet, http.MethodPatch, http.MethodPost, http.MethodHead, http.MethodDelete, http.MethodOptions},
		AllowHeaders:     []string{"Content-Type", "X-XSRF-TOKEN", "Accept", "Origin", "X-Requested-With", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	router.GET("/", handlers.RedirectToRoot)
	router.GET("/ping", handlers.Ping)

	// https://api.github.com/repos/OWNER/REPO/releases
	router.GET("/project/:owner/:repo/releases", handler.GetProjectReleaseNotes)

	// GET single project
	router.GET("/project", handler.GetProjectById)
	// POST single projet
	router.POST("/project", handler.GetAllProjects)
	// DELETE single project
	router.DELETE("/project", handler.DeleteProject)

	// TODO: should it be separate functions? or cast everything to list?
	// GET all projects
	router.GET("/projects", handler.GetAllProjects)
	router.POST("/projects", handler.GetAllProjects)
	// DELETE multiple items by ID
	router.DELETE("/projects", handler.DeleteMultiplyProjects)

	return router
}
