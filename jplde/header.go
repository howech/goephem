// This package lets you load and use jpl de-series ephemeridies
//
// It was first implemented for de-430t, ascii files, but it may
// eventually include support for the others as well.
package jplde

import(
  "io"
  "bufio"
  "os"
  "strconv"
  "strings"
  "errors"
 )

/*
  According to the jpl documentation, the bodies refered to by the various
  ephemerides' columns are:
       Mercury
       Venus
       Earth-Moon barycenter
       Mars 
       Jupiter 
       Saturn
       Uranus
       Neptune
       Pluto
       Moon (geocentric)
       Sun
       Earth Nutations in longitude and obliquity (IAU 1980 model)
       Lunar mantle libration
       Lunar mantle angular velocity
       TT-TDB (at geocenter)

  Note that some of the ephemerides do not include data for all of the 
  bodies. Older versions are missing some of the last bodies entirely.

*/


// The BodyHeader type reords information about the details of the
// record data for the bodies 
type BodyHeader struct {
  IndexStart int
  Coeffs int
  Coords int
  Blocks int
}

// The header type represents (most) of the information contained
// in the ephemeris header file. The constandte get rolled into a
// map, and the Group1050 data gets turned into an array of BodyHeaders
type Header struct {
  KSize int
  NCoeff int

  // Group 1010 
  Version string
  StartDate string
  EndDate string

  // Group 1030
  Start float64
  End float64
  DaysPerRecord int

  // Group 1040/1041
  ConstantCount int
  Constants map[string]float64

  // Group 1050
  BodyHeaders []BodyHeader
}

// The floating point values that show up in ASCII versions of the JPL
// ephemerides use a D character instead of an E character to denote the
// exponent. 
func ParseDouble(s string) (f float64,err error) {
  f, err = strconv.ParseFloat( strings.Replace(s, "D", "E", 1), 64)
  return f,err
}

func ReadHeaderFromReader(headerReader io.Reader) (*Header,error) {
  scanner := bufio.NewScanner(headerReader)
  scanner.Split(bufio.ScanWords)

  var h Header
  var names []string

  for scanner.Scan() {
    switch scanner.Text() {
    case "KSIZE=":
      scanner.Scan()
      h.KSize, _ = strconv.Atoi(scanner.Text())
    case "NCOEFF=":
      scanner.Scan()
      h.NCoeff, _ = strconv.Atoi(scanner.Text())
    case "GROUP":
      scanner.Scan();
      switch scanner.Text() {

      case "1010":
        scanner.Split(bufio.ScanLines)
        scanner.Scan()
        h.Version = scanner.Text()
        scanner.Scan()
        h.StartDate = scanner.Text()
        scanner.Scan()
        h.EndDate =scanner.Text()
        scanner.Split(bufio.ScanWords)

      case "1030":
        var err error
        scanner.Scan()
        h.Start, err = ParseDouble(scanner.Text())
        if err != nil {
          return nil,err
        }
        scanner.Scan()
        h.End, err = ParseDouble(scanner.Text())
        if err != nil {
          return nil,err
        }

        scanner.Scan()
        dpr,err := ParseDouble(scanner.Text())
        if err != nil {
          return nil,err
        }
        h.DaysPerRecord = int(dpr)
        
      case "1040":
        var err error
        scanner.Scan()
        constants, err := strconv.Atoi(scanner.Text())
        if err != nil {
          return nil, err
        }
        names = make( []string, constants, constants)
        for i:=0;i<constants;i++ {
          scanner.Scan()
          names[i] = scanner.Text()
        }
      case "1041":
        var err error
        h.Constants = make(map[string]float64)
        scanner.Scan()
        constants, err := strconv.Atoi(scanner.Text())
        if err != nil {
          return nil,err
        }
        for i:=0;i<constants;i++ {
          scanner.Scan()
           h.Constants[names[i]],err = ParseDouble(scanner.Text())
          if err != nil {
            return nil, err
          }
        }
      case "1050":
        // de430t has 15 'bodies' that it tracks. Other ephemera from jpl
        // have fewer, but there is nothing in the text of the header that explicitly
        // tells you how many bodies it represents. We assume that there are less
        // than 32, and se determine where the end of the first row is by
        // watching for the values to decrease. This will of course fail for an
        // ephemeris for only one body...
        var err error
        offsets := make([]int, 0,15)

        thisOffset := 0
        lastOffset := 0
        columnCount := 0

        for {
          scanner.Scan()
          thisOffset, err = strconv.Atoi(scanner.Text())
          if err != nil {
            return nil,err
          }
          if thisOffset < lastOffset {
            break
          }
          lastOffset = thisOffset
          
          // Grow the offsets array to fit the new element
          offsets = append(offsets, thisOffset)
          columnCount++
        }
        
        coefs := make([]int, columnCount, columnCount)
        sets := make([]int, columnCount, columnCount)
        coords := make([]int, columnCount, columnCount)


        //We already have the first of the coefficient counts
        coefs[0] = thisOffset        
        for i:=1; i< columnCount; i++ {
          scanner.Scan()
          coefs[i], err = strconv.Atoi(scanner.Text())
          if err != nil {
            return nil,err
          }
        }

        for i:=0; i< columnCount; i++ {
          scanner.Scan()
          sets[i], err = strconv.Atoi(scanner.Text())
          if err != nil {
            return nil,err
          }

        }

        //We have to infer the number of coordinates for each set from
        //the offsets
        for i:=0; i<columnCount-1; i++ {
          if sets[i] == 0 || coefs[i] == 0 {
            coords[i] = 0
          } else {
            coords[i] = (offsets[i+1]-offsets[i])/(sets[i]*coefs[i])
          }
        }

        // For the last body, we have to use the total number of coeffs in the
        // Ephemeris
        if h.NCoeff == 0 {
          return nil,errors.New("Header file is missing an NCOEFF= declaration.")
        }
        
        i := columnCount-1
        if sets[i] == 0 || coefs[i] == 0 {
          coords[i] = 0
        } else {
          coords[i] = (h.NCoeff+1-offsets[i])/(sets[i]*coefs[i])
        }
        
        h.BodyHeaders = make([]BodyHeader,columnCount,columnCount)
        for i:=0; i< columnCount; i++ {
          h.BodyHeaders[i].IndexStart = offsets[i]
          h.BodyHeaders[i].Coeffs = coefs[i]
          h.BodyHeaders[i].Coords = coords[i]
          h.BodyHeaders[i].Blocks = sets[i]
        }
        
      }
    }
  }
  return &h, nil
}

func ReadHeaderFromFile(filename string) (*Header, error) {
  f,err := os.Open(filename)
  if err !=nil {
    return nil,err
  }
  defer f.Close()
  
  return ReadHeaderFromReader(f)
}


