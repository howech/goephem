// Packag chebyshev provides support for interpolating Chebyshev
// polynomial interpolations.
package chebyshev

//import "math"

// Chephshev interpolating polynomials are represented as a slice of
// float64s representing the coefficients of the orthogonal
// polynomials.
type Chebyshev []float64

var (
  T0 Chebyshev = Chebyshev {1}
  T1 Chebyshev = Chebyshev {0,1}
  T2 Chebyshev = Chebyshev {0,0,1}
  T3 Chebyshev = Chebyshev {0,0,0,1}
  T4 Chebyshev = Chebyshev {0,0,0,0,1}
  T5 Chebyshev = Chebyshev {0,0,0,0,0,1}
  T6 Chebyshev = Chebyshev {0,0,0,0,0,0,1}
  T7 Chebyshev = Chebyshev {0,0,0,0,0,0,0,1}
  T8 Chebyshev = Chebyshev {0,0,0,0,0,0,0,0,1}
  T9 Chebyshev = Chebyshev {0,0,0,0,0,0,0,0,0,1}

  Zero Chebyshev = Chebyshev {}
  One Chebyshev = T0
  X Chebyshev = T1
)

// integer absolute value
func abs(i int) int {
  if i >= 0 {
    return i
  } else {
    return -i
  }
}

// reslices the chebychev coeffiecients to remove high order zeros
func (a Chebyshev) normalize() Chebyshev {
  order := len(a)
  for  order > 0 && a[order-1] == 0.0 {
    order--
  }
  return a[0:order]
}

// Returns true if the polynomial has all zero coefficents
func (a Chebyshev) IsZero() bool {
  for _,x := range a {
    if x != 0.0 {
      return false
    }
  }
  return true
}

// Returns true if all of the polynomial coefficients are the same in the two arguments
func (a Chebyshev) Equals(b Chebyshev) bool {
  if(len(a) != len(b)) {
    return false
  }
  for i,x := range a {
    if x != b[i] {
      return false
    }
  }
  return true
}

// Returns the (normalized) sum of the two polynomials
func (a Chebyshev) Add(b Chebyshev) Chebyshev {
  l := len(a)
  if len(b) > l {
    l = len(b)
  }

  c :=  make(Chebyshev,l,l)
  for i := range a {
    c[i] += a[i] 
  }

  for i := range b {
    c[i] += b[i] 
  }

  return c.normalize()
}

// Returns the (normalized) difference of the two polynomials
func (a Chebyshev) Subtract(b Chebyshev) Chebyshev {
  l := len(a)
  if len(b) > l {
    l = len(b)
  }

  c :=  make(Chebyshev,l,l)
  for i := range a {
    c[i] += a[i] 
  }

  for i := range b {
    c[i] -= b[i] 
  }

  return c.normalize()
}

// Returns a polynomial scaled by the float64 argument
func (a Chebyshev) Times(b float64) Chebyshev {
  c:= make(Chebyshev, len(a), len(a))

  for i,x := range a {
    c[i] += b * x
  }

  return c.normalize()
}

// Returns the product of the two polynomials
func (a Chebyshev) Multiply(b Chebyshev) Chebyshev {
  // When you multiply two chebyshev polynomials, the order of the result
  // will be the sum of the orders of the factors 
  l := len(a) + len(b) - 1 

  c:= make(Chebyshev, l, l)
  for i,ai := range a {
    if ai == 0.0 {
      continue
    }
    for j,bj := range b {
      if bj== 0.0 {
        continue
      }
      
      // Tn * Tm = 1/2 ( T(n+m) + T(|n-m|) )
      p := ai*bj / 2
      c[i+j] += p
      c[abs(i-j)] += p
    }
  }

  return c.normalize()
}

// Uses Clenshaw's algorithm to interpolate the polynomial at a point
// x, which must be in the range -1<x<1.
func (a Chebyshev) Interpolate(x float64) float64 {
  // Range check
  if x < -1 || 1 < x {
    panic("interpolate argument out of range [-1,1]")
  }
  // Special cases for low order polynomials
  switch len(a) {
  case 0:
    return 0
  case 1:
    return a[0]
  }

  // Using clenshaw
  b0,b1 := 0.0, 0.0
  x2 := 2.0 * x
  
  for i:= len(a) - 1; i > 0; i-- {
    b1, b0 = b0, a[i] + x2 * b0 - b1
  }

  return a[0] + x * b0 - b1
}
