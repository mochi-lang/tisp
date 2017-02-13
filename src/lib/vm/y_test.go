package vm

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestY(t *testing.T) {
	for _, n := range []float64{0, 1, 2, 3, 4, 5, 6, 42, 100} {
		n1 := lazyFactorial(NewNumber(n))
		n2 := strictFactorial(n)

		t.Logf("%d: %f == %f?\n", int(n), n1, n2)

		assert.Equal(t, n1, n2)
	}

	for _, ts := range [][]*Thunk{
		{NewNumber(7)},
		{NewNumber(13), NewString("foobarbaz")},
		{NewNumber(42), Nil, Nil},
	} {
		t.Log(lazyFactorial(ts...))
	}
}

func strictFactorial(n float64) float64 {
	if n == 0 {
		return 1
	}

	return n * strictFactorial(n-1)
}

func lazyFactorial(ts ...*Thunk) float64 {
	return float64(App(App(Y, lazyFactorialImpl), ts...).Eval().(NumberType))
}

var lazyFactorialImpl = NewLazyFunction(func(ts ...*Thunk) Object {
	// fmt.Println(len(ts))

	return App(If,
		App(Equal, ts[1], NewNumber(0)),
		NewNumber(1),
		App(Mul, NewList(
			ts[1],
			App(ts[0], append([]*Thunk{App(Sub, ts[1], NewList(NewNumber(1)))}, ts[2:]...)...))))
})

func TestYsSingleF(t *testing.T) {
	fs := App(Ys, lazyFactorialImpl)

	for _, n := range []float64{0, 1, 2, 3, 4, 5, 6, 42, 100} {
		n1 := float64(App(App(First, fs), NewNumber(n)).Eval().(NumberType))
		n2 := strictFactorial(n)

		t.Logf("%d: %f == %f?\n", int(n), n1, n2)

		assert.Equal(t, n1, n2)
	}
}

func TestYsMultipleFs(t *testing.T) {
	evenWithExtraArg := NewLazyFunction(func(ts ...*Thunk) Object {
		n := ts[3]

		return App(If,
			App(Equal, n, NewNumber(0)),
			True,
			App(ts[1], App(Sub, n, NewList(NewNumber(1)))))
	})

	odd := NewLazyFunction(func(ts ...*Thunk) Object {
		n := ts[2]

		return App(If,
			App(Equal, n, NewNumber(0)),
			False,
			App(ts[0], Nil, App(Sub, n, NewList(NewNumber(1)))))
	})

	fs := App(Ys, evenWithExtraArg, odd)

	for _, n := range []float64{0, 1, 2, 3, 4, 5, 6, 42, 100, 121, 256, 1023} {
		b1 := bool(App(App(First, fs), NewString("unused"), NewNumber(n)).Eval().(BoolType))
		b2 := bool(App(App(First, App(Rest, fs)), NewNumber(n)).Eval().(BoolType))

		t.Logf("n = %v, even? %v, odd? %v\n", n, b1, b2)

		rem := int(n) % 2
		assert.Equal(t, b1, rem == 0)
		assert.Equal(t, b2, rem != 0)
	}
}
