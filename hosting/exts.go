package hosting

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/picapica360/w3go/database/orm"
)

// AddService add a service to server.
func (h *Host) AddService(fn func()) {
	h.servicesFn = append(h.servicesFn, fn)
}

// UseMiddleware add a middleware to server pipe.
func (h *Host) UseMiddleware(fn func() gin.HandlerFunc) {
	h.middlewaresFn = append(h.middlewaresFn, fn)
}

// AddEndpoint add web endpoint.
func (h *Host) AddEndpoint(fn func(c Context)) {
	h.endpointFn = fn
}

// AddDatabase add a database to context.
func (h *Host) AddDatabase(name, conf string) {
	h.C.databases[name] = orm.NewDBString(conf)
}

// AddPProf only listening for pprof
// note: must import _ "net/http/pprof" package.
func (h *Host) AddPProf() {
	h.AddService(func() {
		go func() {
			port := h.Conf.App.PProfPort
			if port == 0 {
				port = defaultPProfPort
			}
			http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
		}()
	})
}

// AddHealth add health check api.
// note: check the database, redis, ES etc.
func (h *Host) AddHealth() {
	h.C.Router.GET("/health", func(c *gin.Context) {
		type errModel struct {
			Name string
			Err  error
		}
		var errs []errModel
		if h.C.DB != nil {
			if err1 := h.C.DB.DB().PingContext(context.TODO()); err1 != nil {
				errs = append(errs, errModel{"database", err1})
			}
		}

		if len(errs) > 0 {
			c.JSON(http.StatusBadRequest, errs)
		} else {
			c.JSON(http.StatusOK, nil)
		}
	})
}
