package lru_cache

import "fmt"

const SIZE = 5

type Node struct {
	Val   string
	Left  *Node
	Right *Node
}

type Queue struct {
	Head   *Node
	Tail   *Node
	Length int
}

type Cache struct {
	//来确认谁最近最近没使用
	Queue Queue
	//来确认有哪些节点
	storage map[string]*Node
}

func NewCache() Cache {
	return Cache{Queue: NewQueue(), storage: map[string]*Node{}}
}
func NewQueue() Queue {
	head := &Node{}
	tail := &Node{}
	head.Right = tail
	tail.Left = head
	return Queue{Head: head, Tail: tail}
}

//删除缓存和存储库中的节点
func (c *Cache) Remove(n *Node) *Node {
	left := n.Left
	right := n.Right
	left.Right = right
	right.Left = left
	c.Queue.Length--
	delete(c.storage, n.Val)
	return n
}

func (c *Cache) Add(n *Node) {
	tmp := c.Queue.Head.Right
	c.Queue.Head.Right = n
	n.Left = c.Queue.Head
	n.Right = tmp
	tmp.Left = n
	c.Queue.Length++
	if c.Queue.Length > SIZE {
		c.Remove(c.Queue.Tail.Left)
	}
}
func (c *Cache) Check(str string) {
	node := &Node{}
	if val, ok := c.storage[str]; ok {
		node = c.Remove(val)
	} else {
		node = &Node{Val: str}

	}
	c.Add(node)
	c.storage[str] = node

}
func (q *Queue) Display() {
	node := q.Head.Right
	for i := 0; i < q.Length; i++ {
		fmt.Printf("{%s}", node.Val)
		if i < q.Length-1 {
			fmt.Printf("<-->")
		}
		node = node.Right
	}
	fmt.Println("]")
}
func (c *Cache) Display() {
	c.Queue.Display()
}
