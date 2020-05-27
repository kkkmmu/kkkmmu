The terminal driver (part of the operating system) establishes the relationship between special characters and signals. Your terminal settings, e.g., using stty, are what it uses to decide what (if anything) to do with characters that you type. You can reassign those special characters as needed with a few caveats:

only one special character per function
only single-byte characters are used
controlC and controlD are conventional: while a few applications may hardcode these values, the terminal driver does not require that.

The terminal driver is software, not part of your terminal. For some keyboards you may find different assignments of special characters more convenient than others (and for different operating systems, a few choices of the default values for the special characters may differ).

Further reading:

11.Special Characters (POSIX *General Terminal Interface)
13.General Terminal Interface
13.stty - set the options for a terminal
