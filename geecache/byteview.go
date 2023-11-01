package geecache

type ByteView struct {
	b []byte
}

func (bv ByteView) Len() int {
	return len(bv.b)
}

func (bv ByteView) ByteSlice() []byte {
	v := make([]byte, bv.Len())
	copy(v, bv.b)
	return v
}

func (bv ByteView) String() string {
	return string(bv.b)
}
