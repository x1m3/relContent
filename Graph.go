package main

import (
	"sync"
	"sort"
	"os"
	"log"
	"bufio"
	"encoding/json"
)

type Graph struct{
	sync.Mutex
	nodes map[string]Edges
}

func newGraph() (*Graph){
	graph := new(Graph)
	graph.nodes = make(map[string]Edges)
	return graph
}

func (graph *Graph) relatedNodes(id string) VertexSlice  {
var list VertexSlice
	graph.Lock()

	edges, ok := graph.nodes[id]
	if (ok) {
		for _,vertex := range edges.List {
			list = append(list,vertex)
		}
	}
	graph.Unlock()
	sort.Sort(list)
	return list
}


func (graph *Graph) addNodeIfNotExists(id string) {
var ok bool
	graph.Lock()
	_, ok = graph.nodes[id]
	if (!ok) {
		newNode := newEdges(id)
		graph.nodes[id]=*newNode
	}
	graph.Unlock()
}

func (graph *Graph) annotateRelation(node1Id string, node2Id string, weight int) {
var vertex Vertex
var ok bool
	graph.addNodeIfNotExists(node1Id)
	graph.addNodeIfNotExists(node2Id)
	graph.Lock()
	node1,_ := graph.nodes[node1Id]
	vertex, ok = node1.List[node2Id]
	if !ok {
		newVertex := newVertex(node2Id, weight)
		node1.List[node2Id] = *newVertex
	} else {
		vertex.weight += weight
		node1.List[node2Id] = vertex
	}
	graph.Unlock()
}

func (graph *Graph) annotateRelationBidirectional(node1Id string, node2Id string, weight int) {
	graph.annotateRelation(node1Id, node2Id, weight)
	graph.annotateRelation(node2Id, node1Id, weight)
}

/*
 * TODO: Inject a repo instead of reading directly from file
 */
func (graph *Graph) loadData(config *Config) {

	file, err := os.Open(config.Global.DatabaseFilename)
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	type UserLikeTemp struct {
		UserId       string  `json:"user_id"`
		ContentIds   []string `json:"content_ids"`
		LastViewed   string   `json:"lastViewed"`
		Consolidated bool     `json:"consolidated"`
	}
	var session UserLikeTemp

	i:=0
	for scanner.Scan() {
		line := scanner.Text()
		err := json.Unmarshal([]byte(line), &session)
		if (err!=nil) {
			log.Fatal(err)
		}
		for _,id1:= range session.ContentIds {
			for _,id2:= range session.ContentIds {
				if id1!=id2 {
					graph.annotateRelationBidirectional(id1, id2, 1)
				}
			}
		}
		i++
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

type Edges struct {
	id string
	List map[string]Vertex
}

// Edges Construct
func newEdges(id string) (*Edges){
	edge := new(Edges)
	edge.id = id
	edge.List = make(map[string]Vertex)
	return edge
}

type Vertex struct{
	edgeId string
	weight int
}

// Vertex construct
func newVertex(edgeId string, weight int) (*Vertex) {
	vertex := new(Vertex)
	vertex.edgeId= edgeId
	vertex.weight= weight
	return vertex;
}



type VertexSlice []Vertex

func (slice VertexSlice) Len() int {
	return len (slice)
}

func (slice VertexSlice) Less(i,j int) bool {
	return slice[i].weight > slice[j].weight
}

func (slice VertexSlice)  Swap(i,j int) {
	slice[i], slice[j] = slice[j], slice[i]
}


