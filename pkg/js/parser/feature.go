package parser

type Feature uint64

const (
	FEAT_NONE              Feature = 0
	FEAT_STRICT            Feature = 1 << 1
	FEAT_GLOBAL_ASYNC      Feature = 1 << 2
	FEAT_LET_CONST         Feature = 1 << 3 // from es6
	FEAT_SPREAD            Feature = 1 << 4 // from es6
	FEAT_BINDING_PATTERN   Feature = 1 << 5 // from es6
	FEAT_BINDING_REST_ELEM Feature = 1 << 6 // from es6
	FEAT_MODULE            Feature = 1 << 7 // from es6
	FEAT_IMPORT_DEC        Feature = 1 << 7 // from es6
	FEAT_EXPORT_DEC        Feature = 1 << 7 // from es6
	FEAT_META_PROPERTY     Feature = 1 << 8 // from es6

	FEAT_POW                      Feature = 1 << 9  // from es7
	FEAT_BINDING_REST_ELEM_NESTED Feature = 1 << 10 // from es7

	FEAT_ASYNC_AWAIT Feature = 1 << 17 // from es8

	FEAT_BAD_ESCAPE_IN_TAGGED_TPL Feature = 1 << 11 // from es9
	FEAT_ASYNC_GENERATOR          Feature = 1 << 12 // from es9
	FEAT_ASYNC_ITERATION          Feature = 1 << 13 // from es9

	FEAT_OPT_CATCH_PARAM Feature = 1 << 14 // from es10

	FEAT_CLASS_PRV      Feature = 1 << 15 // from es11
	FEAT_OPT_EXPR       Feature = 1 << 16 // from es11
	FEAT_NULLISH        Feature = 1 << 17 // from es11
	FEAT_BIGINT         Feature = 1 << 18 // from es11
	FEAT_DYNAMIC_IMPORT Feature = 1 << 19 // from es11

	FEAT_NUM_SEP      Feature = 1 << 20 // from es12
	FEAT_LOGIC_ASSIGN Feature = 1 << 21 // from es12

	FEAT_CLASS_PUB_FIELD  Feature = 1 << 22 // from es13
	FEAT_CLASS_PRIV_FIELD Feature = 1 << 23 // from es13

	// not found where in the spec says that the flags of regexp is neened to check
	// even though it's implemented in some other parsers, so flag `FEAT_CHK_REGEXP_FLAGS`
	// is opt-in in mole
	FEAT_CHK_REGEXP_FLAGS   Feature = 1 << 24
	FEAT_REGEXP_UNICODE     Feature = 1 << 25 // from es6
	FEAT_REGEXP_STICKY      Feature = 1 << 26 // from es6
	FEAT_REGEXP_DOT_ALL     Feature = 1 << 27 // from es8
	FEAT_REGEXP_HAS_INDICES Feature = 1 << 28 // from es10

)

func (f Feature) On(flag Feature) Feature {
	return f | flag
}

func (f Feature) Off(flag Feature) Feature {
	return f & ^flag
}
