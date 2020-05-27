#include <stdio.h>
#include <stdlib.h>
#include <string.h>
int check_authentication(char *password) {
	int auth_flag = 0;
	char password_buffer[16];
	strcpy(password_buffer, password);
	if(strcmp(password_buffer, "brillig") == 0)
		auth_flag = 1;
	if(strcmp(password_buffer, "outgrabe") == 0)
		auth_flag = 1;
	return auth_flag;
}
int main(int argc, char *argv[]) {
	if(argc < 2) {
		printf("Usage: %s <password>\n", argv[0]);
		exit(0);
	}
	if(check_authentication(argv[1])) {
		printf("\n-=-=-=-=-=-=-=-=-=-=-=-=-=-\n");
		printf(" Access Granted.\n");
		printf("-=-=-=-=-=-=-=-=-=-=-=-=-=-\n");
	} else {
		printf("\nAccess Denied.\n");
	}
}

/* gcc auth_overflow.c
 * ./a.out AAAAAAAAAAAAAAAAAAAAAAAAAAAAA
 * ./a.out 11111111111111111111111111111
 * gdb -q a.out
 * (gdb) break 9
 * Breakpoint 1 at 0x8048421: file auth_overflow.c, line 9.
 * (gdb) break 16
 * Breakpoint 2 at 0x804846f: file auth_overflow.c, line 16.
 * (gdb)
 * The GDB debugger is started with the -q option to suppress the welcome
 * banner, and breakpoints are set on lines 9 and 16. When the program is run,
 * execution will pause at these breakpoints and give us a chance to examine
 * memory.
 * (gdb) run AAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
 * Starting program: /home/reader/booksrc/auth_overflow AAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
 * Breakpoint 1, check_authentication (password=0xbffff9af 'A' <repeats 30 times>) at
 * auth_overflow.c:9
 * 9 strcpy(password_buffer, password);
 * (gdb) x/s password_buffer
 * 0xbffff7a0: ")????o??????)\205\004\b?o??p???????"
 * (gdb) x/x &auth_flag
 * 0xbffff7bc: 0x00000000
 * (gdb) print 0xbffff7bc - 0xbffff7a0
 * $1 = 28
 * (gdb) x/16xw password_buffer
 * 0xbffff7a0: 0xb7f9f729 0xb7fd6ff4 0xbffff7d8 0x08048529
 * 0xbffff7b0: 0xb7fd6ff4 0xbffff870 0xbffff7d8 0x00000000
 * 0xbffff7c0: 0xb7ff47b0 0x08048510 0xbffff7d8 0x080484bb
 * 0xbffff7d0: 0xbffff9af 0x08048510 0xbffff838 0xb7eafebc
 * (gdb)
 */
