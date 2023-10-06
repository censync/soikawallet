package tui

import (
	"fmt"
	"github.com/censync/go-i18n"
	"github.com/censync/soikawallet/api/dto"
	"github.com/censync/soikawallet/config"
	"github.com/censync/soikawallet/config/dict"
	"github.com/censync/soikawallet/config/version"
	"github.com/censync/soikawallet/service/tui/tmainframe"
	"github.com/censync/soikawallet/service/tui/twidget/statusview"
	"github.com/censync/soikawallet/service/wallet"
	"github.com/censync/soikawallet/service/wallet/protected_key"
	"github.com/censync/soikawallet/types/event_bus"
	"github.com/censync/tview"
	"github.com/gdamore/tcell/v2"
	"sync"
)

type Tui struct {
	app *tview.Application
	tr  *i18n.Translator

	mainFrame     *tmainframe.TMainFrame
	uiEvents      *event_bus.EventBus
	w3Events      *event_bus.EventBus
	layout        *tview.Flex
	style         *tview.Theme
	isVerboseMode bool
	wg            *sync.WaitGroup
	stopped       bool
}

func NewTui(cfg *config.Config, wg *sync.WaitGroup, uiEvents, w3Events *event_bus.EventBus) *Tui {
	/*style := &tview.Theme{
		PrimitiveBackgroundColor:    tcell.ColorLightYellow,
		ContrastBackgroundColor:     tcell.ColorOrange,
		MoreContrastBackgroundColor: tcell.ColorGreen,
		BorderColor:                 tcell.ColorOrchid,
		TitleColor:                  tcell.ColorDarkOrange,
		GraphicsColor:               tcell.ColorDarkOrange,
		PrimaryTextColor:            tcell.ColorDarkOrange,
		SecondaryTextColor:          tcell.ColorOrangeRed,
		TertiaryTextColor:           tcell.ColorGreen,
		InverseTextColor:            tcell.ColorBlue,
		ContrastSecondaryTextColor:  tcell.ColorDarkBlue,
		DisabledBackgroundColor:     tcell.ColorDarkSlateGray,
		DisabledTextColor:           tcell.ColorLightGray,
	}*/

	style := &tview.Styles

	tui := &Tui{
		app:      tview.NewApplication(),
		tr:       dict.GetTr("en"),
		uiEvents: uiEvents,
		w3Events: w3Events,
		style:    style,
		wg:       wg,
	}
	tui.mainFrame = tmainframe.Init(tui.uiEvents, tui.w3Events, dict.GetTr("en"), tui.style)
	tui.layout = tui.initLayout()
	return tui
}
func (t *Tui) initLayout() *tview.Flex {

	labelTitle := tview.NewTextView().
		SetDynamicColors(true).
		SetWrap(false).
		SetTextAlign(tview.AlignRight).
		SetText(fmt.Sprintf("[darkcyan]Soika Wallet[black] v%s", version.VERSION))

	labelTitle.SetBackgroundColor(tcell.ColorDarkGrey).
		SetBorderPadding(0, 0, 0, 2)

	// Label instance id
	labelInstanceId := tview.NewTextView().
		SetDynamicColors(true).
		SetWrap(false).
		SetTextAlign(tview.AlignLeft).
		SetText(`[darkcyan]ID:[black] not initialized`)

	labelInstanceId.SetBackgroundColor(tcell.ColorDarkGrey).
		SetBorderPadding(0, 0, 2, 0)

	// Label notice
	labelNotice := tview.NewTextView().
		SetDynamicColors(true).
		SetWrap(false).
		SetTextAlign(tview.AlignCenter)

	// TODO: Move to api
	if ok, err := protected_key.IsMemoryProtected(); !ok {
		t.uiEvents.Emit(event_bus.EventWalletNoticeMessage, fmt.Sprintf("[Core] Memory protection error: %s", err))
	}

	labelNotice.SetBackgroundColor(tcell.ColorDarkGrey)

	layoutStatus := statusview.NewStatusView()
	layoutStatus.SetDynamicColors(true)
	layoutStatus.SetBorder(true).
		SetTitle("Status").
		SetTitleAlign(tview.AlignLeft)

	layoutHeader := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(labelInstanceId, 0, 1, false).
		AddItem(labelNotice, 0, 1, false).
		AddItem(labelTitle, 0, 1, false)

	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(layoutHeader, 1, 1, false).
		AddItem(t.mainFrame.Layout(), 0, 6, true).
		AddItem(layoutStatus, 6, 1, false)

	// main TUI event loop
	go func() {
		for {
			select {
			case event := <-t.uiEvents.Events():
				switch event.Type() {
				case event_bus.EventLog:
					layoutStatus.Log(event.String())
				case event_bus.EventLogInfo:
					layoutStatus.Info(event.String())
				case event_bus.EventLogSuccess:
					layoutStatus.Success(event.String())
				case event_bus.EventLogWarning:
					layoutStatus.Warn(event.String())
				case event_bus.EventLogError:
					layoutStatus.Error(event.String())
				case event_bus.EventWalletInitialized:
					layoutStatus.Info("Wallet updated: " + event.String())
					labelInstanceId.SetText(fmt.Sprintf("[darkcyan]ID:[black] %s", event.String()))
					t.w3Events.Emit(event_bus.EventW3WalletAvailable, event.String())
				case event_bus.EventWalletNoticeMessage:
					labelNotice.SetText("[red]Memory protection not available in this system")
					layoutStatus.Error(event.String())
				case event_bus.EventDrawForce:
					t.app.Draw()
				case event_bus.EventShowModal:
					t.app.SetRoot(event.Data().(*tview.Modal), false)
				case event_bus.EventUpdateCurrencies:
					go func() {
						currencies := wallet.API().UpdateFiatCurrencies()
						if currencies != nil {
							layoutStatus.Success(fmt.Sprintf("Currencies loaded: %v", currencies))
						} else {
							layoutStatus.Error("Cannot retrieve currencies data")
						}
					}()
				case event_bus.EventW3Connect:
					t.mainFrame.State().SwitchToPage("w3_confirm_connect", event.Data())
				case event_bus.EventW3RequestAccounts:
					t.mainFrame.State().SwitchToPage("w3_request_accounts", event.Data())
				case event_bus.EventW3ReqCallGetBlockByNumber:
					go func() {
						req, ok := event.Data().(*dto.RequestCallGetBlockByNumberDTO)
						if !ok {
							layoutStatus.Error("Cannot parse w3 request")
							return
						}
						result, err := wallet.API().ExecuteRPC(&dto.ExecuteRPCRequestDTO{
							InstanceId: req.InstanceId,
							//Origin:     req.Origin,
							// RemoteAddr: "",
							ChainKey: req.ChainKey,
							Method:   req.Method,
							Params:   nil,
						})
						if err != nil {
							layoutStatus.Error("Cannot execute w3 call")
							return
						}
						t.mainFrame.State().EmitW3(event_bus.EventW3CallGetBlockByNumber, &dto.ResponseGetBlockByNumberDTO{
							InstanceId: req.InstanceId,
							Data:       result,
						})
					}()
				// Internal
				case event_bus.EventW3InternalConnections:
					layoutStatus.Info("Got connections")
					t.mainFrame.State().SwitchToPage("w3_connections", event.Data())
				case event_bus.EventQuit:
					// graceful shutdown
					// TODO: Uncomment on release
					/*modalConfirmQuit := tview.NewModal().
					SetText("Do you want to quit the application?").
					AddButtons([]string{"Quit", "Cancel"}).
					SetDoneFunc(func(buttonIndex int, buttonLabel string) {
						if buttonIndex == 0 {
							t.app.Stop()
						} else {
							t.app.SetRoot(layout, false)
						}
					})
					t.app.SetRoot(modalConfirmQuit, false)*/
					t.Stop()
					return
				default:
					//layoutStatus.Error()
					layoutStatus.SetText(fmt.Sprintf("unhandled event: %d", event.Type()))
				}
				t.app.ForceDraw()
			}
		}
	}()

	return layout
}

