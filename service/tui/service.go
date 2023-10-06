package tui

import (
	"github.com/censync/soikawallet/config"
	"github.com/censync/soikawallet/service/tui/api_web3"
	"github.com/censync/soikawallet/types/event_bus"
	"sync"
)

var provider *Service

type IServiceProvider interface {
	Start() error
	Stop()
}

type Service struct {
	uiEvents              *event_bus.EventBus
	w3Events              *event_bus.EventBus
	web3ConnectionService *api_web3.Web3Connection
	tuiService            *Tui
}

func NewTUIServiceProvider(cfg *config.Config, wg *sync.WaitGroup) *Service {
	uiEvents := event_bus.NewEventBus()
	w3Events := event_bus.NewEventBus()
	return &Service{
		web3ConnectionService: api_web3.NewWeb3Connection(cfg, wg, uiEvents, w3Events),
		tuiService:            NewTui(cfg, wg, uiEvents, w3Events),
	}
}

func (p *Service) Web3Connection() IServiceProvider {
	return p.web3ConnectionService
}

func (p *Service) UI() IServiceProvider {
	return p.tuiService
}
