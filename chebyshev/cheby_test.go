package chebyshev

import (
  "testing"
  "fmt"
)

func TestAdd(t *testing.T) {
  a,b:=Chebyshev {coefs: []float64{1} }, Chebyshev {coefs: []float64{0,1}}

  c:= a.add(b)
  fmt.Println(c)
  switch {
  case len(c.coefs) != 2:
    t.Error("result not the right length")
  case c.coefs[0] != 1.0:
    t.Error("zero order coef wrong")
  case c.coefs[1] != 1.0:
    t.Error("first order coef wrong")
  }
}
