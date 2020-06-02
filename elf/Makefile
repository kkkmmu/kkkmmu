DIR := .

SRC := $(wildcard $(DIR)/*.c)
INC := $(wildcard $(DIR)/*.h)
#OBJ := $(SRC:$(DIR)/%.c=$(DIR)/%.o)
OBJ := $(patsubst $(DIR)/%.c, $(DIR)/%.o, $(SRC))

EXE := $(patsubst $(DIR)/%.c, $(DIR)/%, $(SRC))

CFLAGS := -Wall
LDLIBS := -lm -lpthread -lreadline -lutil
LDFLAGS :=

all: $(EXE)

$(EXE):$(OBJ)
	$(CC) $(LDFLAGS) $@.o $(LDLIBS) -o $@

$(DIR)/%.o: $(DIR)/%.c
	$(CC) $(CFLAGS) -c $^ 

clean:
	rm -rf $(EXE) $(OBJ)

.PHONY: all clean $(EXE)
