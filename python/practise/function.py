import os
import io
import glob

def fn_test(msg):
	print(msg + " is printed")
	print(msg == True)
	print(os.path)
	print(os.name)
	print(os.curdir)

def read_file(name):
	file = io.open(name)
	print(file.name)
	lines = file.readlines()
	for line in lines:
		print(line)
	file.close()

def diction():
	a = {"1" : "One", "2" : "Two"}
	print(a["1"])
	print(a["2"])
	print("2" in a)

	for k in a:
		print (k)

	print({f: os.stat(f) for f in glob.glob('../*.py')})

def glob_func(path):
	fs = glob.glob(path+"/*.py")
	print(fs)
	print([os.path.realpath(f) for f in glob.glob(path + "/*.py")])
	print([(os.stat(f).st_size, os.path.realpath(f)) for f in glob.glob(path + "/*.py")])

def generator(data):
	for index in range(len(data) -1, -1, -1):
		yield data[index]

class First:
	"""This is first class in python"""
	i = 12345

	def get_i(self):
		return  self.i

class Second:
	"""That is the second class in python"""

	def __init__(self):
		self.x = 100
		self.y = "hello from second"
		self.z = [x ** x for x in range(1,10)]

class Third:
	"""It is the third class in python"""
	def __init__(self, p1=None, p2=[x *x for x in range(1,5)], p3=dict({"Hello" : "World"}), p4=100):
		self.u = p1
		self.v = p2
		self.w = p3
		self.x = p4

	def get_u(self):
		return self.u
	def set_u(self, u):
		self.u = u

class Iter:
	"""Iteratoer"""
	def __init__(self, data):
		self.data = data
		self.index = len(data)
	def __iter__(self):
		return self

	def __next__(self):
		if self.index == 0:
			raise StopIteration
		self.index = self.index -1
		return self.data[self.index]

if __name__ == "__main__":
	fn_test("hello")
	read_file("function.py")
	diction()
	glob_func("..")

	f = First()
	print(f.i)
	print(f.get_i)
	print(f.get_i())
	print(f.__doc__)

	s = Second()
	print(s)
	print(s.x)
	print(s.y)
	print(s.z)

	t1 = Third()
	print(t1.u)
	print(t1.v)
	print(t1.w)
	print(t1.x)

	t = Third("test", "world", 1, '1')
	print(t.u)
	print(t.v)
	print(t.w)
	print(t.x)

	print(t.get_u())
	print(t.set_u(1000))
	print(t.get_u())
	t.y = 10000
	print(t.y)
	print(isinstance(t1, Third))
	print(isinstance(t, Third))

	try:
		i = Iter([x * x for x in range (1,10)])
		it = iter(i)
		print(next(it))
		print(next(it))
		print(next(it))
		print(next(it))
		print(next(it))
		print(next(it))
		print(next(it))
		print(next(it))
		print(next(it))
		print(next(it))
		print(next(it))
		print(next(it))
		print(next(it))
	except StopIteration:
		pass
	i2 = Iter([x * x for x in range (1,10)])
	for ii in iter(i2):
		print(ii)


	for char in generator("hello world"):
		print(char)

