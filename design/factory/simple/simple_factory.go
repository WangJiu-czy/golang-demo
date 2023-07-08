package simple

import "fmt"

//简单工厂模式

type Girl interface {
	weight()
}

type FatGirl struct {
}

func (FatGirl) weight() {
	fmt.Println("80kg")
}

type ThinGirl struct {
}

func (ThinGirl) weight() {
	fmt.Println("50kg")
}

type GirlFactory struct {
}

func (*GirlFactory) CreateGirl(like string) Girl {

	switch like {
	case "fat":
		return &FatGirl{}
	case "thin":
		return &ThinGirl{}
	default:
		return nil
	}
}
