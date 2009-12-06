package unquot

import (
  "hex";
)

func String(s []byte) []byte {
  if s[0] == '<' {
    return hex.Decode(string(s[1 : len(s)-1]))
  }
  if s[0] != '(' {
    return []byte{}
  }
  r := make([]byte, len(s));
  q := 0;
  for p := 1; p < len(s)-1; p++ {
    if s[p] == '\\' {
      p++;
      switch s[p] {
      case 13:
        if s[p+1] == 10 {
          p++
        }
        q--;
      case 10:
        q--
      case 'n':
        r[q] = 10
      case 'r':
        r[q] = 13
      case 't':
        r[q] = 9
      case 'b':
        r[q] = 8
      case 'f':
        r[q] = 12
      case '0', '1', '2', '3', '4', '5', '6', '7':
        a := s[p] - '0';
        if s[p+1] >= '0' && s[p+1] <= '7' {
          p++;
          a = (a << 3) + (s[p] - '0');
          if s[p+1] >= '0' && s[p+1] <= '7' {
            p++;
            a = (a << 3) + (s[p] - '0');
          }
        }
        r[q] = a;
      default:
        r[q] = s[p]
      }
    } else {
      r[q] = s[p]
    }
    q++;
  }
  return r[0:q];
}
