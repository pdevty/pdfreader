package crush

// "crush" bytes into bits - variable length.

type BitT struct {
  s    []byte;
  p, b int;
}

var mask = [9]int{0, 1, 3, 7, 15, 31, 63, 127, 255}

func (x *BitT) Get(n int) (r int) {
  if x.b == 0 {
    x.b = 8;
    x.p++;
  }
  if x.b >= n {
    x.b -= n;
    r = int(x.s[x.p]>>uint8(x.b)) & mask[n];
  } else {
    n -= x.b;
    r = x.Get(x.b) << uint8(n);
    r += x.Get(n);
  }
  return;
}

func NewBits(s []byte) *BitT {
  r := new(BitT);
  r.s = s;
  r.b = 8;
  return r;
}
