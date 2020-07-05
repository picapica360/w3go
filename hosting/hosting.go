package hosting

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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
	cliEnv  string
	cliPort int
)

func init() {
	flag.StringVar(&cliEnv, "env", "", `set the service rumtime environment, like 'development','test' or 'production'`)
	flag.IntVar(&cliPort, "port", 0, `set the net server run port`)

	flag.Parse()
}

// Host service host.
type Host struct {
	C Context // Context

	Conf *config.AppConfig

	endpointFn func(Context)

	servicesFn    []func()                 // services, lazy loading.
	middlewaresFn []func() gin.HandlerFunc // middlewares, lazy loading.
}

// Context host context, with Router and DbContext.
type Context struct {
	Router *gin.Engine
	DB     *gorm.DB

	databases map[string]*gorm.DB
}

// GetDB get the special database from context.
// It will return nil if not exists.
func (c *Context) GetDB(name string) *gorm.DB {
	if k, ok := c.databases[name]; ok {
		return k
	}
	return nil
}

// Run startup the service.
func (h *Host) Run() {
	env.SetEnv(cliEnv) // override

	// config
	config.Init()
	h.Conf = config.Conf()

	h.C.Router = httpd.Default()
	// TODO: set database.
	h.C.databases = make(map[string]*gorm.DB)
	defer func() {
		if len(h.C.databases) > 0 {
			for _, db := range h.C.databases {
				db.Close()
			}
		}
	}()
	// TODO: add all services.

	if len(h.middlewaresFn) > 0 {
		for _, fn := range h.middlewaresFn {
			h.C.Router.Use(fn())
		}
	}
	if len(h.servicesFn) > 0 {
		for _, fn := range h.servicesFn {
			fn()
		}
	}
	if h.endpointFn != nil {
		h.endpointFn(h.C)
	}

	// startup server.
	port := cliPort
	if port == 0 {
		// extract port from config. when not found in config, use default.
		if h.Conf.App.Port > 0 {
			port = h.Conf.App.Port
		} else {
			port = defaultHostPort
		}
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
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
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
