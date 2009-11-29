package strm

// string math - PROTOTYPE! This is not lossless and produces dust...

import (
  "fmt";
  "strconv";
)

func Neg(s string) string {
  if s[0] == '-' {
    return s[1:len(s)]
  }
  return fmt.Sprintf("-%s", s);
}

func Add(a, b string) string {
  av, _ := strconv.Atof(a);
  bv, _ := strconv.Atof(b);
  return fmt.Sprintf("%f", av+bv);
}

func Sub(a, b string) string {
  av, _ := strconv.Atof(a);
  bv, _ := strconv.Atof(b);
  return fmt.Sprintf("%f", av-bv);
}
