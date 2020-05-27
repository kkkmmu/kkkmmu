
m:     file format elf64-x86-64


Disassembly of section .init:

0000000000400400 <_init>:
  400400:	48 83 ec 08          	sub    $0x8,%rsp
  400404:	48 8b 05 ed 0b 20 00 	mov    0x200bed(%rip),%rax        # 600ff8 <__gmon_start__>
  40040b:	48 85 c0             	test   %rax,%rax
  40040e:	74 05                	je     400415 <_init+0x15>
  400410:	e8 4b 00 00 00       	callq  400460 <.plt.got>
  400415:	48 83 c4 08          	add    $0x8,%rsp
  400419:	c3                   	retq   

Disassembly of section .plt:

0000000000400420 <.plt>:
  400420:	ff 35 e2 0b 20 00    	pushq  0x200be2(%rip)        # 601008 <_GLOBAL_OFFSET_TABLE_+0x8>
  400426:	ff 25 e4 0b 20 00    	jmpq   *0x200be4(%rip)        # 601010 <_GLOBAL_OFFSET_TABLE_+0x10>
  40042c:	0f 1f 40 00          	nopl   0x0(%rax)

0000000000400430 <write@plt>:
  400430:	ff 25 e2 0b 20 00    	jmpq   *0x200be2(%rip)        # 601018 <write@GLIBC_2.2.5>
  400436:	68 00 00 00 00       	pushq  $0x0
  40043b:	e9 e0 ff ff ff       	jmpq   400420 <.plt>

0000000000400440 <strlen@plt>:
  400440:	ff 25 da 0b 20 00    	jmpq   *0x200bda(%rip)        # 601020 <strlen@GLIBC_2.2.5>
  400446:	68 01 00 00 00       	pushq  $0x1
  40044b:	e9 d0 ff ff ff       	jmpq   400420 <.plt>

0000000000400450 <__libc_start_main@plt>:
  400450:	ff 25 d2 0b 20 00    	jmpq   *0x200bd2(%rip)        # 601028 <__libc_start_main@GLIBC_2.2.5>
  400456:	68 02 00 00 00       	pushq  $0x2
  40045b:	e9 c0 ff ff ff       	jmpq   400420 <.plt>

Disassembly of section .plt.got:

0000000000400460 <.plt.got>:
  400460:	ff 25 92 0b 20 00    	jmpq   *0x200b92(%rip)        # 600ff8 <__gmon_start__>
  400466:	66 90                	xchg   %ax,%ax

Disassembly of section .text:

0000000000400470 <_start>:
  400470:	31 ed                	xor    %ebp,%ebp
  400472:	49 89 d1             	mov    %rdx,%r9
  400475:	5e                   	pop    %rsi
  400476:	48 89 e2             	mov    %rsp,%rdx
  400479:	48 83 e4 f0          	and    $0xfffffffffffffff0,%rsp
  40047d:	50                   	push   %rax
  40047e:	54                   	push   %rsp
  40047f:	49 c7 c0 20 06 40 00 	mov    $0x400620,%r8
  400486:	48 c7 c1 b0 05 40 00 	mov    $0x4005b0,%rcx
  40048d:	48 c7 c7 8b 05 40 00 	mov    $0x40058b,%rdi
  400494:	e8 b7 ff ff ff       	callq  400450 <__libc_start_main@plt>
  400499:	f4                   	hlt    
  40049a:	66 0f 1f 44 00 00    	nopw   0x0(%rax,%rax,1)

00000000004004a0 <deregister_tm_clones>:
  4004a0:	b8 4f 10 60 00       	mov    $0x60104f,%eax
  4004a5:	55                   	push   %rbp
  4004a6:	48 2d 48 10 60 00    	sub    $0x601048,%rax
  4004ac:	48 83 f8 0e          	cmp    $0xe,%rax
  4004b0:	48 89 e5             	mov    %rsp,%rbp
  4004b3:	77 02                	ja     4004b7 <deregister_tm_clones+0x17>
  4004b5:	5d                   	pop    %rbp
  4004b6:	c3                   	retq   
  4004b7:	b8 00 00 00 00       	mov    $0x0,%eax
  4004bc:	48 85 c0             	test   %rax,%rax
  4004bf:	74 f4                	je     4004b5 <deregister_tm_clones+0x15>
  4004c1:	5d                   	pop    %rbp
  4004c2:	bf 48 10 60 00       	mov    $0x601048,%edi
  4004c7:	ff e0                	jmpq   *%rax
  4004c9:	0f 1f 80 00 00 00 00 	nopl   0x0(%rax)

