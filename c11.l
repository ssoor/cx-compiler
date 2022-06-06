%{
    #include <stdio.h>
    #include "y.tab.h"

    static int node(int token);

    static int check_type(void);
    extern int sym_type(const char *);  /* returns type from symbol table */
    #define sym_type(identifier) IDENTIFIER /* with no symbol table, fake it */

    int prev_yystate = 0;
    static void switch_yystate(const int newyystate);
%}

%option yylineno

%e  1019
%p  2807
%n  371
%k  284
%a  1213
%o  1117

octal_number   [0-7]
decimal_number   [0-9]
NZ  [1-9]
L   [a-zA-Z_]
A   [a-zA-Z_0-9]
hexadecimal_number   [a-fA-F0-9]
hexadecimal_prefix  (0[xX])
e_number   ([Ee][+-]?{decimal_number}+)
p_number   ([Pp][+-]?{decimal_number}+)
FS  (f|F|l|L)
IS  (((u|U)(l|L|ll|LL)?)|((l|L|ll|LL)(u|U)?))
CP  (u|U|L)
SP  (u8|u|U|L)
ES  (\\(['"\?\\abfnrtv]|[0-7]{1,3}|x[a-fA-F0-9]+))
whitespace_separates  [ \t\v\n\f]

%x block_comment

%%
<block_comment>{
    [^*]*        { node(0); /* eat anything that's not a '*' */ }
    "*"[^*/]*   { node(0); /* eat up '*'s not followed by '/'s */ }
    "*/"        { switch_yystate(prev_yystate);  node(IGNORE); }
}

"/*"                                    { node(0); switch_yystate(block_comment); }
"//".*                                    {  node(IGNORE); /* consume //-comment */ }

"#".*                                    {  node(IGNORE); }

"auto"					{ return node(AUTO); }
"static"				{ return node(STATIC); }
"extern"				{ return node(EXTERN); }
"register"				{ return node(REGISTER); }

"const"					{ return node(CONST); }
"volatile"				{ return node(VOLATILE); }

"signed"				{ return node(SIGNED); }
"unsigned"				{ return node(UNSIGNED); }

"restrict"				{ return node(RESTRICT); }


"char"					{ return node(CHAR); }
"int"					{ return node(INT); }
"long"					{ return node(LONG); }
"short"					{ return node(SHORT); }
"float"					{ return node(FLOAT); }
"double"				{ return node(DOUBLE); }
"_Complex"              { return node(COMPLEX); }

"void"					{ return node(VOID); }
"enum"					{ return node(ENUM); }
"union"					{ return node(UNION); }
"struct"				{ return node(STRUCT); }

"inline"				{ return node(INLINE); }
"return"				{ return node(RETURN); }

"if"					{ return node(IF); }
"else"					{ return node(ELSE); }
"goto"					{ return node(GOTO); }
"for"					{ return node(FOR); }
"break"					{ return node(BREAK); }
"continue"				{ return node(CONTINUE); }
"do"					{ return node(DO); }
"while"					{ return node(WHILE); }
"switch"				{ return node(SWITCH); }
"case"					{ return node(CASE); }
"default"				{ return node(DEFAULT); }

"sizeof"				{ return node(IDENTIFIER); } // node(SIZEOF);
"typedef"				{ return node(TYPEDEF); }

"_Alignas"              { return node(ALIGNAS); }
"_Alignof"              { return node(ALIGNOF); }
"_Atomic"               { return node(ATOMIC); }
"_Bool"                 { return node(BOOL); }
"_Generic"              { return node(GENERIC); }
"_Noreturn"             { return node(NORETURN); }
"_Static_assert"        { return node(STATIC_ASSERT); }
"_Thread_local"         { return node(THREAD_LOCAL); }
"__func__"              { return node(FUNC_NAME); }
"_Imaginary"            { return node(IMAGINARY); }

{L}{A}*					{ return node(check_type()); }

{hexadecimal_prefix}{hexadecimal_number}+{IS}?				    { return node(I_CONSTANT); }
{NZ}{decimal_number}*{IS}?				    { return node(I_CONSTANT); }
"0"{octal_number}*{IS}?				    { return node(I_CONSTANT); }
{CP}?"'"([^'\\\n]|{ES})+"'"		{ return node(I_CONSTANT); }

{decimal_number}+{e_number}{FS}?				    { return node(F_CONSTANT); }
{decimal_number}*"."{decimal_number}+{e_number}?{FS}?			{ return node(F_CONSTANT); }
{decimal_number}+"."{e_number}?{FS}?			    { return node(F_CONSTANT); }
{hexadecimal_prefix}{hexadecimal_number}+{p_number}{FS}?			    { return node(F_CONSTANT); }
{hexadecimal_prefix}{hexadecimal_number}*"."{hexadecimal_number}+{p_number}{FS}?			{ return node(F_CONSTANT); }
{hexadecimal_prefix}{hexadecimal_number}+"."{p_number}{FS}?			    { return node(F_CONSTANT); }

({SP}?\"([^"\\\n]|{ES})*\"{whitespace_separates}*)+	{ return node(STRING_LITERAL); }

"..."					{ return node(ELLIPSIS); }

">>="					{ return node(RIGHT_ASSIGN); }
"<<="					{ return node(LEFT_ASSIGN); }
"+="					{ return node(ADD_ASSIGN); }
"-="					{ return node(SUB_ASSIGN); }
"*="					{ return node(MUL_ASSIGN); }
"/="					{ return node(DIV_ASSIGN); }
"%="					{ return node(MOD_ASSIGN); }
"&="					{ return node(AND_ASSIGN); }
"^="					{ return node(XOR_ASSIGN); }
"|="					{ return node(OR_ASSIGN); }

">>"					{ return node(RIGHT_OP); }
"<<"					{ return node(LEFT_OP); }
"&"					{ return node('&'); }
"|"					{ return node('|'); }
"^"					{ return node('^'); }

"++"					{ return node(INC_OP); }
"--"					{ return node(DEC_OP); }
"~"					{ return node('~'); }


"<="					{ return node(LE_OP); }
">="					{ return node(GE_OP); }
"=="					{ return node(EQ_OP); }
"!="					{ return node(NE_OP); }
"<"					{ return node('<'); }
">"					{ return node('>'); }

"&&"					{ return node(AND_OP); }
"||"					{ return node(OR_OP); }
"!"					{ return node('!'); }

"-"					{ return node('-'); }
"+"					{ return node('+'); }
"*"					{ return node('*'); }
"/"					{ return node('/'); }
"%"					{ return node('%'); }

";"					{ return node(';'); }
","					{ return node(','); }
":"					{ return node(':'); }
"="					{ return node('='); }
"?"					{ return node('?'); }
"."					{ return node('.'); }
"->"					{ return node(PTR_OP); }

"("					{ return node('('); }
")"					{ return node(')'); }
("{"|"<%")				{ return node('{'); }
("}"|"%>")				{ return node('}'); }
("["|"<:")				{ return node('['); }
("]"|":>")				{ return node(']'); }


{whitespace_separates}+					{ node(0); /* whitespace separates tokens */ }
.					{ node(0); /* discard bad characters */ }

%%

/*这一函数在文件（或输入）的末尾调用。 如果函数的返回值是1，就停止解析。
 因此它可以用来解析多个文件。 代码可以写在第三段，这就能够解析多个文件。
 方法是使用 yyin 文件指针指向不同的文件，直到所有的文件都被解析。 
 最后，yywrap() 可以返回 1 来表示解析的结束。*/
int yywrap(void)        /* called at end of input */
{
    return 1;           /* terminate now */
}

int column = 0;
int node(int token)
{
	int i;
	char * newval;

	for (i = 0; yytext[i] != '\0'; i++)
		if (yytext[i] == '\n')
            column = 0;
		else if (yytext[i] == '\t')
			column += 8 - (column % 8);
		else
			column++;
    
    // printf("%d\n",token);
	ECHO;
    // printf("\n");
    return token;
}

static int echo_input(void) {
    int c = input();
    do { if (fputc(c, yyout)) {} } while (0);

    return c;
}

static void switch_yystate(int newyystate){
    prev_yystate = YYSTATE;
    BEGIN(newyystate);
}

static int check_type(void)
{
    switch (sym_type(yytext))
    {
    case TYPEDEF_NAME:                /* previously defined */
        return TYPEDEF_NAME;
    case ENUMERATION_CONSTANT:        /* previously defined */
        return ENUMERATION_CONSTANT;
    default:                          /* includes undefined */
        return IDENTIFIER; }
}