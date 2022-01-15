package parser

import "fmt"

type LexerError struct {
	msg  string
	file string
	line int
	col  int
}

func newLexerError(msg, file string, line, col int) *LexerError {
	return &LexerError{
		msg:  msg,
		file: file,
		line: line,
		col:  col,
	}
}

func (e *LexerError) Error() string {
	return fmt.Sprintf("%s at %s(%d:%d)", e.msg, e.file, e.line, e.col)
}

type ParserError struct {
	msg  string
	file string
	line int
	col  int
}

func newParserError(msg, file string, line, col int) *ParserError {
	return &ParserError{
		msg:  msg,
		file: file,
		line: line,
		col:  col,
	}
}

func (e *ParserError) Error() string {
	return fmt.Sprintf("%s at %s(%d:%d)", e.msg, e.file, e.line, e.col)
}

const (
	ERR_UNEXPECTED_CHAR                            = "Unexpected character"
	ERR_UNEXPECTED_TOKEN                           = "Unexpected token"
	ERR_TPL_UNEXPECTED_TOKEN_TYPE                  = "Unexpected token `%s`"
	ERR_UNTERMINATED_COMMENT                       = "Unterminated comment"
	ERR_UNTERMINATED_REGEXP                        = "Unterminated regular expression"
	ERR_UNTERMINATED_STR                           = "Unterminated string constant"
	ERR_INVALID_REGEXP_FLAG                        = "Invalid regular expression flag"
	ERR_IDENT_AFTER_NUMBER                         = "Identifier directly after number"
	ERR_INVALID_NUMBER                             = "Invalid number"
	ERR_TPL_EXPECT_NUM_RADIX                       = "Expected number in radix %s"
	ERR_LEGACY_OCTAL_IN_STRICT_MODE                = "Octal literals are not allowed in strict mode"
	ERR_TPL_LEGACY_OCTAL_ESCAPE_IN                 = "Octal escape sequences are not allowed in template strings"
	ERR_LEGACY_OCTAL_ESCAPE_IN_STRICT_MODE         = "Octal escape sequences are not allowed in strict mode"
	ERR_EXPECTING_UNICODE_ESCAPE                   = "Expecting Unicode escape sequence \\uXXXX"
	ERR_CODEPOINT_OUT_OF_BOUNDS                    = "Code point out of bounds"
	ERR_BAD_ESCAPE_SEQ                             = "Bad character escape sequence"
	ERR_BAD_RUNE                                   = "Bad character"
	ERR_UNTERMINATED_TPL                           = "Unterminated template"
	ERR_INVALID_UNICODE_ESCAPE                     = "Invalid Unicode escape"
	ERR_ILLEGAL_RETURN                             = "Illegal return"
	ERR_ILLEGAL_BREAK                              = "Illegal break"
	ERR_DUP_LABEL                                  = "Label `%s` already declared"
	ERR_UNDEF_LABEL                                = "Undefined label `%s`"
	ERR_ILLEGAL_CONTINUE                           = "Illegal continue"
	ERR_MULTI_DEFAULT                              = "Multiple default clauses"
	ERR_ASSIGN_TO_RVALUE                           = "Assigning to rvalue"
	ERR_INVALID_META_PROP                          = "The only valid meta property for new is `new.target`"
	ERR_META_PROP_OUTSIDE_FN                       = "`new.target` can only be used in functions"
	ERR_DUP_BINDING                                = "Must have a single binding"
	ERR_TPL_BINDING_RESERVED_WORD                  = "Invalid binding `%s`"
	ERR_AWAIT_AS_DEFAULT_VALUE                     = "Await expression cannot be a default value"
	ERR_AWAIT_IN_FORMAL_PARAMS                     = "Await expression can't be used in parameter"
	ERR_TPL_ASSIGN_TO_RESERVED_WORD_IN_STRICT_MODE = "Assigning to `%s` in strict mode"
	ERR_FOR_IN_LOOP_HAS_INIT                       = "for-in loop variable declaration may not have an initializer"
	ERR_FOR_OF_LOOP_HAS_INIT                       = "for-of loop variable declaration may not have an initializer"
	ERR_STRICT_DIRECTIVE_AFTER_NOT_SIMPLE          = "Illegal 'use strict' directive in function with non-simple parameter list"
	ERR_DUP_PARAM_NAME                             = "Parameter name clash"
	ERR_TRAILING_COMMA                             = "Unexpected trailing comma"
	ERR_REST_ELEM_MUST_LAST                        = "Rest element must be last element"
	ERR_DELETE_LOCAL_IN_STRICT                     = "Deleting local variable in strict mode"
	ERR_REDEF_PROP                                 = "Redefinition of property"
	ERR_ILLEGAL_NEWLINE_AFTER_THROW                = "Illegal newline after throw"
	ERR_CONST_DEC_INIT_REQUIRED                    = "Const declarations require an initialization value"
	ERR_TPL_FORBIDED_LEXICAL_NAME                  = "%s is disallowed as a lexically bound name"
	ERR_GETTER_SHOULD_NO_PARAM                     = "Getter must not have any formal parameters"
	ERR_SETTER_SHOULD_ONE_PARAM                    = "Setter must have exactly one formal parameter"
	ERR_ESCAPE_IN_KEYWORD                          = "Keyword must not contain escaped characters"
	ERR_ID_DUP_DEF                                 = "Identifier `%s` has already been declared"
	ERR_WITH_STMT_IN_STRICT                        = "Strict mode code may not include a with statement"
	ERR_CLASS_NAME_REQUIRED                        = "Class name is required"
	ERR_SHORTHAND_PROP_ASSIGN_NOT_IN_DESTRUCT      = "Shorthand property assignments are valid only in destructuring patterns"
	ERR_REST_ARG_NOT_SIMPLE                        = "Invalid rest operator's argument"
	ERR_REST_ARG_NOT_BINDING_PATTERN               = "Binding pattern is not permitted as rest operator's argument"
	ERR_REST_IN_SETTER                             = "Setter cannot use rest params"
	ERR_INVALID_PAREN_ASSIGN_PATTERN               = "Invalid parenthesized assignment pattern"
	ERR_OBJ_PATTERN_CANNOT_FN                      = "Object pattern can't contain getter or setter"
	ERR_INVALID_DESTRUCTING_TARGET                 = "Invalid destructuring assignment target"
	ERR_REST_CANNOT_SET_DEFAULT                    = "Rest elements cannot have a default value"
	ERR_MALFORMED_ARROW_PARAM                      = "Malformed arrow function parameter list"
	ERR_AWAIT_OUTSIDE_ASYNC                        = "Cannot use keyword 'await' outside an async function"
	ERR_AWAIT_AS_NAME_IN_ASYNC                     = "Can not use 'await' as identifier inside an async function"
	ERR_EXPORT_NOT_DEFINED                         = "Export `%s` is not defined"
	ERR_DUP_EXPORT                                 = "Duplicate export `%s`"
	ERR_FN_IN_SINGLE_STMT_CTX                      = "function declarations can't appear in single-statement context"
	ERR_STATIC_PROP_PROTOTYPE                      = "Classes can't have a static field named `prototype`"
	ERR_YIELD_CANNOT_BE_DEFAULT_VALUE              = "Yield expression cannot be a default value"
	ERR_YIELD_IN_FORMAL_PARAMS                     = "Yield expression can't be used in parameter"
	ERR_SUPER_CALL_OUTSIDE_CTOR                    = "super() call outside constructor of a subclass"
	ERR_SUPER_OUTSIDE_CLASS                        = "'super' is only allowed in object methods and classes"
	ERR_CTOR_CANNOT_HAVE_MODIFIER                  = "Constructor can't have get/set modifier"
	ERR_CTOR_CANNOT_BE_GENERATOR                   = "Constructor can't be a generator"
	ERR_CTOR_CANNOT_BE_ASYNC                       = "Constructor can't be a async"
	ERR_CTOR_CANNOT_BE_Field                       = "Classes can't have a field named `constructor`"
	ERR_CTOR_DUP                                   = "Duplicate constructor in the same class"
	ERR_COMPUTE_PROP_MISSING_INIT                  = "A computed property name must have property initialization"
	ERR_IMPORT_EXPORT_SHOULD_AT_TOP_LEVEL          = "'import' and 'export' may only appear at the top level"
	ERR_COMPLEX_BINDING_MISSING_INIT               = "Complex binding patterns require an initialization value"
	ERR_LHS_OF_FOR_OF_CANNOT_ASYNC                 = "The left-hand side of a for-of loop may not be 'async'"
	ERR_UNARY_OPERATOR_IMMEDIATELY_BEFORE_POW      = "Unary operator used immediately before exponentiation expression"
	ERR_TPL_ID_ALREADY_DEF                         = "Identifier `%s` has already been declared"
	ERR_UNEXPECTED_PVT_FIELD                       = "Unexpected private field"
	ERR_DELETE_PVT_FIELD                           = "Private fields can not be deleted"
	ERR_TPL_ALONE_PVT_FIELD                        = "Private field `%s` must be declared in an enclosing class"
	ERR_OPT_EXPR_IN_NEW                            = "Invalid optional chain from new expression"
	ERR_OPT_EXPR_IN_TAG                            = "Invalid tagged template on optional chain"
	ERR_NULLISH_MIXED_WITH_LOGIC                   = "Cannot use unparenthesized `??` within logic expressions"
	ERR_NUM_SEP_BEGIN                              = "Numeric separator is not allowed at the first of digits"
	ERR_NUM_SEP_END                                = "Numeric separator is not allowed at the last of digits"
	ERR_NUM_SEP_DUP                                = "Only one underscore is allowed as numeric separator"
	ERR_NUM_SEP_IN_LEGACY_OCTAL                    = "Numeric separator is not allowed in legacy octal numeric literals"
	ERR_ILLEGAL_IMPORT_PROP                        = "The only valid meta property for import is `import.meta`"
	ERR_META_PROP_CONTAINS_ESCAPE                  = "Meta property can not contain escaped characters"
	ERR_DYNAMIC_IMPORT_CANNOT_NEW                  = "Cannot use new with `import()`"
	ERR_UNTERMINATED_JSX_CONTENTS                  = "Unterminated JSX contents."
	ERR_TPL_UNBALANCED_JSX_TAG                     = "Expected corresponding JSX closing tag for <%s>."
	ERR_JSX_ADJACENT_ELEM_SHOULD_BE_WRAPPED        = "Adjacent JSX elements must be wrapped in an enclosing tag"
	ERR_TPL_JSX_HTML_UNESCAPED_ENTITY              = "Unexpected `%s`, HTML entity "
	ERR_TPL_JSX_UNDEFINED_HTML_ENTITY              = "Undefined HTML entity `%s`"

	ERR_THIS_CANNOT_BE_OPTIONAL      = "The `this` parameter cannot be optional"
	ERR_ILLEGAL_PARAMETER_MODIFIER   = "A parameter property is only allowed in a constructor implementation"
	ERR_CTOR_CANNOT_WITH_TYPE_PARAMS = "Type parameters cannot appear on a constructor declaration"
	ERR_FN_SIG_MISSING_IMPL          = "Function implementation is missing or not immediately following the declaration"
	ERR_TPL_INVALID_FN_IMPL_NAME     = "Function implementation name must be `%s`"
	ERR_TPL_USE_TYP_AS_VALUE         = "`%s` only refers to a type, but is being used as a value here"
	ERR_ASYNC_IN_AMBIENT             = "`async` modifier cannot be used in an ambient context"
	ERR_INIT_NOT_ALLOWED             = "Initializers are not allowed in ambient contexts"
	ERR_UNEXPECTED_TYPE_ANNOTATION   = "Unexpected type annotation"
	ERR_ABSTRACT_MIXED_WITH_STATIC   = "`static` modifier cannot be used with `abstract` modifier"
	ERR_BARE_ABSTRACT_PROPERTY       = "Abstract methods can only appear within an abstract class"
	ERR_ABSTRACT_METHOD_WITH_IMPL    = "Method cannot have an implementation because it's marked abstract"
	ERR_ABSTRACT_PROP_WITH_INIT      = "Property cannot have an initializer because it's marked abstract"
	ERR_OVERRIDE_METHOD_DYNAMIC_NAME = "Method overload name must refer to an expression whose type is a literal type"
)
