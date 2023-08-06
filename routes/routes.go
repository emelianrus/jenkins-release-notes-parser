package routes

import (
	"net/http"

	"github.com/emelianrus/jenkins-release-notes-parser/handlers"
	"github.com/emelianrus/jenkins-release-notes-parser/pkg/pluginManager"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// func preloadPluginManager(pm *pluginManager.PluginManager, redis *redisStorage.RedisStorage) {
// 	watcherData, _ := redis.GetPluginListData()
// 	for name, version := range watcherData {
// 		pm.AddPluginWithVersion(name, version)
// 	}
// }

func SetupRouter() *gin.Engine {
	router := gin.Default()

	pm := pluginManager.NewPluginManager()

	// preloadPluginManager(&pm, redis)

	handler := handlers.ProjectService{
		// Redis:         redis,
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
	router.GET("/watcher-list", handler.GetPluginList)

	router.GET("/api/stats", handler.GetApiStats)

	// TODO: ============== plugin-manager routes ==============

	// should return plugins + versions + core version
	router.GET("/plugin-manager/get-data", handler.GetPluginsData)

	router.POST("/plugin-manager/add-new-plugin", handler.AddNewPlugin)
	router.DELETE("/plugin-manager/delete-plugin", handler.DeletePlugin)
	router.GET("/plugin-manager/check-deps", handler.CheckDeps)
	router.GET("/plugin-manager/get-fixed-deps-diff", handler.GetFixedDepsDiff)

	router.GET("/plugin-changes/download-file", handler.DownloadFilePluginChanges)
	router.GET("/plugin-manager/download-file", handler.DownloadFilePluginManager)

	router.POST("/plugin-manager/rescan", handler.RescanProjectNow)
	router.POST("/plugin-manager/get-manifest-attrs", handler.GetManifestAttrs)

	router.POST("/plugin-manager/edit-core-version", handler.EditCoreVersion)
	router.GET("/plugin-manager/get-core-version", handler.GetCoreVersion)

	router.POST("/plugin-manager/check-versions", func(ctx *gin.Context) {})
	router.POST("/plugin-manager/resolve-deps", func(ctx *gin.Context) {})

	router.GET("/add-plugin-list/get-data", handler.GetPluginList)
	router.POST("/add-plugin-list/add-plugins", handler.AddPluginsFile)

	router.GET("/add-updated-plugins/get-data", handler.GetUpdatedPluginList)
	router.POST("/add-updated-plugins/edit-data", handler.AddUpdatedPluginList)

	// ==========================================================

	return router
}
