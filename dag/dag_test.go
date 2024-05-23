package dag

import (
	"reflect"
	"testing"
)

// 创建DAG节点
var nodeMap map[string]*Node

func setNodeMap() {
	nodeMap = map[string]*Node{
		"nodeA": &Node{name: "A", status: StatusInit},
		"nodeB": &Node{name: "B", status: StatusInit},
		"nodeC": &Node{name: "C", status: StatusInit},
		"nodeD": &Node{name: "D", status: StatusInit},
		"nodeE": &Node{name: "E", status: StatusInit},
		"nodeF": &Node{name: "F", status: StatusInit},
		"nodeG": &Node{name: "G", status: StatusInit},
		"nodeH": &Node{name: "H", status: StatusInit},
		"nodeI": &Node{name: "I", status: StatusInit},
		"nodeJ": &Node{name: "J", status: StatusInit},
		"nodeK": &Node{name: "K", status: StatusInit},
		"nodeL": &Node{name: "L", status: StatusInit},
		"nodeM": &Node{name: "M", status: StatusInit},
		"nodeN": &Node{name: "N", status: StatusInit},
		"nodeO": &Node{name: "O", status: StatusInit},
		"nodeP": &Node{name: "P", status: StatusInit},
		"nodeQ": &Node{name: "Q", status: StatusInit},
		"nodeR": &Node{name: "R", status: StatusInit},
		"nodeS": &Node{name: "S", status: StatusInit},
		"nodeT": &Node{name: "T", status: StatusInit},
		"nodeU": &Node{name: "U", status: StatusInit},
		"nodeV": &Node{name: "V", status: StatusInit},
		"nodeW": &Node{name: "W", status: StatusInit},
		"nodeX": &Node{name: "X", status: StatusInit},
		"nodeY": &Node{name: "Y", status: StatusInit},
		"nodeZ": &Node{name: "Z", status: StatusInit},
	}
}

func setUpChain() *Node {
	setNodeMap()
	// 构建节点关系-链式
	// A->B->C->D->E
	nodeMap["nodeA"].AddChild(nodeMap["nodeB"])
	nodeMap["nodeB"].AddChild(nodeMap["nodeC"])
	nodeMap["nodeC"].AddChild(nodeMap["nodeD"])
	nodeMap["nodeD"].AddChild(nodeMap["nodeE"])
	return nodeMap["nodeA"]
}

func setUpBranch() *Node {
	setNodeMap()
	// 构建节点关系-分支
	// A->B->D->H
	// A->B->E->I
	// A->C->F->J
	// A->C->F->K
	// A->C->G->L
	// A->C->G->M
	nodeMap["nodeA"].AddChild(nodeMap["nodeB"])
	nodeMap["nodeA"].AddChild(nodeMap["nodeC"])
	nodeMap["nodeB"].AddChild(nodeMap["nodeD"])
	nodeMap["nodeB"].AddChild(nodeMap["nodeE"])
	nodeMap["nodeC"].AddChild(nodeMap["nodeF"])
	nodeMap["nodeC"].AddChild(nodeMap["nodeG"])
	nodeMap["nodeD"].AddChild(nodeMap["nodeH"])
	nodeMap["nodeE"].AddChild(nodeMap["nodeI"])
	nodeMap["nodeF"].AddChild(nodeMap["nodeJ"])
	nodeMap["nodeF"].AddChild(nodeMap["nodeK"])
	nodeMap["nodeG"].AddChild(nodeMap["nodeL"])
	nodeMap["nodeG"].AddChild(nodeMap["nodeM"])
	return nodeMap["nodeA"]
}

func setUpBranchAgg() *Node {
	setNodeMap()
	// 构建节点关系-分支+聚合+分支+聚合
	// A->B->D->F->H
	// A->C->D->G->H
	nodeMap["nodeA"].AddChild(nodeMap["nodeB"])
	nodeMap["nodeA"].AddChild(nodeMap["nodeC"])
	nodeMap["nodeB"].AddChild(nodeMap["nodeD"])
	nodeMap["nodeC"].AddChild(nodeMap["nodeD"])
	nodeMap["nodeD"].AddChild(nodeMap["nodeF"])
	nodeMap["nodeD"].AddChild(nodeMap["nodeG"])
	nodeMap["nodeF"].AddChild(nodeMap["nodeH"])
	nodeMap["nodeG"].AddChild(nodeMap["nodeH"])
	return nodeMap["nodeA"]
}

