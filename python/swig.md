在Python代码中调用C/C++代码
2018-01-22 16:35:09 +08  字数：1675  标签： Python C
在Python代码中调用C/C++代码，这需求是比较常见的。 毕竟，当代软件世界的基础设施领域是C语言的天下，很多开发工作不可避免地要与它打交道。 而C++则是家族的嫡长子，也占据了一些不可替代的位置。

从调用的角度细分，还有两种方式： 一是单向调用，只从Python中调用C/C++代码；二是双向调用，C/C++也会回调Python的内容。 此外，从C/C++中调用Python代码，也偶尔出现，因为Python的某些标准库功能非常实用，比如http.server。 大部分情况下，都是从Python单向调用C/C++。

从C/C++代码的编译情况，还能细分成两类： 一是已编译，直接使用C/C++的动态链接库；二是未编译，需要在setup.py中编译成可导入模块。 前者往往用在使用某些著名的库，比如libc.so； 而后者则是往往出现在使用C/C++来提升Python运行效率、或者需要与C/C++交换复杂的数据结构。

本文着重介绍ctypes与SWIG，也会提一下其它方案。

ctypes ¶
对于直接调用C语言的动态链接库，ctypes就是正解。 遗憾的是，它不支持C++。

>>> import ctypes
>>> libc = ctypes.CDLL('libc.so.6')
>>> libc.printf
<_FuncPtr object at 0x7ffba0284a70>
>>> libc.printf(b'Hello world!\n')
Hello world!
13
>>> libc.time()
1516097322
printf的显示中，Hello world!是打印输出，而13则是printf函数返回值，代表13个打印字符。 详细教程，参考《ctypes tutorial》。 （注意，以上代码直接在python交互环境执行。如果是使用IPython，那么printf的打印结果可能需要退出IPython时才能看到。）

setup.py中打包so ¶
如果不是libc.so这种已安装或易安装的库，而是比较冷门，或者私有的库，就需要和package一起打包。 比如，对于以下结构的一个项目，需要把libxxx.so作为数据打包。

.
├── pkg_name
│   ├── __init__.py
│   └── xxx_wrapper
│       ├── __init__.py
│       └── libxxx.so
└── setup.py
在setup.py中，需要添加以下配置。

setup(
    ...
    package_data={'pkg_name': ['xxx_wrapper/libxxx.so']},
)
SWIG ¶
SWIG是个帮助使用C或者C++编写的软件能与其它各种高级编程语言进行嵌入联接的开发工具。 SWIG能应用于各种不同类型的语言包括常用脚本编译语言例如Perl, PHP, Python, Tcl, Ruby and PHP。

写一个interface文件，即可自动生成一个C/C++与一个Python的包装文件。 让连接C/C++与Python的这个抽象层级，完全不用手写。

比如，下面的example.i文件，暴露了四个接口。

%module example

%{
#include "example.h"
%}

extern double My_variable;
extern int fact(int n);
extern int my_mod(int x, int y);
extern char *get_time();
上层的Python可以直接调用这四个接口，而无需考虑中间的转换过程。

import example

print('My_varaiable: %s' % example.cvar.My_variable)
print('fact(5): %s' % example.fact(5))
print('my_mod(7,3): %s' % example.my_mod(7,3))
print('get_time(): %s' % example.get_time())
SWIG与ctypes相比，虽然使用上麻烦一些，但接口更容易使用，也支持C++。 最重要的是，可以很方便地让C/C++代码在setup.py中编译，更容易跨平台。

from setuptools import Extension, setup

EXAMPLE_EXT = Extension(
    name='_example',
    sources=[
        'src/example/example.c',
        'src/example/example.i',
    ],
)

setup(
    ...
    ext_modules=[EXAMPLE_EXT],
)
setuptools中的Extension源于distutils，可以识别*.i文件并自动做出转换。 这样编译的Wheel文件，还自带架构与版本等编译信息，比如：example-0.1.0-cp36-cp36m-linux_x86_64.whl。

同类思想的产物还有SIP，用于PyQt，参考《Using SIP》。

当然，直接通过CPython的API来写这个包装层级的代码，也不是不可以。 有时也很省事，但总是不省心。

dl ¶
如果是老旧的项目，或者老旧的人，可能还会使用dl。 建议不要用了，还是老老实实听官方文档的指示吧。

Deprecated since version 2.6: The dl module has been removed in Python 3. Use the ctypes module instead.

Boost.Python ¶
再标志一个坑，那就是Boost.Python。 说它是坑，并不是说像dl那样不好用，而是因为太复杂，会消耗过多的开发与学习时间。

这里有另一篇笔记《编译运行Boost.Python的HelloWorld》，记录着跑通一个它的HelloWorld是一件多么不容易的事。

参考 ¶
更多的调用C/C++代码的方法，可以参考：

《Integrating Python With Other Languages - Python Wiki》
《Wrapping C/C++ for Python》
《Calling C/C++ from Python? - Stack Overflow》
《Building C and C++ Extensions》
相关笔记
《在setup.py中配置SWIG模块》
《解决setup.py编译C++代码的-Wstrict-prototypes警告》
《编译运行Boost.Python的HelloWorld》
《用SWIG向Python提供C++里STL的容器》
《编译运行SWIG的example代码样例》
