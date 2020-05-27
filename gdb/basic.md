1. set print pretty on
2. print var
3. print *var
4. info locals
5. info args
6. info frame
7. info registers
8. info all-registers
9. ptype var
10. p &((struct foo *)0)->member 获取字段偏移量
11. gdb --pid=28498 -ex "generate-core-file /coredump/nsm.core" -ex "detach" -ex "quit" 用这个命令来生成当前正在运行进程的corefile
