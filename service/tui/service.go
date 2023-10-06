package tui

import (
	"github.com/censync/soikawallet/config"
	"github.com/censync/soikawallet/service/tui/api_web3"
	"github.com/censync/soikawallet/service/tui/events"
	"sync"
)

var provider *Service

type IServiceProvider interface {
	Start() error
	Stop()
}

type Service struct {
	uiEvents              *events.EventBus
	w3Events              *events.EventBus
	web3ConnectionService *api_web3.Web3Connection
	tuiService            *Tui
}

func NewTUIServiceProvider(cfg *config.Config, wg *sync.WaitGroup) *Service {
	uiEvents := events.NewEventBus()
	w3Events := events.NewEventBus()
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