func setUpMultiBranchAgg() *Node {
	setNodeMap()
	// 构建节点关系-多聚合+分支
	// A->D->E->H->N->Q->T->V
	// A->D->E->I->N->Q->T->V
	// B->D->F->J->N->Q->T->V
	// B->D->F->K->O->R->T->V
	// C->D->G->L->O->S->T->V
	// C->D->G->M->P->T->V
	nodeMap["nodeA"].AddChild(nodeMap["nodeD"])
	nodeMap["nodeB"].AddChild(nodeMap["nodeD"])
	nodeMap["nodeC"].AddChild(nodeMap["nodeD"])
	nodeMap["nodeD"].AddChild(nodeMap["nodeE"])
	nodeMap["nodeD"].AddChild(nodeMap["nodeF"])
	nodeMap["nodeD"].AddChild(nodeMap["nodeG"])
	nodeMap["nodeE"].AddChild(nodeMap["nodeH"])
	nodeMap["nodeE"].AddChild(nodeMap["nodeI"])
	nodeMap["nodeF"].AddChild(nodeMap["nodeJ"])
	nodeMap["nodeF"].AddChild(nodeMap["nodeK"])
	nodeMap["nodeG"].AddChild(nodeMap["nodeL"])
	nodeMap["nodeG"].AddChild(nodeMap["nodeM"])
	nodeMap["nodeH"].AddChild(nodeMap["nodeN"])
	nodeMap["nodeI"].AddChild(nodeMap["nodeN"])
	nodeMap["nodeJ"].AddChild(nodeMap["nodeN"])
	nodeMap["nodeK"].AddChild(nodeMap["nodeO"])
	nodeMap["nodeL"].AddChild(nodeMap["nodeO"])
	nodeMap["nodeM"].AddChild(nodeMap["nodeP"])
	nodeMap["nodeN"].AddChild(nodeMap["nodeQ"])
	nodeMap["nodeO"].AddChild(nodeMap["nodeR"])
	nodeMap["nodeO"].AddChild(nodeMap["nodeS"])
	nodeMap["nodeP"].AddChild(nodeMap["nodeT"])
	nodeMap["nodeQ"].AddChild(nodeMap["nodeT"])
	nodeMap["nodeR"].AddChild(nodeMap["nodeT"])
	nodeMap["nodeS"].AddChild(nodeMap["nodeT"])
	nodeMap["nodeT"].AddChild(nodeMap["nodeV"])
	return nodeMap["nodeA"]
}
func TestNode_DriveChain(t *testing.T) {
	initToSucceed := func(nodes []*Node) {
		for _, v := range nodes {
			v.status = StatusSucceed
		}
	}
	type fields struct {
		node   *Node
		afterF func([]*Node)
	}
	tests := []struct {
		name   string
		fields fields
		want   []*Node
		want1  bool
	}{
		{name: "t1-1", fields: fields{node: setUpChain(), afterF: initToSucceed}, want: []*Node{nodeMap["nodeA"]}, want1: false},
		{name: "t1-1", fields: fields{node: nodeMap["nodeA"], afterF: initToSucceed}, want: []*Node{nodeMap["nodeB"]}, want1: false},
		{name: "t1-1", fields: fields{node: nodeMap["nodeA"], afterF: initToSucceed}, want: []*Node{nodeMap["nodeC"]}, want1: false},
		{name: "t1-1", fields: fields{node: nodeMap["nodeA"], afterF: initToSucceed}, want: []*Node{nodeMap["nodeD"]}, want1: false},
		{name: "t1-1", fields: fields{node: nodeMap["nodeA"], afterF: initToSucceed}, want: []*Node{nodeMap["nodeE"]}, want1: false},
		{name: "t1-1", fields: fields{node: nodeMap["nodeA"], afterF: initToSucceed}, want: []*Node{}, want1: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := tt.fields.node
			got, got1 := n.Drive()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Node.Drive() = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Node.Drive() = %v, want %v", got1, tt.want1)
			}
			tt.fields.afterF(got)
		})
	}
}

func TestNode_DriveBranch(t *testing.T) {
	initToSucceed := func(nodes []*Node) {
		for _, v := range nodes {
			v.status = StatusSucceed
		}
	}
	type fields struct {
		node   *Node
		afterF func([]*Node)
	}
	tests := []struct {
		name   string
		fields fields
		want   []*Node
		want1  bool
	}{
		{name: "t2-1", fields: fields{node: setUpBranch(), afterF: initToSucceed}, want: []*Node{nodeMap["nodeA"]}, want1: false},
		{name: "t2-2", fields: fields{node: nodeMap["nodeA"], afterF: initToSucceed}, want: []*Node{nodeMap["nodeB"], nodeMap["nodeC"]}, want1: false},
		{name: "t2-3", fields: fields{node: nodeMap["nodeA"], afterF: initToSucceed}, want: []*Node{nodeMap["nodeD"], nodeMap["nodeE"], nodeMap["nodeF"], nodeMap["nodeG"]}, want1: false},
		{name: "t2-4", fields: fields{node: nodeMap["nodeA"], afterF: initToSucceed}, want: []*Node{nodeMap["nodeH"], nodeMap["nodeI"], nodeMap["nodeJ"], nodeMap["nodeK"], nodeMap["nodeL"], nodeMap["nodeM"]}, want1: false},
		{name: "t2-5", fields: fields{node: nodeMap["nodeA"], afterF: initToSucceed}, want: []*Node{}, want1: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := tt.fields.node
			got, got1 := n.Drive()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Node.Drive() = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Node.Drive() = %v, want %v", got1, tt.want1)
			}
			for _, v := range got {
				v.status = StatusSucceed
			}
		})
	}
}

