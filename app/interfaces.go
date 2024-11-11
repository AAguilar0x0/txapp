package app

type Initializer func(services ServiceProvider) (Lifecycle, error)

type Lifecycle interface {
	Run()
	Close()
}
