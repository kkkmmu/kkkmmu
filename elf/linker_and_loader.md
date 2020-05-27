The basic job of any linker or loader is simple: 
	it binds more abstract names to more concrete names, which permits programmers to write code using the more abstract names. That is, it takes a name written by a programmer such as getline and binds it to "the location 612 bytes from the beginning of the executable code in module iosys." Or it may take a more abstract numeric address such as "the location 450 bytes beyond the beginning of the static data for this module" and bind it to a numeric address.

Program loading: 
	Copy a program from secondary storage (which since about 1968 invariably means a disk) into main memory so it’s ready to run. In some cases loading just involves copying the data from disk to memory, in others it involves allocating storage, setting protection bits, or arranging for virtual memory to map virtual addresses to disk pages.

Relocation: 
	Compilers and assemblers generally create each file of object code with the program addresses starting at zero, but few computers let you load your program at location zero. If a program is created from multiple subprograms, all the subprograms have to be loaded at non-overlapping addresses. Relocation is the process of assigning load addresses to the various parts of the program, adjusting the code and data in the program to reflect the assigned addresses. In many systems, relocation happens more than once. It’s quite common for a linker to create a program from multiple subprograms, and create one linked output program that starts at zero, with the various subprograms relocated to locations within the big program. Then when the program is loaded, the system picks the actual load address and the linked program is relocated as a whole to the load address.

Symbol resolution: 
	When a program is built from multiple subprograms, the references from one subprogram to another are made using symbols; a main program might use a square root routine called sqrt, and the math library defines sqrt. A linker resolves the symbol by noting the location assigned to sqrt in the library, and patching the caller’s object code to so the call instruction refers to that location.
