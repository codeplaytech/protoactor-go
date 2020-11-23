package cluster

import "github.com/AsynkronIT/protoactor-go/actor"

var (
	defaultCluster *Cluster = nil
)

// Start the default instance of cluster.
func Start(c *Config) {
	defaultCluster = New(actor.DefaultSystem, c)
	defaultCluster.Start()
}

func StartClient(c *Config) {
	defaultCluster = New(actor.DefaultSystem, c)
	defaultCluster.StartClient()
}

// Shutdown the default instance of cluster.
func Shutdown(graceful bool) {
	if defaultCluster == nil {
		plog.Error("no default cluster")
		return
	}
	// plog.Error("default cluster is stoping.")
	defaultCluster.Shutdown(graceful)
	defaultCluster = nil
	// plog.Error("default cluster is stoped.")
}