func (t *Tui) Start() error {

	if t.isVerboseMode {
		t.uiEvents.Emit(event_bus.EventLog, "Verbose mode enabled")

		var (
			prevX, prevY           int
			prevFrameX, prevFrameY int
		)

		t.app.SetAfterDrawFunc(func(screen tcell.Screen) {
			x, y := screen.Size()
			if x != prevX || y != prevY {
				prevX = x
				prevY = y
				t.uiEvents.Emit(event_bus.EventLog, fmt.Sprintf("Resolution: %dx%d", x, y))
			}

			x1, y1, x2, y2 := t.mainFrame.Layout().GetItem(1).GetRect()

			if x2 != prevFrameX || y2 != prevFrameY {
				prevFrameX = x2
				prevFrameY = y2
				t.uiEvents.Emit(event_bus.EventLog, fmt.Sprintf("Frame: %dx%d, %dx%d", x1, y1, x2, y2))
			}
		})
	}

	t.app.SetRoot(t.layout, true).
		EnableMouse(true).Run()
	return nil
	//}()
}

func (t *Tui) Stop() {
	if t.isVerboseMode {
		fmt.Println("[TUI] Stopping")
	}
	if t.stopped {
		return
	}
	defer t.wg.Done()

	t.stopped = true
	t.uiEvents.Stop()
	t.app.Stop()
}
