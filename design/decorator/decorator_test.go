package decorator

import (
	"fmt"
	"testing"
)

func TestDecorator(t *testing.T) {
	lw := &laowang{}
	jk := &Jacket{}
	jk.p = lw
	jk.show()
	hat := &Hat{}
	hat.p = jk
	hat.show()

	fmt.Println("cost:", hat.cost())

}
