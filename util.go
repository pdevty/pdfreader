package util

func StrToByte(a string) []byte {
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
}
