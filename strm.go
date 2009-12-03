package strm

// string math

import (
  "bignum";
)

func operand(s string) (r int64, f int) {
  if len(s) < 1 { return 0, 1 }
  sig := s[0] == '-';
  p := 0;
  if sig {
    p++
  }
  for p < len(s) {
    if s[p] == '.' {
      f = 1
    } else {
      f *= 10;
      r *= 10;
      r += int64(s[p] - '0');
    }
    p++;
  }
  if sig {
    r = -r
  }
  if f == 0 {
    f = 1
  }
  return;
}

func Int64(s string, f int) int64 {
  ra, fa := operand(s);
  for fa < f {
    fa *= 10;
    ra *= 10;
  }
  return ra / int64(fa/f);
}

func Int(s string, f int) int { return int(Int64(s, f)) }

func twop(a, b string) (ra, rb int64, f int) {
  ra, f = operand(a);
  rb, fb := operand(b);
  for fb < f {
    fb *= 10;
    rb *= 10;
  }
  for f < fb {
    f *= 10;
    ra *= 10;
  }
  return;
}

func String(a int64, f int) string {
  buf := make([]byte, 128);
  p := 0;
  if a < 0 {
    buf[p] = '-';
    p++;
    a = -a;
  }
  var fu func(c int64);
  step := 1;
  fu = func(c int64) {
    s := step;
    step *= 10;
    if c > 9 || step <= f {
      fu(c / 10)
    }
    buf[p] = '0' + byte(c%10);
    p++;
    if f == s && f != 1 {
      buf[p] = '.';
      p++;
    }
  };
  fu(a);
  return string(buf[0:p]);
}

func Mul(a, b string) string {
  ra, rb, f := twop(a, b);
  ar := bignum.Rat(ra, int64(f));
  br := bignum.Rat(rb, int64(f));
  i, n := ar.Mul(br).Value();
  nv := n.Value();
  d := uint64(1);
  for d%nv != 0 {
    d *= 10
  }
  i = i.Mul1(int64(d / nv));
  if uint64(f) < d {
    i = i.Div(bignum.Int(int64(d / uint64(f))));
    d = uint64(f);
  }
  return String(i.Value(), int(d));
}

func Add(a, b string) string {
  ra, rb, f := twop(a, b);
  return String(ra+rb, f);
}

func Sub(a, b string) string {
  ra, rb, f := twop(a, b);
  return String(ra-rb, f);
}

func Neg(a string) string {
  if a[0] == '-' {
    return a[1:len(a)]
  }
  ra, f := operand(a);
  return String(-ra, f);
}
