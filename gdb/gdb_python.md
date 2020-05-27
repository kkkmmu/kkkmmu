1. gdb & source gdb_PY_EXTENSION.py
2. 在gdb模式下输入pi可以切换到python解释器模式.
3. 在gdb python 解释器模式后import gdb.
4. print(gdb.PYTHONDIR);
5. gdb.execute("set pagination off")
6. gdb.execute("set print pretty")
7. gdb.execute("target remote <ip>:<port>")
8. gdb.execute('set logging file %s' % (LOG_FILE))
9. gdb.execute("set logging on")
10. gdb.execute("show script-extension")
11. gdb.execute("show sysroot")
12. gdb.execute("show solib-search-path")
13. gdb.execute("info local")
14. gdb.execute("info threads")
15. gdb.execute("info args")
16. gdb.execute("bt")
17.
	def stop_handler(event):
		    print('EVENT: %s' % (event))
			    gdb.execute("info local")
			       gdb.execute("info threads")
			       gdb.execute("info args")
			       gdb.execute("bt")
			       if event.breakpoint.location == "func3":
			           gdb.write("Special bp hit\n")
			           gdb.execute("p func3_var")
			           gdb.execute("p *func3_ptr")
			       # don't stop, continue
			       gdb.execute("c")

18. gdb.events.stop.connect(stop_handler)
10. gdb.events.stop.disconnect(stop_handler)
11. def exit_handler (event):
	    print "event type: exit"
	    print "exit code: %d" % (event.exit_code)

   gdb.events.exited.connect (exit_handler)
12. GDB provides a general event facility so that Python code can be notified of various state changes, particularly changes that occur in the inferior.
	An event is just an object that describes some state change. The type of the object and its attributes will vary depending on the details of the change. All the existing events are described below.
	In order to be notified of an event, you must register an event handler with an event registry. An event registry is an object in the gdb.events module which dispatches particular events. A registry provides methods to register and unregister event handlers:
	Function: EventRegistry.connect (object)
		Add the given callable object to the registry. This object will be called when an event corresponding to this registry occurs.
	Function: EventRegistry.disconnect (object)
		Remove the given object from the registry. Once removed, the object will no longer receive notifications of events.
