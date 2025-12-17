package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"howett.net/plist"
)

type AppInfo struct {
	Name     string
	BundleID string
	IconPath string
	Path     string
}

type Settings struct {
	SelectedApps []string `json:"selected_apps"`
}

func loadSettings() (*Settings, error) {
	// For dynamic, scan the folder
	selectedApps := getTranscribBundleIDs()
	return &Settings{SelectedApps: selectedApps}, nil
}

func getTranscribApps() []AppInfo {
	dir := "/Applications/Transcrib"
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	var apps []AppInfo
	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".app") {
			plistPath := filepath.Join(dir, entry.Name(), "Contents", "Info.plist")
			file, err := os.Open(plistPath)
			if err != nil {
				continue
			}
			var info map[string]interface{}
			decoder := plist.NewDecoder(file)
			err = decoder.Decode(&info)
			file.Close()
			if err != nil {
				continue
			}
			var app AppInfo
			if name, ok := info["CFBundleName"].(string); ok {
				app.Name = name
			} else {
				app.Name = strings.TrimSuffix(entry.Name(), ".app")
			}
			if bundleID, ok := info["CFBundleIdentifier"].(string); ok {
				app.BundleID = bundleID
				app.Path = filepath.Join(dir, entry.Name())
				if iconFile, ok := info["CFBundleIconFile"].(string); ok {
					if !strings.HasSuffix(iconFile, ".icns") {
						iconFile += ".icns"
					}
					app.IconPath = filepath.Join(dir, entry.Name(), "Contents", "Resources", iconFile)
				}
				apps = append(apps, app)
			}
		}
	}
	return apps
}

func getTranscribBundleIDs() []string {
	apps := getTranscribApps()
	var ids []string
	for _, app := range apps {
		ids = append(ids, app.BundleID)
	}
	return ids
}

func saveSettings(settings *Settings) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	settingsPath := filepath.Join(home, ".macos_app_switcher_settings.json")
	file, err := os.Create(settingsPath)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	return encoder.Encode(settings)
}
