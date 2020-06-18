package controller

import (
	"github.com/Shitovdm/GitRsync/src/component/gui"
	"github.com/gin-gonic/gin"
	"net/http"
)

// DocsController struct describes docs section controller.
type DocsController struct{}

// Index describes docs index page.
func (ctrl DocsController) Index(c *gin.Context) {

	menu := gui.GetMenu(c)
	templateParams := gin.H{"menu": menu}
	templateParams["title"] = "Docs"

	c.HTML(http.StatusOK, "docs/index", templateParams)
}
