package biofabric

import (
	"fmt"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	log "github.com/inconshreveable/log15"
)

const (
	Reset             = "\033[0m"
	Red               = "\033[31m"
	Green             = "\033[32m"
	neutralEdgeTop    = "▣"
	neutralEdgeBottom = "▣"
	negativeEdge      = "x"
	emptyEdgeTop      = "◯"
	emptyEdgeBottom   = "◯"
	inwardEdgeTop     = "▼"
	inwardEdgeBottom  = "▲"
	outwardEdgeTop    = "▲"
	outwardEdgeBottom = "▼"
	onlyEdgeLine      = "┃"
	onlyRowLine       = "─"
	cellCrossed       = "┃\u0336"
	cellCrossed2      = "╂"
	empty             = " "
)

var topEdgeLookup = map[EdgeMode]string{
	NULLEdgeMode:     emptyEdgeTop,
	NEUTRALEdgeMode:  neutralEdgeTop,
	INWARDEdgeMode:   inwardEdgeTop,
	OUTWARDEdgeMode:  outwardEdgeTop,
	NEGATIVEEdgeMode: negativeEdge,
}

var bottomEdgeLookup = map[EdgeMode]string{
	NULLEdgeMode:     emptyEdgeBottom,
	NEUTRALEdgeMode:  neutralEdgeBottom,
	INWARDEdgeMode:   inwardEdgeBottom,
	OUTWARDEdgeMode:  outwardEdgeBottom,
	NEGATIVEEdgeMode: negativeEdge,
}

/**

  Nodes are represented as one-dimensional horizontal line segments, one per row.
  Edges are represented as one-dimensional vertical line segments, one per column, terminating at the two rows associated with the endpoint nodes.
  Edges are drawn darker than nodes; this has the effect of emphasizing the links and making them appear to float in front of the nodes.
  Both ends of a link are represented as a tiny square. This provides sufficient contrast to make the ends of the link stand out even at large scales. In the case of directed edges, the appropriate end is tagged with an arrowhead.
  Edges are unambiguously represented and never overlap.
  In networks that have multiple edges between the same nodes, i.e. representing different types of relationships, all edges show up clearly.
  As nodes are represented as horizontal lines, there is no requirement that all edges converge upon a single point, allowing for complete flexibility in where a link is drawn.
  Links can originate, and terminate, anywhere along the length of the node segment. This flexibility introduces the powerful ability to create sets of links that share some semantic property and depict them as discrete groups in the visualization.
  The addition of a new edge just increases the width of the visualization, and does not degrade the existing presentation in any fashion.
  The visualization technique produces a distinct edge wedge for each node, created by the close-set juxtaposition of the parallel links, that provides clear visual cues about how the node is connected, and how it compares to other similar nodes.
  A set of 32 colors is used, not randomly, but in a repeating cycle to render node and edge segments. Colors are not used to apply semantic meaning to network elements, but are crucial for providing a framework that allows the user to visually trace features over long distances. Also, the use of cycling insures that antialiased rendering will produce larger-scale color patterns that provide useful visual cues even when individual links cannot be discerned.
  Note that the traditional technique overloads the two-dimensional plane, using the same space to represent both nodes and edges. BioFabric effectively segregates the plane into two separate one-dimensional spaces, and assigns each space to either nodes or edges; the imposition of orthogonality and the use of judicious rendering allow the user to visually distinguish the two. Thus, BioFabric can provide additional clarity of the network structure while using the same underlying two-dimensional resource.
*/

type fabricNode[N any, NPtr node[N]] struct {
	node     NPtr
	numLinks int
}

func (fn *fabricNode[N, NPtr]) addLink() {
	fn.numLinks += 1
}

type Fabric[N any, NPtr node[N]] struct {
	web web[N, NPtr]
	// nodes     []*N
	collapsed bool
	fLog      log.Logger
}

