%{
	package main

	import (
		"log"
		"os"
		"fmt"
		"encoding/json"
		"strconv"
		"reflect"
	)
%}

%union {
	lval lexVal
	strs []string
	fun function
	typ typedecl
	etyp enumdecl
	tdef typedefine
	val vardecl
	vals []vardecl
	varr vardeclarrtype
	varrs []vardeclarrtype
	vref varref
	vrefField varrefField
	vrefFields []varrefField
	fref funcref
	bval varblockdecl
	bvalfield varblocksubfield
	bvalfields []varblocksubfield
	fself funcself
	valitem vardeclvaldecl
	valitems []vardeclvaldecl
	pkgname varpkgname
	varname vardeclnametype
	varnames []vardeclnametype
	typetyp typedecltype

	cblock codeblock
	
	expr expression
	exprs expressions
	op opexpr
	call callexpr

	stmt statement
	stmts []statement
	ifstmt ifstmt
	forstmt forstmt
	retstmt returnstmt
	casestmt casestmt
	whilestmt whilestmt
	switchstmt switchstmt
}

%token	<lval> COMMENT BLOCK_COMMENT LINEEND IGNORE

%token	<lval> IDENTIFIER I_CONSTANT F_CONSTANT STRING_LITERAL FUNC_NAME SIZEOF
%token	<lval> PTR_OP INC_OP DEC_OP LEFT_OP RIGHT_OP LE_OP GE_OP EQ_OP NE_OP
%token	<lval> AND_OP OR_OP MUL_ASSIGN DIV_ASSIGN MOD_ASSIGN ADD_ASSIGN
%token	<lval> SUB_ASSIGN LEFT_ASSIGN RIGHT_ASSIGN AND_ASSIGN
%token	<lval> XOR_ASSIGN OR_ASSIGN
%token	<lval> TYPEDEF_NAME ENUMERATION_CONSTANT

%token	<lval> TYPEDEF EXTERN STATIC AUTO REGISTER INLINE
%token	<lval> CONST RESTRICT VOLATILE
%token	<lval> BOOL CHAR SHORT INT LONG SIGNED UNSIGNED FLOAT DOUBLE
%token	<lval> COMPLEX  
%token	<lval> ELLIPSIS

%token	<lval> VOID ENUM STRUCT UNION IMAGINARY

%token	<lval> CASE DEFAULT IF ELSE SWITCH WHILE DO FOR GOTO CONTINUE BREAK RETURN

%token	<lval> ALIGNAS ALIGNOF ATOMIC GENERIC NORETURN STATIC_ASSERT THREAD_LOCAL

%token  <lval> '(' ')' '{' '}' ':' ',' '?' '*'

// token
%type <lval> variable_storage_qualifier
%type <lval> type_ref_qualifier type_qualifier type_specifier
%type <lval> unary_op bit_op math_op logical_op assignment_op pointer_op

// variable
%type <val>  variable_declaration variable_name_havetype
%type <vals> params_declaration params_declaration_items
%type <varr> variable_declaration_name_arrayop
%type <varrs> variable_declaration_name_arrayops
%type <bval> value_blockdecl_expression
%type <bvalfield> value_blockdecl_expression_item
%type <bvalfields> value_blockdecl_expression_items
%type <strs> pointer_ops
%type <fref> variable_declaration_func
 

// type
%type <tdef> typedef_declaration
%type <typ> type_declaration type_block_specifier  type_ref
%type <typetyp> type_block_qualifier
%type <vals> type_block_body type_body_items

// enum
%type <etyp> enum_declaration
%type <valitem> enum_body_item
%type <valitems> type_enum_body  enum_body_items

// function
%type <fun> function_declaration  
%type <fref> function_qualifier
%type <fself> function_qualifier_self_variable
/* %type <fref> function_declaration_qualifier_block */

// reference
%type <valitem> variable_declaration_name 
%type <varname> variable_declaration_name_uninit variable_ref_name
%type <varnames> variable_declaration_names_uninit
%type <valitems> variable_declaration_names 
%type <pkgname> variable_ref_pkgname
%type <vrefField> variable_subfield_op variable_ref_subfield
%type <vrefFields> variable_ref_subfields

