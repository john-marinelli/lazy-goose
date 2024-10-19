package dstructs

type PartFunc func(pg *Vertex) int

type PartitionedGraph struct {
	Verticies map[string]*Vertex
	HowToPart PartFunc
}

type Vertex struct {
	Level int
	Id    string
	Edges []*Edge
}

type Edge struct {
	Vertex *Vertex
}

type AdjList map[string][]string

func NewPartitionedGraph(howToPart PartFunc, adjList AdjList) (*PartitionedGraph, error) {
	pg := new(PartitionedGraph)
	pg.HowToPart = howToPart
	for k := range adjList {
		pg.Verticies[k] = &Vertex{
			Id: k,
		}
	}

	for k, v := range adjList {
		for _, n := range v {
			pg.Verticies[k].Edges = append(pg.Verticies[k].Edges, &Edge{
				Vertex: pg.Verticies[n],
			})
		}
	}

	return pg, nil
}

func howToPart(v *Vertex) int {
	l := 0
	var descend func(v *Vertex, currL int)
	var result []int

	descend = func(v *Vertex, currL int) {
		for _, e := range v.Edges {
			l = e.Vertex.Level
			descend(e.Vertex, currL+1)
		}
		result = append(result, currL)
	}

	currV := v
	for len(currV.Edges) > 0 {

	}

	return l
}
