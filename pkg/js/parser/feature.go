package parser

type Feature uint64

const (
	FEAT_NONE                     Feature = 0
	FEAT_STRICT                           = 1 << 1
	FEAT_GLOBAL_ASYNC                     = 1 << 2
	FEAT_LET_CONST                        = 1 << 3 // from es6
	FEAT_SPREAD                           = 1 << 4
	FEAT_BINDING_PATTERN                  = 1 << 5
	FEAT_BINDING_REST_ELEM                = 1 << 6
	FEAT_BINDING_REST_ELEM_NESTED         = 1 << 7 // from es7
	FEAT_MODULE                           = 1 << 8 // from es6
	FEAT_META_PROPERTY                    = 1 << 9 // from es6
	FEAT_CHK_REGEXP_FLAGS                 = 1 << 10
	FEAT_REGEXP_DOT_ALL                   = 1 << 11 // from es8
	FEAT_REGEXP_HAS_INDICES               = 1 << 12
	FEAT_REGEXP_UNICODE                   = 1 << 13 // from es6
	FEAT_REGEXP_STICKY                    = 1 << 14 // from es6
)

func (f Feature) On(flag Feature) Feature {
	return f | flag
}

func (f Feature) Off(flag Feature) Feature {
	return f & ^flag
}
