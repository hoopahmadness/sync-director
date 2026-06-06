package biofabric

const (
	NULLEdgeMode     = EdgeMode("")
	NEUTRALEdgeMode  = EdgeMode("neutral")
	INWARDEdgeMode   = EdgeMode("inward")
	OUTWARDEdgeMode  = EdgeMode("outward")
	NEGATIVEEdgeMode = EdgeMode("negative")
)

type EdgeMode string

type EdgeAttributes struct {
	Mode      EdgeMode
	Highlight bool
}

type edge struct {
	topRow     int
	bottomRow  int
	attributes map[int]EdgeAttributes
}
