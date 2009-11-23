package hex

// hex decoder for PDF

var deco [256]byte

func init() {
  for i := 0; i <= 255; i++ {
    deco[i] = 255
  }
  for i := '0'; i <= '9'; i++ {
    deco[i] = byte(i) - '0'
  }
  for i := 'A'; i <= 'F'; i++ {
    deco[i] = byte(i) - 'A' + 10
  }
  for i := 'a'; i <= 'f'; i++ {
    deco[i] = byte(i) - 'a' + 10
  }
}

func Decode(s string) []byte {
  r := make([]byte, (len(s)+1)/2);
  q := 0;
  for p := 0; p < len(s); p++ {
    if c := deco[s[p]]; c != 255 {
      if q%2 == 0 {
        c <<= 4
      }
      r[q/2] += c;
      q++;
    } else if s[p] > 32 {
      if s[p] == '>' {
        break
      }
      return []byte{};
    }
  }
  return r[0 : (q+1)/2];
}
