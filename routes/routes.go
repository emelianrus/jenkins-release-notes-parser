package routes

import (
	"net/http"

	"github.com/emelianrus/jenkins-release-notes-parser/handlers"
	"github.com/emelianrus/jenkins-release-notes-parser/pkg/pluginManager"
	"github.com/emelianrus/jenkins-release-notes-parser/storage/redisStorage"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func preloadPluginManager(pm *pluginManager.PluginManager, redis *redisStorage.RedisStorage) {

	watcherData, _ := redis.GetWatcherData()
	for name, version := range watcherData {
		pm.AddPluginWithVersion(name, version)
		// pm.AddPlugin(pluginManager.NewPluginWithVersion(name, version))
	}
}

func SetupRouter(redis *redisStorage.RedisStorage) *gin.Engine {
	router := gin.Default()

	pm := pluginManager.NewPluginManager()
	preloadPluginManager(&pm, redis)

	// pm.FixPluginDependencies()

	// for _, v := range pm {
	// 	fmt.Println(v.Name)
	// }

	handler := handlers.ProjectService{
		Redis:         redis,
		PluginManager: pm,
	}

	router.Use(cors.New(cors.Config{
		// TODO: should not be "*"
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{http.MethodGet, http.MethodPatch, http.MethodPost, http.MethodHead, http.MethodDelete, http.MethodOptions},
		AllowHeaders:     []string{"Content-Type", "X-XSRF-TOKEN", "Accept", "Origin", "X-Requested-With", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// =================== helpers ===================
	router.GET("/", handlers.RedirectToRoot)
	router.GET("/ping", handlers.Ping)
	router.GET("/redis/status", handler.RedisStatus)

	// =================== routes ===================

	router.POST("/watcher-list", handler.EditWatcherList)
	router.GET("/watcher-list", handler.GetWatcherList)

	router.GET("/potential-updates", handler.GetPotentialUpdates)

	// GET all projects
	router.GET("/projects", handler.GetAllProjects)
	// https://api.github.com/repos/OWNER/REPO/releases
	router.GET("/project/:owner/:repo/releases", handler.GetProjectReleaseNotes)

	router.GET("/api/stats", handler.GetApiStats)

	// TODO: ============== plugin-manager routes ==============

	// should return plugins + versions + core version
	router.GET("/plugin-manager/get-data", func(c *gin.Context) {
		type pluginManagerData struct {
			Plugins     map[string]*pluginManager.Plugin
			CoreVersion string
		}

		data := pluginManagerData{
			Plugins:     pm.GetPlugins(),
			CoreVersion: pm.GetCoreVersion(),
		}
		c.JSON(http.StatusOK, data)
	})

	router.POST("/plugin-manager/add-new-plugin", handler.AddNewPlugin)
	router.DELETE("/plugin-manager/delete-plugin", handler.DeletePlugin)

	router.POST("/plugin-manager/rescan", handler.RescanProjectNow)

	router.POST("/plugin-manager/edit-core-version", handler.EditCoreVersion)
	router.GET("/plugin-manager/get-core-version", handler.GetCoreVersion)

	router.POST("/plugin-manager/check-versions", func(ctx *gin.Context) {})
	router.POST("/plugin-manager/resolve-deps", func(ctx *gin.Context) {})
	// ==========================================================

	return router
}
