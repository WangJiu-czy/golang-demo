package simple

import "testing"

func TestSimple(t *testing.T) {
	girlFactory := new(GirlFactory)
	girl := girlFactory.CreateGirl("fat")
	girl.weight()
	girl = girlFactory.CreateGirl("thin")
	girl.weight()
}
