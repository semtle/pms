package widgets

import (
	"github.com/ambientsound/pms/index"
	"github.com/ambientsound/pms/songlist"
	"github.com/ambientsound/pms/version"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

type StyleMap map[string]tcell.Style

type UI struct {
	// UI elements
	App    *views.Application
	Layout *views.BoxLayout

	Topbar        *views.TextBar
	Playbar       *PlaybarWidget
	Columnheaders *ColumnheadersWidget
	Multibar      *MultibarWidget
	Songlist      *SongListWidget

	// Data resources
	Index           *index.Index
	defaultSongList *songlist.SongList

	// TCell
	view views.View
	widget
}

func NewUI() *UI {
	ui := &UI{}

	ui.App = &views.Application{}

	ui.Topbar = views.NewTextBar()
	ui.Playbar = NewPlaybarWidget()
	ui.Columnheaders = NewColumnheadersWidget()
	ui.Multibar = NewMultibarWidget()
	ui.Songlist = NewSongListWidget()

	ui.Multibar.Watch(ui)
	ui.Songlist.Watch(ui)
	ui.Playbar.Watch(ui)

	styles := StyleMap{
		"album":    tcell.StyleDefault.Foreground(tcell.ColorTeal),
		"artist":   tcell.StyleDefault.Foreground(tcell.ColorYellow),
		"cursor":   tcell.StyleDefault.Background(tcell.ColorWhite).Foreground(tcell.ColorBlack),
		"date":     tcell.StyleDefault.Foreground(tcell.ColorGreen),
		"elapsed":  tcell.StyleDefault.Foreground(tcell.ColorGreen),
		"header":   tcell.StyleDefault.Foreground(tcell.ColorGreen).Bold(true),
		"switches": tcell.StyleDefault.Foreground(tcell.ColorTeal),
		//"time":     tcell.StyleDefault.Foreground(tcell.ColorTeal),
		"time":  tcell.StyleDefault.Foreground(tcell.ColorDarkMagenta),
		"title": tcell.StyleDefault.Foreground(tcell.ColorWhite).Bold(true),
		//"title":  tcell.StyleDefault.Background(tcell.ColorBlue).Foreground(tcell.ColorWhite).Bold(true),
		//"title":  tcell.StyleDefault.Foreground(tcell.ColorWhite).Bold(true),
		"topbar": tcell.StyleDefault.Foreground(tcell.ColorYellow).Bold(true),
		"track":  tcell.StyleDefault.Foreground(tcell.ColorGreen),
		"volume": tcell.StyleDefault.Foreground(tcell.ColorGreen),
		"year":   tcell.StyleDefault.Foreground(tcell.ColorGreen),
	}

	// Styles for widgets that don't have their own class yet.
	ui.SetStyleMap(styles)

	ui.Columnheaders.SetStyleMap(styles)
	ui.Playbar.SetStyleMap(styles)
	ui.Songlist.SetStyleMap(styles)

	ui.Topbar.SetStyle(ui.Style("topbar"))
	ui.Topbar.SetLeft(version.ShortName(), ui.Style("topbar"))
	ui.Topbar.SetRight(version.Version(), ui.Style("topbar"))

	ui.Multibar.SetDefaultText("Type to search.")

	ui.CreateLayout()
	ui.App.SetRootWidget(ui)

	return ui
}

func (ui *UI) CreateLayout() {
	ui.Layout = views.NewBoxLayout(views.Vertical)
	ui.Layout.AddWidget(ui.Topbar, 0)
	ui.Layout.AddWidget(ui.Playbar, 0)
	ui.Layout.AddWidget(ui.Columnheaders, 0)
	ui.Layout.AddWidget(ui.Songlist, 2)
	ui.Layout.AddWidget(ui.Multibar, 0)
	ui.Layout.SetView(ui.view)
}

func (ui *UI) SetIndex(i *index.Index) {
	ui.Index = i
}

func (ui *UI) SetDefaultSonglist(s *songlist.SongList) {
	ui.defaultSongList = s
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

func (ui *UI) HandleEvent(ev tcell.Event) bool {
	switch ev := ev.(type) {

	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyCtrlC:
			fallthrough
		case tcell.KeyCtrlD:
			ui.App.Quit()
			return true
		case tcell.KeyCtrlL:
			ui.App.Refresh()
			return true
		}

	case *EventListChanged:
		ui.App.Update()
		ui.Topbar.SetCenter(" "+ui.Songlist.Name()+" ", ui.Style("title"))
		ui.Columnheaders.SetColumns(ui.Songlist.Columns())
		return true

	case *EventInputChanged:
		term := ui.Multibar.GetRuneString()
		ui.runIndexSearch(term)
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
	ui.Multibar.SetRight(str, tcell.StyleDefault)
}

func (ui *UI) runIndexSearch(term string) {
	if ui.Index == nil {
		return
	}
	if len(term) == 0 {
		ui.Songlist.SetCursor(0)
		ui.Songlist.SetSongList(ui.defaultSongList)
		return
	}
	if len(term) == 1 {
		return
	}
	results, err := ui.Index.Search(term)
	if err == nil {
		ui.Songlist.SetCursor(0)
		ui.Songlist.SetSongList(results)
		return
	}
}