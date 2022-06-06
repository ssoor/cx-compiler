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
	fun function
	typ typedecl
	val vardecl
	vals []vardecl
	stmt statement
	stmts []statement
	ret returnstmt
	expr expression
	exprs []expression
	varr vararrdecl
	vref varref
	vrefField varrefField
	tref typeref
	nref nameref
	fref funcref
	op opexpr
	call callexpr
	typetyp typedecltype
}

%token	COMMENT BLOCK_COMMENT

%token	<lval>	IDENTIFIER

%token	<lval> IGNORE
%token  <lval> '(' ')' '{' '}' ':' ',' '?'

%token	<lval> I_CONSTANT F_CONSTANT STRING_LITERAL FUNC_NAME SIZEOF
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


// token
%type <lval> variable_storage_qualifier
%type <lval> type_ref_qualifier type_qualifier type_declarator
%type <lval> pointer_ops unary_op bit_op math_op logical_op assignment_op

// variable
%type <val>  variable_declaration variable_declaration_item variable_declaration_item_normal
%type <varr>  variable_declaration_item_array
%type <vals> params_declaration params_declaration_items
 

// type
%type <typ> type_declaration 
%type <typetyp> type_block_qualifier
%type <val>  type_declaration_block_item
%type <vals> type_declaration_block type_declaration_block_items

// function
%type <fun> function_declaration  
%type <fref> function_declaration_qualifier function_qualifier
/* %type <val> variable_declaration_item_function
%type <fref> function_declaration_qualifier_block */

// reference
%type <tref> type_ref 
%type <nref> name_ref
%type <vref> variable_ref
%type <vrefField> variable_ref_sub_field

// expression
%type <op> op_expression
%type <call> call_expression
%type <lval> constant_expressions
%type <vref> variable_ref_expression
%type <expr> expression expression_statement case_expression  variable_decl_expression 
%type <exprs> expressions params_assignment

// statement
%type <ret> return_statement
%type <stmt> block_statement if_statement for_statement while_statement switch_statement 
%type <stmts> block_declaration block_statements

%%

global_statements
	: global_statement
	| global_statement global_statements
	;

global_statement
	: ';'
	| IGNORE
	| COMMENT
	| BLOCK_COMMENT
	| type_declaration ';'    		{ dump_lval($1); }
	| variable_declaration ';'   	{ dump_lval($1); }
	| function_declaration   		{ dump_lval($1); }
	;

function_declaration
	: function_declaration_qualifier block_declaration { $$.Typ = $1; $$.Block = $2; }
	;

function_declaration_qualifier
	: function_qualifier
	| INLINE function_declaration_qualifier
	| variable_storage_qualifier function_declaration_qualifier
	;

function_qualifier
	: type_ref name_ref params_declaration { $$.Retval = $1; $$.Name = $2; $$.Params = $3; }
	| type_ref '(' variable_declaration_item_normal ')' IDENTIFIER params_declaration { $$.Retval = $1; $$.Name.Name = $5.text; $$.Params = append(append($$.Params,$3),$6...); $$.Typ = functionSelf; } // 新增语法 - 定义成员方法
	;

type_declaration
	: type_block_qualifier IDENTIFIER { $$.Typ = $1; $$.Name = $2.text; }
	| type_block_qualifier type_declaration_block { $$.Typ = $1; $$.Fields = $2; }
	| type_block_qualifier IDENTIFIER type_declaration_block  { $$.Typ = $1; $$.Name = $2.text; $$.Fields = $3; }
	| TYPEDEF type_ref IDENTIFIER  { $$.Ref = $2; $$.Name = $3.text; $$.Typ = typedeclRef; }
	| TYPEDEF type_declaration IDENTIFIER  { $$ = $2; $$.Name = $3.text; }
	;

type_declaration_block
	: '{' '}'
	| '{' type_declaration_block_items '}'  { $$ = $2; }
	;

type_declaration_block_items
	: type_declaration_block_item  { $$ = append($$, $1); }
	| type_declaration_block_items type_declaration_block_item  { $$ = append($1, $2);  }
	;

type_declaration_block_item
	: variable_declaration_item ';'	 { $$ = $1; } /* for anonymous struct/union */
	;

variable_declaration
	: variable_declaration_item { $$ = $1; }
	| variable_declaration_item '=' expression { $$ = $1; $$.Value = $3; }
	| variable_storage_qualifier variable_declaration { $$ = $2; $$.Storage = $1.text; }
	;

