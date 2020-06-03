Telling the time
Generating a unix timestamp in NASM is easy with the sys_time function of the linux kernel. Simply pass OPCODE 13 to the kernel with no arguments and you are returned the Unix Epoch in the EAX register.

That is the number of seconds that have elapsed since January 1st 1970 UTC
