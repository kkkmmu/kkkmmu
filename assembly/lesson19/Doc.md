Exccute Command
Firstly, some background
The EXEC family of functions replace the currently running process with a new process, that executes the command you specified when calling it. We will be using the sys_execve function in this lesson to replace our program's running process with a new process that will execute the linux program /bin/echo to print out “Hello World!”.

Naming convention
The naming convention used for this family of functions is exec (execute) followed by one or more of the following letters.

E - An array of pointers to environment variables is explicitly passed to the new process image.
L - Command-line arguments are passed individually to the function.
P - Uses the PATH environment variable to find the file named in the path argument to be executed.
V - Command-line arguments are passed to the function as an array of pointers.
Writing our program
The V & E at the end of our function name means we will need to pass our arguments in the following format: The first argument is a string containing the command to execute, then an array of arguments to pass to that command and then another array of environment variables that the new process will use. As we are calling a simple command we won't pass any special environment variables to the new process and instead will pass 0h (null).

Both the command arguments and the environment arguments need to be passed as an array of pointers (addresses to memory). That's why we define our strings first and then define a simple null-terminated struct (array) of the variables names. This is then passed to sys_execve. We call the function and the process is replaced by our command and output is returned to the terminal.
