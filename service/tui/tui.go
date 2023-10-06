package tui

import (
	"fmt"
	"github.com/censync/go-i18n"
	"github.com/censync/soikawallet/api/dto"
	"github.com/censync/soikawallet/config"
	"github.com/censync/soikawallet/config/dict"
	"github.com/censync/soikawallet/config/version"
	"github.com/censync/soikawallet/service/core"
	"github.com/censync/soikawallet/service/tui/events"
	"github.com/censync/soikawallet/service/tui/tmainframe"
	"github.com/censync/soikawallet/service/tui/twidget/statusview"
	"github.com/censync/soikawallet/types/protected_key"
	"github.com/censync/tview"
	"github.com/gdamore/tcell/v2"
	"sync"
)

type Tui struct {
	app *tview.Application
	tr  *i18n.Translator

	mainFrame     *tmainframe.TMainFrame
	uiEvents      *events.EventBus
	w3Events      *events.EventBus
	layout        *tview.Flex
	style         *tview.Theme
	isVerboseMode bool
	wg            *sync.WaitGroup
	stopped       bool
}

func NewTui(cfg *config.Config, wg *sync.WaitGroup, uiEvents, w3Events *events.EventBus) *Tui {
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
		t.uiEvents.Emit(events.EventWalletNoticeMessage, fmt.Sprintf("[Core] Memory protection error: %s", err))
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
				case events.EventLog:
					layoutStatus.Log(event.String())
				case events.EventLogInfo:
					layoutStatus.Info(event.String())
				case events.EventLogSuccess:
					layoutStatus.Success(event.String())
				case events.EventLogWarning:
					layoutStatus.Warn(event.String())
				case events.EventLogError:
					layoutStatus.Error(event.String())
				case events.EventWalletInitialized:
					layoutStatus.Info("Wallet updated: " + event.String())
					labelInstanceId.SetText(fmt.Sprintf("[darkcyan]ID:[black] %s", event.String()))
					t.w3Events.Emit(events.EventW3WalletAvailable, event.String())
				case events.EventWalletNoticeMessage:
					labelNotice.SetText("[red]Memory protection not available in this system")
					layoutStatus.Error(event.String())
				case events.EventDrawForce:
					t.app.Draw()
				case events.EventShowModal:
					t.app.SetRoot(event.Data().(*tview.Modal), false)
				case events.EventUpdateCurrencies:
					go func() {
						currencies := core.API().UpdateFiatCurrencies()
						if currencies != nil {
							layoutStatus.Success(fmt.Sprintf("Currencies loaded: %v", currencies))
						} else {
							layoutStatus.Error("Cannot retrieve currencies data")
						}
					}()
				case events.EventW3Connect:
					t.mainFrame.State().SwitchToPage("w3_confirm_connect", event.Data())
				case events.EventW3RequestAccounts:
					t.mainFrame.State().SwitchToPage("w3_request_accounts", event.Data())
				case events.EventW3ReqCallGetBlockByNumber:
					go func() {
						req, ok := event.Data().(*dto.RequestCallGetBlockByNumberDTO)
						if !ok {
							layoutStatus.Error("Cannot parse w3 request")
							return
						}
						result, err := core.API().ExecuteRPC(&dto.ExecuteRPCRequestDTO{
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
						t.mainFrame.State().EmitW3(events.EventW3CallGetBlockByNumber, &dto.ResponseGetBlockByNumberDTO{
							InstanceId: req.InstanceId,
							Data:       result,
						})
					}()
				// Internal
				case events.EventW3InternalConnections:
					layoutStatus.Info("Got connections")
					t.mainFrame.State().SwitchToPage("w3_connections", event.Data())
				case events.EventQuit:
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
		t.uiEvents.Emit(events.EventLog, "Verbose mode enabled")

		var (
			prevX, prevY           int
			prevFrameX, prevFrameY int
		)

		t.app.SetAfterDrawFunc(func(screen tcell.Screen) {
			x, y := screen.Size()
			if x != prevX || y != prevY {
				prevX = x
				prevY = y
				t.uiEvents.Emit(events.EventLog, fmt.Sprintf("Resolution: %dx%d", x, y))
			}

			x1, y1, x2, y2 := t.mainFrame.Layout().GetItem(1).GetRect()

			if x2 != prevFrameX || y2 != prevFrameY {
				prevFrameX = x2
				prevFrameY = y2
				t.uiEvents.Emit(events.EventLog, fmt.Sprintf("Frame: %dx%d, %dx%d", x1, y1, x2, y2))
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
