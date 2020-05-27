list = [1, 2, 3, 4]

it = iter(list)
print(next(it))
print(next(it))
print(next(it))

it = iter(list)
for x in it:
    print(x)


class Fibs:
    def __init__(self):
        self.a = 0
        self.b = 1

    def __next__(self):
        self.a, self.b = self.b, self.a+self.b
        return self.a

    next = __next__

    def __iter__(self):
        return self


fibs = Fibs()

for f in fibs:
    if f > 1000:
        print(f)
        break


class TestIterator:
    value = 0

    def __next__(self):
        self.value += 1
        if self.value > 10:
            raise StopIteration
        return self.value

    next = __next__

    def __iter__(self):
        return self


ti = TestIterator()
# print(list(ti))


class MyNumber:
    def __iter__(self):
        self.a = 1000
        return self

    def __next__(self):
        if self.a <= 2000:
            x = self.a
            self.a += 1
            return x
        else:
            raise stopIteration

    next = __next__


myclass = MyNumber()
myiter = iter(myclass)
print(next(myiter))
print(next(myiter))
print(next(myiter))
