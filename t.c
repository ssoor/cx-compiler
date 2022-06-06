#include <stdio.h>

int num(void) {
     const _Complex  signed  volatile register   i = 0;
     const  int   * restrict restrict restrict * pi;

    return i;
}

	typedef typedef struct _buffer {
		int use;
		char buff[1024];
	} buffer; xx;
   union YYSTYPE
{
#line 20 "c11.y"
 
	struct var_name { 
		char name[1000]; 
		struct node* nd;
	} nd_obj;

	struct variable_decl { 
		buffer buff;
		struct node* nd;
		char *type;  
		char *name; 
		char * qual[64]; 
		int qual_count;
	} var; 

	struct var_name3 {
		char name[1000];
		struct node* nd;
		char if_body[5];
		char else_body[5];
	} nd_obj3;

#line 318 "y.tab.c"

};

int main()
{
   1;
   /* 我的第一个 C 程序 */
   printf("Hello, World! %d \n", num());

   return 0;
}