package Controllers

import (
	"github.com/Shitovdm/git-rsync/src/Components/Interface"
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