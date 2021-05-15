package hello

import (
	"context"

	"github.com/soumitra003/goframework/config"
)

type ModuleHello struct {
	config *config.Config
}

//New creates module instance
func New(config config.Config) *ModuleHello {
	md := &ModuleHello{config: &config}
	return md
}

func (h *ModuleHello) Init(ctx context.Context, config config.Config) {

}
