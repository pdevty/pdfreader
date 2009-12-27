package pfb

import (
  "bytes"
  "hex"
)
// Decoder for pfb fonts.

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
      r = bytes.Add(r, hex.Encode(b[6:l]))
    }
    b = b[l:]
  }
  return r
}
