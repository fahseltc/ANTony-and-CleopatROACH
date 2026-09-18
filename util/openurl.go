//go:build !js

package util

import (
	"log/slog"
	"os/exec"
	"runtime"
)

// OpenURL opens the given URL in the user's default web browser.
//
// This is the desktop (non-WASM) implementation. It shells out to the
// platform's standard "open a thing" command. The browser build has its own
// implementation in openurl_js.go that uses window.open.
func OpenURL(url string) {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		// rundll32 avoids cmd.exe quoting pitfalls with URLs containing & etc.
		cmd = "rundll32"
		args = []string{"url.dll,FileProtocolHandler", url}
	case "darwin":
		cmd = "open"
		args = []string{url}
	default: // linux, bsd, ...
		cmd = "xdg-open"
		args = []string{url}
	}

	if err := exec.Command(cmd, args...).Start(); err != nil {
		slog.Warn("failed to open URL", "url", url, "err", err)
	}
}