// expression
%type <op> op_expression assignment_statement
%type <call> call_expression
%type <lval> constant_expression
%type <vref> variable_ref_expression
%type <expr> expression expression_statement
%type <exprs> expressions params_assignment

// statement
%type <cblock> codeblock_declaration
%type <stmts> codeblock_statements 
%type <stmt> global_statement codeblock_statement      block_statement brack_statement


%type <ifstmt> if_statement
%type <forstmt> for_statement
%type <whilestmt> while_statement
%type <retstmt> return_statement
%type <casestmt> case_statement
%type <switchstmt> switch_statement

%%

global_statements
	: global_statement 		 				{ ST().AddStmt($1); }				
	| global_statements global_statement	{ ST().AddStmt($2); }	
	;

global_statement
	: ';'
	| IDENTIFIER ';'				{ $$.Stmt = $1; $$.Typ = stmtNone; }
	| typedef_declaration ';'		{ $$.Stmt = $1; $$.Typ = stmtTypeDef; }
	| function_declaration   		{ $$.Stmt = $1; $$.Typ = stmtFuncDecl; }
	| enum_declaration  ';' 		{ $$.Stmt = $1; $$.Typ = stmtEnumDecl; }
	| type_declaration  ';' 		{ $$.Stmt = $1; $$.Typ = stmtTypeDecl; }
	| variable_declaration  ';' 	{ $$.Stmt = $1; $$.Typ = stmtVarDecl; }
	;

/* global_statement
	: ';'
	| IDENTIFIER ';'				{ $$.Stmt = $1; $$.Typ = stmtNone; ST().AddStmt($$); }
	| typedef_declaration ';'		{ $$.Stmt = $1; $$.Typ = stmtTypeDef; ST().AddStmt($$); }
	| function_declaration   		{ $$.Stmt = $1; $$.Typ = stmtFuncDecl; ST().AddStmt($$); }
	| enum_declaration  ';' 		{ $$.Stmt = $1; $$.Typ = stmtEnumDecl; ST().AddStmt($$); }
	| type_declaration  ';' 		{ $$.Stmt = $1; $$.Typ = stmtTypeDecl; ST().AddStmt($$); }
	| variable_declaration  ';' 	{ $$.Stmt = $1; $$.Typ = stmtVarDecl; ST().AddStmt($$); }
	; */

codeblock_declaration
	: codeblock_statement  					{ $$.Body = $1; }
	;

codeblock_statements
	: codeblock_statement 					{ $$ = append($$, $1); }
	| codeblock_statements codeblock_statement  { $$ = append($1, $2); }
	;

codeblock_statement
	: ';'  { $$.Typ = stmtNone; }
	| brack_statement { $$ = $1; }
	| block_statement { $$ = $1; }
	| if_statement { $$.Stmt = $1; $$.Typ = stmtIf; }
	| for_statement { $$.Stmt = $1; $$.Typ = stmtFor }
	| case_statement { $$.Stmt = $1; $$.Typ = stmtCase }
	| while_statement { $$.Stmt = $1; $$.Typ = stmtWhile }
	| switch_statement { $$.Stmt = $1; $$.Typ = stmtSwitch; }
	| return_statement { $$.Stmt = $1; $$.Typ = stmtReturn; }
	| expression_statement { $$.Stmt = $1; $$.Typ = stmtExpr; }
	| assignment_statement  { $$.Stmt = $1; $$.Typ = stmtVarDecl; }
	| enum_declaration  ';' { $$.Stmt = $1; $$.Typ = stmtEnumDecl; }
	| type_declaration  ';' { $$.Stmt = $1; $$.Typ = stmtTypeDecl; }
	| variable_declaration  ';' { $$.Stmt = $1; $$.Typ = stmtVarDecl; }
	;

brack_statement
	: BREAK ';'  { $$.Typ = stmtBreak; }
	| codeblock_statement BREAK ';'   {  $$.Stmt = $1; $$.Typ = stmtBreak; }
	;

