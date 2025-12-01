package router

import "github.com/codecrafters-io/shell-starter-go/internal/handlers"

type Handler interface {
	Run([]string) (string, error)
}

type Routes map[string]Handler

type Router struct {
	routes Routes
}

func New() Router {
	return Router{
		routes: make(Routes),
	}
}

func (r Router) Handle(command string, hf Handler) {
	r.routes[command] = hf
}

func (r Router) Run(op string, context []string) (string, error) {
	h, ok := r.routes[op]
	if !ok {
		eh := handlers.NewExternal(op)
		return eh.Run(context)
	}

	return h.Run(context)
}
