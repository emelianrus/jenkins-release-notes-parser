package main

import (
	"html/template"
	"log"
	"net/http"
)

// func constractPageData() {
// 	// getJenkinsInstance
// 	// getListOfPlugins
// 	// checkIsPluginInCache
// 	// createStructure
// }

func StartWeb(redisclient *Redis) {
	log.Println("Starting server")
	data := getPlugins(redisclient)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Retrieve releases from Redis
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		tmpl := template.Must(template.ParseFiles("templates/release-notes.html"))
		err := tmpl.Execute(w, data)
		if err != nil {
			log.Println(err)
		}
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
