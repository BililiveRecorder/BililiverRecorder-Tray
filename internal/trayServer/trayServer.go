package trayserver

import (
	"context"
	"fmt"
	"io/fs"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type TrayServer struct {
	Server   *http.Server
	Router   *gin.Engine
	Start    func()
	Shutdown func(*sync.WaitGroup)
}

func Create(assets fs.FS, port int, host string) *TrayServer {
	server := &TrayServer{}
	router := gin.Default()
	assetsFS := http.FS(assets)
	router.StaticFS("/", assetsFS)
	router.NoRoute(func(c *gin.Context) {
		c.Request.URL.Path = "/"
		router.HandleContext(c)
	})
	server.Server = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", host, port),
		Handler: router,
	}
	server.Router = router

	server.Start = func() {
		go func() {
			log.Info("Starting server on port ", port)
			err := server.Server.ListenAndServe()
			if err != nil {
				log.Error(err)
			}
		}()
	}
	server.Shutdown = func(wg *sync.WaitGroup) {
		wg.Add(1)
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := server.Server.Shutdown(ctx); err != nil {
				log.Fatal("Server Shutdown:", err)
			}
			log.Println("Server exiting")
			wg.Done()
		}()
	}
	return server
}