block_statement
	:  '{'  '}' { $$.Stmt = NewSymbolTable(); $$.Typ = stmtBlock; }
	|  '{' codeblock_statements '}' { $$.Stmt = NewBodySymbolTable($2...); $$.Typ = stmtBlock; }
	;

function_declaration
	: function_qualifier codeblock_declaration 								{ $$.Typ = $1; $$.Body = $2; }
	| INLINE function_qualifier codeblock_declaration 						{ $$.Typ = $2; $$.Body = $3; }
	| variable_storage_qualifier function_qualifier codeblock_declaration 	{ $$.Typ = $2; $$.Body = $3; }
	;
	
enum_declaration
	: ENUM type_enum_body   { $$.Fields = $2; }
	;

type_declaration
	: type_block_specifier						{ $$ = $1; }
	| type_block_specifier variable_declaration_names { $$ = $1; $$.Values = append($$.Values, $2...); }
	;

typedef_declaration
	: TYPEDEF type_ref IDENTIFIER  				{ $$.Name = $3.text; $$.Typ = $2; }
	| TYPEDEF enum_declaration IDENTIFIER  	{ $$.Name = $3.text; $$.Typ = $2; }
	| TYPEDEF type_block_specifier IDENTIFIER  	{ $$.Name = $3.text; $$.Typ = $2; }
	;

variable_declaration
	: type_ref variable_declaration_names 								{ $$.Values = append($$.Values, $2...); $$.RefType = $1; $$.Typ = varType; }
	| variable_declaration_func 										{ $$.RefFunc = $1; $$.Typ = varFunc; }
	| variable_storage_qualifier type_ref variable_declaration_names 	{ $$.Values = append($$.Values, $3...); $$.RefType = $2; $$.Storage = $1.text; $$.Typ = varType; }
	;

variable_declaration_func
	: type_ref '(' variable_declaration_name_uninit ')' params_declaration { $$.Ref.Name = $3; $$.Retval = $1; $$.Params = $5; $$.Typ = functionRef; } // 新增语法 - 定义成员方法
	| variable_declaration_func '=' expression	{ $$ = $1; $$.Ref.Value = $3; }
	;

type_enum_body
	: '{' '}'
	| '{' enum_body_items '}'  { $$ = $2; }
	;

enum_body_items
	: enum_body_item ',' 						{ $$ = append($$, $1); }
	| enum_body_items enum_body_item ',' 	{ $$ = append($1, $2);  }
	;

enum_body_item
	: IDENTIFIER 	 			{ $$.Name.Name = $1.text; }
	| IDENTIFIER '=' expression { $$.Name.Name = $1.text; $$.Value = $3; }
	;

if_statement
	: IF '(' expression ')' codeblock_declaration { $$.Expr = $3; $$.Body = $5; }
	| IF '(' expression ')' codeblock_declaration ELSE codeblock_declaration  { $$.Expr = $3; $$.Body = $5; $$.Else = $7; }
	;

for_statement
	: FOR '(' expression_statement expression_statement ')' codeblock_declaration { $$.Init = $3; $$.Expr = $4; $$.Body = $6; $$.Incr.Typ = exprNone;}
	| FOR '(' expression_statement expression_statement expression ')' codeblock_declaration  { $$.Init = $3; $$.Expr = $4; $$.Incr = $5; $$.Body = $7; }
	| FOR '(' variable_declaration ';'  expression_statement ')' codeblock_declaration { $$.Init = $3; $$.Expr = $5; $$.Body = $7; $$.Incr.Typ = exprNone;}
	| FOR '(' variable_declaration ';'  expression_statement expression ')' codeblock_declaration { $$.Init = $3; $$.Expr = $5; $$.Incr = $6; $$.Body = $8; }
	;

while_statement
	: WHILE '(' expression ')' codeblock_declaration 			{ $$.Expr = $3; $$.Body = $5; }
	| DO codeblock_declaration WHILE '(' expression ')' ';'  	{ $$.Expr = $5; $$.Body = $2; $$.Do = true; }
	;

switch_statement
	: SWITCH '(' expression ')' codeblock_declaration { $$.Expr = $3; $$.Body = $5; }
	;

return_statement
	: RETURN ';'
	| RETURN expression ';' { $$.Expr = $2; }
	;

