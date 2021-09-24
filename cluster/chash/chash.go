package chash

// ConsistentHash ...
type ConsistentHash interface {
	Get(string) string
	Size() int
	Dump() string
}
