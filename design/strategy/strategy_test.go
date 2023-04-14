package strategy

import (
	"fmt"
	"testing"
)

func TestStrategy(t *testing.T) {
	op := Operator{}
	op.setStrategy(&add{})
	res := op.calculate(1, 2)
	fmt.Println("add: ", res)

	op.setStrategy(&reduce{})
	res = op.calculate(34, 4)
	fmt.Println("reduce: ", res)

}
