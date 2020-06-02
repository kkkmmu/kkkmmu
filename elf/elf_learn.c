#include "stdio.h"
int func0(int a, char b);
int func1(int a, char b);
int func2(int a);
int func3(char a);

int add(int a, int b, int *c)
{
	*c = a + b;

	return *c;
}

int sub(int a, int b, int *c)
{
	*c = a - b;
	return *c;
}

int factorial(int n)
{
	if (n == 1)
		return n;

	return n * factorial(n-1);
}

int func(int a, char b)
{
	func1(a, b);

	return a - b;
}

int func1(int a, char b)
{
	func2(a);
	func3(b);

	return a+b;
}

int func2(int a)
{
	int c = 0;
	c = a * a;
	printf("%d\n", c);
	return c;
}

int func3(char b)
{
	int c = 0;
	c = b * b;
	printf("%d\n", c);
	return c;
}

int main(int argc, char **argv)
{
	int x = 6;
	int y = 5;
	int z = 0;

	add(x, y, &z);
	printf("add: x = %d, y = %d, z = %d\n", x, y, z);

	sub(x, y, &z);
	printf("sub: x = %d, y = %d, z = %d\n", x, y, z);

	z = func1(x, y);
	printf("func: x = %d, y = %d, z = %d\n", x, y, z);

	z = factorial(x);
	printf("factorial (%d) = %d\n", x, z);

	return 0;
}
