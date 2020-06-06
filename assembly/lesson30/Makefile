DIR := .

SRC := $(wildcard $(DIR)/*.asm)
OBJ := $(patsubst $(DIR)/%.asm, $(DIR)/%.o, $(SRC))

EXE := $(patsubst $(DIR)/%.asm, $(DIR)/%, $(SRC))

CFLAGS := -Wall
LDLIBS := -lm -lpthread -lreadline -lutil
LDFLAGS :=

all: $(EXE)

$(EXE):$(OBJ)
	$(LD) -m elf_i386 -o $@ $@.o

$(DIR)/%.o: $(DIR)/%.asm
	nasm -f elf -o $@ $^ 

clean:
	rm -rf $(EXE) $(OBJ) *.txt

.PHONY: all clean $(EXE)
