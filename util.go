package util

import "fmt"

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

import "xchar"

var wrongUniCode = xchar.Utf8(-1)

// util.Bytes() is a dup of string.Bytes()
func Bytes(a string) []byte {
  r := make([]byte, len(a));
  for k := 0; k < len(a); k++ {
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

func (st *StackT) Index(p int) []byte { return st.st[st.sp-p] }

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
  Index(p int) []byte;
}


func set(o []byte, q string) int {
  for k := 0; k < len(q); k++ {
    o[k] = q[k]
  }
  return len(q);
}

func ToXML(s []byte) []byte {
  l := len(s);
  for k := range s {
    switch s[k] {
    case '<', '>':
      l += 3
    case '&':
      l += 4
    case 0, 1, 2, 3, 4, 5, 6, 7, 8,
      11, 12, 14, 15, 16, 17, 18, 19, 20,
      21, 22, 23, 24, 25, 26, 27, 28, 29, 30,
      31:
      l += len(wrongUniCode) - 1
    }
  }
  r := make([]byte, l);
  p := 0;
  for k := range s {
    switch s[k] {
    case '<':
      p += set(r[p:p+4], "&lt;")
    case '>':
      p += set(r[p:p+4], "&gt;")
    case '&':
      p += set(r[p:p+5], "&amp;")
    case 10, 9, 13:
      r[p] = s[k];
      p++;
    default:
      if s[k] < 32 {
        p += copy(r[p:], wrongUniCode)
      } else {
        r[p] = s[k];
        p++;
      }
    }
  }
  return r;
}

type OutT struct {
  Content []byte;
}

func (t *OutT) Out(f string, args ...) {
  p := fmt.Sprintf(f, args);
  q := len(t.Content);
  if cap(t.Content)-q < len(p) {
    n := make([]byte, cap(t.Content)+(len(p)/512+2)*512);
    copy(n, t.Content);
    t.Content = n[0:q];
  }
  t.Content = t.Content[0 : q+len(p)];
  for k := 0; k < len(p); k++ {
    t.Content[q+k] = p[k]
  }
}
