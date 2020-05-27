Go 函数选项模式
CODE TALKS 2018-02-01   go, translation
本文译自 Functional Options Pattern in Go 版权@归原文所有.
Golang 开发者遇到的许多问题之一是尝试将一个函数的参数设置为可选. 这是一个非常常见的用例, 有些对象应该使用一些基本的默认设置来开箱即用, 并且你偶尔可能需要提供一些更详细的配置.

在很多语言中这很容易; 在 C 族语言中, 可以使用不同数量的参数提供相同函数的多个版本; 在像 PHP 这样的语言中, 可以给参数一个默认值，并在调用方法时忽略它们. 但是在 Golang 中, 这两种方式你哪个也用不了. 那么你如何创建一个函数, 用户可以指定一些额外的配置?

有很多可能的方法可以做到这一点, 但是大多数都不能满足要求, 或者需要在服务端的代码中进行额外的检查和验证, 或者通过传递额外的客户端不关心的参数来为客户端做额外的工作.

我将介绍一些不同的方案, 并说明为什么每个都不是最理想的, 然后我们将建立我们自己的方式来作为最终干净的解决方案: 函数式选项模式 (Functional Options Pattern).

我们来看一个例子. 比方说, 我们有一些名为 StuffClient 的服务, 它有一些东西, 并有两个配置选项(timeout 和 retries):

type StuffClient interface {
	DoStuff() error
}
type stuffClient struct {
	conn    Connection
	timeout int
	retries int
}
结构体 stuffClient 是私有的, 所以我们应该为它提供一些构造器:

func NewStuffClient(conn Connection, timeout, retries int) StuffClient {
	return &stuffClient{
		conn:    conn,
		timeout: timeout,
		retries: retries,
	}
}
嗯, 但是现在我们每次调用 NewStuffClient 时都要提供 timeout 和 retries. 而大多数时候我们只想使用默认值. 我们不能用不同的参数数量来定义多个版本的 NewStuffClient, 否则我们会得到一个类似 “NewStuffClient redeclared in this blockt” 的编译错误.

一个方案是创建另一个不同名称的构造器:

func NewStuffClient(conn Connection) StuffClient {
	return &stuffClient{
		conn:    conn,
		timeout: DEFAULT_TIMEOUT,
		retries: DEFAULT_RETRIES,
	}
}
func NewStuffClientWithOptions(conn Connection, timeout, retries int) StuffClient {
	return &stuffClient{
		conn:    conn,
		timeout: timeout,
		retries: retries,
	}
}
但是, 这有点蹩脚. 我们可以做得比这更好. 如果我们传入了一个配置对象呢:

type StuffClientOptions struct {
	Retries int //number of times to retry the request before giving up
	Timeout int //connection timeout in seconds
}
func NewStuffClient(conn Connection, options StuffClientOptions) StuffClient {
	return &stuffClient{
		conn:    conn,
		timeout: options.Timeout,
		retries: options.Retries,
	}
}
但是, 这也不是很好. 现在我们总是需要创建 StuffClientOptions, 并且即使我们不想指定任何选项也要传递它. 而且我们也没有自动填充默认值, 除非我们在代码中添加了一堆检查, 或者我们可以传入一个 DefaultStuffClientOptions 变量 (也不好, 因为它可以在一个地方修改导致其他地方有问题).

那么解决方案是什么? 解决这个难题最好的方法是使用函数式选项模式, 利用 Go 对闭包的方便支持. 让我们继续我们上面定义的 StuffClientOptions, 但是我们会添加一些东西:

type StuffClientOption func(*StuffClientOptions)
type StuffClientOptions struct {
	Retries int //number of times to retry the request before giving up
	Timeout int //connection timeout in seconds
}
func WithRetries(r int) StuffClientOption {
	return func(o *StuffClientOptions) {
		o.Retries = r
	}
}
func WithTimeout(t int) StuffClientOption {
	return func(o *StuffClientOptions) {
		o.Timeout = t
	}
}
泥土般芬芳, 对吧? 究竟发生了什么? 基本上我们有我们的结构定义我们的 StuffClient 的可用选项. 另外现在我们定义了一个名为 StuffClientOption 的东东(这次是单数), 它只是一个接受我们的选项 struct 作为参数的函数. 我们已经定义了另外一些函数 WithRetries 和 WithTimeout, 它们返回一个闭包. 现在魔法降临:

