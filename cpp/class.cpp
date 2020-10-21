#include <iostream>

using namespace std;

class Base 
{
	public:
		Base() { std::cout << "Create base class " << endl; }
		Base(int i) : _i(i) { std::cout << "Create base class with int  " << i << endl; }
		Base(double i) : _d(i) { std::cout << "Create base class with double " << i << endl; }
		Base(std::string i) : _s(i) { std::cout << "Create base class with string " << i << endl; }
		~Base() { std::cout << "Destroy base class " << endl; }
		int get_int() { return _i; }
		double get_double() { return _d; }
		std::string get_string();

	private:
		int _i;
		double _d;
		std::string _s;
};

class Sub : public Base 
{
	public:
		Sub() { cout << "Create Sub class " << endl; }
		Sub(int i) : Base(i) { cout << "Create Sub class with int " << i << endl; }
		Sub(double i) : Base(i) { cout << "Create Sub class with double " << i << endl; }
		Sub(std::string i) : Base(i) { cout << "Create Sub class with string " << i << endl; }
		~Sub() { cout << "Destroy sub class " << endl; }
};

class Child : public Base 
{
	public:
		Child() { cout << "Create Child class " << endl; }
		Child(int i) : Base(i) { cout << "Create Child class with int " << i << endl; }
		Child(double i) : Base(i) { cout << "Create Child class with double " << i << endl; }
		Child(std::string i) : Base(i) { cout << "Create Child class with string " << i << endl; }
		~Child() { cout << "Destroy sub class " << endl; }
};

std::string Base::get_string()
{
	return _s;
}

int main()
{
	class Base b;
	class Base c(1);
	class Base d(1.11);
	class Base s(std::string("Hello Base"));

	class Sub sb;
	class Sub sc(2);
	class Sub sd(2.22);
	class Sub ss(std::string("Hello Sub"));

	class Child cb;
	class Child cc(3);
	class Child cd(3.33);
	class Child cs(std::string("Hello Child"));

	class Base& rb = sb;
	class Base& rc = sc;
	class Base& rd = sd;
	class Base& rs = ss;

	class Base& rb1 = cb;
	class Base& rc1 = cc;
	class Base& rd1 = cd;
	class Base& rs1 = cs;



	cout << c.get_int() << " " << d.get_double() << " " << s.get_string() << endl;
	cout << sc.get_int() << " " << sd.get_double() << " " << ss.get_string() << endl;
	cout << rc.get_int() << " " << rd.get_double() << " " << rs.get_string() << endl;
	cout << rc1.get_int() << " " << rd1.get_double() << " " << rs1.get_string() << endl;
}
