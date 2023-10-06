package service

import (
	"github.com/censync/soikawallet/config"
	"github.com/censync/soikawallet/service/api_web3"
	"github.com/censync/soikawallet/service/event_bus"
	"github.com/censync/soikawallet/service/tui"
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
	tuiService            *tui.Tui
}

func NewServiceProvider(cfg *config.Config, wg *sync.WaitGroup) *Service {
	uiEvents := event_bus.NewEventBus()
	w3Events := event_bus.NewEventBus()
	return &Service{
		web3ConnectionService: api_web3.NewWeb3Connection(cfg, wg, uiEvents, w3Events),
		tuiService:            tui.NewTui(cfg, wg, uiEvents, w3Events),
	}
}

func (p *Service) Web3Connection() IServiceProvider {
	return p.web3ConnectionService
}

func (p *Service) UI() IServiceProvider {
	return p.tuiService
}
