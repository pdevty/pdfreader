package crush
/*
"crush" bytes into bits - variable length.

Copyright (c) 2009 Helmar Wodtke

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

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
