package product

import "github.com/americanas-go/config"

const (
	root      = "pkg.domain.repository.product"
	dataStore = root + ".dataStore"
)

func init() {
	config.Add(dataStore, "mock", "product datastore mock/xpto")
}

func DataStore() string {
	return config.String(dataStore)
}