expression_statement
	: ';'{ $$.Typ = exprNone; }
	| expression ';' { $$ = $1;}
	;

assignment_statement
	: expression assignment_op expression ';'				{ $$.Op = $2.text; $$.L = $1; $$.R = $3; }
	;

expressions
	: expression { $$ = append($$, $1); }
	| expressions ',' expression  { $$ = append($1, $3); }
	;

expression
	: op_expression { $$.Expr = $1; $$.Typ = exprOp; }
	| call_expression { $$.Expr = $1; $$.Typ = exprCall; }
	| variable_ref_expression { $$.Expr = $1; $$.Typ = exprVar; }
	| value_blockdecl_expression { $$.Expr = $1; $$.Typ = exprVar; }
	| constant_expression { $$.Expr = $1.text; $$.Typ = exprConstant; }
	| '(' expression ')'  { $$.Expr = $2; $$.Typ = exprParenthese; }
	| '(' type_ref ')' expression
	;

op_expression
	: expression unary_op 								{ $$.Op = $2.text; $$.L = $1; }
	| unary_op expression  								{ $$.Op = $1.text; $$.R = $2; }
	| expression bit_op expression  	 				{ $$.Op = $2.text; $$.L = $1; $$.R = $3; }
	| expression math_op expression  	 				{ $$.Op = $2.text; $$.L = $1; $$.R = $3; }
	| expression logical_op expression   				{ $$.Op = $2.text; $$.L = $1; $$.R = $3; }
	| expression '?' expression ':' expression  		{ $$.Op = $2.text; $$.L = $1; $$.R = $3; $$.There = $5 }
	;

case_statement
	: DEFAULT ':' { $$.Default = true;}
	| CASE IDENTIFIER ':' { $$.Expr = $2.text;}
	| CASE expressions ':' { $$.Expr = $2; }
	| DEFAULT ':' codeblock_declaration { $$.Default = true; $$.Body = $3; }
	| CASE IDENTIFIER ':' codeblock_declaration { $$.Expr = $2.text; $$.Body = $4; }
	| CASE expressions ':' codeblock_declaration { $$.Expr = $2; $$.Body = $4; }
	;

call_expression
	: variable_ref_expression params_assignment { $$.Var = $1; $$.Params = $2; } 
	;

constant_expression
	: I_CONSTANT			{ $$ = $1; }	/* includes character_constant */
	| F_CONSTANT			{ $$ = $1; }
	| STRING_LITERAL		{ $$ = $1; }
	| ENUMERATION_CONSTANT	{ $$ = $1; }	/* after it has been defined as such */
	;

type_ref 
	: type_ref_qualifier 			{ $$.Ref = append($$.Ref, $1.text); $$.Typ = typedeclRef; }
	| IDENTIFIER pointer_ops		{ $$.Ref = append($$.Ref, $1.text); $$.Ref = append($$.Ref, $2...); $$.Typ = typedeclRef; }
	| type_ref type_ref_qualifier 	{ $$.Ref = append($1.Ref, $2.text); $$.Typ = typedeclRef; }
	;

variable_ref_expression
	: variable_ref_pkgname 										{ $$.Parent = $1; }
	| expression variable_ref_subfields								{ $$.Parent = $1; $$.Fields = append($$.Fields, $2...); }
	;

value_blockdecl_expression
	: '[' value_blockdecl_expression_items ']' { $$.Fields = $2; $$.Arr = true; }
	| '[' value_blockdecl_expression_items ',' ']'
	| '{' value_blockdecl_expression_items '}' { $$.Fields = $2; $$.Arr = false; }
	| '{' value_blockdecl_expression_items ',' '}'
	;

value_blockdecl_expression_items
	: value_blockdecl_expression_item			{ $$ = append($$,$1); }
	| value_blockdecl_expression_item		{ $$ = append($$,$1); }
	| value_blockdecl_expression_items ',' value_blockdecl_expression_item { $$ = append($1,$3); }
	;

value_blockdecl_expression_item
	: expression 						{ $$.Value = $1; }
	| '.' IDENTIFIER '=' expression 	{ $$.Name = $2.text; $$.Value = $4; }
	;

