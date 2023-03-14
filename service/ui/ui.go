package ui

import (
	"fmt"
	"github.com/censync/go-i18n"
	"github.com/censync/soikawallet/config/dict"
	"github.com/censync/soikawallet/config/version"
	h "github.com/censync/soikawallet/service/ui/handler"
	"github.com/censync/soikawallet/service/ui/walletframe"
	"github.com/censync/soikawallet/service/ui/widgets/statusview"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Tui struct {
	app *tview.Application
	tr  *i18n.Translator

	frame         *walletframe.WalletFrame
	tbus          h.TBus
	layout        *tview.Flex
	style         *tview.Theme
	isVerboseMode bool
}

func Init() *Tui {
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
		app:   tview.NewApplication(),
		tr:    dict.GetTr("en"),
		tbus:  make(h.TBus, 100),
		style: style,
	}
	tui.frame = walletframe.Init(&tui.tbus, dict.GetTr("en"), tui.style)
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

	labelUUID := tview.NewTextView().
		SetDynamicColors(true).
		SetWrap(false).
		SetTextAlign(tview.AlignLeft).
		SetText(`[darkcyan]UUID:[black] 00000000-0000-0000-0000-000000000000`)

	labelUUID.SetBackgroundColor(tcell.ColorDarkGrey).
		SetBorderPadding(0, 0, 2, 0)

	layoutStatus := statusview.NewStatusView()
	layoutStatus.SetDynamicColors(true)
	layoutStatus.SetBorder(true).
		SetTitle("Status").
		SetTitleAlign(tview.AlignLeft)

	layoutHeader := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(labelUUID, 0, 1, false).
		AddItem(labelTitle, 0, 1, false)

	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(layoutHeader, 1, 1, false).
		AddItem(t.frame.Layout(), 0, 6, true).
		AddItem(layoutStatus, 6, 1, false)

	go func() {
		for {
			select {
			case event := <-t.tbus:
				switch event.Type() {
				case h.EventLog:
					layoutStatus.Log(event.String())
				case h.EventLogInfo:
					layoutStatus.Info(event.String())
				case h.EventLogSuccess:
					layoutStatus.Success(event.String())
				case h.EventLogWarning:
					layoutStatus.Warn(event.String())
				case h.EventLogError:
					layoutStatus.Error(event.String())
				case h.EventUpdatedWallet:
					layoutStatus.Info("Wallet updated: " + event.String())
					labelUUID.SetText(fmt.Sprintf("[darkcyan]UUID:[black] %s", event.String()))
				case h.EventDrawForce:
					t.app.Draw()
				case h.EventShowModal:
					t.app.SetRoot(event.Data().(*tview.Modal), false)
				case h.EventQuit:
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
					t.app.Stop()
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

func (t *Tui) Run(verbose bool) {
	// Run the application
	t.isVerboseMode = false

	if t.isVerboseMode {
		t.tbus.Emit(h.EventLog, "Verbose mode enabled")

		var (
			prevX, prevY           int
			prevFrameX, prevFrameY int
		)

		t.app.SetAfterDrawFunc(func(screen tcell.Screen) {
			x, y := screen.Size()
			if x != prevX || y != prevY {
				prevX = x
				prevY = y
				t.tbus.Emit(h.EventLog, fmt.Sprintf("Resolution: %dx%d", x, y))
			}

			x1, y1, x2, y2 := t.frame.Layout().GetItem(1).GetRect()

			if x2 != prevFrameX || y2 != prevFrameY {
				prevFrameX = x2
				prevFrameY = y2
				t.tbus.Emit(h.EventLog, fmt.Sprintf("Frame: %dx%d, %dx%d", x1, y1, x2, y2))
			}
		})
	}

	if err := t.app.SetRoot(t.layout, true).
		EnableMouse(true).
		Run(); err != nil {
		panic(err)
	}
}
