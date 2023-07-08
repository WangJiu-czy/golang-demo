package observer

import "testing"

func TestObserver(t *testing.T) {
	cA := &CustomerA{}
	cB := &CustomerB{}

	offic := &NewsOffice{}

	//模拟订阅
	offic.addCustomer(cA)
	offic.addCustomer(cB)

	//更新内容
	offic.newspaperCome()

}
