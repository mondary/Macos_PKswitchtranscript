package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// EnterButton is a standard Fyne button that also triggers on Enter/Return.
type EnterButton struct {
	widget.Button
}

func NewEnterButton(label string, tapped func()) *EnterButton {
	b := &EnterButton{}
	b.Text = label
	b.OnTapped = tapped
	b.ExtendBaseWidget(b)
	return b
}

func NewEnterButtonWithIcon(label string, icon fyne.Resource, tapped func()) *EnterButton {
	b := &EnterButton{}
	b.Text = label
	b.Icon = icon
	b.OnTapped = tapped
	b.ExtendBaseWidget(b)
	return b
}

func (b *EnterButton) TypedKey(ev *fyne.KeyEvent) {
	switch ev.Name {
	case fyne.KeySpace, fyne.KeyReturn, fyne.KeyEnter:
		b.Tapped(nil)
	default:
		b.Button.TypedKey(ev)
	}
}
