package main

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/nfnt/resize"
	"github.com/vcaesar/screenshot"
)

var fyneApp fyne.App
var window fyne.Window
var settings *Settings
var apps []AppInfo

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func main() {
	var err error
	settings, err = loadSettings()
	if err != nil {
		log.Fatal(err)
	}

	apps = getTranscribApps()

	// Filter apps based on selected
	if len(settings.SelectedApps) > 0 {
		var filtered []AppInfo
		for _, app := range apps {
			if contains(settings.SelectedApps, app.BundleID) {
				filtered = append(filtered, app)
			}
		}
		apps = filtered
	}

	bounds := screenshot.GetDisplayBounds(0)
	screenWidth := bounds.Dx()
	maxWidth := int(float64(screenWidth) * 0.8)
	iconSize := 128
	windowWidth := len(apps) * iconSize
	if windowWidth > maxWidth {
		iconSize = maxWidth / len(apps)
		if iconSize < 64 {
			iconSize = 64
		}
		windowWidth = len(apps) * iconSize
	}

	fyneApp = app.New()
	rect := image.Rect(0, 0, 128, 128)
	img := image.NewRGBA(rect)
	for y := 0; y < 128; y++ {
		for x := 0; x < 128; x++ {
			img.Set(x, y, color.RGBA{0, 123, 255, 255})
		}
	}
	var buf bytes.Buffer
	png.Encode(&buf, img)
	appIcon := fyne.NewStaticResource("appicon", buf.Bytes())

	window = fyneApp.NewWindow("App Switcher")
	window.SetIcon(appIcon)
	window.SetOnClosed(func() {
		window.Hide()
	})

	var buttons []fyne.CanvasObject
	for _, app := range apps {
		app := app // capture
		var button *widget.Button
		var icon fyne.Resource
		if app.IconPath != "" {
			tempPng := filepath.Join(os.TempDir(), app.BundleID+".png")
			cmd := exec.Command("sips", "-s", "format", "png", app.IconPath, "--out", tempPng)
			err := cmd.Run()
			if err == nil {
				data, err := os.ReadFile(tempPng)
				if err == nil {
					img, _, err := image.Decode(bytes.NewReader(data))
					if err == nil {
						resized := resize.Resize(uint(iconSize), uint(iconSize), img, resize.Lanczos3)
						var buf bytes.Buffer
						png.Encode(&buf, resized)
						icon = fyne.NewStaticResource(app.BundleID, buf.Bytes())
					}
				}
				os.Remove(tempPng)
			}
		}
		if icon != nil {
			button = widget.NewButtonWithIcon("", icon, func() {
				exec.Command("open", app.Path).Run()
				for _, other := range apps {
					if other.BundleID != app.BundleID {
						exec.Command("osascript", "-e", "tell application \""+other.Name+"\" to quit").Run()
					}
				}
				window.Hide()
			})
		} else {
			button = widget.NewButton(app.Name, func() {
				exec.Command("open", app.Path).Run()
				for _, other := range apps {
					if other.BundleID != app.BundleID {
						exec.Command("osascript", "-e", "tell application \""+other.Name+"\" to quit").Run()
					}
				}
				window.Hide()
			})
		}
		buttons = append(buttons, button)
	}

	moveFocus := func(delta int) {
		if len(buttons) == 0 {
			return
		}
		focused := window.Canvas().Focused()
		currentIndex := -1
		if btn, ok := focused.(*widget.Button); ok {
			for i, b := range buttons {
				if b.(*widget.Button) == btn {
					currentIndex = i
					break
				}
			}
		}
		if currentIndex == -1 {
			currentIndex = 0
		}
		newIndex := (currentIndex + delta + len(buttons)) % len(buttons)
		window.Canvas().Focus(buttons[newIndex].(*widget.Button))
	}

	window.SetContent(container.NewCenter(container.NewHBox(buttons...)))

	openSettings := func() {
		settingsWindow := fyneApp.NewWindow("Settings")
		allApps := getTranscribApps()
		var checkboxes []*widget.Check
		for _, app := range allApps {
			cb := widget.NewCheck(app.Name, nil)
			cb.SetChecked(contains(settings.SelectedApps, app.BundleID))
			checkboxes = append(checkboxes, cb)
		}
		saveBtn := widget.NewButton("Save", func() {
			settings.SelectedApps = nil
			for i, cb := range checkboxes {
				if cb.Checked {
					settings.SelectedApps = append(settings.SelectedApps, allApps[i].BundleID)
				}
			}
			saveSettings(settings)
			settingsWindow.Close()
		})
		content := container.NewVBox()
		for _, cb := range checkboxes {
			content.Add(cb)
		}
		content.Add(saveBtn)
		settingsWindow.SetContent(content)
		settingsWindow.Resize(fyne.NewSize(300, 400))
		settingsWindow.Show()
	}
	window.Canvas().SetOnTypedKey(func(ev *fyne.KeyEvent) {
		if ev.Name == fyne.KeyReturn || ev.Name == fyne.KeyEnter {
			if focused := window.Canvas().Focused(); focused != nil {
				if btn, ok := focused.(*widget.Button); ok {
					btn.OnTapped()
				}
			}
		} else if ev.Name == fyne.KeyLeft {
			moveFocus(-1)
		} else if ev.Name == fyne.KeyRight {
			moveFocus(1)
		} else if ev.Name == fyne.KeyS {
			openSettings()
		}
	})

	window.Resize(fyne.NewSize(float32(windowWidth), float32(iconSize)))
	window.SetFixedSize(true)
	window.CenterOnScreen()
	if len(buttons) > 0 {
		window.Canvas().Focus(buttons[0].(*widget.Button))
	}
	window.Show()

	fyneApp.Run()
}
