package jplde

import (
	"testing"
)

func TestReadHeaderFromFile(t *testing.T) {
  header, err := ReadHeaderFromFile("./header.430t")
  if err != nil {
    t.Error(err)
  }

  t.Log(*header)
  
  if header.Constants["CLIGHT"] != 2.99792458E+05 {
    t.Error("Unexpected value for CLIGHT", header.Constants["CLIGHT"],2.99792458E+05 )
  }
}
  
