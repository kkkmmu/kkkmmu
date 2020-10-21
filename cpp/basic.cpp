#include <iostream>
#include <algorithm>
#include <vector>
#include <memory>
#include <thread>

using namespace std;

template<typename T, typename U> auto add(T x, U y){
	return x + y;
};

int add(int a)
{
	return a + 1;
}

int main()
{
	auto x = 1;
	auto y = 2;
	decltype (x+y) z;

	if (std::is_same<decltype(x), int>::value)
		std::cout << "type x == int" << std::endl;

	if (std::is_same<decltype(z), decltype(x)>::value)
		std::cout << "type z == type x" << std::endl;

	std::cout <<add(10, 10.0) << endl;
	std::cout <<add<int, int>(10, 10.0) << endl;
	std::cout <<add<int, const char*>(10, "hello") << endl;

	std::vector<int> vec = {1, 3, 5, 7, 9};
	if (auto itr = std::find(vec.begin(), vec.end(), 3); itr != vec.end())
		*itr = 4;
	for (auto element : vec)
		std::cout << element << std::endl;

	for (auto &element : vec) 
		element += 1;

	for (auto element : vec)
		cout << element << endl;

	auto pointer = std::make_shared<int>(10);
	auto pointer2 = pointer;
	auto pointer3 = pointer;

	int *p = pointer.get();

	cout << "Pointer.use_count() = " << pointer.use_count() << std::endl;
	cout << "Pointer2.use_count() = " << pointer2.use_count() << std::endl;
	cout << "Pointer3.use_count() = " << pointer3.use_count() << std::endl;

	pointer.reset();

	cout << "Pointer.use_count() = " << pointer.use_count() << std::endl;
	cout << "Pointer2.use_count() = " << pointer2.use_count() << std::endl;
	cout << "Pointer3.use_count() = " << pointer3.use_count() << std::endl;
	
	pointer2.reset();
	cout << "Pointer.use_count() = " << pointer.use_count() << std::endl;
	cout << "Pointer2.use_count() = " << pointer2.use_count() << std::endl;
	cout << "Pointer3.use_count() = " << pointer3.use_count() << std::endl;

	std::unique_ptr<int> upointer = std::make_unique<int>(10);
	//std::unique_ptr<int> upointer2 = upointer;
	std::unique_ptr<int> upointer2 = move(upointer);

	std::thread t1([]() { cout << "Hello world 1" << endl; });
	std::thread t2([]() { cout << "Hello world 2" << endl; });

	t1.join();
	t2.join();

	clog << "Hello world!" << endl;

	int i  = -42;
	unsigned int j = 10;
	cout << i + j << endl;
	cout << j + i << endl;
	return 0;
}
