package chebyshev

import (
	"math"
	"testing"
)

func TestNormalize(t *testing.T) {
	a := Chebyshev{0, 0, 0, 0}
	b := Chebyshev{1, 0, 0, 0}
	c := Chebyshev{0, 1, 0, 0}
	d := Chebyshev{0, 0, 0, 1}

	if len(a.normalize()) != 0 {
		t.Error("Error normalizing null Chebyshev")
	}
	if len(b.normalize()) != 1 {
		t.Error("Error normalizing zero order Chebyshev")
	}
	if len(c.normalize()) != 2 {
		t.Error("Error normalizing first order Chebyshev")
	}
	if len(d.normalize()) != 4 {
		t.Error("Error normalizing already normalized Chebyshev")
	}

}

func TestEqual(t *testing.T) {
	a := Chebyshev{}
	b := Chebyshev{0, 0, 0, 0}.normalize()
	c := Chebyshev{0, 0, 1}

	if !a.Equals(b) || !b.Equals(a) {
		t.Error("Failed to compare null Chebyshevs")
	}

	if !c.Equals(T2) || !T2.Equals(c) {
		t.Error("Failed to compare T2")
	}

	if T0.Equals(T1) || T1.Equals(T2) || T2.Equals(T3) {
		t.Error("Equals failed to differentiate distinct polynomials")
	}

}

func TestAdd(t *testing.T) {
	a, b, e := T0, T1, Chebyshev{-1, -1}

	expected := Chebyshev{1, 1}

	c := a.Add(b)
	d := b.Add(a)

	if !expected.Equals(c) || !expected.Equals(d) {
		t.Error("Failed Adding T0 to T1")
	}

	if !T0.Equals(T0.Add(Zero)) || !T0.Equals(Zero.Add(T0)) {
		t.Error("Failed Adding Zero to T0")
	}

	g := c.Add(e)
	if !g.IsZero() {
		t.Error("Addative inverse failed")
	}

}

func TestSubtract(t *testing.T) {
	a, b := T0, T1

	expected := Chebyshev{1, -1}

	c := a.Subtract(b)
	d := b.Subtract(a)

	if !expected.Equals(c) {
		t.Error("Failed subtractinb T1 to T0")
	}

	if !c.Add(d).IsZero() {
		t.Error("Addative inverse failure")
	}

	if !T0.Equals(T0.Subtract(Zero)) {
		t.Error("Failed subtracting Zero from T0")
	}
}

func TestTimes(t *testing.T) {
	if !T0.Times(0).IsZero() {
		t.Error("Failed scalar zero multiplication test")
	}
	if !T0.Times(1).Equals(T0) {
		t.Error("Failed scalar unit multiplication test")
	}
	if !T0.Times(-1).Add(T0).IsZero() {
		t.Error("Failed -1 scalar multiply test")
	}
	if !T0.Times(7).Equals(Chebyshev{7}) {
		t.Error("Failed to multipy by scalar 7")
	}

}

func TestMultiply(t *testing.T) {
	a := Chebyshev{5, 0, 1}
	b := Chebyshev{3, 2}

	if !(a.Multiply(Zero).IsZero() && b.Multiply(Zero).IsZero() && Zero.Multiply(a).IsZero() && Zero.Multiply(b).IsZero()) {
		t.Error("Polynomial multiplication zero failure")
	}

	if !(a.Multiply(One).Equals(a) && One.Multiply(a).Equals(a) && b.Multiply(One).Equals(b)) {
		t.Error("Polynomial unit multiplication failure")
	}

	if !a.Multiply(b).Equals(b.Multiply(a)) {
		t.Error("Polynomial multiplication commutation error")
	}

	if !T1.Multiply(T2).Equals(Chebyshev{0, 0.5, 0, 0.5}) {
		t.Error("Failed Multiplying T1 by T2")
	}
}

func TestInterpolate(t *testing.T) {
	myT0 := func(x float64) float64 {
		return 1
	}
	myT1 := func(x float64) float64 {
		return x
	}
	myT2 := func(x float64) float64 {
		return 2*x*x - 1
	}
	myT5 := func(x float64) float64 {
		return 16*x*x*x*x*x - 20*x*x*x + 5*x
	}

	tests := []struct {
		c Chebyshev
		f func(float64) float64
	}{
		{T0, myT0}, {T1, myT1}, {T2, myT2}, {T5, myT5}}

	for _, v := range tests {
		c, f := v.c, v.f
		for x := -1.0; x <= 1.0; x += 0.1 {
			if math.Abs(c.Interpolate(x)-f(x)) > 0.00001 {
				t.Error("Interpolation error:", x, c, c.Interpolate(x), f(x))
			}
		}
	}
}
