//go:build windows

package companion

import (
	"os"
	"path/filepath"
)

const windowsAppName = "Praxis"

func defaultDesktopDataDir() string {
	appData := os.Getenv("LOCALAPPDATA")
	if appData == "" {
		appData = os.Getenv("APPDATA")
	}
	return filepath.Join(appData, windowsAppName)
}

func platformProvider() Provider {
	return desktopProvider{secrets: newPipeSecretsClient()}
}
