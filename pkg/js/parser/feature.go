package parser

type Feature uint64

const (
	FEAT_NONE              Feature = 0
	FEAT_STRICT            Feature = 1 << iota
	FEAT_GLOBAL_ASYNC      Feature = 1 << iota
	FEAT_LET_CONST         Feature = 1 << iota // from es6
	FEAT_SPREAD            Feature = 1 << iota // from es6
	FEAT_BINDING_PATTERN   Feature = 1 << iota // from es6
	FEAT_BINDING_REST_ELEM Feature = 1 << iota // from es6
	FEAT_MODULE            Feature = 1 << iota // from es6
	FEAT_IMPORT_DEC        Feature = 1 << iota // from es6
	FEAT_EXPORT_DEC        Feature = 1 << iota // from es6
	FEAT_META_PROPERTY     Feature = 1 << iota // from es6

	FEAT_POW                      Feature = 1 << iota // from es7
	FEAT_BINDING_REST_ELEM_NESTED Feature = 1 << iota // from es7

	FEAT_ASYNC_AWAIT Feature = 1 << iota // from es8

	FEAT_BAD_ESCAPE_IN_TAGGED_TPL Feature = 1 << iota // from es9
	FEAT_ASYNC_GENERATOR          Feature = 1 << iota // from es9
	FEAT_ASYNC_ITERATION          Feature = 1 << iota // from es9

	FEAT_OPT_CATCH_PARAM Feature = 1 << iota // from es10
	FEAT_JSON_SUPER_SET  Feature = 1 << iota // from es10

	FEAT_CLASS_PRV        Feature = 1 << iota // from es11
	FEAT_OPT_EXPR         Feature = 1 << iota // from es11
	FEAT_NULLISH          Feature = 1 << iota // from es11
	FEAT_BIGINT           Feature = 1 << iota // from es11
	FEAT_DYNAMIC_IMPORT   Feature = 1 << iota // from es11
	FEAT_EXPORT_ALL_AS_NS Feature = 1 << iota // from es11

	FEAT_NUM_SEP      Feature = 1 << iota // from es12
	FEAT_LOGIC_ASSIGN Feature = 1 << iota // from es12

	FEAT_CLASS_PUB_FIELD  Feature = 1 << iota // from es13
	FEAT_CLASS_PRIV_FIELD Feature = 1 << iota // from es13

	FEAT_JSX Feature = 1 << iota

	// not found where in the spec says that the flags of regexp is neened to check
	// even though it's implemented in some other parsers, so flag `FEAT_CHK_REGEXP_FLAGS`
	// is opt-in in mole
	FEAT_CHK_REGEXP_FLAGS   Feature = 1 << iota
	FEAT_REGEXP_UNICODE     Feature = 1 << iota // from es6
	FEAT_REGEXP_STICKY      Feature = 1 << iota // from es6
	FEAT_REGEXP_DOT_ALL     Feature = 1 << iota // from es8
	FEAT_REGEXP_HAS_INDICES Feature = 1 << iota // from es10

)

func (f Feature) On(flag Feature) Feature {
	return f | flag
}

func (f Feature) Off(flag Feature) Feature {
	return f & ^flag
}
