package controllers

import (
	"github.com/Shitovdm/GitRsync/src/components/interface"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AboutController struct{}

func (ctrl AboutController) Index(c *gin.Context) {
	menu := Interface.GetMenu(c)
	templateParams := gin.H{"menu": menu}
	templateParams["title"] = "About"
	c.HTML(http.StatusOK, "about/index", templateParams)
}
