package parser

type Feature uint64

const (
	FEAT_NONE                     Feature = 0
	FEAT_STRICT                   Feature = 1 << 1
	FEAT_GLOBAL_ASYNC             Feature = 1 << 2
	FEAT_LET_CONST                Feature = 1 << 3 // from es6
	FEAT_SPREAD                   Feature = 1 << 4
	FEAT_BINDING_PATTERN          Feature = 1 << 5
	FEAT_BINDING_REST_ELEM        Feature = 1 << 6
	FEAT_BINDING_REST_ELEM_NESTED Feature = 1 << 7 // from es7
	FEAT_MODULE                   Feature = 1 << 8 // from es6
	FEAT_META_PROPERTY            Feature = 1 << 9 // from es6
	FEAT_CHK_REGEXP_FLAGS         Feature = 1 << 10
	FEAT_REGEXP_DOT_ALL           Feature = 1 << 11 // from es8
	FEAT_REGEXP_HAS_INDICES       Feature = 1 << 12
	FEAT_REGEXP_UNICODE           Feature = 1 << 13 // from es6
	FEAT_REGEXP_STICKY            Feature = 1 << 14 // from es6
	FEAT_ASYNC_ITERATION          Feature = 1 << 15 // from es9
	FEAT_ASYNC_AWAIT              Feature = 1 << 16 // from es8
	FEAT_ASYNC_GENERATOR          Feature = 1 << 17 // from es9
	FEAT_POW                      Feature = 1 << 18 // from es7
)

func (f Feature) On(flag Feature) Feature {
	return f | flag
}

func (f Feature) Off(flag Feature) Feature {
	return f & ^flag
}
