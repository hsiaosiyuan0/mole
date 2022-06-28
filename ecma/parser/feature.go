package parser

type Feature uint64

const (
	FEAT_NONE   Feature = 0
	FEAT_STRICT Feature = 1 << iota
	FEAT_GLOBAL_ASYNC
	FEAT_LET_CONST         // from es6
	FEAT_SPREAD            // from es6
	FEAT_BINDING_PATTERN   // from es6
	FEAT_BINDING_REST_ELEM // from es6
	FEAT_MODULE            // from es6
	FEAT_IMPORT_DEC        // from es6
	FEAT_EXPORT_DEC        // from es6
	FEAT_META_PROPERTY     // from es6

	FEAT_POW                      // from es7
	FEAT_BINDING_REST_ELEM_NESTED // from es7

	FEAT_ASYNC_AWAIT // from es8

	FEAT_BAD_ESCAPE_IN_TAGGED_TPL // from es9
	FEAT_ASYNC_GENERATOR          // from es9
	FEAT_ASYNC_ITERATION          // from es9

	FEAT_OPT_CATCH_PARAM // from es10
	FEAT_JSON_SUPER_SET  // from es10

	FEAT_CLASS_PRV        // from es11
	FEAT_OPT_EXPR         // from es11
	FEAT_NULLISH          // from es11
	FEAT_BIGINT           // from es11
	FEAT_DYNAMIC_IMPORT   // from es11
	FEAT_EXPORT_ALL_AS_NS // from es11

	FEAT_NUM_SEP      // from es12
	FEAT_LOGIC_ASSIGN // from es12

	FEAT_CLASS_PUB_FIELD  // from es13
	FEAT_CLASS_PRIV_FIELD // from es13

	FEAT_JSX
	FEAT_JSX_NS

	// not found where in the spec says that the flags of regexp is needed to check
	// even though it's implemented in some other parsers, so flag `FEAT_CHK_REGEXP_FLAGS`
	// is opt-in in mole
	FEAT_CHK_REGEXP_FLAGS
	FEAT_REGEXP_UNICODE     // from es6
	FEAT_REGEXP_STICKY      // from es6
	FEAT_REGEXP_DOT_ALL     // from es8
	FEAT_REGEXP_HAS_INDICES // from es10

	FEAT_TS
	FEAT_DTS

	FEAT_DECORATOR
)

func (f Feature) On(flag Feature) Feature {
	return f | flag
}

func (f Feature) Off(flag Feature) Feature {
	return f & ^flag
}

func (f Feature) Turn(flag Feature, on bool) Feature {
	if on {
		return f.On(flag)
	}
	return f.Off(flag)
}
