File Handling - Delete
Deleting a file on linux is achieved by calling sys_unlink.

sys_unlink expects 1 argument - the filename in EBX. The sys_unlink opcode is then loaded into EAX and the kernel is called to delete the file.

Note: A file 'readme.txt' has been included in the code folder for this lesson. This file will be deleted after running the program.
