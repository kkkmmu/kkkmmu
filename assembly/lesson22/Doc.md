File Handling - Create
Firstly, some background
File Handling in Linux is achieved through a small number of system calls related to creating, updating and deleting files. These functions require a file descriptor which is a unique, non-negative integer that identifies the file on the system.

Writing our program
We begin the tutorial by creating a file using sys_creat. We will then build upon our program in each of the following file handling lessons, adding code as we go. Eventually we will have a full program that can create, update, open, close and delete files.

sys_creat expects 2 arguments - the file permissions in ECX and the filename in EBX. The sys_creat opcode is then loaded into EAX and the kernel is called to create the file. The file descriptor of the created file is returned in EAX. This file descriptor can then be used for all other file handling functions.
