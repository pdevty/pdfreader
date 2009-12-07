package cmap

import (
  "fancy";
  "util";
  "strm";
  "ps";
  "strconv";
  "utf8";
)

var (
  Others []string;
  oc     int;
)

func init() {
  Others = make([]string, 1024);
  Others[1] = "&lt;";
  Others[2] = "&gt;";
  Others[3] = "&amp;";
  oc = 4;
}

// FIXME/STUB: Everything is "too simple" here...
func Read(f fancy.Reader) []int {
  r := make([]int, 256);
  for k := range r {
    r[k] = k
  }
  r['<'] = -1;
  r['>'] = -2;
  r['&'] = -3;
  if f == nil {
    return r
  }
  c := util.Bytes("0");
  for {
    t, _ := ps.Token(f);
    if string(t) == "" {
      break
    }
    if string(t) == "beginbfchar" {
      cc := strm.Int(string(c), 1);
      for k := 0; k < cc; k++ {
        ch, _ := ps.Token(f);
        tr, _ := ps.Token(f);
        ci, _ := strconv.Btoi64(string(ch[1:len(ch)-1]), 16);
        ti, _ := strconv.Btoi64(string(tr[1:len(tr)-1]), 16);
        switch ti {
        case '<':
          ti = -1
        case '>':
          ti = -2
        case '&':
          ti = -3
        }
        r[int(ci)] = int(ti);
      }
      break;
    }
    c = t;
  }
  return r;
}

func Decode(s []byte, m []int) []byte {
  r := make([]byte, len(s)*6);
  p := 0;
  for k := range s {
    if m[s[k]] < 0 {
      q := Others[-m[s[k]]];
      for l := range q {
        r[p] = q[l];
        p++;
      }
    } else {
      if m[s[k]] != 0 { // FIXME, WRONG ASSUMPTION, for now this fixes some CID-Fonts.
        p += utf8.EncodeRune(m[s[k]], r[p:len(r)])
      }
    }
  }
  return r[0:p];
}
