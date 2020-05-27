1. View all data/code in every section of an ELF file:
   objdump -D <elf_object>

2. View only program code in an ELF file:
   objdump -d <elf_object>

3. View all symbols:
   objdump -tT <elf_object>

4.To copy the .data section from an ELF object to a file, use this line:
   objcopy –only-section=.data <infile> <outfile>

5.This is the strace command used to trace a basic program:
   strace /bin/ls -o ls.out

6. The strace command used to attach to an existing process is as follows:
   strace -p <pid> -o daemon.out

7. The initial output will show you the file descriptor number of each system call that
   takes a file descriptor as an argument, such as this:
   SYS_read(3, buf, sizeof(buf));
  
8. If you want to see all of the data that was being read into file descriptor 3, you can run the following command:
   strace -e read=3 /bin/ls
   You may also use -e write=fd to see written data. 

9. To retrieve a section header table:
   readelf -S <object>

10. To retrieve a program header table:
   readelf -l <object>

11. To retrieve a symbol table:
   readelf -s <object>

12. To retrieve the ELF file header data:
   readelf -e <object>

13. To retrieve relocation entries:
   readelf -r <object>

14. To retrieve a dynamic segment:
   readelf -d <object>
 
15. /proc/<pid>/maps
   /proc/<pid>/maps fle contains the layout of a process image by showing each memory mapping. This includes the executable, shared libraries, stack, heap, VDSO,
 and more. This fle is critical for being able to quickly parse the layout of a process address space.

16. The LD_PRELOAD environment variable
   The LD_PRELOAD environment variable can be set to specify a library path that should be dynamically linked before any other libraries. This has the effect of allowing functions and symbols from the preloaded library to override the ones from the other libraries that are linked afterwards. This essentially allows you to perform runtime patching by redirecting shared library functions. As we will see in later chapters, this technique can be used to bypass anti-debugging code and for userland rootkits.

17. The LD_SHOW_AUXV environment variable This environment variable tells the program loader to display the program's auxiliary vector during runtime. The auxiliary vector is information that is placed on the program's stack (by the kernel's ELF loading routine), with information that is passed to the dynamic linker with certain information about the program. 

18. Linker scripts Linker scripts are a point of interest to us because they are interpreted by the linker and help shape a program's layout with regard to sections, memory, and symbols. The default linker script can be viewed with ld -verbose.

19. ELF fle types
   An ELF fle may be marked as one of the following types:
      • ET_NONE: This is an unknown type. It indicates that the file type is unknown, or has not yet been defined.
	  • ET_REL: This is a relocatable file. ELF type relocatable means that the file is marked as a relocatable piece of code or sometimes called an object file.
	            Relocatable object files are generally pieces of Position independent code (PIC) that have not yet been linked into an executable. You will often see
	            .o files in a compiled code base. These are the files that hold code and data suitable for creating an executable file.
	  • ET_EXEC: This is an executable file. ELF type executable means that the file is marked as an executable file. These types of files are also called programs
	            and are the entry point of how a process begins running.
	  • ET_DYN: This is a shared object. ELF type dynamic means that the file is marked as a dynamically linkable object file, also known as shared libraries.
             	These shared libraries are loaded and linked into a program's process image at runtime.
	  • ET_CORE: This is an ELF type core that marks a core file. A core file is a dump of a full process image during the time of a program crash or when the
	            process has delivered an SIGSEGV signal (segmentation violation). GDB can read these files and aid in debugging to determine what caused the program
             	to crash.

20. man elf

21. ELF program headers 
	ELF program headers are what describe segments within a binary and are necessary for program loading. Segments are understood by the kernel during load time and
describe the memory layout of an executable on disk and how it should translate to memory. The program header table can be accessed by referencing the offset found
in the initial ELF header member called e_phoff (program header table offset), as shown in the ElfN_Ehdr structure.
     Program headers describe the segments of an executable file (shared libraries included) and what type of segment it is (that is, what type of data or code it is reserved for). 
	 Program headers are primarily there to describe the layout of a program for when it is executing and in memory. 


