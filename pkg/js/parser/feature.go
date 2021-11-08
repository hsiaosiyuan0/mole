package parser

type Feature int

const (
	FEAT_NONE         = 0
	FEAT_STRICT       = 1 << 1
	FEAT_GLOBAL_ASYNC = 1 << 2
)

func (f Feature) On(flag Feature) Feature {
	return f & flag
}

func (f Feature) Off(flag Feature) Feature {
	return f & ^flag
}