// returns a string representation of the web with default settings
/**
┰
┃ box drawings heavy vert
┃
┃
┃
▣
┸
┸─────┰ box drawings light horizontal
┸ box drawings up heavy and horizontal light
┰ box drawings down heavy and...
╂ box drawing vert heavy horiz light

first guy ───▣──▼──▲
             ┃  ┃  ┃
second guy ──▣──╂──╂──┰
                ┃  ┃  ┃
third guy ──────▼──╂──╂─
                   ┃  ┃
fourth guy ────────▲──┸

*/
func (f Fabric[N, NPtr]) View() string {
	renderedStr := "boop\n"
	listNodes := f.createSortedList()
	for _, node := range listNodes {
		renderedStr += "\n"
		renderedStr += node.node.Name()
		renderedStr += "\n"
	}
	edges := f.generateEdges(listNodes)
	grid := f.generateGrid(listNodes, edges)
	names := f.justify(listNodes)
	for rowCount, rowRunes := range grid {
		// f.fLog.Info("above the thing", "rowcount", rowCount, "runes", string(rowRunes), "names", names)
		rowStr := names[rowCount] + string(rowRunes) + "\n"
		renderedStr += rowStr
	}
	renderedStr += "boopt 2"
	return renderedStr
}

// Generates an array of arrays of runes that will be joined into one long string
// each rune represents one "cell" of the biofabric
// Nodes are represented as rows; Cells in a row are either an edge terminator (┰, ┸, ▣, ▼, ▲)
// or a row line with or without intersection (─, ╂)
// or a space, which is used before and after the first and last edge terminators in that row
// Also there may be lines between each node which will mostly be empty space
// Edges are represented as columns; cells in a column are either the edge terminators
// or a column line with our without intersetction (┃, ╂)
// or a space or horizontal line, which is used above and/or below the edge value
// Also there may be lines between edges which will mostly be empty space
func (f Fabric[N, NPtr]) generateGrid(nodes []*fabricNode[N, NPtr], edges []edge) [][]rune {
	grid := [][]rune{}
	for nodeCount := 0; nodeCount < len(nodes); nodeCount++ {
		nodeLine := []rune{}
		spacerLine := []rune{}

		for edgeCount := 0; edgeCount < len(edges); edgeCount++ {
			cell := empty
			spacer01 := empty // spacers are the other 3 cells that combine with the main cell to form a 2x2 grid
			spacer10 := empty
			spacer11 := empty

			thisEdge := edges[edgeCount]
			// thisNode := nodes[nodeCount]
			firstEdge, lastEdge := getEdgeIntersections(nodeCount, edges)
			cellTopOfEdge := nodeCount == thisEdge.topRow
			cellBottomOfEdge := nodeCount == thisEdge.bottomRow
			cellWithinEdge := nodeCount > thisEdge.topRow && nodeCount < thisEdge.bottomRow

			cellWithinRow := edgeCount >= firstEdge && edgeCount <= lastEdge-1
			switch {
			case cellTopOfEdge:
				cell = topEdgeLookup[thisEdge.attributes[nodeCount].Mode]
				spacer10 = onlyEdgeLine
				if cellWithinRow {
					spacer01 = onlyRowLine
				} else {
					spacer01 = empty
				}
			case cellBottomOfEdge:
				cell = bottomEdgeLookup[thisEdge.attributes[nodeCount].Mode]
				spacer10 = empty
				if cellWithinRow {
					spacer01 = onlyRowLine
				} else {
					spacer01 = empty
				}
			case cellWithinEdge && cellWithinRow:
				cell = cellCrossed
				spacer01 = onlyRowLine
				spacer10 = onlyEdgeLine

			case cellWithinEdge: // but outside main row
				cell = onlyEdgeLine
				spacer01 = empty
				spacer10 = onlyEdgeLine

			case cellWithinRow: // but not inside an edge
				cell = onlyRowLine
				spacer01 = onlyRowLine
				spacer10 = empty
			}

			nodeLine = append(nodeLine, []rune(cell)...)
			if !f.collapsed {
				nodeLine = append(nodeLine, []rune(spacer01)[0])
				spacerLine = append(spacerLine, []rune(spacer10)[0], []rune(spacer11)[0])
			}
		}
		grid = append(grid, nodeLine)
		if !f.collapsed {
			grid = append(grid, spacerLine)
		}
	}
	return grid
}

