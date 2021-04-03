package memviz

import (
	"io"
	"strings"
)

type cell struct {
	port string
	name string
}

type field struct {
	cells []cell
}

type cnode struct {
	id     nodeID
	name   string
	fields []field
}

func createNode(id nodeID, name string) *cnode {
	node := cnode{
		id:   id,
		name: name,
	}
	return &node
}

func (s *cnode) addFieldInlined(port string, name, summary string) {
	s.fields = append(s.fields, field{
		cells: []cell{
			cell{
				port: port,
				name: name,
			},
			cell{
				name: summary,
			},
		},
	})
}

func (s *cnode) addField(port string, name string) {
	s.fields = append(s.fields, field{
		cells: []cell{
			cell{
				port: port,
				name: name,
			},
		},
	})
}

func (s *cnode) addFields(port1 string, name1 string, port2 string, name2 string) {
	s.fields = append(s.fields, field{
		cells: []cell{
			cell{
				port: port1,
				name: name1,
			},
			cell{
				port: port2,
				name: name2,
			},
		},
	})
}

func (m *mapper) addConnection(fromNode nodeID, port string, toNode nodeID) {
	m.connections = append(m.connections, connection{
		fromNode: fromNode,
		fromPort: port,
		toNode:   toNode,
		toPort:   portTitle,
	})
}

func (m *mapper) addNode(node *cnode) {
	m.nodes = append(m.nodes, node)
}

func (m *mapper) write(w io.Writer) {
	m.optimize()
	Mrecord(w, m.nodes, m.connections, m.comment)
}

func (m *mapper) optimize() {

	if Options().CollapsePointerNodes || Options().CollapseSingleSliceNodes {

		direct := make(map[nodeID][]nodeID)
		reverse := make(map[nodeID][]nodeID)

		for _, conn := range m.connections {
			direct[conn.fromNode] = append(direct[conn.fromNode], conn.toNode)
			reverse[conn.toNode] = append(reverse[conn.toNode], conn.fromNode)
		}

		access := make(map[nodeID]*cnode)
		for _, node := range m.nodes {
			access[node.id] = node
		}

		singleUse := make([]nodeID, 0)
		remap := make(map[nodeID]nodeID)

		for _, node := range m.nodes {
			from := node.id

			if len(direct[from]) == 1 {
				to := direct[from][0]
				// if len(reverse[to]) == 1 {
				if len(node.fields) <= 1 { // == 0
					parts := strings.Split(node.name, ".")
					suffix := parts[len(parts)-1]

					toName := access[to].name
					if toName == node.name || toName == suffix {

						singleUse = append(singleUse, from)
						remap[from] = to
					}
				}
				// }
			}
		}

		var connections []connection
		for _, conn := range m.connections {
			from := conn.fromNode
			to := conn.toNode

			if newTo, exists := remap[to]; exists {
				conn.toNode = newTo

				if len(access[to].fields) == 0 {
					conn.style = 1
				} else {
					conn.style = 2
				}

			} else if _, exists := remap[from]; exists {
				continue
			}

			connections = append(connections, conn)
		}
		m.connections = connections

		for _, id := range singleUse {
			delete(access, id)
		}

		var nodes []*cnode
		for _, value := range access {
			nodes = append(nodes, value)
		}
		m.nodes = nodes
	}
}
