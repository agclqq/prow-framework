package dag

import (
	"errors"
	"fmt"
	"strings"
)

type status int
type FvsFunc func(n *Node) bool

const (
	StatusInit status = iota
	StatusRunning
	StatusSucceed
	StatusFailed
)

// Node 节点
// 一个节点有三种状态：StatusInit、StatusRunning、StatusSucceed
// 一个节点有多个父节点和多个子节点
type Node struct { //
	name     string
	status   status
	parents  []*Node
	children []*Node
}

// AddParent 添加父节点到当前节点
func (n *Node) AddParent(parent *Node) {
	n.parents = append(n.parents, parent)
}

// AddChild 添加子节点到当前节点
func (n *Node) AddChild(child *Node) {
	for c := range n.children {
		if n.children[c].name == child.name {
			return
		}
	}
	n.children = append(n.children, child)
	child.AddParent(n)
}

// InsertBefore 在指定节点之前插入一个节点
// 指定节点必须没有父节点，否则不知道要插入到哪里
func (n *Node) InsertBefore(node *Node) error {
	if node == nil {
		return fmt.Errorf("node is nil")
	}
	if len(node.parents) != 0 {
		//如果node有父节点，则不插入
		return errors.New("node has parent, should use InsertBetween")
	}
	node.AddParent(n)
	n.AddChild(node)
	return nil
}

// InsertAfter 在指定节点之后插入一个节点
// 指定节点必须没有子节点，否则不知道要插入到哪里
func (n *Node) InsertAfter(node *Node) error {
	if node == nil {
		return fmt.Errorf("node is nil")
	}
	if len(node.children) != 0 {
		return errors.New("node has children, should use InsertBetween")
	}
	node.AddChild(n)
	n.AddParent(node)
	return nil
}

// InsertBetween 在两个节点之间插入一个节点
func (n *Node) InsertBetween(front, behind *Node) error {
	if front == nil && behind == nil {
		return errors.New("front and behind are nil")
	}
	if front == nil {
		// 如果front为nil，则插入到behind之前
		return n.InsertBefore(behind)
	}
	if behind == nil {
		// 如果behind为nil，则插入到front之后
		return n.InsertAfter(front)
	}
	for i := 0; i < len(front.children); i++ {
		if front.children[i].name == behind.name {
			n.children = append(n.children, behind)
			front.children = append(front.children[:i], append([]*Node{n}, front.children[i:]...)...)
			return nil
		}
	}
	return nil
}

// CanEnterRunning 检查当前节点是否可以进入running状态
func (n *Node) CanEnterRunning() bool {
	if n.status != StatusInit {
		return false // 如果节点不在notStart状态，则不能进入running状态
	}
	// 检查所有父节点是否都已完成
	for _, parent := range n.parents {
		if parent.status != StatusSucceed {
			return false
		}
	}
	return true
}

// DFS 深度优先遍历
func (n *Node) DFS(f FvsFunc) {
	// 处理当前节点
	if !f(n) {
		return // 如果f返回false，则退出遍历
	}
	for _, child := range n.children {
		child.DFS(f)
	}
}

// BFS 广度优先遍历
func (n *Node) BFS(f FvsFunc) {
	queue := []*Node{n}
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		if !f(current) {
			continue // 如果f返回false，则退出当前分支的遍历
		}
		for _, child := range current.children {
			queue = append(queue, child)
		}
	}
}

// EnterRunning 将节点的状态从notStart更改为running（如果符合条件）
func (n *Node) EnterRunning() {
	if n.CanEnterRunning() {
		n.status = StatusRunning
		fmt.Printf("%s is entering the StatusRunning status.\n", n.name)
	}
}

func (n *Node) EnterFinished() {
	if n.status == StatusRunning {
		n.status = StatusSucceed
		fmt.Printf("%s is entering the StatusSucceed status.\n", n.name)
	}
}

// Verify 验证DAG是否合法
func (n *Node) Verify() bool {
	return false
}

