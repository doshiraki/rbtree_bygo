package main

import (
	"fmt"
	"os"
)

//RBNodeDir is
type RBNodeDir int

const (

	//RBNodeLeft 左方向
	RBNodeLeft RBNodeDir = iota

	//RBNodeRight 右方向
	RBNodeRight

	//RBNodeHere ここ
	RBNodeHere
)

//RBNode is
type RBNode struct {
	isRed    bool
	index    int
	parent   *RBNode
	children [2]*RBNode
}

//RBCursor is
type RBCursor struct {
	node *RBNode
	less func(i int, j int) bool
}

func (cur *RBCursor) clone() *RBCursor {
	wrk := *cur
	return &wrk
}

func (cur *RBCursor) root() *RBCursor {

	if cur.node != nil {
		for ; cur.node.parent != nil; cur.node = cur.node.parent {

		}
	}
	return cur
}

func (cur *RBCursor) move(dir RBNodeDir) *RBCursor {
	rev := dir ^ 1
	if next := cur.node.children[dir]; next != nil {
		for next.children[rev] != nil {
			next = next.children[rev]
		}
		cur.node = next
	} else {
		find := false
		for cur.node.parent != nil {
			now := cur.node
			cur.node = cur.node.parent
			if cur.node.children[rev] == now {
				find = true
				break
			}
		}
		if !find {
			cur.node = nil
		}
	}
	return cur
}
func (cur *RBCursor) find(index int) (*RBCursor, RBNodeDir) {
	cur.root()
	dir := RBNodeHere
	for {
		if cur.less(index, cur.node.index) {
			dir = RBNodeLeft
			next := cur.node.children[dir]
			if next == nil {
				break
			}
			cur.node = next

		} else if cur.less(cur.node.index, index) {
			dir = RBNodeRight
			next := cur.node.children[dir]
			if next == nil {
				break
			}
			cur.node = next
		} else {
			dir = RBNodeHere
			break
		}
	}
	return cur, dir
}

func (cur *RBCursor) add(index int) *RBCursor {
	newNode := &RBNode{index: index, isRed: true}
	cur.root()
	if cur.node != nil {
		pnode, dir := cur.find(newNode.index)
		if dir != RBNodeHere {
			newNode.parent = pnode.node
			newNode.parent.children[dir] = newNode

			wrk := cur.clone()
			wrk.node = newNode
			wrk.opt()
		} else {
			pnode.node.index = newNode.index
			newNode = pnode.node
		}
	}
	cur.node = newNode
	return cur
}

func (node *RBNode) sibiling(child *RBNode) *RBNode {
	dir := RBNodeLeft
	if node.children[RBNodeLeft] == child {
		dir = RBNodeRight
	}
	return node.children[dir]
}

func (node *RBNode) flip(dir RBNodeDir) {
	curGranPa := node.parent
	newParent := node.children[dir]
	newParent.parent, node.parent = node.parent, newParent

	node.children[dir] = newParent.children[dir^1]
	if node.children[dir] != nil {
		node.children[dir].parent = node
	}
	newParent.children[dir^1] = node

	if curGranPa != nil {
		dir = RBNodeLeft
		if curGranPa.children[dir] != node {
			dir ^= 1
		}
		curGranPa.children[dir] = newParent
	}

}
func (cur *RBCursor) opt() {
	node := cur.node
	for node != nil && node.isRed {
		parent := node.parent
		if parent == nil {
			break
		} else if !parent.isRed {
			break
		} else if parent.parent == nil {
			parent.isRed = false
		} else {
			grandparent := parent.parent
			parentsibiling := grandparent.sibiling(parent)
			if parentsibiling != nil && parentsibiling.isRed {
				grandparent.isRed = true
				parent.isRed = false
				parentsibiling.isRed = false
				node = grandparent
			} else {
				dir := RBNodeLeft
				if grandparent.children[dir] != parent {
					dir ^= 1
				}

				if parent.children[dir] != node {
					parent.flip(dir ^ 1)
				}

				grandparent.flip(dir)
				grandparent.parent.isRed = false
				grandparent.isRed = true
				break
			}
		}
	}
}

func (cur *RBCursor) end(dir RBNodeDir) *RBCursor {
	cur.root()
	if cur.node != nil {
		for ; cur.node.children[dir] != nil; cur.node = cur.node.children[dir] {

		}
	}
	return cur
}
func main() {
	x := []int{7, 5, 9, 6, 8, 9, 10, 11, 12, 13, 14}
	tree := RBCursor{less: func(i, j int) bool {
		return x[i] < x[j]
	}}

	{
		for i := 0; i < len(x); i++ {
			tree.add(i)
		}
		for cur := tree.clone().end(RBNodeLeft); cur.node != nil; cur.move(RBNodeRight) {
			fmt.Println(x[cur.node.index])
		}
		for cur := tree.clone().end(RBNodeRight); cur.node != nil; cur.move(RBNodeLeft) {
			fmt.Println(x[cur.node.index])
		}
	}
	os.Exit(0)
}
