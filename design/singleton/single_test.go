package singleton

import "testing"

func TestGoSingleton(t *testing.T) {
	instance1 := GetIns()
	t.Logf("%p", instance1)
	instance2 := GetIns()
	t.Logf("%p", instance2)

}
