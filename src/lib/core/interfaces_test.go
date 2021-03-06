package core

import "testing"

func testEqual(ts ...*Thunk) bool {
	return bool(PApp(Equal, ts...).Eval().(BoolType))
}

func testLess(t1, t2 *Thunk) bool {
	return t1.Eval().(ordered).less(t2.Eval().(ordered))
}

func TestXFailLess(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fail()
		}
	}()

	less(NewNumber(42).Eval(), NewError("you", "failed.").Eval())
}
