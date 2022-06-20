/* Hello World Example

   This example code is in the Public Domain (or CC0 licensed, at your option.)

   Unless required by applicable law or agreed to in writing, this
   software is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
   CONDITIONS OF ANY KIND, either express or implied.
*/
#include <stdio.h>
#include <string.h>

// 单行注释

/*
多行注释
*/

typedef int int_t; // 声明类型

union union_t
{ // 声明联合体
    int a;
    int_t b;
    int_t *c, d, e;
    // auto int_t d;
    // static int_t e;
    // extern int_t f;
    volatile int_t g;
    volatile int_t g1[1];
    volatile int_t g2[1];
    volatile int_t g3[1][1];
};

union union_t union_v1 = {0};              // 声明联合体变量
union union_t union_v2 = {.a = 0, .b = 1}; // 声明联合体变量

struct struct_t
{ // 声明结构体
    int a;
    int_t b;
    int_t *c, d, e;
    // auto int_t d;
    // static int_t e;
    // extern int_t f;
    volatile int_t g;
    volatile int_t g1[1];
    volatile int_t g2[1];
    volatile int_t g3[1][1];
};

struct struct_t struct_v1 = {0};              // 声明结构体变量
struct struct_t struct_v2 = {.a = 0, .b = 1}; // 声明结构体变量

struct_v8;

enum
{ // 声明枚举
    enum_a,
    enum_b = 2,
    enum_c = 2 + 3,
    enum_d = enum_a + enum_b,
};

// void empv; error: storage size of ‘empv’ isn’t known
char strc;
int numberi;
long numberl;
short numbers;
float numberf;
double numberd;

char strc_arr[] = {0};
int numberi_arr[] = {0};
long numberl_arr[] = {0};
short numbers_arr[] = {0};
float numberf_arr[] = {0};
double numberd_arr[] = {0};

char strc_arr1[1] = {0};
int numberi_arr1[1] = {0};
long numberl_arr1[1] = {0};
short numbers_arr1[1] = {0};
float numberf_arr1[1] = {0};
double numberd_arr1[1] = {0};

char strc_arr2[1][1] = {0};
int numberi_arr2[1][1] = {0};
long numberl_arr2[1][1] = {0};
short numbers_arr2[1][1] = {0};
float numberf_arr2[1][1] = {0};
double numberd_arr2[1][1] = {0};

// void empv_a, empv_b; error: storage size of ‘empv_a’ isn’t known
char strc_a, strc_b;
int numberi_a, numberi_b;
long numberl_a, numberl_b;
short numbers_a, numbers_b;
float numberf_a, numberf_b;
double numberd_a, numberd_b;

// void empv1 = 0; // error: variable ‘empv1’ has initializer but incomplete type
char strc1 = 0;
int numberi1 = 0;
long numberl1 = 0;
short numbers1 = 0;
float numberf1 = 0.0;
double numberd1 = 0.0;

// void empv3 = {0}; // error: variable ‘empv3’ has initializer but incomplete type
char strc3 = {0};
int numberi3 = {0};
long numberl3 = {0};
short numbers3 = {0};
float numberf3 = {0.0};
double numberd3 = {0.0};

// void empv2 = 1; // error: variable ‘empv2’ has initializer but incomplete type
char strc2 = 1;
int numberi2 = 1;
long numberl2 = 1;
short numbers2 = 1;
float numberf2 = 1.1;
double numberd2 = 1.1;

// void empv4 = {1}; // error: variable ‘empv4’ has initializer but incomplete type
char strc4 = {1};
int numberi4 = {1};
long numberl4 = {1};
short numbers4 = {1};
float numberf4 = {1.1};
double numberd4 = {1.1};

// void empv_c = 0, empv_d = {0}; // error: variable ‘empv_c’ has initializer but incomplete type
char strc_c = 0, strc_d = {0};
int numberi_c = 0, numberi_d = {0};
long numberl_c = 0, numberl_d = {0};
short numbers_c = 0, numbers_d = {0};
float numberf_c = 0, numberf_d = {0};
double numberd_c = 0, numberd_d = {0};

// long void empvl; // error: both ‘long’ and ‘void’ in declaration specifiers
// long char strlc; // error: both ‘long’ and ‘char’ in declaration specifiers
long int numberli;
long long numberll;
// long short numberls; // error: both ‘long’ and ‘short’ in declaration specifiers
// long float numberlf; // error: both ‘long’ and ‘float’ in declaration specifiers
long double numberld;

// signed void empsv; // error: both ‘signed’ and ‘void’ in declaration specifiers
signed char strsc;
signed int numbersi;
signed long numbersl;
signed short numberss;
// signed float numbersf; // error: both ‘signed’ and ‘float’ in declaration specifiers
// signed double numbersd; // error: both ‘signed’ and ‘double’ in declaration specifiers