variable_declaration_item
	: variable_declaration_item_normal { $$ = $1; $$.Typ = varType; }
	| variable_declaration_item_normal variable_declaration_item_array { $$ = $1; $$.Arr = $2; $$.Typ = varType; }
	;

variable_declaration_item_normal
	: type_ref IDENTIFIER { $$.RefType = $1; $$.Name = $2.text; }
	;

variable_declaration_item_array
	: '[' ']' { $$.Arr = true; $$.Count = "-1";  }
	| '[' I_CONSTANT ']' { $$.Arr = true; $$.Count = $2.text; }
	;

/* variable_declaration_item_function
	: function_declaration_qualifier_block IDENTIFIER { $$.RefFunc.Typ = $1; $$.Name = $2.text; }
	;

function_declaration_qualifier_block
	: type_ref params_declaration { $$.Retval = $1; $$.Params = $2; } // 新增语法 - 匿名函数
	; */

block_declaration
	: '{'  '}'
	| '{' block_statements '}' { $$ = append($$, $2...) }
	| block_statement   { $$ = append($$, $1) }
	| block_statement BREAK   { $$ = append($$, $1) }
	;

block_statement
	: ';'
	| if_statement { $$ = $1; $$.Typ = stmtIf; }
	| for_statement { $$ = $1; $$.Typ = stmtFor }
	| while_statement { $$ = $1; $$.Typ = stmtWhile }
	| switch_statement { $$ = $1; $$.Typ = stmtSwitch }
	| return_statement { $$.Stmt = $1; $$.Typ = stmtReturn; dump_lval($1); }
	| expression_statement { $$.Stmt = $1; $$.Typ = stmtExpr; dump_lval($1); }
	| type_declaration  ';' { $$.Stmt = $1; $$.Typ = stmtTypeDecl; dump_lval($1); }
	| variable_declaration  ';' { $$.Stmt = $1; $$.Typ = stmtVarDecl; dump_lval($1); }
	;

block_statements
	: block_statement { $$ = append($$, $1) }
	| block_statements block_statement  { $$ = append($1, $2) }
	;

if_statement
	: IF '(' expression ')' block_declaration
	| IF '(' expression ')' block_declaration ELSE block_declaration
	;

for_statement
	: FOR '(' expression_statement expression_statement ')' block_declaration
	| FOR '(' expression_statement expression_statement expression ')' block_declaration
	| FOR '(' variable_declaration ';' expression_statement ')' block_declaration
	| FOR '(' variable_declaration ';' expression_statement expression ')' block_declaration
	;

while_statement
	: WHILE '(' expression ')' block_declaration
	| DO block_declaration WHILE '(' expression ')' ';'
	;

switch_statement
	: SWITCH '(' expression ')' block_declaration
	;

return_statement
	: RETURN ';'
	| RETURN expression ';' { $$.Expr = $2; }
	;

expression_statement
	: expression ';' { $$ = $1;}
	;

expressions
	: expression { $$ = append($$, $1); }
	| expressions ',' expression  { $$ = append($1, $3); }
	;

expression
	: op_expression { $$.Expr = $1; $$.Typ = exprOp; }
	| case_expression
	| call_expression { $$.Expr = $1; $$.Typ = exprCall; }
	| variable_ref_expression { $$.Expr = $1; $$.Typ = exprVar; }
	| variable_decl_expression { $$.Expr = $1; $$.Typ = exprVar; }
	| constant_expressions { $$.Expr = $1.text; $$.Typ = exprConstant; }
	| '(' expression ')'  { $$.Expr = $2; $$.Typ = exprParenthese; }
	| '(' type_ref ')' expression
	;

op_expression
	: expression unary_op 								{ $$.Op = $2.text; $$.L = $1; }
	| unary_op expression  								{ $$.Op = $1.text; $$.R = $2; }
	| expression bit_op expression  	 				{ $$.Op = $2.text; $$.L = $1; $$.R = $3; }
	| expression math_op expression  	 				{ $$.Op = $2.text; $$.L = $1; $$.R = $3; }
	| expression logical_op expression   				{ $$.Op = $2.text; $$.L = $1; $$.R = $3; }
	| expression assignment_op expression 				{ $$.Op = $2.text; $$.L = $1; $$.R = $3; }
	| expression '?' expression ':' expression  		{ $$.Op = $2.text; $$.L = $1; $$.R = $3; $$.There = $5 }
	;

call_expression
	: name_ref params_assignment { $$.Name = $1; $$.Params = $2; } 
	;

