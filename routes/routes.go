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
	// TODO: remove
	pm.AddPluginWithVersion("blueocean", "1.25.5")
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

	// ============== plugin-manager routes ==============

	// get plugin manager data (plugins + coreversion)
	router.GET("/plugin-manager/get-data", handler.GetPluginsData)
	// adds new plugin to plugin manager
	router.POST("/plugin-manager/add-new-plugin", handler.AddNewPlugin)
	// changes core version in plugin manager
	router.POST("/plugin-manager/edit-core-version", handler.EditCoreVersion)
	// download formated file with plugin list from plugin manager
	router.GET("/plugin-manager/download-file", handler.DownloadFilePluginManager)
	// read plugin manifest data
	router.POST("/plugin-manager/get-manifest-attrs", handler.GetManifestAttrs)
	// delete plugin from plugin manager
	router.DELETE("/plugin-manager/delete-plugin", handler.DeletePlugin)
	// ============== END ^ plugin-manager routes ==============

	// ============== plugin-changes routes ==============
	router.GET("/plugin-manager/check-deps-with-update", handler.CheckDeps)
	router.GET("/plugin-manager/check-deps-without-update", func(ctx *gin.Context) {})

	router.GET("/plugin-manager/get-fixed-deps-diff", handler.GetVersionsDiff)
	router.POST("/plugin-manager/get-release-notes-diff", handler.GetReleaseNotesDiff)

	router.GET("/plugin-changes/download-file", handler.DownloadFilePluginChanges)

	router.POST("/plugin-manager/rescan", handler.RescanProjectNow)

	// router.GET("/plugin-manager/get-core-version", handler.GetCoreVersion)

	// router.POST("/plugin-manager/check-versions", func(ctx *gin.Context) {})
	// router.POST("/plugin-manager/resolve-deps", func(ctx *gin.Context) {})

	router.GET("/add-plugin-list/get-data", handler.GetPluginList)
	router.POST("/add-plugin-list/add-plugins", handler.AddPluginsFile)

	router.GET("/add-updated-plugins/get-data", handler.GetUpdatedPluginList)
	router.POST("/add-updated-plugins/edit-data", handler.AddUpdatedPluginList)

	// ==========================================================

	return router
}