func TestNode_DriveBranchAgg(t *testing.T) {
	initToSucceed := func(nodes []*Node) {
		for _, v := range nodes {
			v.status = StatusSucceed
		}
	}
	type fields struct {
		node   *Node
		afterF func([]*Node)
	}
	tests := []struct {
		name   string
		fields fields
		want   []*Node
		want1  bool
	}{
		{name: "t3-1", fields: fields{node: setUpBranchAgg(), afterF: initToSucceed}, want: []*Node{nodeMap["nodeA"]}, want1: false},
		{name: "t3-2", fields: fields{node: nodeMap["nodeA"], afterF: initToSucceed}, want: []*Node{nodeMap["nodeB"], nodeMap["nodeC"]}, want1: false},
		{name: "t3-3", fields: fields{node: nodeMap["nodeA"], afterF: initToSucceed}, want: []*Node{nodeMap["nodeD"]}, want1: false},
		{name: "t3-4", fields: fields{node: nodeMap["nodeA"], afterF: initToSucceed}, want: []*Node{nodeMap["nodeF"], nodeMap["nodeG"]}, want1: false},
		{name: "t3-5", fields: fields{node: nodeMap["nodeA"], afterF: initToSucceed}, want: []*Node{nodeMap["nodeH"]}, want1: false},
		{name: "t3-6", fields: fields{node: nodeMap["nodeA"], afterF: initToSucceed}, want: []*Node{}, want1: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := tt.fields.node
			got, got1 := n.Drive()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Node.Drive() = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Node.Drive() = %v, want %v", got1, tt.want1)
			}
			for _, v := range got {
				v.status = StatusSucceed
			}
		})
	}
}

func TestNode_DriveMultiBranchAgg(t *testing.T) {
	initToSucceed := func(nodes []*Node) {
		for _, v := range nodes {
			v.status = StatusSucceed
		}
	}
	type fields struct {
		node   *Node
		afterF func([]*Node)
	}
	tests := []struct {
		name   string
		fields fields
		want   []*Node
		want1  bool
	}{
		{name: "t3-1", fields: fields{node: setUpMultiBranchAgg(), afterF: initToSucceed}, want: []*Node{nodeMap["nodeA"]}, want1: false},
		{name: "t3-2", fields: fields{node: nodeMap["nodeB"], afterF: initToSucceed}, want: []*Node{nodeMap["nodeB"]}, want1: false},
		{name: "t3-3", fields: fields{node: nodeMap["nodeC"], afterF: initToSucceed}, want: []*Node{nodeMap["nodeC"]}, want1: false},
		{name: "t3-4", fields: fields{node: nodeMap["nodeA"], afterF: initToSucceed}, want: []*Node{nodeMap["nodeD"]}, want1: false},
		{name: "t3-5", fields: fields{node: nodeMap["nodeA"], afterF: initToSucceed}, want: []*Node{nodeMap["nodeE"], nodeMap["nodeF"], nodeMap["nodeG"]}, want1: false},
		{name: "t3-5", fields: fields{node: nodeMap["nodeA"], afterF: initToSucceed}, want: []*Node{nodeMap["nodeH"], nodeMap["nodeI"], nodeMap["nodeJ"], nodeMap["nodeK"], nodeMap["nodeL"], nodeMap["nodeM"]}, want1: false},
		{name: "t3-5", fields: fields{node: nodeMap["nodeA"], afterF: initToSucceed}, want: []*Node{nodeMap["nodeN"], nodeMap["nodeO"], nodeMap["nodeP"]}, want1: false},
		{name: "t3-5", fields: fields{node: nodeMap["nodeA"], afterF: initToSucceed}, want: []*Node{nodeMap["nodeQ"], nodeMap["nodeR"], nodeMap["nodeS"]}, want1: false},
		{name: "t3-5", fields: fields{node: nodeMap["nodeA"], afterF: initToSucceed}, want: []*Node{nodeMap["nodeT"]}, want1: false},
		{name: "t3-5", fields: fields{node: nodeMap["nodeA"], afterF: initToSucceed}, want: []*Node{nodeMap["nodeV"]}, want1: false},
		{name: "t3-6", fields: fields{node: nodeMap["nodeA"], afterF: initToSucceed}, want: []*Node{}, want1: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := tt.fields.node
			got, got1 := n.Drive()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Node.Drive() = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Node.Drive() = %v, want %v", got1, tt.want1)
			}
			for _, v := range got {
				v.status = StatusSucceed
			}
		})
	}
}
