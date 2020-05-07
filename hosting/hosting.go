package hosting

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	"github.com/picapica360/w3go/config"
	"github.com/picapica360/w3go/config/env"
	"github.com/picapica360/w3go/logs"
	httpd "github.com/picapica360/w3go/net/http"
)

const (
	defaultHostPort  = 5000
	defaultPProfPort = 6060
)

var (
	cliEnv *string
)

func init() {
	cliEnv = flag.String("env", "", `set the service rumtime environment, like 'development','test' or 'production'`)
	flag.Parse()
}

// Host service host.
type Host struct {
	C Context // Context

	conf *config.AppConfig

	middlewareFn []func() gin.HandlerFunc
}

// Context host context, with Router and DbContext.
type Context struct {
	Router *gin.Engine
	DB     *gorm.DB
}

// ConfigureServices configure the services.
func (h *Host) ConfigureServices() {

}

// Configure configure the middleware.
func (h *Host) Configure(fn func(*Host)) {

}

// Run startup the service.
func (h *Host) Run() {
	env.SetEnv(*cliEnv) // override

	// config
	config.Init()
	h.conf = config.Conf()

	h.C.Router = httpd.Default()

	port := 1000
	if port == 0 {
		port = defaultHostPort
	}

	// h.C.Engine.Run(fmt.Sprintf(":%d", port))
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      h.C.Router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	// listen serve
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logs.Infof("[server] listen error %v", err)
		}
	}()

	shutdown(srv)
}

func shutdown(srv *http.Server) {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	sig := <-quit
	logs.Infof("[server] get a signal %s, stop the process", sig.String())

	fmt.Println("[server] shutdown server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil && err != http.ErrServerClosed {
		// note: do not call Fatal, because it will call os.Exit(),
		// 	and the 'defer func' (include caller) will not be executed.
		logs.Infof("[server] server shutdown error: %v", err)
	}
	fmt.Println("[server] server exited")
}