// unsigned void empuv; // error: both ‘unsigned’ and ‘void’ in declaration specifiers
unsigned char struc;
unsigned int numberui;
unsigned long numberul;
unsigned short numberus;
// unsigned float numberuf; // error: both ‘unsigned’ and ‘float’ in declaration specifiers
// unsigned double numberud; // error: both ‘unsigned’ and ‘double’ in declaration specifiers

const int const_1 = 212;
const int const_3 = 0xFeeL;
const int const_4 = 0777;
const int const_2 = 215u;
const int const_5 = 032U;

const char const_6 = 'a';
const char const_7[] = "合法的\n";
const char const_8[] = "合法的跨行\
                            合法的跨行结束";

// auto int storagea; // error: file-scope declaration of ‘stdina’ specifies ‘auto’
extern int storagee;
static int storages;
// register int storager; // error: register name not specified for ‘storager’
volatile int storagev;

int op_1 = 1 + 1;
int op_2 = 1 - 1;
int op_3 = 1 * 1;
int op_4 = 1 / 1;
int op_5 = 1 % 1;
// int op_6_1 = op_5++; // error: initializer element is not constant
// int op_6_2 = ++op_5; // error: initializer element is not constant
// int op_7_1 = --op_5; // error: initializer element is not constant
// int op_7_2 = op_5--; // error: initializer element is not constant
int op_8 = 1 == 1;
int op_9 = 1 != 1;
int op_10 = 1 > 1;
int op_11 = 1 >= 1;
int op_12 = 1 < 1;
int op_13 = 1 <= 1;
int op_14 = 1 && 1;
int op_15 = 1 || 1;
int op_16 = !1;
int op_17 = 1 & 1;
int op_18 = 1 | 1;
int op_19 = 1 ^ 1;
int op_20 = ~1;
int op_21 = 1 << 1;
int op_22 = 1 >> 1;
int op_23 = 1;
// int op_24 += 1;
// int op_25 -= 1;
// int op_26 *= 1;
// int op_27 /= 1;
// int op_28 %= 1;
// int op_29 <<= 1;
// int op_30 >>= 1;
// int op_31 &= 1;
// int op_32 |= 1;
// int op_33 ^= 1;
int *op_33 = &op_23;
// int op_34 = *op_33; // error: initializer element is not constant
int op_35 = 1 ? 0 : 1;
int op_36 = sizeof(op_35);
// int op_37 = fun1(); //  error: initializer element is not constant

numberi; // warning: data definition has no type or storage class

int funa(int p1, int p2)
{
    int op_1 = 1 + 1;
    int op_2 = 1 - 1;
    int op_3 = 1 * 1;
    int op_4 = 1 / 1;
    int op_5 = 1 % 1;
    int op_6_1 = op_5++; // error: initializer element is not constant
    int op_6_2 = ++op_5; // error: initializer element is not constant
    int op_7_1 = --op_5; // error: initializer element is not constant
    int op_7_2 = op_5--; // error: initializer element is not constant
    int op_8 = 1 == 1;
    int op_9 = 1 != 1;
    int op_10 = 1 > 1;
    int op_11 = 1 >= 1;
    int op_12 = 1 < 1;
    int op_13 = 1 <= 1;
    int op_14 = 1 && 1;
    int op_15 = 1 || 1;
    int op_16 = !1;
    int op_17 = 1 & 1;
    int op_18 = 1 | 1;
    int op_19 = 1 ^ 1;
    int op_20 = ~1;
    int op_21 = 1 << 1;
    int op_22 = 1 >> 1;
    int op_23 = 1;
    // int op_24 += 1;
    // int op_25 -= 1;
    // int op_26 *= 1;
    // int op_27 /= 1;
    // int op_28 %= 1;
    // int op_29 <<= 1;
    // int op_30 >>= 1;
    // int op_31 &= 1;
    // int op_32 |= 1;
    // int op_33 ^= 1;
    int *op_33 = &op_23;
    int op_34 = *op_33; // error: initializer element is not constant
    int op_35 = 1 ? 0 : 1;
    int op_36 = sizeof(op_35);
    int op_37 = fun1(); //  error: initializer element is not constant

    op_1 = 1 + 1;
    op_2 = 1 - 1;
    op_3 = 1 * 1;
    op_4 = 1 / 1;
    op_5 = 1 % 1;
    op_6_1 = op_5++; // error: initializer element is not constant
    op_6_2 = ++op_5; // error: initializer element is not constant
    op_7_1 = --op_5; // error: initializer element is not constant
    op_7_2 = op_5--; // error: initializer element is not constant
    op_8 = 1 == 1;
    op_9 = 1 != 1;
    op_10 = 1 > 1;
    op_11 = 1 >= 1;
    op_12 = 1 < 1;
    op_13 = 1 <= 1;
    op_14 = 1 && 1;
    op_15 = 1 || 1;
    op_16 = !1;
    op_17 = 1 & 1;
    op_18 = 1 | 1;
    op_19 = 1 ^ 1;
    op_20 = ~1;
    op_21 = 1 << 1;
    op_22 = 1 >> 1;
    op_23 = 1;
    op_23 += 1;
    op_23 -= 1;
    op_23 *= 1;
    op_23 /= 1;
    op_23 %= 1;
    op_23 <<= 1;
    op_23 >>= 1;
    op_23 &= 1;
    op_23 |= 1;
    op_23 ^= 1;
    op_33 = &op_23;
    op_34 = *op_33; // error: initializer element is not constant
    op_35 = 1 ? 0 : 1;
    op_36 = sizeof(op_35);
    op_37 = fun1(); //  error: initializer element is not constant

    struct struct_t struct_v1;
    struct struct_t *struct_v2;

    struct_v1.a = struct_v1.a;
    ((struct_v2)->a) = struct_v2->a;

    auto int storagea;
    extern int storagee;
    static int storages;
    register int storager;
    volatile int storagev;

    op_1 = op_2 + op_3 - op_4 * op_5 / op_8 % op_9;
}

