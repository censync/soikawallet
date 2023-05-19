package service

import (
	"github.com/censync/soikawallet/config"
	"github.com/censync/soikawallet/service/api_web3"
	"github.com/censync/soikawallet/service/internal/event_bus"
	"github.com/censync/soikawallet/service/ui"
	"sync"
)

var provider *ServiceProvider

type IServiceProvider interface {
	Start() error
	Stop()
}

type ServiceProvider struct {
	uiEvents              *event_bus.EventBus
	w3Events              *event_bus.EventBus
	web3ConnectionService *api_web3.Web3Connection
	tuiService            *ui.Tui
}

func NewServiceProvider(cfg *config.Config, wg *sync.WaitGroup) *ServiceProvider {
	uiEvents := event_bus.NewEventBus()
	w3Events := event_bus.NewEventBus()
	return &ServiceProvider{
		web3ConnectionService: api_web3.NewWeb3Connection(cfg, wg, uiEvents, w3Events),
		tuiService:            ui.NewTui(cfg, wg, uiEvents, w3Events),
	}
}

func (p *ServiceProvider) Web3Connection() IServiceProvider {
	return p.web3ConnectionService
}

func (p *ServiceProvider) UI() IServiceProvider {
	return p.tuiService
}
