
a.out:     file format elf32-i386


Disassembly of section .init:

080482f0 <_init>:
 80482f0:	53                   	push   %ebx
 80482f1:	83 ec 08             	sub    $0x8,%esp
 80482f4:	e8 a7 00 00 00       	call   80483a0 <__x86.get_pc_thunk.bx>
 80482f9:	81 c3 07 1d 00 00    	add    $0x1d07,%ebx
 80482ff:	8b 83 fc ff ff ff    	mov    -0x4(%ebx),%eax
 8048305:	85 c0                	test   %eax,%eax
 8048307:	74 05                	je     804830e <_init+0x1e>
 8048309:	e8 52 00 00 00       	call   8048360 <__libc_start_main@plt+0x10>
 804830e:	83 c4 08             	add    $0x8,%esp
 8048311:	5b                   	pop    %ebx
 8048312:	c3                   	ret    

Disassembly of section .plt:

08048320 <printf@plt-0x10>:
 8048320:	ff 35 04 a0 04 08    	pushl  0x804a004
 8048326:	ff 25 08 a0 04 08    	jmp    *0x804a008
 804832c:	00 00                	add    %al,(%eax)
	...

08048330 <printf@plt>:
 8048330:	ff 25 0c a0 04 08    	jmp    *0x804a00c
 8048336:	68 00 00 00 00       	push   $0x0
 804833b:	e9 e0 ff ff ff       	jmp    8048320 <_init+0x30>

08048340 <__stack_chk_fail@plt>:
 8048340:	ff 25 10 a0 04 08    	jmp    *0x804a010
 8048346:	68 08 00 00 00       	push   $0x8
 804834b:	e9 d0 ff ff ff       	jmp    8048320 <_init+0x30>

08048350 <__libc_start_main@plt>:
 8048350:	ff 25 14 a0 04 08    	jmp    *0x804a014
 8048356:	68 10 00 00 00       	push   $0x10
 804835b:	e9 c0 ff ff ff       	jmp    8048320 <_init+0x30>

Disassembly of section .plt.got:

08048360 <.plt.got>:
 8048360:	ff 25 fc 9f 04 08    	jmp    *0x8049ffc
 8048366:	66 90                	xchg   %ax,%ax

Disassembly of section .text:

08048370 <_start>:
 8048370:	31 ed                	xor    %ebp,%ebp
 8048372:	5e                   	pop    %esi
 8048373:	89 e1                	mov    %esp,%ecx
 8048375:	83 e4 f0             	and    $0xfffffff0,%esp
 8048378:	50                   	push   %eax
 8048379:	54                   	push   %esp
 804837a:	52                   	push   %edx
 804837b:	68 70 85 04 08       	push   $0x8048570
 8048380:	68 10 85 04 08       	push   $0x8048510
 8048385:	51                   	push   %ecx
 8048386:	56                   	push   %esi
 8048387:	68 82 84 04 08       	push   $0x8048482
 804838c:	e8 bf ff ff ff       	call   8048350 <__libc_start_main@plt>
 8048391:	f4                   	hlt    
 8048392:	66 90                	xchg   %ax,%ax
 8048394:	66 90                	xchg   %ax,%ax
 8048396:	66 90                	xchg   %ax,%ax
 8048398:	66 90                	xchg   %ax,%ax
 804839a:	66 90                	xchg   %ax,%ax
 804839c:	66 90                	xchg   %ax,%ax
 804839e:	66 90                	xchg   %ax,%ax

080483a0 <__x86.get_pc_thunk.bx>:
 80483a0:	8b 1c 24             	mov    (%esp),%ebx
 80483a3:	c3                   	ret    
 80483a4:	66 90                	xchg   %ax,%ax
 80483a6:	66 90                	xchg   %ax,%ax
 80483a8:	66 90                	xchg   %ax,%ax
 80483aa:	66 90                	xchg   %ax,%ax
 80483ac:	66 90                	xchg   %ax,%ax
 80483ae:	66 90                	xchg   %ax,%ax

080483b0 <deregister_tm_clones>:
 80483b0:	b8 23 a0 04 08       	mov    $0x804a023,%eax
 80483b5:	2d 20 a0 04 08       	sub    $0x804a020,%eax
 80483ba:	83 f8 06             	cmp    $0x6,%eax
 80483bd:	76 1a                	jbe    80483d9 <deregister_tm_clones+0x29>
 80483bf:	b8 00 00 00 00       	mov    $0x0,%eax
 80483c4:	85 c0                	test   %eax,%eax
 80483c6:	74 11                	je     80483d9 <deregister_tm_clones+0x29>
 80483c8:	55                   	push   %ebp
 80483c9:	89 e5                	mov    %esp,%ebp
 80483cb:	83 ec 14             	sub    $0x14,%esp
 80483ce:	68 20 a0 04 08       	push   $0x804a020
 80483d3:	ff d0                	call   *%eax
 80483d5:	83 c4 10             	add    $0x10,%esp
 80483d8:	c9                   	leave  
 80483d9:	f3 c3                	repz ret 
 80483db:	90                   	nop
 80483dc:	8d 74 26 00          	lea    0x0(%esi,%eiz,1),%esi

