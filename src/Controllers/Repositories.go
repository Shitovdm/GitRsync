package Controllers

import (
	"encoding/json"
	"fmt"
	"github.com/Shitovdm/git-repo-exporter/src/Components/Helpers"
	"github.com/Shitovdm/git-repo-exporter/src/Components/Interface"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type RepositoriesController struct{}

var WebSocketUpgrade = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (ctrl RepositoriesController) Index(c *gin.Context) {
	menu := Interface.GetMenu(c)
	templateParams := gin.H{"menu": menu}
	templateParams["title"] = "Repositories"

	_, conferr := Helpers.GetAppConfig()
	if conferr != nil {
		c.Redirect(http.StatusTemporaryRedirect, "/setup/")
	}

	c.HTML(http.StatusOK, "repositories/index", templateParams)
}

func (ctrl RepositoriesController) Add(c *gin.Context) {
	wsHandler(c.Writer, c.Request)
}

func (ctrl RepositoriesController) Remove(c *gin.Context) {
	wsHandler(c.Writer, c.Request)
}

// Actual startAll form processing
func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := WebSocketUpgrade.Upgrade(w, r, nil)
	if err != nil {
		log.Println(fmt.Sprintf("failed to set websocket upgrade: %+v", err))
		return
	}

	_, message, _ := conn.ReadMessage()

	var formData map[string]string
	err = json.Unmarshal(message, &formData)
	if err != nil {
		log.Println("Failed to decode json from form posted")
	}

	byteData, _ := json.Marshal(formData)
	log.Println(string(byteData))

	return
}
