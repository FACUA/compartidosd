package app

import (
	"facua.org/compartidosd/client/fs"
)

// Start starts the application daemon
func Start() {
	fs.Init()
	Tick()
}
