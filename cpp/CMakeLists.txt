CMAKE_MINIMUM_REQUIRED(VERSION 3.5.1)

PROJECT (cpp)
INCLUDE_DIRECTORIES(include)
#ADD_SUBDIRECTORY(src)

SET(CMAKE_CXX_STANDARD 17)
SET(CMAKE_CXX_STANDARD_REQUIRED True)

ADD_EXECUTABLE(cpp_17 cpp_17.cpp)
