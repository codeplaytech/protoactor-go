package remote

import "github.com/AsynkronIT/protoactor-go/actor"

var (
	defaultRemote *Remote = nil
)

// Start the default instance of remote
func Start(c Config) {
	defaultRemote = NewRemote(actor.DefaultSystem, c)
	for kind, props := range c.Kinds {
		defaultRemote.Register(kind, props)
	}
	defaultRemote.Start()
}

// Shutdown the default instance of remote
func Shutdown(graceful bool) {
	if defaultRemote == nil {
		plog.Error("default instance was nil")
		return
	}
	defaultRemote.Shutdown(graceful)
	defaultRemote = nil
}