variable_declaration_names
	: variable_declaration_name { $$ = append($$, $1); }
	| variable_declaration_names ',' variable_declaration_name { $$ = append($1, $3); }
	;

variable_declaration_name
	: variable_declaration_name_uninit { $$.Name = $1 }
	| variable_declaration_name_uninit '=' expression { $$.Name = $1; $$.Value = $3; }
	;

variable_declaration_name_uninit
	: IDENTIFIER { $$.Name = $1.text; }
	| pointer_ops variable_declaration_name_uninit { $$ = $2; $$.Ops = append($$.Ops, $1...); }
	| variable_declaration_name_uninit variable_declaration_name_arrayops { $$ = $1;  $$.Arr = $2; }
	;

variable_declaration_names_uninit
	: variable_declaration_name_uninit { $$ = append($$, $1); }
	| variable_declaration_names_uninit ',' variable_declaration_name_uninit { $$ = append($1, $3); }
	;

variable_declaration_name_arrayops
	: variable_declaration_name_arrayop 									{ $$ = append($$, $1); }
	| variable_declaration_name_arrayops variable_declaration_name_arrayop  { $$ = append($1, $2); }
	;

variable_declaration_name_arrayop
	: '[' ']' { $$.Arr = true; $$.Count = "-1";  }
	| '[' I_CONSTANT ']' { $$.Arr = true; $$.Count = $2.text; }
	;

variable_storage_qualifier
	: AUTO
	| EXTERN
	| STATIC
	| REGISTER
	;

type_body_items
	: variable_declaration ';'  						{ $$ = append($$, $1); }
	| type_body_items variable_declaration ';'  	{ $$ = append($1, $2);  }
	;

type_ref_qualifier 
	: type_qualifier    { $$ = $1; }
	| type_specifier	{ $$ = $1; }
	;

type_specifier
	: variable_ref_pkgname
	| enum_declaration
	| type_block_specifier
	| type_pointer_specifier
	| type_keyword_normal_specifier
	| type_keyword_havesign_specifier
	;
	
type_block_specifier
	: type_block_qualifier IDENTIFIER   				{ $$.Typ = $1; $$.Name = $2.text; }
	| type_block_qualifier type_block_body   			{ $$.Typ = $1; $$.Fields = $2; }
	| type_block_qualifier IDENTIFIER type_block_body   { $$.Typ = $1; $$.Name = $2.text; $$.Fields = $3; }
	;
	
type_block_qualifier
	: UNION { $$ = typedeclUnion }
	| STRUCT { $$ = typedeclStruct }
	;

type_block_body
	: '{' '}'
	| '{' type_body_items '}'  { $$ = $2; }
	;

type_pointer_specifier
	: pointer_ops
	| type_pointer_specifier type_pointer_access_qualifier
	;

type_pointer_access_qualifier
	: RESTRICT
	| type_access_qualifier
	;

type_keyword_normal_specifier
	: VOID 
	| BOOL
	| IMAGINARY	  	/* non-mandated extension */
	;
	
type_keyword_havesign_specifier
	: CHAR
	| INT
	| LONG
	| SHORT
	| FLOAT
	| DOUBLE
	| COMPLEX
	;

type_qualifier
	: type_sign_qualifier
	| type_access_qualifier
	;
	
type_sign_qualifier
	: SIGNED
	| UNSIGNED
	;

type_access_qualifier
	: CONST
	| ATOMIC
	| VOLATILE
	;

function_qualifier
	: type_ref variable_ref_pkgname params_declaration { $$.Retval = $1; $$.Name = $2; $$.Params = $3; }
	| type_ref '(' function_qualifier_self_variable ')' IDENTIFIER params_declaration { $$.Retval = $1; $$.Name.Name.Name = $5.text;  $$.Self = $3; $$.Params = $6; $$.Typ = functionSelf; } // 新增语法 - 定义成员方法
	;

function_qualifier_self_variable
	: type_ref IDENTIFIER { $$.Typ = $1; $$.Name = $2.text; }
	;

params_assignment
	: '(' ')'
	| '(' expressions ')' { $$ = $2 }
	;