080483e0 <register_tm_clones>:
 80483e0:	b8 20 a0 04 08       	mov    $0x804a020,%eax
 80483e5:	2d 20 a0 04 08       	sub    $0x804a020,%eax
 80483ea:	c1 f8 02             	sar    $0x2,%eax
 80483ed:	89 c2                	mov    %eax,%edx
 80483ef:	c1 ea 1f             	shr    $0x1f,%edx
 80483f2:	01 d0                	add    %edx,%eax
 80483f4:	d1 f8                	sar    %eax
 80483f6:	74 1b                	je     8048413 <register_tm_clones+0x33>
 80483f8:	ba 00 00 00 00       	mov    $0x0,%edx
 80483fd:	85 d2                	test   %edx,%edx
 80483ff:	74 12                	je     8048413 <register_tm_clones+0x33>
 8048401:	55                   	push   %ebp
 8048402:	89 e5                	mov    %esp,%ebp
 8048404:	83 ec 10             	sub    $0x10,%esp
 8048407:	50                   	push   %eax
 8048408:	68 20 a0 04 08       	push   $0x804a020
 804840d:	ff d2                	call   *%edx
 804840f:	83 c4 10             	add    $0x10,%esp
 8048412:	c9                   	leave  
 8048413:	f3 c3                	repz ret 
 8048415:	8d 74 26 00          	lea    0x0(%esi,%eiz,1),%esi
 8048419:	8d bc 27 00 00 00 00 	lea    0x0(%edi,%eiz,1),%edi

08048420 <__do_global_dtors_aux>:
 8048420:	80 3d 20 a0 04 08 00 	cmpb   $0x0,0x804a020
 8048427:	75 13                	jne    804843c <__do_global_dtors_aux+0x1c>
 8048429:	55                   	push   %ebp
 804842a:	89 e5                	mov    %esp,%ebp
 804842c:	83 ec 08             	sub    $0x8,%esp
 804842f:	e8 7c ff ff ff       	call   80483b0 <deregister_tm_clones>
 8048434:	c6 05 20 a0 04 08 01 	movb   $0x1,0x804a020
 804843b:	c9                   	leave  
 804843c:	f3 c3                	repz ret 
 804843e:	66 90                	xchg   %ax,%ax

08048440 <frame_dummy>:
 8048440:	b8 10 9f 04 08       	mov    $0x8049f10,%eax
 8048445:	8b 10                	mov    (%eax),%edx
 8048447:	85 d2                	test   %edx,%edx
 8048449:	75 05                	jne    8048450 <frame_dummy+0x10>
 804844b:	eb 93                	jmp    80483e0 <register_tm_clones>
 804844d:	8d 76 00             	lea    0x0(%esi),%esi
 8048450:	ba 00 00 00 00       	mov    $0x0,%edx
 8048455:	85 d2                	test   %edx,%edx
 8048457:	74 f2                	je     804844b <frame_dummy+0xb>
 8048459:	55                   	push   %ebp
 804845a:	89 e5                	mov    %esp,%ebp
 804845c:	83 ec 14             	sub    $0x14,%esp
 804845f:	50                   	push   %eax
 8048460:	ff d2                	call   *%edx
 8048462:	83 c4 10             	add    $0x10,%esp
 8048465:	c9                   	leave  
 8048466:	e9 75 ff ff ff       	jmp    80483e0 <register_tm_clones>

0804846b <add>:
 804846b:	55                   	push   %ebp
 804846c:	89 e5                	mov    %esp,%ebp
 804846e:	8b 55 08             	mov    0x8(%ebp),%edx
 8048471:	8b 45 0c             	mov    0xc(%ebp),%eax
 8048474:	01 c2                	add    %eax,%edx
 8048476:	8b 45 10             	mov    0x10(%ebp),%eax
 8048479:	89 10                	mov    %edx,(%eax)
 804847b:	8b 45 10             	mov    0x10(%ebp),%eax
 804847e:	8b 00                	mov    (%eax),%eax
 8048480:	5d                   	pop    %ebp
 8048481:	c3                   	ret    