var defaultStuffClientOptions = StuffClientOptions{
	Retries: 3,
	Timeout: 2,
}
func NewStuffClient(conn Connection, opts ...StuffClientOption) StuffClient {
	options := defaultStuffClientOptions
	for _, o := range opts {
		o(&options)
	}
	return &stuffClient{
		conn:    conn,
		timeout: options.Timeout,
		retries: options.Retries,
	}
}
我们已经定义了一个额外的未导出(unexposed)变量, 包含我们的默认选项, 我们现在调整了我们的构造函数, 而不是接受一个可变参数. 然后, 我们遍历 StuffClientOption (单数) 列表, 并对每一项应用返回的闭包到我们的选项变量.

现在我们要做的就是这样:

x := NewStuffClient(Connection{})
fmt.Println(x) // prints &{{} 2 3}
x = NewStuffClient(
	Connection{},
	WithRetries(1),
)
fmt.Println(x) // prints &{{} 2 1}
x = NewStuffClient(
	Connection{},
	WithRetries(1),
	WithTimeout(1),
)
fmt.Println(x) // prints &{{} 1 1}
这看起来相当不错而且可用. 而关于它的好的部分是, 我们可以随时添加新的选项, 只需要对代码进行非常少量的更改. 把这些都组合起来就是这样:

var defaultStuffClientOptions = StuffClientOptions{
	Retries: 3,
	Timeout: 2,
}
type StuffClientOption func(*StuffClientOptions)
type StuffClientOptions struct {
	Retries int //number of times to retry the request before giving up
	Timeout int //connection timeout in seconds
}
func WithRetries(r int) StuffClientOption {
	return func(o *StuffClientOptions) {
		o.Retries = r
	}
}
func WithTimeout(t int) StuffClientOption {
	return func(o *StuffClientOptions) {
		o.Timeout = t
	}
}
type StuffClient interface {
	DoStuff() error
}
type stuffClient struct {
	conn    Connection
	timeout int
	retries int
}
type Connection struct {}
func NewStuffClient(conn Connection, opts ...StuffClientOption) StuffClient {
	options := defaultStuffClientOptions
	for _, o := range opts {
		o(&options)
	}
	return &stuffClient{
		conn:    conn,
		timeout: options.Timeout,
		retries: options.Retries,
	}
}
func (c stuffClient) DoStuff() error {
	return nil
}
如果你想自己尝试一下, 去 Go Playground 吧.

但是可以通过删除 StuffClientOptions 结构并直接将选项应用到我们的 StuffClient 来进一步简化.

var defaultStuffClient = stuffClient{
	retries: 3,
	timeout: 2,
}
type StuffClientOption func(*stuffClient)
func WithRetries(r int) StuffClientOption {
	return func(o *stuffClient) {
		o.retries = r
	}
}
func WithTimeout(t int) StuffClientOption {
	return func(o *stuffClient) {
		o.timeout = t
	}
}
type StuffClient interface {
	DoStuff() error
}
type stuffClient struct {
	conn    Connection
	timeout int
	retries int
}
type Connection struct{}
func NewStuffClient(conn Connection, opts ...StuffClientOption) StuffClient {
	client := defaultStuffClient
	for _, o := range opts {
		o(&client)
	}

	client.conn = conn
	return client
}
func (c stuffClient) DoStuff() error {
	return nil
}
可以在这里尝试一下. 在我们的示例中, 我们只是直接将配置应用到我们的结构中, 在中间有一个额外的配置结构是没有意义的. 但是请注意, 在许多情况下, 您可能仍然想使用前面示例中的 config 结构体; 例如, 如果你的构造器使用配置选项来执行一些操作, 但不把它们存储到结构中, 或者传递到其他地方. 配置结构变体是更通用的实现.

感谢 Rob Pike 和 Dave Cheney 推广这种设计模式
