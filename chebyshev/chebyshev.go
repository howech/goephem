package chebyshev

type Chebyshev []float64


func (a Chebyshev) add(b Chebyshev) Chebyshev {
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

  return c
}