00000000004004d0 <register_tm_clones>:
  4004d0:	b8 48 10 60 00       	mov    $0x601048,%eax
  4004d5:	55                   	push   %rbp
  4004d6:	48 2d 48 10 60 00    	sub    $0x601048,%rax
  4004dc:	48 c1 f8 03          	sar    $0x3,%rax
  4004e0:	48 89 e5             	mov    %rsp,%rbp
  4004e3:	48 89 c2             	mov    %rax,%rdx
  4004e6:	48 c1 ea 3f          	shr    $0x3f,%rdx
  4004ea:	48 01 d0             	add    %rdx,%rax
  4004ed:	48 d1 f8             	sar    %rax
  4004f0:	75 02                	jne    4004f4 <register_tm_clones+0x24>
  4004f2:	5d                   	pop    %rbp
  4004f3:	c3                   	retq   
  4004f4:	ba 00 00 00 00       	mov    $0x0,%edx
  4004f9:	48 85 d2             	test   %rdx,%rdx
  4004fc:	74 f4                	je     4004f2 <register_tm_clones+0x22>
  4004fe:	5d                   	pop    %rbp
  4004ff:	48 89 c6             	mov    %rax,%rsi
  400502:	bf 48 10 60 00       	mov    $0x601048,%edi
  400507:	ff e2                	jmpq   *%rdx
  400509:	0f 1f 80 00 00 00 00 	nopl   0x0(%rax)

0000000000400510 <__do_global_dtors_aux>:
  400510:	80 3d 2a 0b 20 00 00 	cmpb   $0x0,0x200b2a(%rip)        # 601041 <_edata>
  400517:	75 11                	jne    40052a <__do_global_dtors_aux+0x1a>
  400519:	55                   	push   %rbp
  40051a:	48 89 e5             	mov    %rsp,%rbp
  40051d:	e8 7e ff ff ff       	callq  4004a0 <deregister_tm_clones>
  400522:	5d                   	pop    %rbp
  400523:	c6 05 17 0b 20 00 01 	movb   $0x1,0x200b17(%rip)        # 601041 <_edata>
  40052a:	f3 c3                	repz retq 
  40052c:	0f 1f 40 00          	nopl   0x0(%rax)

0000000000400530 <frame_dummy>:
  400530:	48 83 3d e8 08 20 00 	cmpq   $0x0,0x2008e8(%rip)        # 600e20 <__JCR_END__>
  400537:	00 
  400538:	74 1e                	je     400558 <frame_dummy+0x28>
  40053a:	b8 00 00 00 00       	mov    $0x0,%eax
  40053f:	48 85 c0             	test   %rax,%rax
  400542:	74 14                	je     400558 <frame_dummy+0x28>
  400544:	55                   	push   %rbp
  400545:	bf 20 0e 60 00       	mov    $0x600e20,%edi
  40054a:	48 89 e5             	mov    %rsp,%rbp
  40054d:	ff d0                	callq  *%rax
  40054f:	5d                   	pop    %rbp
  400550:	e9 7b ff ff ff       	jmpq   4004d0 <register_tm_clones>
  400555:	0f 1f 00             	nopl   (%rax)
  400558:	e9 73 ff ff ff       	jmpq   4004d0 <register_tm_clones>

