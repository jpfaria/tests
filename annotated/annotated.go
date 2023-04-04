package annotated

import (
	"context"
	"net/http"
)

type Loren struct {
	Name string
}

// LorenMethod Lorem ipsum dolor sit amet, consectetur adipiscing elit
// @A App xpto
// @A HandlerType HTTP
// @A Path /foo
// @A Path /
// @A Method POST
// @A Consume application/json
// @A Consume application/yaml
// @A Produce application/json
// @A Param query foo bool true tiam sed efficitur purus
// @A Param query bar string true tiam sed efficitur purus
// @A Param path foo string tiam sed efficitur purus
// @A Param path bar string tiam sed efficitur purus
// @A Param header foo string true tiam sed efficitur purus
// @A Param header bar string true tiam sed efficitur purus
// @A Body github.com/americanas-go/inject/examples/simple.Request
// @A Response 201 github.com/americanas-go/inject/examples/simple.Response tiam sed efficitur purus, at lacinia magna
// @A provider index=0 name=xpto
//
// ctx lorem inpsum
func LorenMethod(ctx context.Context, req *http.Request) (resp *Loren, err error) {
	return &Loren{
		Name: "Hello world",
	}, nil
}