func getEdgeIntersections(nodeNumber int, edges []edge) (firstEdge int, lastEdge int) {
	intersections := []int{}
	for edgeNum, edge := range edges {
		topNode := edge.topRow
		bottomNode := edge.bottomRow
		intersects := nodeNumber == topNode || nodeNumber == bottomNode
		if intersects {
			intersections = append(intersections, edgeNum)
		}
	}
	firstEdge = intersections[0]
	lastEdge = intersections[len(intersections)-1]
	return
}

func (f Fabric[N, NPtr]) generateEdges(rows []*fabricNode[N, NPtr]) []edge {
	edges := []edge{}
	chartedPairs := map[Pair[N, NPtr]]bool{}
	for topIndex, topNode := range rows {
		for bottomIndex, bottomNode := range rows {
			if topIndex == bottomIndex {
				continue
			}
			var thisPair Pair[N, NPtr]
			pairings := f.web.GetRelatedPairs(topNode.node)
			for _, pair := range pairings {
				if _, alreadyChartedPair := chartedPairs[pair]; pair.Other(topNode.node) == bottomNode.node && !alreadyChartedPair {
					thisPair = pair
					chartedPairs[pair] = true
				}
			}
			if thisPair == nil {
				continue
			}
			newEdge := edge{
				topRow:    topIndex,
				bottomRow: bottomIndex,
				attributes: map[int]EdgeAttributes{
					topIndex:    thisPair.GetNodeAttributes(topNode.node),
					bottomIndex: thisPair.GetNodeAttributes(bottomNode.node),
				},
			}
			edges = append(edges, newEdge)
		}
	}
	return edges
}

func (f Fabric[N, NPtr]) justify(nodes []*fabricNode[N, NPtr]) []string {
	maxLength := 0
	names := []string{}
	justifiedLines := []string{}
	extraSpaces := !f.collapsed
	for _, node := range nodes {
		name := node.node.Name()
		amnt := node.numLinks
		name = name + fmt.Sprintf(" (%d)", amnt)
		names = append(names, name)
		if len(name) > maxLength {
			maxLength = len(name)
		}
	}
	for _, name := range names {
		requiredExtra := maxLength - len(name)
		name += strings.Repeat(empty, requiredExtra)
		justifiedLines = append(justifiedLines, name)
		if extraSpaces {
			justifiedLines = append(justifiedLines, strings.Repeat(empty, maxLength))
		}
	}
	return justifiedLines
}

func (f Fabric[N, NPtr]) createSortedList() []*fabricNode[N, NPtr] {
	fabricNodesMap := map[string]*fabricNode[N, NPtr]{}
	fabricNodesArray := []*fabricNode[N, NPtr]{}
	for pair := range f.web.GetAllPairs() {
		for _, item := range pair.GetItems() {
			fabNode, exists := fabricNodesMap[item.Name()]
			if !exists {
				fabNode = &fabricNode[N, NPtr]{item, 0}
				fabricNodesMap[item.Name()] = fabNode
				fabricNodesArray = append(fabricNodesArray, fabNode)
			}
			fabNode.addLink()
		}
	}

	sort.Stable(NodesByPopularity[N, NPtr](fabricNodesArray))
	return fabricNodesArray
}

func (f Fabric[N, NPtr]) Init() tea.Cmd {
	return nil
}

func (f Fabric[N, NPtr]) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return nil, nil
}

func NewFabric[N any, NPtr node[N]](w web[N, NPtr], logger log.Logger) Fabric[N, NPtr] {
	return Fabric[N, NPtr]{
		web:       w,
		collapsed: false,
		fLog:      logger,
	}
}
