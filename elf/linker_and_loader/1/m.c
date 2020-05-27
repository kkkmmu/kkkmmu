extern void a(char *);
extern int external;

int main(int argc, char **argv)
{
	static char string[] = "Hello world\n";
	external = 1;
	a(string);
	sleep(1000);
}
