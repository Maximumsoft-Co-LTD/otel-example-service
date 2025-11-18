package port

type Router interface {
	Start() error
	RegisterRoutes() error
}
