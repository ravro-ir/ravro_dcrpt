package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var (
	logBgColor   = color.NRGBA{R: 0x11, G: 0x11, B: 0x11, A: 0xff}
	logTextColor = color.NRGBA{R: 0xf5, G: 0xf5, B: 0xf5, A: 0xff}
)

type terminalLog struct {
	widget.BaseWidget
	lines []string
}

func newTerminalLog() *terminalLog {
	l := &terminalLog{}
	l.ExtendBaseWidget(l)
	return l
}

func (l *terminalLog) Append(message string) {
	l.lines = append(l.lines, message)
	l.Refresh()
}

func (l *terminalLog) Clear() {
	l.lines = nil
	l.Refresh()
}

func (l *terminalLog) CreateRenderer() fyne.WidgetRenderer {
	bg := canvas.NewRectangle(logBgColor)
	content := container.NewVBox()
	scroll := container.NewScroll(content)
	stack := container.NewStack(bg, container.NewPadded(scroll))
	r := &terminalLogRenderer{
		log:     l,
		bg:      bg,
		content: content,
		scroll:  scroll,
		stack:   stack,
	}
	r.layoutLines()
	return r
}

type terminalLogRenderer struct {
	log     *terminalLog
	bg      *canvas.Rectangle
	content *fyne.Container
	scroll  *container.Scroll
	stack   *fyne.Container
}

func (r *terminalLogRenderer) layoutLines() {
	objs := make([]fyne.CanvasObject, 0, len(r.log.lines))
	for _, line := range r.log.lines {
		t := canvas.NewText(line, logTextColor)
		t.TextSize = 13
		t.TextStyle = fyne.TextStyle{Monospace: true}
		objs = append(objs, t)
	}
	r.content.Objects = objs
	r.content.Refresh()
}

func (r *terminalLogRenderer) Destroy() {}

func (r *terminalLogRenderer) Layout(size fyne.Size) {
	r.stack.Resize(size)
}

func (r *terminalLogRenderer) MinSize() fyne.Size {
	return fyne.NewSize(200, 180)
}

func (r *terminalLogRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.stack}
}

func (r *terminalLogRenderer) Refresh() {
	r.bg.FillColor = logBgColor
	r.layoutLines()
	r.bg.Refresh()
	r.stack.Refresh()
	r.scroll.ScrollToBottom()
}
