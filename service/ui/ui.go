package ui

import (
	"fmt"
	"github.com/censync/go-i18n"
	"github.com/censync/soikawallet/config"
	"github.com/censync/soikawallet/config/dict"
	"github.com/censync/soikawallet/config/version"
	"github.com/censync/soikawallet/service/internal/event_bus"
	"github.com/censync/soikawallet/service/ui/walletframe"
	"github.com/censync/soikawallet/service/ui/widgets/statusview"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"sync"
)

type Tui struct {
	app *tview.Application
	tr  *i18n.Translator

	frame         *walletframe.WalletFrame
	events        *event_bus.EventBus
	layout        *tview.Flex
	style         *tview.Theme
	isVerboseMode bool
	wg            *sync.WaitGroup
	stopped       bool
}

func NewTui(cfg *config.Config, wg *sync.WaitGroup, events *event_bus.EventBus) *Tui {
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
		app:    tview.NewApplication(),
		tr:     dict.GetTr("en"),
		events: events,
		style:  style,
		wg:     wg,
	}
	tui.frame = walletframe.Init(tui.events, dict.GetTr("en"), tui.style)
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

	labelInstanceId := tview.NewTextView().
		SetDynamicColors(true).
		SetWrap(false).
		SetTextAlign(tview.AlignLeft).
		SetText(`[darkcyan]ID:[black] not initialized`)

	labelInstanceId.SetBackgroundColor(tcell.ColorDarkGrey).
		SetBorderPadding(0, 0, 2, 0)

	layoutStatus := statusview.NewStatusView()
	layoutStatus.SetDynamicColors(true)
	layoutStatus.SetBorder(true).
		SetTitle("Status").
		SetTitleAlign(tview.AlignLeft)

	layoutHeader := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(labelInstanceId, 0, 1, false).
		AddItem(labelTitle, 0, 1, false)

	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(layoutHeader, 1, 1, false).
		AddItem(t.frame.Layout(), 0, 6, true).
		AddItem(layoutStatus, 6, 1, false)

	// main TUI event loop
	go func() {
		for {
			select {
			case event := <-t.events.Events():
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
				case event_bus.EventDrawForce:
					t.app.Draw()
				case event_bus.EventShowModal:
					t.app.SetRoot(event.Data().(*tview.Modal), false)
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
				//t.app.ForceDraw()
			}
		}
	}()

	return layout
}

func (t *Tui) Start() error {
	// Start the application
	//go func() {

	if t.isVerboseMode {
		t.events.Emit(event_bus.EventLog, "Verbose mode enabled")

		var (
			prevX, prevY           int
			prevFrameX, prevFrameY int
		)

		t.app.SetAfterDrawFunc(func(screen tcell.Screen) {
			x, y := screen.Size()
			if x != prevX || y != prevY {
				prevX = x
				prevY = y
				t.events.Emit(event_bus.EventLog, fmt.Sprintf("Resolution: %dx%d", x, y))
			}

			x1, y1, x2, y2 := t.frame.Layout().GetItem(1).GetRect()

			if x2 != prevFrameX || y2 != prevFrameY {
				prevFrameX = x2
				prevFrameY = y2
				t.events.Emit(event_bus.EventLog, fmt.Sprintf("Frame: %dx%d, %dx%d", x1, y1, x2, y2))
			}
		})
	}

	t.app.SetRoot(t.layout, true).
		EnableMouse(true).Run()
	return nil
	//}()
}

func (t *Tui) Stop() {
	if t.stopped {
		return
	}
	defer t.wg.Done()

	t.stopped = true
	t.events.Stop()
	t.app.Stop()
}
