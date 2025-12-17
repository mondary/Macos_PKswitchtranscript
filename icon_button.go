package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var _ fyne.Tappable = (*IconButton)(nil)

// IconButton is a tappable icon-only control that is intentionally not focusable,
// so it won't participate in Tab focus traversal.
type IconButton struct {
	widget.BaseWidget
	Icon    fyne.Resource
	OnTap   func()
	minSize fyne.Size
}

func NewIconButton(icon fyne.Resource, minSize fyne.Size, onTap func()) *IconButton {
	b := &IconButton{
		Icon:    icon,
		OnTap:   onTap,
		minSize: minSize,
	}
	b.ExtendBaseWidget(b)
	return b
}

func (b *IconButton) MinSize() fyne.Size {
	b.ExtendBaseWidget(b)
	if b.minSize.Width == 0 && b.minSize.Height == 0 {
		return b.BaseWidget.MinSize()
	}
	return b.minSize
}

func (b *IconButton) Tapped(*fyne.PointEvent) {
	if b.OnTap != nil {
		b.OnTap()
	}
}

func (b *IconButton) CreateRenderer() fyne.WidgetRenderer {
	b.ExtendBaseWidget(b)

	th := b.Theme()
	v := fyne.CurrentApp().Settings().ThemeVariant()

	background := canvas.NewRectangle(th.Color(theme.ColorNameButton, v))
	background.CornerRadius = th.Size(theme.SizeNameInputRadius)

	icon := canvas.NewImageFromResource(b.Icon)
	icon.FillMode = canvas.ImageFillContain

	objects := []fyne.CanvasObject{background, icon}
	return &iconButtonRenderer{
		objects:    objects,
		button:     b,
		background: background,
		icon:       icon,
	}
}

type iconButtonRenderer struct {
	objects    []fyne.CanvasObject
	button     *IconButton
	background *canvas.Rectangle
	icon       *canvas.Image
}

func (r *iconButtonRenderer) Layout(size fyne.Size) {
	r.background.Resize(size)

	th := r.button.Theme()
	padding := th.Size(theme.SizeNameInnerPadding)
	pos := fyne.NewPos(padding, padding)
	iconSize := size.Subtract(fyne.NewSize(padding*2, padding*2))
	if iconSize.Width < 0 {
		iconSize.Width = 0
	}
	if iconSize.Height < 0 {
		iconSize.Height = 0
	}
	r.icon.Move(pos)
	r.icon.Resize(iconSize)
}

func (r *iconButtonRenderer) MinSize() fyne.Size {
	return r.button.MinSize()
}

func (r *iconButtonRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *iconButtonRenderer) Destroy() {}

func (r *iconButtonRenderer) Refresh() {
	th := r.button.Theme()
	v := fyne.CurrentApp().Settings().ThemeVariant()

	r.background.FillColor = th.Color(theme.ColorNameButton, v)
	r.background.CornerRadius = th.Size(theme.SizeNameInputRadius)
	r.background.Refresh()

	r.icon.Resource = r.button.Icon
	r.icon.Refresh()
}
