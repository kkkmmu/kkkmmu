#include <unistd.h>
#include <stdio.h>

int main()
{
	daemon(1, 0);
	while(1)
	{
		printf("Hello world\n");
	}

	return 0;
}
