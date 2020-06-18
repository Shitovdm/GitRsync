package controller

import (
	"github.com/Shitovdm/GitRsync/src/component/gui"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AboutController struct{}

func (ctrl AboutController) Index(c *gin.Context) {
	menu := gui.GetMenu(c)
	templateParams := gin.H{"menu": menu}
	templateParams["title"] = "About"
	c.HTML(http.StatusOK, "about/index", templateParams)
}
