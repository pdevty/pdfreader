package pfb

import (
  "bytes"
)
// Decoder for pfb fonts.


const _hex = "0123456789ABCDEF"

func hexenc(i []byte) []byte {
  r := make([]byte, len(i)*2)
  for k := range i {
    r[k*2] = _hex[i[k]>>4]
    r[k*2+1] = _hex[i[k]&15]
  }
  return r
}

func Decode(b []byte) []byte {
  r := make([]byte, len(b)*2)[0:0]
  for {
    if b[0] != 128 {
      break
    }
    if b[1] == 3 {
      break
    }
    l := int(b[2]) + (int(b[3]) << 8) + (int(b[4]) << 16) + 6
    if b[1] == 1 {
      r = bytes.Add(r, b[6:l])
    } else {
      r = bytes.Add(r, hexenc(b[6:l]))
    }
    b = b[l:]
  }
  return r
}
