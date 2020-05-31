CMAKE_MINIMUM_REQUIRED(VERSION 3.16.3) 

PROJECT(awsome_program)

message("root='${CMAKE_SOURCE_DIR}'")

SET(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -std=c++2a -Wall")

#FIND_PACKAGE(OpenCV REQUIRED)

#ADD_SUBDIRECTORY(.)
SET(SRC useful_functions.cpp utilies.cpp stuf_that_should_be_in_opencv.cpp main.cpp)
foreach(src ${SRC})
	message("${src}\n")
endforeach(src)

ADD_EXECUTABLE(awsome_program ${SRC})
#TARGET_LINK_LIBRARIES(awsome_program /usr/lib/opencv_core)
