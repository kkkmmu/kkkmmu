//c里面是通过malloc/free手动管理堆内存的.而Rust则有多种方式，其中最常用的一种就是Box，通过Box::new可以在堆上申请一块内存.不像C里面一样堆上空间需要调用free来释放，rust中是在编译期编译器借助lifetime对内存生命期进行分析，在生命期结束时自动插入free。当前Rust底层即Box背后是调用jemalloc来做内存管理的,所以堆上空间是不需要手动释放的.
//大多数带GC的面向对象语言里的对象都是借助Box来实现的，比如常见的动态语言python/ruby/javascript等，其宣称的一切皆对象里面所谓的对象基本上就是Boxed
//Value。
//
//Boxed值相对于Unboxed值，内存占用空间会大些，同时访问值的时候也需要先进行Unbox，即对指针进行解引用再获取真正存储的值，所以内存访问开销也会大些。既然Boxed值既浪费空间又浪费时间，为什么还要这么做呢？因为通过Box，所有对象看起来就像是以相同大小存储的，因为只需要存储一个指针就够了，应用程序可以同等看待各种值，而不用去管实际存储是多大的值，如果申请和释放相应的资源。
//

//通过Box::new()申请一块堆内存，并返回该内催的地址。
//
//
//Rc和Arc
//  Rust建立在所有权之上的这一套机制，它要求一个资源同一时刻有且只能有一个拥有该资源的绑定或者&mut
//  引用，这在大部分情况下保证了内存安全。但是这样的设计是相当严格的，在另外一些情况下，它限制了程序的书写，无法实现某些功能。因此Rust在std扩中提供了额外的措施来补充所有权机制，以应对更广泛的场景。
//  默认情况下Rust中，对一个资源，同一时刻，有且只有一个所有权拥有者。Rc和Arc使用引用计数的方法，让程序在同一时刻，实现同一资源的多个所有权拥有者，多个拥有者共享该资源。
//
//  Rc: 
//      Rc用于同一个线程内部，通过use::std::rc::Rc来引入，它有以下特点：
//          1. 用Rc包装起来的类型对象，是immutable的，即不可变的。即你无法修改Rc<T>中的T对象。
//          2. 一旦最后一个拥有者消失，则该资源会被自动回收，这个生命周期是在编译期就确定下来的。
//          3. Rc只能用于同一个线程内部，不能用于线程之间的共享对象（不能跨线程传递）.
//          4. Rc实际上是一个指针，它不影响包裹对象的方法调用形式。
//
//      fn main() {
//          use std::rc::Rc;
//
//          let five = Rc::new(5);
//          let five2 = five.clone();
//          let five3 = five.clone();
//      }
//
//  Rc Weak:
//      Weak通过 use std::rc::Weak来引入。
//      Rc是一个引用计数指针，而Weak是一个指针，但不能增加引用计数，是Rc的Weak版本。它有以下几个特点：
//          1. 可访问，但不拥有，不增加引用计数，因此，不会对资源回收管理造成影响。
//          2. 可由Rc<T>调用downgrade方法而转换成Weak<T>;
//          3. Weak<T>可以使用upgrade方法转换成Option<Rc<T>>，如果资源已被释放，则Option值为None；
//          4. 常用于解决循环引用的问题。
//
//
//  Arc:
//      Arc是原子引用计数，是Rc的多线程版本。Arc通过sdt::sync::Arc引入。
//      它的特点：
//          1. Arc可以跨线程传递，用于跨线程共享一个对象；
//          2. 用Arc包裹起来的类型对象，对可变性没有要求。
//          3. 一旦最后一个拥有者消失，则资源会被自动回收，这个生命周期是在编译期就确定下来的。
//          4. Arc实际上是一个指针，它不影响包裹对象的方法调用方式（即不存在先解包再调用这一说）.
//          5. Arc 对于多线程的共享状态几乎是必须的。
//
//  Arc Weak:
//      与Rc Weak类似，只是是多线程版本的