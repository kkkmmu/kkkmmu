__attribute__((section("FOO"))) int global = 100;
__attribute__((section("BAR"))) int function(int var) { return 1; }
int main()
{
	function(global);
}

/*
 * gcc -c attribute.c
 * objdump -h attribute.o
 * objdump -h attribute.o        
 * attribute.o:     file format elf64-x86-64
 *
 * Sections:
 * Idx Name          Size      VMA               LMA               File off  Algn
 * 0 .text         00000013  0000000000000000  0000000000000000  00000040  2**0
 *                 CONTENTS, ALLOC, LOAD, RELOC, READONLY, CODE
 * 1 .data         00000000  0000000000000000  0000000000000000  00000053  2**0
 *                 CONTENTS, ALLOC, LOAD, DATA
 * 2 .bss          00000000  0000000000000000  0000000000000000  00000053  2**0
 *                 ALLOC
 * 3 FOO           00000004  0000000000000000  0000000000000000  00000054  2**2
 *                 CONTENTS, ALLOC, LOAD, DATA
 * 4 BAR           0000000e  0000000000000000  0000000000000000  00000058  2**0
 *                 CONTENTS, ALLOC, LOAD, READONLY, CODE
 * 5 .comment      0000002e  0000000000000000  0000000000000000  00000066  2**0
 *                 CONTENTS, READONLY
 * 6 .note.GNU-stack 00000000  0000000000000000  0000000000000000  00000094  2**0
 *                 CONTENTS, READONLY
 * 7 .eh_frame     00000058  0000000000000000  0000000000000000  00000098  2**3
 *                 CONTENTS, ALLOC, LOAD, RELOC, READONLY, DATA
 *
 * hexdump -C attribute.o 
 * 00000000  7f 45 4c 46 02 01 01 00  00 00 00 00 00 00 00 00  |.ELF............|
 * 00000010  01 00 3e 00 01 00 00 00  00 00 00 00 00 00 00 00  |..>.............|
 * 00000020  00 00 00 00 00 00 00 00  18 03 00 00 00 00 00 00  |................|
 * 00000030  00 00 00 00 40 00 00 00  00 00 40 00 0e 00 0d 00  |....@.....@.....|
 * 00000040  55 48 89 e5 8b 05 00 00  00 00 89 c7 e8 00 00 00  |UH..............|
 * 00000050  00 5d c3 00 64 00 00 00  55 48 89 e5 89 7d fc b8  |.]..d...UH...}..|
 * 00000060  01 00 00 00 5d c3 00 47  43 43 3a 20 28 47 4e 55  |....]..GCC: (GNU|
 * 00000070  29 20 34 2e 38 2e 35 20  32 30 31 35 30 36 32 33  |) 4.8.5 20150623|
 * 00000080  20 28 52 65 64 20 48 61  74 20 34 2e 38 2e 35 2d  | (Red Hat 4.8.5-|
 * 00000090  32 38 29 00 00 00 00 00  14 00 00 00 00 00 00 00  |28).............|
 * 000000a0  01 7a 52 00 01 78 10 01  1b 0c 07 08 90 01 00 00  |.zR..x..........|
 * 000000b0  1c 00 00 00 1c 00 00 00  00 00 00 00 0e 00 00 00  |................|
 * 000000c0  00 41 0e 10 86 02 43 0d  06 49 0c 07 08 00 00 00  |.A....C..I......|
 * 000000d0  1c 00 00 00 3c 00 00 00  00 00 00 00 13 00 00 00  |....<...........|
 * 000000e0  00 41 0e 10 86 02 43 0d  06 4e 0c 07 08 00 00 00  |.A....C..N......|
 * objdump -d  -j BAR attribute.o  
 * attribute.o:     file format elf64-x86-64
 * Disassembly of section BAR:
 * 0000000000000000 <function>:
 *   0:   55                      push   %rbp
 *   1:   48 89 e5                mov    %rsp,%rbp
 *   4:   89 7d fc                mov    %edi,-0x4(%rbp)
 *   7:   b8 01 00 00 00          mov    $0x1,%eax
 *   c:   5d                      pop    %rbp
 *   d:   c3                      retq   
 */
