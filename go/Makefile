
all:
	lex -o ./lex/lex.yy.c ./lex/c11.l && gcc -I./lex -ll ./lex/lex.yy.c -o cxlex  &&  ./cxlex  < ./lex/simple_ota_example.c > simple_ota_example.lex
# 	go run ./yacc/yacc.go -o expr.go -v expr.output -p "expr" expr.y && go run . simple_ota_example.lex $(DEBUG) 
	go run ./yacc/yacc.go -o expr.go -v expr.output -p "expr" expr.y && go run . simple_ota_example.lex $(DEBUG) > ../simple_ota_example.c
	gcc ../simple_ota_example.c -o simple_ota_example; rm -f simple_ota_example