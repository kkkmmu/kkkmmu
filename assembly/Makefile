DIR := .

SRC := $(wildcard $(DIR)/*.asm)
OBJ := $(patsubst $(DIR)/%.asm, $(DIR)/%.o, $(SRC))

EXE := $(patsubst $(DIR)/%.asm, $(DIR)/%, $(SRC))

CFLAGS := -Wall
LDLIBS := -lm -lpthread -lreadline -lutil
LDFLAGS :=

all: $(EXE)

$(EXE):$(OBJ)
	$(LD) -m elf_x86_64 -o $@ $@.o

$(DIR)/%.o: $(DIR)/%.asm
	nasm -f elf64 -o $@ $^ 

clean:
	rm -rf $(EXE) $(OBJ)

.PHONY: all clean $(EXE)
