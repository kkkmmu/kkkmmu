1. man ascii
1. gcc -g X.c -o X
2. objdump -D X | grep -A20 main.:
2. objdump -S X | grep -A20 main.:
3. objdump -D -M intel X | grep -A20 main.:
3. objdump -S -M intel X | grep -A20 main.:
4. 对于一个Process来说只包含一个Stack, 栈底是固定的，栈顶(ESP)是动态变化的。该Process所调用的每一个函数都在栈上创建相应的栈帧（Frame). 每个函数调用的Frame的基地址(EBP)是上一个函数调用的栈顶。分配于栈上的变量都可以通过栈顶（ESP)指针或者帧底（EBP)加上偏移来寻址。注意栈地址是从高地址向低地址生长的，因此ESP -= x的操作相当于在栈上分配空间, ESP += X的操作相当于在栈上回收空间. 
4. EIP contains a memory address that points to an instruction.
5. 通过比较info r esp, info r ebp, 以及info r eip的值可看出, 程序实际运行的栈的位置和程序本身在内存中存储的位置是不同的. 也就是说代码段（.txt）与栈占用不同的内存空间。
6. 通过gdb: info files可以查看程序在内存中的具体分布。
   通过gdb: info proc mapping可以查看程序的内存映射

6. set disassembly-flavor intel

5. X86 Registers:
	EAX: Accumulator
	ECX: Counter
	EDX: Data
	EBX: Base
	ESP: Stack Pointer
	EBP: Base Pointer
	ESI: Source Index
	EDI: Destination Index
	EIP: Instruction Pointer
	EFLAGS: Flags
5. gdb X
5. set pagination off
5. layout asm
5. info files
5. info proc mappings
5. info proc stat
6. b main
6. run
6. info frame
6. disassemble /m main
6. disassemble /r main
6. disassemble main
6. bt
6. bt full
6. info locals
6. info variables
6. info arg
6. info stack
6. info registers
6. info registers rip
7. nexti
6. disassemble main
8. info registers
6. info registers rip
7. nexti
6. disassemble main
8. info registers
6. info registers rip
6. i r rip
6. x /i $rip
6. x /t $rip
6. x /24w $rip
6. x /24b $rip
6. x /24wx main
6. x /24bx main
6. disassemble /r main
6. x /32x main
6. i r esp ebp eip
6. x/5i $eip
6. x/16xw $esp
