package rbtree

//RBNodeDir is
type RBNodeDir int

const (

	//RBNodeLeft is
	RBNodeLeft RBNodeDir = iota

	//RBNodeRight is
	RBNodeRight

	//RBNodeHere is
	RBNodeHere
)

//RBNode is
type RBNode struct {
	isRed    bool
	Index    int
	parent   *RBNode
	children [2]*RBNode
}

//RBTree is
type RBTree struct {
	Node *RBNode
	less func(i int, j int) bool
}

//RBCursor is
type RBCursor RBTree

//NewTree is
func NewTree(less func(i int, j int) bool) *RBTree {
	tree := &RBTree{less: less}
	return tree
}

//Clone is
func (cur *RBTree) Clone() *RBCursor {
	wrk := RBCursor(*cur)
	return &wrk
}

//Root is
func (cur *RBTree) Root() *RBTree {

	if cur.Node != nil {
		for ; cur.Node.parent != nil; cur.Node = cur.Node.parent {

		}
	}
	return cur
}

//Move is
func (cur *RBCursor) Move(dir RBNodeDir) *RBCursor {
	rev := dir ^ 1
	if next := cur.Node.children[dir]; next != nil {
		for next.children[rev] != nil {
			next = next.children[rev]
		}
		cur.Node = next
	} else {
		find := false
		for cur.Node.parent != nil {
			now := cur.Node
			cur.Node = cur.Node.parent
			if cur.Node.children[rev] == now {
				find = true
				break
			}
		}
		if !find {
			cur.Node = nil
		}
	}
	return cur
}

//Find is
func (cur *RBTree) Find(Index int) (*RBTree, RBNodeDir) {
	dir := RBNodeHere
	if cur.Node == nil {
		return cur, dir
	}
	cur.Root()
	for {
		if cur.less(Index, cur.Node.Index) {
			dir = RBNodeLeft
			next := cur.Node.children[dir]
			if next == nil {
				break
			}
			cur.Node = next

		} else if cur.less(cur.Node.Index, Index) {
			dir = RBNodeRight
			next := cur.Node.children[dir]
			if next == nil {
				break
			}
			cur.Node = next
		} else {
			dir = RBNodeHere
			break
		}
	}
	return cur, dir
}

//Add is
func (cur *RBTree) Add(Index int) *RBTree {
	newNode := &RBNode{Index: Index, isRed: true}
	cur.Root()
	if cur.Node != nil {
		pnode, dir := cur.Find(newNode.Index)
		if dir != RBNodeHere {
			newNode.parent = pnode.Node
			newNode.parent.children[dir] = newNode

			wrk := cur.Clone()
			wrk.Node = newNode
			wrk.opt()
		} else {
			pnode.Node.Index = newNode.Index
			newNode = pnode.Node
		}
	}
	cur.Node = newNode
	return cur
}

func (Node *RBNode) flip(dir RBNodeDir) {
	curGranPa := Node.parent
	newParent := Node.children[dir]
	if curGranPa != nil {
		curGranPa.children[Node.dir()] = newParent
	}
	newParent.parent, Node.parent = Node.parent, newParent

	Node.children[dir] = newParent.children[dir^1]
	if Node.children[dir] != nil {
		Node.children[dir].parent = Node
	}
	newParent.children[dir^1] = Node

}
func (Node *RBNode) dir() RBNodeDir {
	dir := RBNodeLeft
	if Node.parent.children[dir] != Node {
		dir = RBNodeRight
	}
	return dir
}
func (cur *RBCursor) opt() {
	Node := cur.Node
	for Node != nil && Node.isRed {
		parent := Node.parent
		if parent == nil {
			break
		} else if !parent.isRed {
			break
		} else if parent.parent == nil {
			parent.isRed = false
		} else {
			grandparent := parent.parent
			parentsibling := grandparent.children[parent.dir()^1]
			if parentsibling != nil && parentsibling.isRed {
				grandparent.isRed = true
				parent.isRed = false
				parentsibling.isRed = false
				Node = grandparent
			} else {
				dir := parent.dir()

				if parent.children[dir] != Node {
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

//End is
func (cur *RBTree) End(dir RBNodeDir) *RBCursor {
	cur.Root()
	if cur.Node != nil {
		for ; cur.Node.children[dir] != nil; cur.Node = cur.Node.children[dir] {

		}
	}
	return cur.Clone()
}

func (Node *RBNode) cut() {
	if Node.parent != nil {
		Node.parent.children[Node.dir()] = nil
	}
}

//Delete is
func (cur *RBTree) Delete(Index int) (ret bool) {
	//wcur := cur.Root().Clone()
	wcur, dir := cur.Find(Index)
	delNode := wcur.Node
	if delNode == nil || dir != RBNodeHere {
		return
	}
	ret = true

	dir = RBNodeLeft
	if delNode.children[dir] == nil {
		dir = RBNodeRight
		if delNode.children[dir] == nil {
			dir = RBNodeHere
		}
	}

	if delNode.children[RBNodeLeft] != nil &&
		delNode.children[RBNodeRight] != nil {
		var wrk *RBCursor
		if dir != RBNodeHere {
			wrk = wcur.Clone().Move(dir)
			delNode.Index = wrk.Node.Index
			delNode = wrk.Node
			dir = RBNodeLeft
			if delNode.children[dir] == nil {
				dir = RBNodeRight
				if delNode.children[dir] == nil {
					dir = RBNodeHere
				}
			}
		}
	}

	if dir == RBNodeHere {
		if delNode.isRed {
			delNode.cut()
			return
		}
	} else {
		wrk := delNode.children[dir]
		delNode.Index = wrk.Index
		wrk.cut()
		return
	}

	Node := delNode
	for {
		parent := Node.parent
		if parent == nil {
			break
		}

		dir := Node.dir()
		dirOther := dir ^ 1
		sibling := parent.children[dirOther]

		if sibling.isRed {
			//sibling is Red
			parent.flip(dirOther)
			sibling.isRed = false
			parent.isRed = true
			sibling = parent.children[dirOther]
		}
		//sibling is Black

		nephew := sibling.children[dirOther]
		if nephew == nil || !nephew.isRed {
			//far nephew is Black
			nephew = sibling.children[dir]
			if nephew == nil || !nephew.isRed {
				//near nephew is Black
				sibling.isRed = true
				if parent.isRed {
					parent.isRed = false
					break
				} else {
					Node = parent
					continue
				}
			}
			//near nephew is Red and far nephew is Black
			sibling.flip(dir)
			sibling, nephew = nephew, sibling
			sibling.isRed = false
			nephew.isRed = true
		}
		//sibling is Black && far nephew is Red

		saveColor := parent.isRed
		parent.flip(dirOther)
		sibling.isRed = saveColor
		parent.isRed = false
		nephew.isRed = false
		break

	}
	delNode.cut()
	if delNode.parent == nil {
		cur.Node = nil
	}
	return
}
