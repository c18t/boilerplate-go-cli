package inject

import (
	"sync"

	"github.com/samber/do/v2"
)

var (
	injector *do.RootScope
	once     sync.Once
)

func GetInjector() *do.RootScope {
	once.Do(func() {
		injector = AddProvider()
	})
	return injector
}

func AddProvider() *do.RootScope {
	var i = do.New()

	// add your di code
	// do.Provide(i, gateway.SomeGateway)
	// do.Provide(i, repository.SomeRepository)

	return i
}
