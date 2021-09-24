package cluster

import (
	"github.com/cespare/xxhash"
	rdv "github.com/dgryski/go-rendezvous"
)

// RendezvousV2 ...
type RendezvousV2 struct {
	rdv      *rdv.Rendezvous
	size     int
	dumpInfo []string
}

// NewRendezvousV2 ...
func NewRendezvousV2(members []*Member) *RendezvousV2 {
	addrs := make([]string, len(members))
	for i, member := range members {
		addrs[i] = member.Address()
	}
	return &RendezvousV2{
		rdv:      rdv.New(addrs, xxhash.Sum64String),
		size:     len(addrs),
		dumpInfo: addrs,
	}
}

// Get ...
func (r *RendezvousV2) Get(key string) string {
	return r.rdv.Lookup(key)
}

// Size ...
func (r *RendezvousV2) Size() int {
	return r.size
}

// Dump is use for debug
func (r *RendezvousV2) Dump() []string {
	return r.dumpInfo
}
