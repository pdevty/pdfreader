package hex
/*
hex decoder for PDF.

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

func IsHex(c byte) bool { return deco[c] != 255 }
