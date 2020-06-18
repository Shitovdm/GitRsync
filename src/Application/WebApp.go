package application

import (
	rice "github.com/GeertJohan/go.rice"
	"github.com/Shitovdm/GitRsync/src/components/helpers"
	"github.com/Shitovdm/GitRsync/src/components/logger"
	"github.com/Shitovdm/GitRsync/src/controllers"
	"github.com/foolin/goview"
	"github.com/foolin/goview/supports/ginview"
	"github.com/foolin/goview/supports/gorice"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	"html/template"
	"net/http"
	"strings"
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

	helpers.OpenBrowser("http://localhost:8888")

	g.Go(func() error {
		return serverWeb.ListenAndServe()
	})

	if err := g.Wait(); err != nil {
		logger.Error("StartServer", err.Error())
	}
}

func WebRouter() http.Handler {
	e := gin.New()
	riceBox := rice.MustFindBox("../../public/views")
	ginView := ginview.New(goview.Config{
		Root:      riceBox.Name(),
		Extension: ".html",
		Master:    "layouts/master",
		Partials:  []string{},
		Funcs: template.FuncMap{
			"replace": strings.Replace,
			"add":     helpers.Add,
		},
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

	indexController := new(controllers.IndexController)
	e.GET("/", indexController.Index)

	actions := e.Group("/actions")
	{
		actionsController := new(controllers.ActionsController)
		actions.GET("pull/", actionsController.Pull)
		actions.GET("push/", actionsController.Push)
		actions.GET("clear/", actionsController.Clear)
		actions.GET("info/", actionsController.Info)
		actions.GET("block/", actionsController.Block)
		actions.GET("activate/", actionsController.Activate)
	}

	repositories := e.Group("/repositories")
	{
		repositoriesController := new(controllers.RepositoriesController)
		repositories.GET("/", repositoriesController.Index)
		repositories.GET("add/", repositoriesController.Add)
		repositories.GET("edit/", repositoriesController.Edit)
		repositories.GET("remove/", repositoriesController.Remove)
	}

	platforms := e.Group("/platforms")
	{
		platformsController := new(controllers.PlatformsController)
		platforms.GET("/", platformsController.Index)
		platforms.GET("add/", platformsController.Add)
		platforms.GET("edit/", platformsController.Edit)
		platforms.GET("remove/", platformsController.Remove)
	}

	logs := e.Group("/logs")
	{
		logsController := new(controllers.LogsController)
		logs.GET("/", logsController.Index)
		logs.GET("subscribe/", logsController.Subscribe)
		logs.GET("remove/runtime/", logsController.RemoveRuntime)
		logs.GET("remove/all/", logsController.RemoveAll)
	}

	settings := e.Group("/settings")
	{
		settingsController := new(controllers.SettingsController)
		settings.GET("/", settingsController.Index)
		settings.GET("save/", settingsController.Save)
	}

	docs := e.Group("/docs")
	{
		docsController := new(controllers.DocsController)
		docs.GET("/", docsController.Index)
	}

	about := e.Group("/about")
	{
		aboutController := new(controllers.AboutController)
		about.GET("/", aboutController.Index)
	}

	e.NoRoute(func(c *gin.Context) {
		c.HTML(404, "404.html", gin.H{})
	})

	return e
}