case_expression
	: CASE expressions ':' block_declaration
	;

variable_ref_expression
	: variable_ref 											{ $$ = $1 }
	;

variable_decl_expression
	: '[' variable_subfield_expressions ']'
	/* | '[' variable_expression_items ',' ']' */
	| '{' variable_subfield_expressions '}'
	/* | '{' variable_expression_items ',' '}' */
	;

variable_subfield_expressions
	: variable_subfield_expression
	| variable_subfield_expressions ',' variable_subfield_expression
	;

variable_subfield_expression
	: expression
	| '.' IDENTIFIER '=' expression
	;

variable_storage_qualifier
	: AUTO
	| EXTERN
	| STATIC
	| REGISTER
	;

constant_expressions
	: I_CONSTANT			{ $$ = $1; }	/* includes character_constant */
	| F_CONSTANT			{ $$ = $1; }
	| STRING_LITERAL		{ $$ = $1; }
	| ENUMERATION_CONSTANT	{ $$ = $1; }	/* after it has been defined as such */
	;

name_ref
	: IDENTIFIER { $$.Name = $1.text; }
	| IDENTIFIER ':' IDENTIFIER { $$.Name = $3.text; $$.Pkgs = append($$.Pkgs, $1.text); } // 新增语法 - 定义静态方法, 引用指定上下文内容
	;

variable_ref
	: name_ref 											{ $$.Name = $1 }
	| '(' variable_ref ')' 										{ $$ = $2 }
	| pointer_ops variable_ref 									{ $$ = $2; $$.Ops = append($$.Ops, $1.text) }
	| variable_ref variable_ref_sub_field						{ $$ = $1; $$.Fields = append($$.Fields, $2); }
	;

variable_ref_sub_field
	: '.' name_ref 		{ $$.Name = $2; $$.Ptr =false; }
	| PTR_OP name_ref 	{ $$.Name = $2; $$.Ptr =true; }
	;

type_ref 
	: type_ref_qualifier { $$ = append($$, $1.text); }
	| type_ref type_ref_qualifier { $$ = $1; $$ = append($$, $2.text); }
	;

type_ref_qualifier 
	: type_qualifier    { $$ = $1;}
	| type_declarator   { $$ = $1;}
	;

type_declarator
	: IDENTIFIER
	| type_block_qualifier
	| type_pointer_qualifier
	| type_keyword_normal_specifier
	| type_keyword_havesign_specifier
	;

type_pointer_qualifier
	: pointer_ops
	| type_pointer_qualifier type_pointer_access_qualifier
	;

type_pointer_access_qualifier
	: RESTRICT
	| type_access_qualifier
	;
	
type_block_qualifier
	: STRUCT { $$ = typedeclStruct }
	| UNION { $$ = typedeclUnion }
	;

type_keyword_normal_specifier
	: VOID 
	| ENUM
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
	/* | variable_storage_qualifier */
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

params_assignment
	: '(' ')'
	| '(' expressions ')' { $$ = $2 }
	;

params_declaration
	: '(' ')' 
	| '(' params_declaration_items ')' { $$ = $2; }
	;
	
params_declaration_items
	: variable_declaration { $$ = append($$, $1); }
	| params_declaration_items ',' variable_declaration { $$ = append($1, $3);  }
	;

pointer_op
	: '*'
	| '&'
	;

pointer_ops
	: pointer_op
	| pointer_ops pointer_op
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

	exprParse(lexer)
}

func dump_lval(val interface{}) {
	typeName := reflect.TypeOf(val).Name()
	bin,err := json.Marshal(val)
	valueStr := string(bin)
	if err != nil {
		fmt.Printf("%s: %s\n", typeName, err)
	}else {
		fmt.Printf(" >>%s:%s<<", typeName, valueStr)
	}

		return

	/* switch v := val.(type) {
	case lexVal:
		fmt.Printf(" >>%s%s<<",  v.token.String(),valueStr)
	case function:
		fmt.Printf(" >>%s%s<<", v.Typ.String(), valueStr)
	case typedecl:
		fmt.Printf(" >>%s%s<<", v.Typ.String(), valueStr)
	case expression:
		fmt.Printf(" >>%s%s<<", v.Typ.String(), valueStr)
	case statement:
		switch (v.Typ) {
		case stmtExpr:
			dump_lval(v.Stmt)
		default:
			fmt.Printf(" >>%s%s<<", v.Typ.String(), valueStr)
		}
	default:
		fmt.Printf(" >>%s%s<<", typeName, valueStr)
	} */
}