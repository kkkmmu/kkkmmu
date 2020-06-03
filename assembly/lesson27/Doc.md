File Handling - Seek
In this lesson we will open a file and update the file contents at the end of the file using sys_lseek.

Using sys_lseek you can move the cursor within the file by an offset in bytes. The below example will move the cursor to the end of the file, then pass 0 bytes as the offset (so we append to the end of the file and not beyond) before writing a string in that position. Try different values in ECX and EDX to write the content to different positions within the opened file.

sys_lseek expects 3 arguments - the whence argument (table below) in EDX, the offset in bytes in ECX, and the file descriptor in EBX. The sys_lseek opcode is then loaded into EAX and we call the kernel to move the file pointer to the correct offset. We then use sys_write to update the content at that position.

Description	Value
SEEK_SET	beginning of the file	0
SEEK_CUR	current file offset	1
SEEK_END	end of the file	2
