#include "stdio.h"
#include "string.h"
#include "stdlib.h"

int main()
{
	int i = 0;
	char buf[2048];
	char *ptr = NULL;

	for (i = 0; i < 10; i++)
		printf("Hello world\n");

	memcpy(buf, "Hello world1\n", sizeof(buf));
	printf("%s\n", buf);

	ptr = malloc(20);
	if (!ptr)
		return -1;

	memcpy(ptr, "Hello world2\n", strlen("Hello world2\n") + 1);

	return 0;
}

/* at&t syntax
  000000000040051d <main>:
  40051d:       55                      push   %rbp               # save base pointer of previous frame onto stack.
  40051e:       48 89 e5                mov    %rsp,%rbp          # Set ESP value as current frame's base pointer.
  400521:       48 83 ec 10             sub    $0x10,%rsp         # Set the new ESP.
  400525:       c7 45 fc 00 00 00 00    movl   $0x0,-0x4(%rbp)    # i = 0;
  40052c:       c7 45 fc 00 00 00 00    movl   $0x0,-0x4(%rbp)    # i = 0;
  400533:       eb 0e                   jmp    400543 <main+0x26> # goto 400543;
  400535:       bf e0 05 40 00          mov    $0x4005e0,%edi     # move the value in 4005e0 to edi
  40053a:       e8 c1 fe ff ff          callq  400400 <puts@plt>  # call 400400(puts)
  40053f:       83 45 fc 01             addl   $0x1,-0x4(%rbp)    # i += 1
  400543:       83 7d fc 09             cmpl   $0x9,-0x4(%rbp)    # compare i with 9
  400547:       7e ec                   jle    400535 <main+0x18> # if i <= 9 goto 400535
  400549:       c9                      leaveq 
  40054a:       c3                      retq   
  40054b:       0f 1f 44 00 00          nopl   0x0(%rax,%rax,1)
 */

/* intel syntax
  000000000040051d <main>:
  40051d:       55                      push   rbp
  40051e:       48 89 e5                mov    rbp,rsp
  400521:       48 83 ec 10             sub    rsp,0x10
  400525:       c7 45 fc 00 00 00 00    mov    DWORD PTR [rbp-0x4],0x0
  40052c:       c7 45 fc 00 00 00 00    mov    DWORD PTR [rbp-0x4],0x0
  400533:       eb 0e                   jmp    400543 <main+0x26>
  400535:       bf e0 05 40 00          mov    edi,0x4005e0
  40053a:       e8 c1 fe ff ff          call   400400 <puts@plt>
  40053f:       83 45 fc 01             add    DWORD PTR [rbp-0x4],0x1
  400543:       83 7d fc 09             cmp    DWORD PTR [rbp-0x4],0x9
  400547:       7e ec                   jle    400535 <main+0x18>
  400549:       c9                      leave  
  40054a:       c3                      ret    
  40054b:       0f 1f 44 00 00          nop    DWORD PTR [rax+rax*1+0x0]
 */

/*
 (gdb) b main
Breakpoint 1 at 0x400525: file hello.c, line 5.
(gdb) r
Starting program: /home/kkkmmu/useful_script/hack/a.out 

Breakpoint 1, main () at hello.c:5
5               int i = 0;
Missing separate debuginfos, use: debuginfo-install glibc-2.17-222.el7.x86_64
(gdb) info registers
rax            0x40051d 4195613
rbx            0x0      0
rcx            0x400550 4195664
rdx            0x7fffffffb5c8   140737488336328
rsi            0x7fffffffb5b8   140737488336312
rdi            0x1      1
rbp            0x7fffffffb4d0   0x7fffffffb4d0
rsp            0x7fffffffb4c0   0x7fffffffb4c0
r8             0x7ffff7dd5e80   140737351868032
r9             0x0      0
r10            0x7fffffffb020   140737488334880
r11            0x7ffff7a30350   140737348043600
r12            0x400430 4195376
r13            0x7fffffffb5b0   140737488336304
r14            0x0      0
r15            0x0      0
rip            0x400525 0x400525 <main+8>
eflags         0x206    [ PF IF ]
cs             0x33     51
ss             0x2b     43
ds             0x0      0
es             0x0      0
fs             0x0      0
gs             0x0      0
(gdb)
 */

