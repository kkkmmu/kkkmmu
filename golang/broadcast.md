What you are doing is a fan out pattern, that is to say, multiple endpoints are listening to a single input source. The result of this pattern is, only one of these listeners will be able to get the message whenever there's a message in the input source. The only exception is a close of channel. This close will be recognized by all of the listeners, and thus a "broadcast".

But what you want to do is broadcasting a message read from connection, so we could do something like this:

When the number of listeners is known
Let each worker listen to dedicated broadcast channel, and dispatch the message from the main channel to each dedicated broadcast channel.

type worker struct {
	source chan interface{}
	quit chan struct{}
}

func (w *worker) Start() {
	w.source = make(chan interface{}, 10) // some buffer size to avoid blocking
	go func() {
		for {
			select {
				case msg := <-w.source
				// do something with msg
				case <-quit: // will explain this in the last section
				return
			}
		}
	}()
}
And then we could have a bunch of workers:

workers := []*worker{&worker{}, &worker{}}
for _, worker := range workers { worker.Start() }
Then start our listener:

go func() {
	for {
		conn, _ := listener.Accept()
		ch <- conn
	}
}()
And a dispatcher:

go func() {
	for {
		msg := <- ch
		for _, worker := workers {
			worker.source <- msg
		}
	}
}()
When the number of listeners is not known
In this case, the solution given above still works. The only difference is, whenever you need a new worker, you need to create a new worker, start it up, and then push it into workers slice. But this method requires a thread-safe slice, which need a lock around it. One of the implementation may look like as follows:

type threadSafeSlice struct {
	sync.Mutex
	workers []*worker
}

func (slice *threadSafeSlice) Push(w *worker) {
	slice.Lock()
	defer slice.Unlock()

	workers = append(workers, w)
}

func (slice *threadSafeSlice) Iter(routine func(*worker)) {
	slice.Lock()
	defer slice.Unlock()

	for _, worker := range workers {
		routine(worker)
	}
}
Whenever you want to start a worker:

w := &worker{}
w.Start()
threadSafeSlice.Push(w)
And your dispatcher will be changed to:

go func() {
	for {
		msg := <- ch
		threadSafeSlice.Iter(func(w *worker) { w.source <- msg })
	}
}()
Last words: never leave a dangling goroutine
One of the good practices is: never leave a dangling goroutine. So when you finished listening, you need to close all of the goroutines you fired. This will be done via quit channel in worker:

First we need to create a global quit signalling channel:

globalQuit := make(chan struct{})
And whenever we create a worker, we assign the globalQuit channel to it as its quit signal:

worker.quit = globalQuit
Then when we want to shutdown all workers, we simply do:

close(globalQuit)
Since close will be recognized by all listening goroutines (this is the point you understood), all goroutines will be returned. Remember to close your dispatcher routine as well, but I will leave it to you :)
