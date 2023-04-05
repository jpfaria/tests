package finder

import "github.com/americanas-go/config"

const (
	root     = "pkg.domain.service.product.finder"
	provider = root + ".provider"
)

func init() {
	config.Add(provider, "standard", "product finder service standard/xpto")
}

func Provider() string {
	return config.String(provider)
}
