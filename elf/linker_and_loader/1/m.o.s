
m.o:     file format elf64-x86-64


Disassembly of section .text:

0000000000000000 <main>:
   0:	55                   	push   %rbp
   1:	48 89 e5             	mov    %rsp,%rbp
   4:	48 83 ec 10          	sub    $0x10,%rsp
   8:	89 7d fc             	mov    %edi,-0x4(%rbp)
   b:	48 89 75 f0          	mov    %rsi,-0x10(%rbp)
   f:	bf 00 00 00 00       	mov    $0x0,%edi
  14:	e8 00 00 00 00       	callq  19 <main+0x19>
  19:	c9                   	leaveq 
  1a:	c3                   	retq   
