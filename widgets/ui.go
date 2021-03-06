package widgets

import (
	"fmt"
	"strings"

	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/console"
	"github.com/ambientsound/pms/constants"
	"github.com/ambientsound/pms/index"
	"github.com/ambientsound/pms/options"
	"github.com/ambientsound/pms/songlist"
	"github.com/ambientsound/pms/style"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

type UI struct {
	// UI elements
	Screen tcell.Screen
	App    *views.Application
	Layout *views.BoxLayout

	Topbar        *Topbar
	Columnheaders *ColumnheadersWidget
	Multibar      *MultibarWidget
	Songlist      *SonglistWidget

	// Input events
	EventInputCommand chan string
	EventKeyInput     chan *tcell.EventKey

	// Data resources
	api          api.API
	Index        *index.Index
	options      *options.Options // FIXME: use api instead
	searchResult songlist.Songlist

	// TCell
	view views.View
	style.Styled
	views.WidgetWatchers
}

func NewUI(a api.API) *UI {
	var err error

	ui := &UI{}

	ui.Screen, err = tcell.NewScreen()
	if err != nil {
		return nil
	}

	ui.EventInputCommand = make(chan string, 16)
	ui.EventKeyInput = make(chan *tcell.EventKey, 16)

	ui.App = &views.Application{}
	ui.api = a
	ui.options = ui.api.Options()

	ui.Topbar = NewTopbar()
	ui.Columnheaders = NewColumnheadersWidget()
	ui.Multibar = NewMultibarWidget(ui.api, ui.EventKeyInput)
	ui.Songlist = NewSonglistWidget(ui.api)

	ui.Multibar.Watch(ui)
	ui.Songlist.Watch(ui)

	// Set styles
	ui.SetStylesheet(ui.api.Styles())
	ui.Topbar.SetStylesheet(ui.api.Styles())
	ui.Columnheaders.SetStylesheet(ui.api.Styles())
	ui.Songlist.SetStylesheet(ui.api.Styles())
	ui.Multibar.SetStylesheet(ui.api.Styles())

	ui.CreateLayout()
	ui.App.SetScreen(ui.Screen)
	ui.App.SetRootWidget(ui)

	return ui
}

func (ui *UI) CreateLayout() {
	ui.Layout = views.NewBoxLayout(views.Vertical)
	ui.Layout.AddWidget(ui.Topbar, 1)
	ui.Layout.AddWidget(ui.Columnheaders, 0)
	ui.Layout.AddWidget(ui.Songlist, 2)
	ui.Layout.AddWidget(ui.Multibar, 0)
	ui.Layout.SetView(ui.view)
}

func (ui *UI) Refresh() {
	ui.App.Refresh()
}

func (ui *UI) SetIndex(i *index.Index) {
	ui.Index = i
}

func (ui *UI) CurrentSonglistWidget() api.SonglistWidget {
	return ui.Songlist
}

func (ui *UI) Start() {
	ui.App.Start()
}

func (ui *UI) Wait() error {
	return ui.App.Wait()
}

func (ui *UI) Quit() {
	ui.App.Quit()
}

func (ui *UI) Draw() {
	ui.Layout.Draw()
}

func (ui *UI) Resize() {
	ui.CreateLayout()
	ui.Layout.Resize()
	ui.PostEventWidgetResize(ui)
}

func (ui *UI) SetView(v views.View) {
	ui.view = v
	ui.Layout.SetView(v)
}

func (ui *UI) Size() (int, int) {
	return ui.view.Size()
}

func (ui *UI) Title() string {
	index, err := ui.Songlist.SonglistIndex()
	if err == nil {
		return fmt.Sprintf("[%d/%d] %s", index+1, ui.Songlist.SonglistsLen(), ui.Songlist.Name())
	} else {
		return fmt.Sprintf("[...] %s", ui.Songlist.Name())
	}
}

func (ui *UI) UpdateCursor() {
	switch ui.Multibar.Mode() {
	case constants.MultibarModeInput, constants.MultibarModeSearch:
		_, ymax := ui.Screen.Size()
		ui.Screen.ShowCursor(ui.Multibar.Cursor()+1, ymax-1)
	default:
		ui.Screen.HideCursor()
	}
}

func (ui *UI) PostFunc(f func()) {
	ui.App.PostFunc(f)
}

func (ui *UI) HandleEvent(ev tcell.Event) bool {
	switch ev.(type) {

	// If a list was changed, make sure we obtain the correct column widths.
	case *EventListChanged:
		tags := strings.Split(ui.options.StringValue("columns"), ",")
		cols := ui.Songlist.Songlist().Columns(tags)
		ui.Songlist.SetColumns(tags)
		ui.Columnheaders.SetColumns(cols)
		return true

	case *EventInputChanged:
		term := ui.Multibar.RuneString()
		mode := ui.Multibar.Mode()
		switch mode {
		case constants.MultibarModeSearch:
			if err := ui.runIndexSearch(term); err != nil {
				console.Log("Error while searching: %s", err)
			}
		}
		ui.UpdateCursor()
		return true

	case *EventInputFinished:
		term := ui.Multibar.RuneString()
		mode := ui.Multibar.Mode()
		switch mode {
		case constants.MultibarModeInput:
			ui.EventInputCommand <- term
		case constants.MultibarModeSearch:
			if ui.searchResult != nil {
				if ui.searchResult.Len() > 0 {
					ui.Songlist.AddSonglist(ui.searchResult)
				} else {
					ui.searchResult = nil
				}
			}
			ui.showSearchResult()
		}
		ui.Multibar.SetMode(constants.MultibarModeNormal)
		return true

	case *EventScroll:
		ui.refreshPositionReadout()
		return true
	}

	if ui.Layout.HandleEvent(ev) {
		return true
	}

	return false
}

func (ui *UI) refreshPositionReadout() {
	str := ui.Songlist.PositionReadout()
	ui.Multibar.SetRight(str, ui.Style("readout"))
}

func (ui *UI) runIndexSearch(term string) error {
	var err error

	if ui.Index == nil {
		return fmt.Errorf("Search index is not operational")
	}

	ui.searchResult, err = ui.Index.Search(term)

	ui.showSearchResult()

	return err
}

func (ui *UI) showSearchResult() {
	if ui.searchResult != nil {
		ui.Songlist.SetSonglist(ui.searchResult)
	} else if ui.Songlist.FallbackSonglist() != nil {
		ui.Songlist.SetSonglist(ui.Songlist.FallbackSonglist())
	} else {
		ui.Songlist.SetSonglistIndex(0)
	}
}
