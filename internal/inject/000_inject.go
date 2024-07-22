package inject

import "github.com/samber/do/v2"

var Injector = AddProvider()

func AddProvider() *do.RootScope {
	var i = do.New()

	// add your di code
	// do.Provide(i, gateway.SomeGateway)
	// do.Provide(i, repository.SomeRepository)

	return i
}