// Drive 驱动DAG
// 返回节点map和流水线是否终止，是否终止的含义为：是否还有应该执行但还未执行的节点。
// 如果没有失败的节点，且没有未初化、未执行的节点，则流水线终止。
// 如果有失败的节点，则流水线终止。
func (n *Node) Drive() ([]*Node, bool) {
	var readyToRunning = make([]*Node, 0)
	end := true
	n.DFS(func(curr *Node) bool {
		if curr.status == StatusFailed {
			end = true
		}
		if curr.status == StatusRunning {
			end = false
			return false // 如果节点已经在running状态，则不再遍历其子节点
		}
		if curr.status == StatusInit {
			end = false
			if curr.CanEnterRunning() {
				//如果当前节点是聚合节点，先判断一下是否已在readyToRunning中
				for _, r := range readyToRunning {
					if r.name == curr.name {
						return false
					}
				}
				readyToRunning = append(readyToRunning, curr)
				//当前节点可以进入running状态，只返回当前节点
			}
			return false // 如果节点是notStart状态，则不再遍历其子节点
		}
		return true
	})
	return readyToRunning, end
}

// PrintDAG 打印DAG示意图
func (n *Node) PrintDAG() {
	var sb strings.Builder

	// DFS遍历，并构建图形化表示
	var dfs func(node *Node, indent string)
	dfs = func(node *Node, indent string) {
		// 打印当前节点
		sb.WriteString(indent)
		sb.WriteString(node.name + " " + node.status.String())
		sb.WriteString("\n")

		// 遍历子节点
		for _, child := range node.children {
			// 为每个子节点添加箭头并递归调用dfs
			sb.WriteString(indent + "  ")
			sb.WriteString("-> ")
			dfs(child, indent+"  ") // 缩进以表示层次
		}
	}
	dfs(n, "") // 从根节点开始遍历
	// 打印DAG示意图
	fmt.Print(sb.String())
}

func main() {
	// 创建DAG节点
	nodeA := &Node{name: "A", status: StatusInit}
	nodeB := &Node{name: "B", status: StatusInit}
	nodeC := &Node{name: "C", status: StatusInit}
	nodeD := &Node{name: "D", status: StatusInit}
	nodeE := &Node{name: "E", status: StatusInit}

	// 构建节点关系
	nodeA.AddChild(nodeB)
	nodeA.AddChild(nodeC)
	nodeB.AddParent(nodeA)
	nodeC.AddParent(nodeA)
	nodeC.AddChild(nodeD)
	nodeD.AddParent(nodeC)
	nodeB.AddChild(nodeE)
	nodeD.AddChild(nodeE)
	nodeE.AddParent(nodeB)
	nodeE.AddParent(nodeD)

	//打印DAG示意图
	nodeA.PrintDAG()
	//驱动
	fmt.Println("DAG Drive:")
	//for rfr := nodeA.Drive(); len(rfr) > 0; rfr = nodeA.Drive() {
	//	//fmt.Printf("%v\n", rfr)
	//	for _, n := range rfr {
	//		fmt.Println("=====")
	//		fmt.Println("got:" + n.name)
	//		fmt.Println("------")
	//		time.Sleep(3 * time.Second)
	//		n.EnterRunning()
	//		n.EnterFinished()
	//		//fmt.Println("------")
	//	}
	//	nodeA.PrintDAG()
	//}
}

// printNodeStatuses 递归打印所有节点的状态
func printNodeStatuses(node *Node) {
	fmt.Printf("%s: %s\n", node.name, node.status.String())
	for _, child := range node.children {
		printNodeStatuses(child)
	}
}

// String 将status转换为字符串
func (s status) String() string {
	switch s {
	case StatusInit:
		return "StatusInit"
	case StatusRunning:
		return "StatusRunning"
	case StatusSucceed:
		return "StatusSucceed"
	default:
		return "unknown"
	}
}
