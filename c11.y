%{
	#include <string.h>
	#define YYDEBUG 1

    extern int yylex (void);
	extern char *yytext;
	extern void yyerror(const char *);  /* prints grammar violation message */

	extern void dump_expr(const char * name);
	extern void dump_stmt(const char * name);

	typedef struct _buffer {
		int use;
		char buff[1024];
	} buffer;
	
	static buffer g_buff;

	extern void buff_reset(void);
	extern char* buff_str(const char* srcstr);
	extern char* buff_strcat(const char* srcstr, const char* catstr);
%}

%union { 
	struct node_name { 
		char *name; 
		struct node* nd; 
	} node;

	struct node_type {
		char *name; 
		struct node* nd;   
		char *type;
		char *qual[64];
		int qcount;
	} type;
} 

%token	COMMENT BLOCK_COMMENT

%token	<type>	IDENTIFIER

%token	IGNORE

%token	I_CONSTANT F_CONSTANT STRING_LITERAL FUNC_NAME SIZEOF
%token	PTR_OP INC_OP DEC_OP LEFT_OP RIGHT_OP LE_OP GE_OP EQ_OP NE_OP
%token	AND_OP OR_OP MUL_ASSIGN DIV_ASSIGN MOD_ASSIGN ADD_ASSIGN
%token	SUB_ASSIGN LEFT_ASSIGN RIGHT_ASSIGN AND_ASSIGN
%token	XOR_ASSIGN OR_ASSIGN
%token	TYPEDEF_NAME ENUMERATION_CONSTANT

%token	TYPEDEF EXTERN STATIC AUTO REGISTER INLINE
%token	CONST RESTRICT VOLATILE
%token	<type> BOOL CHAR SHORT INT LONG SIGNED UNSIGNED FLOAT DOUBLE
%token	COMPLEX  
%token	   ELLIPSIS

%token	<type> VOID ENUM STRUCT UNION IMAGINARY

%token	CASE DEFAULT IF ELSE SWITCH WHILE DO FOR GOTO CONTINUE BREAK RETURN

%token	ALIGNAS ALIGNOF ATOMIC GENERIC NORETURN STATIC_ASSERT THREAD_LOCAL
%start global_statements

// \s+\{.*\}

%%

global_statements
	: global_statement 						
	| global_statement global_statements 	
	;

global_statement
	: ';'
	| IDENTIFIER ';'
	| enum_declaration ';'
	| type_declaration ';'
	| typedef_declaration ';'
	| variable_declaration ';'
	| function_declaration
	;

block_statement
	: ';'
	| BREAK ';'
	| if_statement
	| for_statement
	| case_statement
	| while_statement
	| switch_statement
	| return_statement
	| expression_statement
	| assignment_statement
	| enum_declaration  ';'
	| type_declaration  ';'
	| variable_declaration  ';'
	;

function_declaration
	: function_qualifier block_declaration
	| INLINE function_qualifier block_declaration
	| variable_storage_qualifier function_qualifier block_declaration
	;
	
enum_declaration
	: ENUM type_enum_body
	;

type_declaration
	: type_block_specifier
	| type_block_specifier variable_declaration_names
	;

typedef_declaration
	: TYPEDEF type_ref IDENTIFIER
	| TYPEDEF enum_declaration IDENTIFIER
	| TYPEDEF type_block_specifier IDENTIFIER
	;

variable_declaration
	: type_ref variable_declaration_names
	| variable_storage_qualifier type_ref variable_declaration_names
	;

type_enum_body
	: '{' '}'
	| '{' enum_body_items '}'
	;

enum_body_items
	: enum_body_item ','
	| enum_body_items enum_body_item ','
	;

enum_body_item
	: IDENTIFIER
	| IDENTIFIER '=' expression
	;

block_declaration
	: '{'  '}'
	| '{' block_statements '}'
	| block_statement
	| block_statement BREAK
	;

block_statements
	: block_statement
	| block_statements block_statement
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
	| RETURN expression ';'
	;

expression_statement
	: expression ';'
	;

assignment_statement
	: expression assignment_op expression ';'
	;

expressions
	: expression
	| expressions ',' expression
	;

expression
	: op_expression
	| call_expression
	| variable_ref_expression
	| value_blockdecl_expression
	| constant_expression
	| '(' expression ')'
	| '(' type_ref ')' expression
	;

op_expression
	: expression unary_op
	| unary_op expression
	| expression bit_op expression
	| expression math_op expression
	| expression logical_op expression
	| expression '?' expression ':' expression
	;

case_statement
	: DEFAULT ':' block_declaration
	| CASE IDENTIFIER ':' block_declaration
	| CASE expressions ':' block_declaration
	;

call_expression
	: variable_ref_expression params_assignment 
	;

constant_expression
	: I_CONSTANT	/* includes character_constant */
	| F_CONSTANT
	| STRING_LITERAL
	| ENUMERATION_CONSTANT	/* after it has been defined as such */
	;

type_ref 
	: IDENTIFIER '*'
	| type_ref_qualifier
	| type_ref type_ref_qualifier
	;

variable_ref_expression
	: variable_ref_pkgname
	| pointer_ops variable_ref_expression
	| expression variable_ref_subfields
	;

