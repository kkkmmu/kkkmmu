File Hanndling - Close
Building upon the previous lesson we will now use sys_close to properly close an open file.

sys_close expects 1 argument - the file descriptor in EBX. We will use the previous lessons code to obtain the file descriptor which we will then load into EBX. The sys_close opcode is then loaded into EAX and the kernel is called to close the file and remove the active file descriptor.
