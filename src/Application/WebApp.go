package Application

import (
	rice "github.com/GeertJohan/go.rice"
	"github.com/Masterminds/sprig"
	"github.com/Shitovdm/git-repo-exporter/src/Components/Helpers"
	"github.com/Shitovdm/git-repo-exporter/src/Controllers"
	"github.com/foolin/goview"
	"github.com/foolin/goview/supports/ginview"
	"github.com/foolin/goview/supports/gorice"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"time"
)

var g errgroup.Group

func StartServer() {
	gin.SetMode(gin.ReleaseMode)
	serverWeb := &http.Server{
		Addr:         ":8888",
		Handler:      WebRouter(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	Helpers.OpenBrowser("http://localhost:8888")

	g.Go(func() error {
		return serverWeb.ListenAndServe()
	})

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}

}

func WebRouter() http.Handler {
	e := gin.New()
	riceBox := rice.MustFindBox("../../public/views")
	ginView := ginview.New(goview.Config{
		Root:         riceBox.Name(),
		Extension:    ".html",
		Master:       "layouts/master",
		Partials:     []string{},
		Funcs:        sprig.FuncMap(),
		DisableCache: true,
	})
	ginView.SetFileHandler(gorice.FileHandler(riceBox))
	e.HTMLRender = ginView

	staticBox := rice.MustFindBox("../../public/assets")
	staticFileServer := http.StripPrefix("/public/assets/", http.FileServer(staticBox.HTTPBox()))

	e.Use(gin.Recovery())
	e.Any("/public/assets/*filepath", gin.WrapF(staticFileServer.ServeHTTP))
	e.Handle(http.MethodGet, "/public/views", func(context *gin.Context) {
		staticFileServer.ServeHTTP(context.Writer, context.Request)
	})

	// Homepage
	indexController := new(Controllers.IndexController)
	e.GET("/", indexController.Index)

	// /repositories
	repositories := e.Group("/repositories")
	{
		repositoriesController := new(Controllers.RepositoriesController)
		repositories.GET("/", repositoriesController.Index)
		repositories.GET("add/", repositoriesController.Add)
		repositories.GET("edit/", repositoriesController.Edit)
		repositories.GET("remove/", repositoriesController.Remove)

	}

	// /platforms
	platforms := e.Group("/platforms")
	{
		platformsController := new(Controllers.PlatformsController)
		platforms.GET("/", platformsController.Index)
		platforms.GET("add/", platformsController.Add)
		platforms.GET("edit/", platformsController.Edit)
		platforms.GET("remove/", platformsController.Remove)
	}

	// /logs
	logs := e.Group("/logs")
	{
		logsController := new(Controllers.LogsController)
		logs.GET("/", logsController.Index)
		logs.GET("append/", logsController.Append)
	}

	// 404
	e.NoRoute(func(c *gin.Context) {
		c.HTML(404, "404.html", gin.H{})
	})

	return e
}