value_blockdecl_expression
	: '[' value_blockdecl_expression_items ']'
	| '[' value_blockdecl_expression_items ',' ']'
	| '{' value_blockdecl_expression_items '}'
	| '{' value_blockdecl_expression_items ',' '}'
	;

value_blockdecl_expression_items
	: value_blockdecl_expression_item
	| value_blockdecl_expression_item
	| value_blockdecl_expression_items ',' value_blockdecl_expression_item
	;

value_blockdecl_expression_item
	: expression
	| '.' IDENTIFIER '=' expression
	;

variable_declaration_names
	: variable_declaration_name
	| variable_declaration_names ',' variable_declaration_name
	;

variable_declaration_name
	: variable_declaration_name_uninit
	| variable_declaration_name_uninit '=' expression
	;

variable_declaration_name_uninit
	: IDENTIFIER
	| IDENTIFIER variable_declaration_name_arrayops
	;

variable_declaration_name_arrayops
	: variable_declaration_name_arrayop
	| variable_declaration_name_arrayops variable_declaration_name_arrayop
	;

variable_declaration_name_arrayop
	: '[' ']'
	| '[' I_CONSTANT ']'
	;

variable_storage_qualifier
	: AUTO
	| EXTERN
	| STATIC
	| REGISTER
	;

type_body_items
	: variable_declaration ';'
	| type_body_items variable_declaration ';'
	;

type_ref_qualifier 
	: type_qualifier
	| type_specifier
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
	: type_block_qualifier IDENTIFIER
	| type_block_qualifier type_block_body
	| type_block_qualifier IDENTIFIER type_block_body
	;
	
type_block_qualifier
	: UNION
	| STRUCT
	;

type_block_body
	: '{' '}'
	| '{' type_body_items '}'
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
	: type_ref variable_ref_pkgname params_declaration
	| type_ref '(' function_qualifier_self_variable ')' IDENTIFIER params_declaration // 新增语法 - 定义成员方法
	;

function_qualifier_self_variable
	: type_ref IDENTIFIER
	;

params_assignment
	: '(' ')'
	| '(' expressions ')'
	;

params_declaration
	: '(' ')' 
	| '(' params_declaration_items ')'
	;
	
params_declaration_items
	: variable_name_havetype
	| params_declaration_items ',' variable_name_havetype
	;
	
variable_name_havetype
	: variable_declaration_name
	| type_ref variable_declaration_name
	| variable_storage_qualifier type_ref variable_declaration_name
	;

variable_ref_pkgname
	: variable_ref_name
	| variable_ref_pkgname ':' variable_ref_name // 新增语法 - 定义静态方法, 引用指定上下文内容
	;

variable_ref_name
	: IDENTIFIER
	;

variable_ref_subfields
	: variable_ref_subfield
	| variable_ref_subfields variable_ref_subfield
	;

variable_ref_subfield
	: variable_subfield_op variable_ref_pkgname
	;

pointer_op
	: '*'
	| '&'
	;

pointer_ops
	: pointer_op
	| pointer_ops pointer_op
	;

variable_subfield_op
	: '.'
	| PTR_OP
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
#include <stdio.h>

extern int column;
extern int yyleng;
extern int yylineno;
extern FILE *yyin, *yyout;

int yydebug = 0;
void yyerror(char const *s)
{
	int i = 0;
	int placen = column;
	fflush(stdout);
	
	if (placen >= yyleng){
		placen -= yyleng;
	}

	if (placen > 0) {
		printf("\n%*s", placen, "");
	} else {
		printf("\n");
	}
	
	for (i=0;i<yyleng;i++) {
		fputc('^', stdout);
	}

	printf(" %s\n", s);
}

extern int yyparse();
int main(int argc, char *argv[])
{
    yyin = fopen(argv[1], "r");
    if(!yyin)
    {
        printf("couldn't open file for reading\n");
        return 1;
    }

	if (argc > 2) {
    	yydebug = 1;
	}

    return yyparse();
}

void dump_expr(const char * name) {
	printf(" >>>%s %s %s<<< ",name, yylval.node.name, name); 
}

void dump_stmt(const char * name) {
	printf("\n>>>%s\n%s\n%s<<<\n", name,yylval.node.name, name); 
}

void dump_type(YYSTYPE* lval) {
	int i;

	printf("!!!!!!!!!!!!!!!!  addr: %p, ", lval);

	for (i = 0; i < lval->type.qcount; i++) {
	  printf("qual: %s, ", lval->type.qual[lval->type.qcount]);
	}
	
	printf("type: %s", lval->type.type);
}

void buff_reset(void){
	buffer *buff = &g_buff;

	buff->use = 0;
}

void* buff_alloc(int size){
	buffer *buff = &g_buff;

	int use = buff->use;
	if (use + size > sizeof(buff->buff)){
		return NULL;
	}

	buff->use = use + size;
	return (void*)&buff->buff[use];

}

char * buff_str(const char* srcstr) {
	int size = strlen(srcstr) + 1;
	char * dststr = (char*)buff_alloc(size);

	memcpy(dststr, srcstr, size);

	return dststr;
}

char * buff_strcat(const char* srcstr, const char* catstr) {
	int size = strlen(srcstr) + strlen(catstr) + 1;
	char * dststr = (char*)buff_alloc(size);
	
	dststr[0] = '\0';
	dststr = strcpy(dststr, srcstr);
	
	return strcat(dststr, catstr);
}