package strategy

//策略模式

// 实现此接口,则为一个策略
type IStraregy interface {
	do(int, int) int
}

// 加
type add struct {
}

func (_ *add) do(a, b int) int {
	return a + b
}

// 减
type reduce struct {
}

func (_ *reduce) do(a, b int) int {
	return a - b
}

// 具体的策略执行者
type Operator struct {
	strategy IStraregy
}

func (o *Operator) setStrategy(straregy IStraregy) {
	o.strategy = straregy

}

// 调用策略方法
func (o *Operator) calculate(a, b int) int {
	return o.strategy.do(a, b)
}
