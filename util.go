package util

/* Some utilities.

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

// util.Bytes() is a dup of string.Bytes()
func Bytes(a string) []byte {
  r := make([]byte, len(a));
  for k := range a {
    r[k] = byte(a[k])
  }
  return r;
}

func JoinStrings(a []string, c byte) []byte {
  if a == nil {
    return []byte{}
  }
  l := 0;
  for k := range a {
    l += len(a[k]) + 1
  }
  r := make([]byte, l);
  q := 0;
  for k := range a {
    for i := 0; i < len(a[k]); i++ {
      r[q] = a[k][i];
      q++;
    }
    r[q] = c;
    q++;
  }
  return r[0 : l-1];
}

func StringArray(i [][]byte) []string {
  r := make([]string, len(i));
  for k := range i {
    r[k] = string(i[k])
  }
  return r;
}

// Stacks

type StackT struct {
  st [][]byte;
  sp int;
}

func (st *StackT) Push(s []byte) {
  st.st[st.sp] = s;
  st.sp++;
}

func (st *StackT) Drop(n int) [][]byte {
  st.sp -= n;
  return st.st[st.sp : st.sp+n];
}

func (st *StackT) Pop() []byte {
  st.sp--;
  return st.st[st.sp];
}

func (st *StackT) Dump() [][]byte { return st.st[0:st.sp] }

func (st *StackT) Depth() int { return st.sp }

func NewStack(n int) *StackT {
  r := new(StackT);
  r.st = make([][]byte, n);
  return r;
}

type Stack interface {
  Push([]byte);
  Pop() []byte;
  Drop(int) (st [][]byte);
  Dump() [][]byte;
  Depth() int;
}


func set(o []byte, q string) {
  for k := range q {
    o[k] = q[k]
  }
}

func ToXML(s []byte) []byte {
  l := len(s);
  for k := range s {
    switch s[k] {
    case '<', '>':
      l += 3
    case '&':
      l += 4
    }
  }
  r := make([]byte, l);
  p := 0;
  for k := range s {
    switch s[k] {
    case '<':
      set(r[p:p+3], "&lt;");
      p += 4;
    case '>':
      set(r[p:p+3], "&gt;");
      p += 4;
    case '&':
      set(r[p:p+4], "&amp;");
      p += 5;
    default:
      r[p] = s[k];
      p++;
    }
  }
  return r;
}
