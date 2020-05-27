#include "unistd.h"
#include "string.h"

int external;

void a(char *s)
{
	write(1, s, strlen(s));
}
