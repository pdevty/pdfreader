package xchar

// encode Utf8

var utconv = []int{
  0x7F, 0x00,
  0x7FF, 0xC0,
  0xFFFF, 0xE0,
  0x1FFFFF, 0xF0,
}

func Utf8(rune int) []byte {
  up := 3;
  out := make([]byte, up+1);
  if rune < 0 || rune > 0x10FFFF {
    rune = 0xFFFD;
  }
  uc := 0;
  r := rune;
  for ; utconv[uc] < rune; uc += 2 {
    out[up] = byte((r & 0x3F) | 0x80);
    up--;
    r >>= 6;
  }
  out[up] = byte(r | utconv[uc+1]);
  return out[up:];
}

func EncodeRune(rune int, out []byte) int {
  r := Utf8(rune);
  copy(out, r);
  return len(r);
}
