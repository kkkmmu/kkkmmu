class Foo(object):
    def __init__(self, val):
        self.__data = val

    @property
    def data(self):
        return self.__data

    @data.setter
    def data(self, val):
        if val < 0 or val > 200:
            raise Exception
        self.__data = val


foo = Foo(100)
# print(foo.__data)
foo.data = 200
foo.data = 300
print(foo.data)
