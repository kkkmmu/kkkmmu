  File Handling - Write
  Building upon the previous lesson we will now use sys_write to write content to a newly created file.

  sys_write expects 3 arguments - the number of bytes to write in EDX, the contents string to write in ECX and the file descriptor in EBX. The sys_write opcode is then loaded into EAX and the kernel is called to write the content to the file. In this lesson we will first call sys_creat to get a file descriptor which we will then load into EBX.
