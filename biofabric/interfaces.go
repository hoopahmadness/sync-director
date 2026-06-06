package biofabric

type web[N any, NPtr node[N]] interface {
	// Given a single node, return all pairs of nodes containing that node
	GetRelatedPairs(NPtr) []Pair[N, NPtr]

	// Return all pairs in this web
	GetAllPairs() map[Pair[N, NPtr]]bool

	// Return the number of nodes inluded in this web.
	// NumNodes() int

	// Return all nodes included in this web
	// ListNodes() []N
}

type Pair[N any, NPtr node[N]] interface {
	Other(node NPtr) NPtr
	GetItems() [2]NPtr
	GetNodeAttributes(node NPtr) EdgeAttributes
}

type node[n any] interface {
	Name() string
	*n
}
