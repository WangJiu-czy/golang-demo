package decorator

import "fmt"

type Person interface {
	cost() int
	show()
}

// 被装饰的对象
type laowang struct {
}

func (*laowang) show() {

	fmt.Println("赤裸裸的老王...")
}
func (*laowang) cost() int {
	return 0
}

// 衣服装饰器
type clothesDecorator struct {
	p Person
}

func (*clothesDecorator) show() {

}
func (*clothesDecorator) cost() int {
	return 0
}

// 夹克
type Jacket struct {
	clothesDecorator
}

func (j *Jacket) show() {
	j.p.show()
	fmt.Println("穿上夹克的老王...")
}
func (j *Jacket) cost() int {
	return j.p.cost() + 10
}

// 帽子
type Hat struct {
	clothesDecorator
}

func (h *Hat) cost() int {
	return h.p.cost() + 5
}
func (h *Hat) show() {
	h.p.show()
	fmt.Println("戴上帽子的老王。。。")
}
