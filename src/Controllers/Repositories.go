package Controllers

import (
	"github.com/Shitovdm/git-repo-exporter/src/Components/Helpers"
	"github.com/Shitovdm/git-repo-exporter/src/Components/Interface"
	"github.com/gin-gonic/gin"
	"net/http"
)

type RepositoriesController struct{}

func (ctrl RepositoriesController) Index(c *gin.Context) {
	menu := Interface.GetMenu(c)
	templateParams := gin.H{"menu": menu}
	templateParams["title"] = "Repositories"
	templateParams["repositories"] = Helpers.GetRepositoriesConfig()
	println(templateParams)

	c.HTML(http.StatusOK, "repositories/index", templateParams)
}

func (ctrl RepositoriesController) Add(c *gin.Context) {
	//Helpers.WsHandler(c.Writer, c.Request)
}

func (ctrl RepositoriesController) Remove(c *gin.Context) {
	//Helpers.WsHandler(c.Writer, c.Request)
}


