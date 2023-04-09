package handlers

import (
	"net/http"

	"github.com/emelianrus/jenkins-release-notes-parser/types"
	"github.com/gin-gonic/gin"
)

func AddProject(c *gin.Context) {
	c.JSON(http.StatusOK, "ASd")
}

func AddMultiplyProjects(c *gin.Context) {
	c.JSON(http.StatusOK, "ASd")
}

func GetProjectById(c *gin.Context) {
	projectName := c.DefaultQuery("name", "")
	c.String(http.StatusOK, "Hello %s", projectName)
}

func GetAllProjects(c *gin.Context) {

	releaseNotes := []types.GitHubReleaseNote{}

	releaseNotes = append(releaseNotes, types.GitHubReleaseNote{
		Name:    "1.1.1",
		TagName: "1.1.1T",
		Body: string(`<h2>🐛 Bug fixes</h2>
		<ul>
			<li>Fix publish release artifact GitHub action (<a class="issue-link js-issue-link" href="#">#526</a>)
			<a href="https://github.com/timja">@timja</a>
			</li>
		</ul>`),
		CreatedAt: "AUG 21",
	})

	releaseNotes = append(releaseNotes, types.GitHubReleaseNote{
		Name:    "2.2.2",
		TagName: "2.2.2T",
		Body: string(`<h2>🐛 Bug fixes222</h2>
		<ul>
			<li>abasdfasdfadf (<a class="issue-link js-issue-link" href="#">#526</a>)
			<a href="https://github.com/timja">@timja</a>
			</li>
		</ul>`),
		CreatedAt: "FEB 15",
	})
	releaseNotes = append(releaseNotes, types.GitHubReleaseNote{
		Name:    "3.3.3",
		TagName: "3.3.3T",
		Body: string(`<h2>🐛 Bug fixes333</h2>
		<ul>
			<li>abasdfasdfadf (<a class="issue-link js-issue-link" href="#">#526</a>)
			<a href="https://github.com/timja">@timja</a>
			</li>
		</ul>`),
		CreatedAt: "MAR 7",
	})

	// c.JSON(200, pessoas)

	c.JSON(http.StatusOK, releaseNotes)
}

func GetProjectsById(c *gin.Context) {
	c.JSON(http.StatusOK, "GetProjectsById")
}

func DeleteProject(c *gin.Context) {
	id := c.DefaultQuery("id", "")
	c.String(http.StatusOK, "Hello %s", id)
}

func DeleteMultiplyProjects(c *gin.Context) {
	var ids []string
	if err := c.BindJSON(&ids); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
		return
	}
	// Your logic to delete the items with the given IDs goes here
	c.JSON(http.StatusOK, gin.H{"message": "Deleted items with IDs", "ids": ids})
}