params_declaration
	: '(' ')' 
	| '(' params_declaration_items ')' { $$ = $2; }
	;
	
params_declaration_items
	: variable_name_havetype { $$ = append($$, $1); }
	| params_declaration_items ',' variable_name_havetype { $$ = append($1, $3);  }
	;
	
variable_name_havetype
	: variable_declaration_name 										{ $$.Values = append($$.Values, $1); $$.Typ = varType; }
	| type_ref variable_declaration_name 								{ $$.Values = append($$.Values, $2); $$.Typ = varType; $$.RefType = $1; }
	| variable_storage_qualifier type_ref variable_declaration_name 	{ $$.Values = append($$.Values, $3); $$.Typ = varType; $$.RefType = $2; $$.Storage = $1.text; }
	;

variable_ref_pkgname
	: variable_declaration_name_uninit { $$.Name = $1; }
	| IDENTIFIER ':' variable_declaration_name_uninit {  $$.Name = $3; $$.Pkgs = append($$.Pkgs, $1.text); } // 新增语法 - 定义静态方法, 引用指定上下文内容
	;

variable_ref_name
	: IDENTIFIER { $$.Name = $1.text; }
	;

variable_ref_subfields
	: variable_ref_subfield 						{ $$ = append($$, $1); }
	| variable_ref_subfields variable_ref_subfield 	{ $$ = append($1, $2); }
	;

variable_ref_subfield
	: variable_subfield_op variable_ref_pkgname { $$.Name = $2; $$.Ptr = $1.Ptr; }
	;

pointer_op
	: '*'
	| '&'
	;

pointer_ops
	: pointer_op 				{ $$ = append($$, $1.text); }
	| pointer_ops pointer_op 	{ $$ = append($1, $2.text); }
	;

variable_subfield_op
	: '.' 		{ $$.Ptr =false; }
	| PTR_OP 	{ $$.Ptr =true; }
	;

unary_op
	: INC_OP
	| DEC_OP
	| '!'
	| '~'
	;

assignment_op
	: '='
	| MUL_ASSIGN
	| DIV_ASSIGN
	| MOD_ASSIGN
	| ADD_ASSIGN
	| SUB_ASSIGN
	| AND_ASSIGN
	| OR_ASSIGN
	| XOR_ASSIGN
	| LEFT_ASSIGN
	| RIGHT_ASSIGN
	;

bit_op
	: LEFT_OP
	| RIGHT_OP
	| '&'
	| '|'
	| '^'
	;

math_op
	: '+'
	| '-'
	| '*'
	| '/'
	| '%'
	;

logical_op
	: EQ_OP
	| NE_OP
	| GE_OP
	| LE_OP
	| '>'
	| '<'
	| AND_OP
	| OR_OP
	/* | '!' */
	;

%%

func main() {
	exprErrorVerbose = true
	lexer,err := newLex(os.Args[1])
	if err != nil {
		log.Fatalf("ReadBytes: %s", err)
	}

	if len(os.Args) > 2 {
		exprDebug, _ = strconv.Atoi(os.Args[2])
	}

	if exprParse(lexer) == 0 {
		fmt.Println(ST().Source())
	}

}

func dump_lval(val interface{}) {
	
	msg := ""
	if str, ok := val.(string); ok {
		msg = str
	} else if strm, ok := val.(fmt.Stringer); ok {
		msg = strm.String()
	} else if val == nil {
		msg = "<nil>"
	} else {
		bin,err := json.Marshal(val)
	
		if err != nil {
			msg = err.Error()
		}else {
			msg = string(bin)
		}
	}

	valueStr := msg

		typeName := reflect.TypeOf(val).Name()
	switch val.(type) {
	case nil:
	case function:
		fmt.Printf("\n>>>>>>>>>>%s\n%s\n",typeName, valueStr)
	case expression:
		/* dump_lval(v.Expr); */
	case varref:
			/* fmt.Printf("<%s>[varref]", valueStr) */
	default:
		/* fmt.Printf(" >> %s", valueStr) */
		/* fmt.Printf(" >> %s [%s]", valueStr, typeName) */
	}
}