package biofabric

type NodesByPopularity[N any, NPtr node[N]] []*fabricNode[N, NPtr]

func (byPop NodesByPopularity[N, NPtr]) Len() int {
	return len(byPop)
}
func (byPop NodesByPopularity[N, NPtr]) Less(i, j int) bool {
	if byPop[i].numLinks == byPop[j].numLinks {
		return byPop[i].node.Name() < byPop[j].node.Name()
	}
	return byPop[i].numLinks > byPop[j].numLinks
}
func (byPop NodesByPopularity[N, NPtr]) Swap(i, j int) {
	byPop[i], byPop[j] = byPop[j], byPop[i]
}
