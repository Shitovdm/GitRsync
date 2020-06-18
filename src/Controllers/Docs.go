package Controllers

import (
	"github.com/Shitovdm/GitRsync/src/Components/Interface"
	"github.com/gin-gonic/gin"
	"net/http"
)

type DocsController struct{}

func (ctrl DocsController) Index(c *gin.Context) {
	menu := Interface.GetMenu(c)
	templateParams := gin.H{"menu": menu}
	templateParams["title"] = "Docs"
	c.HTML(http.StatusOK, "docs/index", templateParams)
}
