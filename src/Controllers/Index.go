package Controllers

import (
	"github.com/Shitovdm/git-repo-exporter/src/Components/Helpers"
	"github.com/Shitovdm/git-repo-exporter/src/Components/Interface"
	"github.com/gin-gonic/gin"
	"net/http"
)

type IndexController struct{}

func (ctrl IndexController) Index(c *gin.Context) {

	// Check we have config set
	_, err := Helpers.GetAppConfig()
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, "/setup")
	}
	// Start working on projects
	menu := Interface.GetMenu(c)
	templateParams := gin.H{"menu": menu}
	templateParams["title"] = "Dashboard"

	repositoriesConfig, _ := Helpers.GetRepositoriesConfig()

	templateParams["repositories"] = repositoriesConfig

	c.HTML(http.StatusOK, "index/index", templateParams)
}
