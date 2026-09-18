//go:build js

package util

import "syscall/js"

// OpenURL opens the given URL in a new browser tab.
//
// This is the WASM/browser implementation. window.open with "_blank" pops a
// new tab; "noopener" is passed for the usual security hygiene. The desktop
// build has its own implementation in openurl.go.
func OpenURL(url string) {
	js.Global().Call("open", url, "_blank", "noopener")
}