08048482 <main>:
 8048482:	8d 4c 24 04          	lea    0x4(%esp),%ecx
 8048486:	83 e4 f0             	and    $0xfffffff0,%esp
 8048489:	ff 71 fc             	pushl  -0x4(%ecx)
 804848c:	55                   	push   %ebp
 804848d:	89 e5                	mov    %esp,%ebp
 804848f:	51                   	push   %ecx
 8048490:	83 ec 24             	sub    $0x24,%esp
 8048493:	89 c8                	mov    %ecx,%eax
 8048495:	8b 40 04             	mov    0x4(%eax),%eax
 8048498:	89 45 e4             	mov    %eax,-0x1c(%ebp)
 804849b:	65 a1 14 00 00 00    	mov    %gs:0x14,%eax
 80484a1:	89 45 f4             	mov    %eax,-0xc(%ebp)
 80484a4:	31 c0                	xor    %eax,%eax
 80484a6:	c7 45 ec 0a 00 00 00 	movl   $0xa,-0x14(%ebp)
 80484ad:	c7 45 f0 0b 00 00 00 	movl   $0xb,-0x10(%ebp)
 80484b4:	c7 45 e8 00 00 00 00 	movl   $0x0,-0x18(%ebp)
 80484bb:	8d 45 e8             	lea    -0x18(%ebp),%eax
 80484be:	50                   	push   %eax
 80484bf:	ff 75 f0             	pushl  -0x10(%ebp)
 80484c2:	ff 75 ec             	pushl  -0x14(%ebp)
 80484c5:	e8 a1 ff ff ff       	call   804846b <add>
 80484ca:	83 c4 0c             	add    $0xc,%esp
 80484cd:	8b 45 e8             	mov    -0x18(%ebp),%eax
 80484d0:	50                   	push   %eax
 80484d1:	ff 75 f0             	pushl  -0x10(%ebp)
 80484d4:	ff 75 ec             	pushl  -0x14(%ebp)
 80484d7:	68 90 85 04 08       	push   $0x8048590
 80484dc:	e8 4f fe ff ff       	call   8048330 <printf@plt>
 80484e1:	83 c4 10             	add    $0x10,%esp
 80484e4:	b8 00 00 00 00       	mov    $0x0,%eax
 80484e9:	8b 55 f4             	mov    -0xc(%ebp),%edx
 80484ec:	65 33 15 14 00 00 00 	xor    %gs:0x14,%edx
 80484f3:	74 05                	je     80484fa <main+0x78>
 80484f5:	e8 46 fe ff ff       	call   8048340 <__stack_chk_fail@plt>
 80484fa:	8b 4d fc             	mov    -0x4(%ebp),%ecx
 80484fd:	c9                   	leave  
 80484fe:	8d 61 fc             	lea    -0x4(%ecx),%esp
 8048501:	c3                   	ret    
 8048502:	66 90                	xchg   %ax,%ax
 8048504:	66 90                	xchg   %ax,%ax
 8048506:	66 90                	xchg   %ax,%ax
 8048508:	66 90                	xchg   %ax,%ax
 804850a:	66 90                	xchg   %ax,%ax
 804850c:	66 90                	xchg   %ax,%ax
 804850e:	66 90                	xchg   %ax,%ax

08048510 <__libc_csu_init>:
 8048510:	55                   	push   %ebp
 8048511:	57                   	push   %edi
 8048512:	56                   	push   %esi
 8048513:	53                   	push   %ebx
 8048514:	e8 87 fe ff ff       	call   80483a0 <__x86.get_pc_thunk.bx>
 8048519:	81 c3 e7 1a 00 00    	add    $0x1ae7,%ebx
 804851f:	83 ec 0c             	sub    $0xc,%esp
 8048522:	8b 6c 24 20          	mov    0x20(%esp),%ebp
 8048526:	8d b3 0c ff ff ff    	lea    -0xf4(%ebx),%esi
 804852c:	e8 bf fd ff ff       	call   80482f0 <_init>
 8048531:	8d 83 08 ff ff ff    	lea    -0xf8(%ebx),%eax
 8048537:	29 c6                	sub    %eax,%esi
 8048539:	c1 fe 02             	sar    $0x2,%esi
 804853c:	85 f6                	test   %esi,%esi
 804853e:	74 25                	je     8048565 <__libc_csu_init+0x55>
 8048540:	31 ff                	xor    %edi,%edi
 8048542:	8d b6 00 00 00 00    	lea    0x0(%esi),%esi
 8048548:	83 ec 04             	sub    $0x4,%esp
 804854b:	ff 74 24 2c          	pushl  0x2c(%esp)
 804854f:	ff 74 24 2c          	pushl  0x2c(%esp)
 8048553:	55                   	push   %ebp
 8048554:	ff 94 bb 08 ff ff ff 	call   *-0xf8(%ebx,%edi,4)
 804855b:	83 c7 01             	add    $0x1,%edi
 804855e:	83 c4 10             	add    $0x10,%esp
 8048561:	39 f7                	cmp    %esi,%edi
 8048563:	75 e3                	jne    8048548 <__libc_csu_init+0x38>
 8048565:	83 c4 0c             	add    $0xc,%esp
 8048568:	5b                   	pop    %ebx
 8048569:	5e                   	pop    %esi
 804856a:	5f                   	pop    %edi
 804856b:	5d                   	pop    %ebp
 804856c:	c3                   	ret    
 804856d:	8d 76 00             	lea    0x0(%esi),%esi

08048570 <__libc_csu_fini>:
 8048570:	f3 c3                	repz ret 

Disassembly of section .fini:

08048574 <_fini>:
 8048574:	53                   	push   %ebx
 8048575:	83 ec 08             	sub    $0x8,%esp
 8048578:	e8 23 fe ff ff       	call   80483a0 <__x86.get_pc_thunk.bx>
 804857d:	81 c3 83 1a 00 00    	add    $0x1a83,%ebx
 8048583:	83 c4 08             	add    $0x8,%esp
 8048586:	5b                   	pop    %ebx
 8048587:	c3                   	ret    
