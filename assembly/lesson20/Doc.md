Process Forking
Firstly, some background
In this lesson we will use sys_fork to create a new process that duplicates our current process. sys_fork takes no arguments - you just call fork and the new process is created. Both processes run concurrently. We can test the return value (in EAX) to test whether we are currently in the parent or child process. The parent process returns a non-negative, non-zero integer. In the child process EAX is zero. This can be used to branch your logic between the parent and child.

In our program we exploit this fact to print out different messages in each process.

Note: Each process is responsible for safely exiting.
