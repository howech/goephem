package chebyshev

type Chebyshev struct {
  coefs []float64
}

func (a Chebyshev) add(b Chebyshev) Chebyshev {
  l := len(a.coefs)
  if len(b.coefs) > l {
    l = len(b.coefs)
  }

  c := Chebyshev{ coefs: make([]float64,l,l) }

  for i := range a.coefs {
    c.coefs[i] += a.coefs[i] 
  }

  for i := range b.coefs {
    c.coefs[i] += b.coefs[i] 
  }

  return c
}
