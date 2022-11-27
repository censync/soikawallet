package state

import (
	"github.com/censync/go-i18n"
	"github.com/censync/soikawallet/service"
	"github.com/censync/soikawallet/service/ui/handler"
	"github.com/censync/soikawallet/service/ui/widgets/extpages"
	"github.com/censync/soikawallet/service/wallet"
)

const (
	ModeIdle = iota
	ModeWithAirGap
	ModeWithoutAirGap

	StatusIdle = iota
	StateInitAirGap
	StateInitLocal
	StateAwaitQR
)

type State struct {
	events     *handler.TBus
	walletMode uint8
	status     uint8

	isInitialised bool
	//wallet *wallet.Wallet
	tr    *i18n.Translator
	pages *extpages.ExtPages
}

func InitState(events *handler.TBus, tr *i18n.Translator) *State {
	return &State{
		events:     events,
		walletMode: ModeIdle,
		status:     StatusIdle,
		tr:         tr,
	}
}

func (s *State) SetWallet(wallet *wallet.Wallet) {
	s.isInitialised = true
	s.events.Emit(handler.EventUpdatedWallet, "xxxx-xxxxx-xxxx") // GetInstanceId()
}

func (s *State) Emit(event handler.EventType, data interface{}) {
	s.events.Emit(event, data)
}

func (s *State) WalletMode() uint8 {
	return s.walletMode
}

func (s *State) Status() uint8 {
	return s.status
}

func (s *State) API() service.WalletAdapter {
	return service.API()
}

func (s *State) Tr() *i18n.Translator {
	return s.tr
}

func (s *State) SetStatus(status uint8) {
	s.status = status
}

func (s *State) SetWalletMode(mode uint8) {
	s.walletMode = mode
}

// Pages
func (s *State) SetPages(pages *extpages.ExtPages) {
	s.pages = pages
}

func (s *State) Pages() *extpages.ExtPages {
	return s.pages
}

func (s *State) SwitchToPage(page string, args ...interface{}) {
	s.pages.SwitchToPage(page, args)
}

// Current page

func (s *State) Params() []interface{} {
	return s.pages.Current().Params()
}
