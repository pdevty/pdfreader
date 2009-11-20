package fancy

import (
  "io";
  "os";
)

type Reader interface {
  ReadAt(buf []byte, pos int64) (n int, err os.Error);
  Read(b []byte) (n int, err os.Error);
  Slice(n int) []byte;
  Seek(off int64, whence int) (ret int64, err os.Error);
  ReadByte() (c byte, err os.Error);
  UnreadByte() os.Error;
  Size() int64;
  Close();
}

// ------------------------------------------------------------------

const (
  _SECTOR_SIZE  = 1024;
  _SECTOR_COUNT = 10;
)

type SecReaderT struct {
  cache     map[int64][]byte;
  age       map[int64]int;
  ticker    int;
  pos, size int64;
  f         io.ReaderAt;
}

func (sr *SecReaderT) access(pos int64) (sl []byte, p int) {
  p = int(pos % _SECTOR_SIZE);
  pos /= _SECTOR_SIZE;
  if s, ok := sr.cache[pos]; ok {
    if sr.age[pos] != sr.ticker {
      sr.ticker++;
      sr.age[pos] = sr.ticker;
    }
    return s, p;
  }
  if len(sr.cache) >= _SECTOR_COUNT {
    a := sr.ticker;
    old := int64(0);
    for k := range sr.age {
      if sr.age[k] < a {
        old = k;
        a = sr.age[k];
      }
    }
    sr.cache[old] = nil, false;
    sr.age[old] = 0, false;
  }
  sr.ticker++;
  sr.cache[pos] = make([]byte, _SECTOR_SIZE);
  sl = sr.cache[pos];
  sr.age[pos] = sr.ticker;
  sr.f.ReadAt(sr.cache[pos], pos*_SECTOR_SIZE);
  return;
}

func (sr *SecReaderT) ReadAt(buf []byte, pos int64) (n int, err os.Error) {
  if pos >= sr.size {
    return 0, os.EOF
  }
  b, p := sr.access(pos);
  for ; p < _SECTOR_SIZE && n < len(buf); p++ {
    buf[n] = b[p];
    n++;
  }
  if secs := (len(buf) - n) / _SECTOR_SIZE; secs > 0 {
    sr.f.ReadAt(buf[n:n+secs*_SECTOR_SIZE], pos+int64(n));
    n += secs * _SECTOR_SIZE;
  }
  if len(buf)-n > 0 {
    b, p = sr.access(pos + int64(n));
    for ; n < len(buf); p++ {
      buf[n] = b[p];
      n++;
    }
  }
  if pos+int64(n) >= sr.size {
    n -= int(pos + int64(n) - sr.size)
  }
  return;
}

func (sr *SecReaderT) Read(b []byte) (n int, err os.Error) {
  n, err = sr.ReadAt(b, sr.pos);
  sr.pos += int64(n);
  return;
}

func (sr *SecReaderT) Seek(off int64, whence int) (ret int64, err os.Error) {
  ret = sr.pos;
  switch whence {
  case 0:
    sr.pos = 0
  case 2:
    sr.pos = sr.size
  }
  sr.pos += off;
  return;
}

func (sr *SecReaderT) ReadByte() (c byte, err os.Error) {
  if sr.pos < sr.size {
    b, p := sr.access(sr.pos);
    c = b[p];
    sr.pos++;
  } else {
    err = os.EOF
  }
  return;
}

func (sr *SecReaderT) UnreadByte() os.Error {
  sr.pos--;
  return nil;
}

func (sr *SecReaderT) Size() int64 { return sr.size }

func (sr *SecReaderT) Close() {}

func (sr *SecReaderT) Slice(n int) []byte {
  return []byte{} // FIXME
}

func SecReader(f io.ReaderAt, size int64) Reader {
  sr := new(SecReaderT);
  sr.f = f;
  sr.size = size;
  return sr;
}

// ------------------------------------------------------------------

type SliceReaderT struct {
  bin []byte;
  pos int64;
}

func (sl *SliceReaderT) ReadAt(b []byte, off int64) (n int, err os.Error) {
  for n := 0; n < len(b); n++ {
    if off >= int64(len(sl.bin)) {
      if n > 0 {
        break
      }
      return n, os.EOF;
    }
    b[n] = sl.bin[off];
    off++;
  }
  return len(b), nil;
}

func (sl *SliceReaderT) Read(b []byte) (n int, err os.Error) {
  n, err = sl.ReadAt(b, sl.pos);
  sl.pos += int64(n);
  return;
}

func (sl *SliceReaderT) Seek(off int64, whence int) (ret int64, err os.Error) {
  ret = sl.pos;
  switch whence {
  case 0:
    sl.pos = 0
  case 2:
    sl.pos = int64(len(sl.bin))
  }
  sl.pos += off;
  return;
}

func (sl *SliceReaderT) Size() int64 { return int64(len(sl.bin)) }

func (sl *SliceReaderT) ReadByte() (c byte, err os.Error) {
  if sl.pos < int64(len(sl.bin)) {
    c = sl.bin[sl.pos];
    sl.pos++;
  } else {
    err = os.EOF
  }
  return;
}

func (sl *SliceReaderT) UnreadByte() os.Error {
  sl.pos--;
  return nil;
}

func (sl *SliceReaderT) Close() {}

func (sl *SliceReaderT) Slice(n int) []byte {
  return []byte{} // FIXME
}

func SliceReader(bin []byte) Reader {
  r := new(SliceReaderT);
  r.bin = bin;
  return r;
}

// ------------------------------------------------------------------