000000000040055d <a>:
  40055d:	55                   	push   %rbp
  40055e:	48 89 e5             	mov    %rsp,%rbp
  400561:	48 83 ec 10          	sub    $0x10,%rsp
  400565:	48 89 7d f8          	mov    %rdi,-0x8(%rbp)
  400569:	48 8b 45 f8          	mov    -0x8(%rbp),%rax
  40056d:	48 89 c7             	mov    %rax,%rdi
  400570:	e8 cb fe ff ff       	callq  400440 <strlen@plt>
  400575:	48 89 c2             	mov    %rax,%rdx
  400578:	48 8b 45 f8          	mov    -0x8(%rbp),%rax
  40057c:	48 89 c6             	mov    %rax,%rsi
  40057f:	bf 01 00 00 00       	mov    $0x1,%edi
  400584:	e8 a7 fe ff ff       	callq  400430 <write@plt>
  400589:	c9                   	leaveq 
  40058a:	c3                   	retq   

000000000040058b <main>:
  40058b:	55                   	push   %rbp
  40058c:	48 89 e5             	mov    %rsp,%rbp
  40058f:	48 83 ec 10          	sub    $0x10,%rsp
  400593:	89 7d fc             	mov    %edi,-0x4(%rbp)
  400596:	48 89 75 f0          	mov    %rsi,-0x10(%rbp)
  40059a:	bf 34 10 60 00       	mov    $0x601034,%edi
  40059f:	e8 b9 ff ff ff       	callq  40055d <a>
  4005a4:	c9                   	leaveq 
  4005a5:	c3                   	retq   
  4005a6:	66 2e 0f 1f 84 00 00 	nopw   %cs:0x0(%rax,%rax,1)
  4005ad:	00 00 00 

00000000004005b0 <__libc_csu_init>:
  4005b0:	41 57                	push   %r15
  4005b2:	41 89 ff             	mov    %edi,%r15d
  4005b5:	41 56                	push   %r14
  4005b7:	49 89 f6             	mov    %rsi,%r14
  4005ba:	41 55                	push   %r13
  4005bc:	49 89 d5             	mov    %rdx,%r13
  4005bf:	41 54                	push   %r12
  4005c1:	4c 8d 25 48 08 20 00 	lea    0x200848(%rip),%r12        # 600e10 <__frame_dummy_init_array_entry>
  4005c8:	55                   	push   %rbp
  4005c9:	48 8d 2d 48 08 20 00 	lea    0x200848(%rip),%rbp        # 600e18 <__init_array_end>
  4005d0:	53                   	push   %rbx
  4005d1:	4c 29 e5             	sub    %r12,%rbp
  4005d4:	31 db                	xor    %ebx,%ebx
  4005d6:	48 c1 fd 03          	sar    $0x3,%rbp
  4005da:	48 83 ec 08          	sub    $0x8,%rsp
  4005de:	e8 1d fe ff ff       	callq  400400 <_init>
  4005e3:	48 85 ed             	test   %rbp,%rbp
  4005e6:	74 1e                	je     400606 <__libc_csu_init+0x56>
  4005e8:	0f 1f 84 00 00 00 00 	nopl   0x0(%rax,%rax,1)
  4005ef:	00 
  4005f0:	4c 89 ea             	mov    %r13,%rdx
  4005f3:	4c 89 f6             	mov    %r14,%rsi
  4005f6:	44 89 ff             	mov    %r15d,%edi
  4005f9:	41 ff 14 dc          	callq  *(%r12,%rbx,8)
  4005fd:	48 83 c3 01          	add    $0x1,%rbx
  400601:	48 39 eb             	cmp    %rbp,%rbx
  400604:	75 ea                	jne    4005f0 <__libc_csu_init+0x40>
  400606:	48 83 c4 08          	add    $0x8,%rsp
  40060a:	5b                   	pop    %rbx
  40060b:	5d                   	pop    %rbp
  40060c:	41 5c                	pop    %r12
  40060e:	41 5d                	pop    %r13
  400610:	41 5e                	pop    %r14
  400612:	41 5f                	pop    %r15
  400614:	c3                   	retq   
  400615:	90                   	nop
  400616:	66 2e 0f 1f 84 00 00 	nopw   %cs:0x0(%rax,%rax,1)
  40061d:	00 00 00 

0000000000400620 <__libc_csu_fini>:
  400620:	f3 c3                	repz retq 

Disassembly of section .fini:

0000000000400624 <_fini>:
  400624:	48 83 ec 08          	sub    $0x8,%rsp
  400628:	48 83 c4 08          	add    $0x8,%rsp
  40062c:	c3                   	retq   
