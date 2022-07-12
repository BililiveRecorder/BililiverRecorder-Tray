package httpServer

import (
	"context"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func Main(assets fs.FS, port int, host string, ch chan os.Signal) {
	router := gin.Default()

	assetsFS := http.FS(assets)
	router.StaticFS("/", assetsFS)
	router.NoRoute(func(c *gin.Context) {
		c.Request.URL.Path = "/"
		router.HandleContext(c)
	})

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", host, port),
		Handler: router,
	}

	go func() {
		log.Info("Starting server on port ", port)
		err := srv.ListenAndServe()
		if err != nil {
			log.Error(err)
		}
	}()
	<-ch
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
