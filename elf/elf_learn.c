#include "stdio.h"

int add(int a, int b, int *c)
{
	*c = a + b;

	return *c;
}

int main(int argc, char **argv)
{
	int x = 10;
	int y = 11;
	int z = 0;

	add(x, y, &z);

	printf("x = %d, y = %d, z = %d\n", x, y, z);

	return 0;
}
