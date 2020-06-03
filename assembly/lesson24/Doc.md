ndling - Open
Building upon the previous lesson we will now use sys_open to obtain the file descriptor of the newly created file. This file descriptor can then be used for all other file handling functions.

sys_open expects 2 arguments - the access mode (table below) in ECX and the filename in EBX. The sys_open opcode is then loaded into EAX and the kernel is called to open the file and return the file descriptor.

sys_open additionally accepts zero or more file creation flags and file status flags in EDX. Click here for more information about the access mode, file creation flags and file status flags.

Description	Value
O_RDONLY	open file in read only mode	0
O_WRONLY	open file in write only mode	1
O_RDWR	open file in read and write mode	2
Note: sys_open returns the file descriptor in EAX. On linux this will be a unique, non-negative integer which we will print using our integer printing function.
