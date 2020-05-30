#include <sys/types.h>
#include <sys/stat.h>
#include <fcntl.h>
#include <signal.h>
#include <time.h>
#include <unistd.h>
#include <stdio.h>
#include <stdlib.h>

int daemonize(int nochdir,int noclose)
{
	pid_t pid;

	pid = fork ();

	if (pid < 0)
	{
		perror ("fork");
		return -1;
	}

	if (pid != 0)
		exit (0);

	pid = setsid();
	if (pid < -1)
	{
		perror ("setsid");
		return -1;
	}

	if (!nochdir)
		chdir ("/");

	if (! noclose)
	{
		int fd;

		fd = open ("/dev/null", O_RDWR, 0);
		if (fd != -1)
		{
			dup2 (fd, STDIN_FILENO);
			dup2 (fd, STDOUT_FILENO);
			dup2 (fd, STDERR_FILENO);
			if (fd > 2)
				close (fd);
		}
	}

	umask (0027);

	return 0;
}

int main()
{
	daemonize(0, 0);
	while(1)
	{
		printf("Hello world	\n");
	}

	return 0;
}