int fun1()
{
    return 1;
}

inline int fun2()
{
    return 1;
}

static int fun3()
{
    return 1;
}

static int fun4(int p1)
{
    return 1;
}

// typedef int (*call[])();
// typedef int (*calls[5])();

int fun_if()
{
    int a = 1;
    if (a)
    {
        return 1;
    }

    if (!a)
    {
        return 1;
    }

    if (a)
    {
        return 1;
    }
    else
    {
        return 0;
    }

    if (a)
    {
        return 1;
    }
    else if (a == 2)
    {
        return 2;
    }

    if (a)
    {
        return 1;
    }
    else
    {
        if (a == 2)
        {
            return 2;
        }
    }
}

int for_for()
{
    int a = 1;

    for (int i = 0; i < 10; i++)
    {
        break;
    }

    for (int i = 0;; i < 10)
    {
        break;
    }

    for (int i = 0; i < 10;)
    {
        break;
    }

    for (; a < 10; a + 0)
    {
        break;
    }

    for (; a < 10;)
    {
        break;
    }

    for (a < 10;;)
    {
        break;
    }

    for (;; a < 10)
    {
        break;
    }

    for (;;)
    {
        break;
    }
}

int for_while()
{
    int a = 1;
    do
    {
        int b = 0;
    } while (a);

    while (a)
    {
        int c = a + 0;
        c++;
    }
}

int fun_switch()
{
    int str = 0, a, b, c, d[1], e, f, g, h, i, j, k, l, m, n, o, p, q, r, s, t, u, v, w, x, y, **z[1][1];
    int (*calls[5])() = {0};

    switch (str)
    {
    case 1:
    case enum_a:
    default:
        break;
    }

    switch (str)
    {
    case 1:
        break;
    case enum_a:
        break;
    default:
        break;
    }

    switch (str)
    {
    case 1:
    {
        break;
    }
    case enum_a:
    {
        break;
    }
    default:
    {
        break;
    }
    }

    switch (str)
    {
    case 1:
        break;
    case enum_a:
        break;
    default:
        break;
    }

    switch (str)
    {
    case 1:
    {
        break;
    }
    case enum_a:
    {
        break;
    }
    default:
    {
        break;
    }
    }

    switch (str)
    {
    case 1:
        calls[0](a, b + c, d[1], sizeof(str) * a, e * f, g - h / i * s);
        break;
    case enum_a:
        calls[0](a, b + c, d[1], sizeof(str) * a, e * f, g - h / i * s);
        break;
    default:
        calls[0](a, b + c, d[1], sizeof(str) * a, e * f, g - h / i * s);
        break;
    }

    switch (str)
    {
    case 1:
    {
        calls[0](a, b + c, d[1], sizeof(str) * a, e * f, g - h / i * s);
        break;
    }
    case enum_a:
    {
        calls[0](a, b + c, d[1], sizeof(str) * a, e * f, g - h / i * s);
        break;
    }
    default:
    {
        calls[0](a, b + c, d[1], sizeof(str) * a, e * f, g - h / i * s);
        break;
    }
    }
}

typedef struct string
{
    int len;
    char ptr[1]; // '\0'
} string;

const char *cat3(string *str, string *str2, string *str3)
{
    const char *(*call1)(string * str) = func(string * str) const char *
    {
        return func(string * str) const char *
        {
            return str;
        }
        (str);
    };

    return call1(str2);
}

int main()
{
    return 0;
}