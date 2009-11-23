package lzw

const (
  _LZW_EOD       = 257;
  _LZW_RESET     = 256;
  _LZW_DICSIZE   = 4096;
  _LZW_STARTBITS = 9;
  _LZW_STARTUTOK = 258;
)

type BitT struct {
  s    []byte;
  p, b int;
}

type lzwDecoder struct {
  bits   *BitT;
  bc, cp int;
  early  bool;
}

var mask = [9]int{0, 1, 3, 7, 15, 31, 63, 127, 255}

func (x *BitT) get(n int) (r int) {
  if x.b == 0 {
    x.b = 8;
    x.p++;
  }
  if x.b >= n {
    x.b -= n;
    r = int(x.s[x.p]>>uint8(x.b)) & mask[n];
  } else {
    n -= x.b;
    r = x.get(x.b) << uint8(n);
    r += x.get(n);
  }
  return;
}

func NewBits(s []byte) *BitT {
  r := new(BitT);
  r.s = s;
  r.b = 8;
  return r;
}

func (lzw *lzwDecoder) reset() {
  lzw.bc = _LZW_STARTBITS;
  lzw.cp = _LZW_STARTUTOK - 1;
}

func newLzwDecoder(s []byte, early bool) (lzw *lzwDecoder) {
  lzw = new(lzwDecoder);
  lzw.bits = NewBits(s);
  lzw.early = early;
  lzw.reset();
  return;
}

func (lzw *lzwDecoder) update() bool {
  if lzw.cp < _LZW_DICSIZE-1 {
    if lzw.early {
      lzw.cp++
    }
    switch lzw.cp {
    case 511:
      lzw.bc = 10
    case 1023:
      lzw.bc = 11
    case 2047:
      lzw.bc = 12
    }
    if !lzw.early {
      lzw.cp++
    }
    return true;
  }
  return false;
}

func (lzw *lzwDecoder) token() (r int) {
  for {
    r = lzw.bits.get(lzw.bc);
    if r != _LZW_RESET {
      break
    }
    lzw.reset();
  }
  return r;
}

func DecodeToSlice(s []byte, out []byte, early bool) (r int) {
  lzw := newLzwDecoder(s, early);
  dict := make([][]byte, _LZW_DICSIZE);
  for i := 0; i <= 255; i++ {
    dict[i] = []byte{byte(i)}
  }
  for c := lzw.token(); c != _LZW_EOD; c = lzw.token() {
    k := r;
    for i := 0; i < len(dict[c]); i++ {
      out[r] = dict[c][i];
      r++;
    }
    if lzw.update() {
      dict[lzw.cp] = out[k : r+1]
    }
  }
  return;
}

func CalculateLength(s []byte, early bool) (r int) {
  lzw := newLzwDecoder(s, early);
  dict := make([]int, _LZW_DICSIZE);
  for i := 0; i <= 255; i++ {
    dict[i] = 1
  }
  for c := lzw.token(); c != _LZW_EOD; c = lzw.token() {
    r += dict[c];
    if lzw.update() {
      dict[lzw.cp] = dict[c] + 1
    }
  }
  return;
}

func Decode(s []byte, early bool) []byte {
  r := make([]byte, CalculateLength(s, early)+1);
  return r[0:DecodeToSlice(s, r, early)];
}
