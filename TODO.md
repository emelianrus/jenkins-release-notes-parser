* set UpdatedPlugins from UI to be able get release notes for predefined plugins
* add get txt file for main plugin manager page
* add tests for plugin-manager
* release notes not working


* replace redis by plugin manager
router.GET("/watcher-list", handler.GetPluginList)
router.GET("/add-plugin-list/get-data", handler.GetPluginList)