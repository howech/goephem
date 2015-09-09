package chebyshev

import (
  "testing"
  "fmt"
)

func TestAdd(t *testing.T) {
  a,b:=Chebyshev {1}, Chebyshev {0,1}

  c:= a.add(b)
  fmt.Println(c)
  switch {
  case len(c) != 2:
    t.Error("result not the right length")
  case c[0] != 1.0:
    t.Error("zero order coef wrong")
  case c[1] != 1.0:
    t.Error("first order coef wrong")
  }
}
