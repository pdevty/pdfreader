package util

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
