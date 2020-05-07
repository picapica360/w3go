package http

import (
	"github.com/gin-gonic/gin"
	"github.com/picapica360/w3go/config/env"
	"github.com/picapica360/w3go/net/http/validator"
)

// Options options.
type Options struct {
	Handlers gin.HandlersChain // middlewares
}

// New a gin engine.
func New(opt *Options) (engine *gin.Engine) {
	engine = makeEngine()
	engine.Use(builtinMiddleware()...) // add builtin middleware
	engine.Use(opt.Handlers...)
	return
}

// Default create gin engine, with middleware.
func Default() (engine *gin.Engine) {
	engine = makeEngine()
	engine.Use(builtinMiddleware()...) // add builtin middleware
	return
}

func makeEngine() (engine *gin.Engine) {
	if env.IsProduction() {
		gin.SetMode("release") // debug or release
	}
	engine = gin.New()
	validator.Register() // register custom validators.
	return
}

func builtinMiddleware() []gin.HandlerFunc {
	var handlers []gin.HandlerFunc
	handlers = append(handlers, gin.Recovery())

	return handlers
}
