
build:
	yacc -d -v c11.y
	lex c11.l 
	gcc -g y.tab.c lex.yy.c