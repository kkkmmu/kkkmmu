# CMAKE generated file: DO NOT EDIT!
# Generated by "Unix Makefiles" Generator, CMake Version 3.8

# Delete rule output on recipe failure.
.DELETE_ON_ERROR:


#=============================================================================
# Special targets provided by cmake.

# Disable implicit rules so canonical targets will work.
.SUFFIXES:


# Remove some rules from gmake that .SUFFIXES does not remove.
SUFFIXES =

.SUFFIXES: .hpux_make_needs_suffix_list


# Produce verbose output by default.
VERBOSE = 1

# Suppress display of executed commands.
$(VERBOSE).SILENT:


# A target that is always out of date.
cmake_force:

.PHONY : cmake_force

#=============================================================================
# Set environment variables for the build.

# The shell in which to execute make rules.
SHELL = /bin/sh

# The CMake executable.
CMAKE_COMMAND = /usr/bin/cmake

# The command to remove a file.
RM = /usr/bin/cmake -E remove -f

# Escaping for special characters.
EQUALS = =

# The top-level source directory on which CMake was run.
CMAKE_SOURCE_DIR = /home/leeweop/kkkmmu/pppp

# The top-level build directory on which CMake was run.
CMAKE_BINARY_DIR = /home/leeweop/kkkmmu/pppp

# Include any dependencies generated for this target.
include .lib/llog/src/CMakeFiles/llog.dir/depend.make

# Include the progress variables for this target.
include .lib/llog/src/CMakeFiles/llog.dir/progress.make

# Include the compile flags for this target's objects.
include .lib/llog/src/CMakeFiles/llog.dir/flags.make

.lib/llog/src/CMakeFiles/llog.dir/llog.c.o: .lib/llog/src/CMakeFiles/llog.dir/flags.make
.lib/llog/src/CMakeFiles/llog.dir/llog.c.o: lib/llog/src/llog.c
	@$(CMAKE_COMMAND) -E cmake_echo_color --switch=$(COLOR) --green --progress-dir=/home/leeweop/kkkmmu/pppp/CMakeFiles --progress-num=$(CMAKE_PROGRESS_1) "Building C object .lib/llog/src/CMakeFiles/llog.dir/llog.c.o"
	cd /home/leeweop/kkkmmu/pppp/.lib/llog/src && /usr/bin/cc $(C_DEFINES) $(C_INCLUDES) $(C_FLAGS) -o CMakeFiles/llog.dir/llog.c.o   -c /home/leeweop/kkkmmu/pppp/lib/llog/src/llog.c

.lib/llog/src/CMakeFiles/llog.dir/llog.c.i: cmake_force
	@$(CMAKE_COMMAND) -E cmake_echo_color --switch=$(COLOR) --green "Preprocessing C source to CMakeFiles/llog.dir/llog.c.i"
	cd /home/leeweop/kkkmmu/pppp/.lib/llog/src && /usr/bin/cc $(C_DEFINES) $(C_INCLUDES) $(C_FLAGS) -E /home/leeweop/kkkmmu/pppp/lib/llog/src/llog.c > CMakeFiles/llog.dir/llog.c.i

.lib/llog/src/CMakeFiles/llog.dir/llog.c.s: cmake_force
	@$(CMAKE_COMMAND) -E cmake_echo_color --switch=$(COLOR) --green "Compiling C source to assembly CMakeFiles/llog.dir/llog.c.s"
	cd /home/leeweop/kkkmmu/pppp/.lib/llog/src && /usr/bin/cc $(C_DEFINES) $(C_INCLUDES) $(C_FLAGS) -S /home/leeweop/kkkmmu/pppp/lib/llog/src/llog.c -o CMakeFiles/llog.dir/llog.c.s

.lib/llog/src/CMakeFiles/llog.dir/llog.c.o.requires:

.PHONY : .lib/llog/src/CMakeFiles/llog.dir/llog.c.o.requires

.lib/llog/src/CMakeFiles/llog.dir/llog.c.o.provides: .lib/llog/src/CMakeFiles/llog.dir/llog.c.o.requires
	$(MAKE) -f .lib/llog/src/CMakeFiles/llog.dir/build.make .lib/llog/src/CMakeFiles/llog.dir/llog.c.o.provides.build
.PHONY : .lib/llog/src/CMakeFiles/llog.dir/llog.c.o.provides

.lib/llog/src/CMakeFiles/llog.dir/llog.c.o.provides.build: .lib/llog/src/CMakeFiles/llog.dir/llog.c.o


# Object files for target llog
llog_OBJECTS = \
"CMakeFiles/llog.dir/llog.c.o"

# External object files for target llog
llog_EXTERNAL_OBJECTS =

.lib/llog/lib/libllog.so.1.1: .lib/llog/src/CMakeFiles/llog.dir/llog.c.o
.lib/llog/lib/libllog.so.1.1: .lib/llog/src/CMakeFiles/llog.dir/build.make
.lib/llog/lib/libllog.so.1.1: .lib/llog/src/CMakeFiles/llog.dir/link.txt
	@$(CMAKE_COMMAND) -E cmake_echo_color --switch=$(COLOR) --green --bold --progress-dir=/home/leeweop/kkkmmu/pppp/CMakeFiles --progress-num=$(CMAKE_PROGRESS_2) "Linking C shared library ../lib/libllog.so"
	cd /home/leeweop/kkkmmu/pppp/.lib/llog/src && $(CMAKE_COMMAND) -E cmake_link_script CMakeFiles/llog.dir/link.txt --verbose=$(VERBOSE)
	cd /home/leeweop/kkkmmu/pppp/.lib/llog/src && $(CMAKE_COMMAND) -E cmake_symlink_library ../lib/libllog.so.1.1 ../lib/libllog.so.1 ../lib/libllog.so

.lib/llog/lib/libllog.so.1: .lib/llog/lib/libllog.so.1.1
	@$(CMAKE_COMMAND) -E touch_nocreate .lib/llog/lib/libllog.so.1

.lib/llog/lib/libllog.so: .lib/llog/lib/libllog.so.1.1
	@$(CMAKE_COMMAND) -E touch_nocreate .lib/llog/lib/libllog.so

# Rule to build all files generated by this target.
.lib/llog/src/CMakeFiles/llog.dir/build: .lib/llog/lib/libllog.so

.PHONY : .lib/llog/src/CMakeFiles/llog.dir/build

.lib/llog/src/CMakeFiles/llog.dir/requires: .lib/llog/src/CMakeFiles/llog.dir/llog.c.o.requires

.PHONY : .lib/llog/src/CMakeFiles/llog.dir/requires

.lib/llog/src/CMakeFiles/llog.dir/clean:
	cd /home/leeweop/kkkmmu/pppp/.lib/llog/src && $(CMAKE_COMMAND) -P CMakeFiles/llog.dir/cmake_clean.cmake
.PHONY : .lib/llog/src/CMakeFiles/llog.dir/clean

.lib/llog/src/CMakeFiles/llog.dir/depend:
	cd /home/leeweop/kkkmmu/pppp && $(CMAKE_COMMAND) -E cmake_depends "Unix Makefiles" /home/leeweop/kkkmmu/pppp /home/leeweop/kkkmmu/pppp/lib/llog/src /home/leeweop/kkkmmu/pppp /home/leeweop/kkkmmu/pppp/.lib/llog/src /home/leeweop/kkkmmu/pppp/.lib/llog/src/CMakeFiles/llog.dir/DependInfo.cmake --color=$(COLOR)
.PHONY : .lib/llog/src/CMakeFiles/llog.dir/depend

