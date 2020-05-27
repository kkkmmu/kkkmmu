1. 	gcc -c m.c
	gcc -c a.c
	readelf -S m.o
	readelf -S a.o
	readelf -a m.o
	readelf -a a.o
	objdump -d m.o
	objdump -d a.o

	gcc -g -c m.c
	gcc -g -c a.c
	objdump -S m.o
	objdump -S a.o
	gcc -g a.c m.c -o m
	readelf -a m
	objdump -d m
	objdump -S m
	readelf -S m
	hexdump -C m
	nm m
	size m
	objdump -s -j .text m
	objdump -s m
	objdump -p m
	objdump -x -s -d m
