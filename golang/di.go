Dependency Injection in Go
I recently built a small project in Go. I’ve been working with Java for the past few years and was immediately struck by the lack of momentum behind Dependency Injection (DI) in the Go ecosystem. I decided to try building my project using Uber’s dig library and was very impressed.

I found that DI helped solve a lot of problems I had encountered in my previous Go applications – overuse of the init function, abuse of globals and complicated application setup.

In this post I’ll give an introduction to DI and then show an example application before and after using a DI framework (via the dig library).

A Brief Overview of DI
Dependency Injection is the idea that your components (usually structs in go) should receive their dependencies when being created. This runs counter to the associated anti-pattern of components building their own dependencies during initialization. Let’s look at an example.

Suppose you have a Server struct that requires a Config struct to implement its behavior. One way to do this would be for the Server to build its own Config during initialization.

type Server struct {
	config *Config
}

func New() *Server {
	return &Server{
		config: buildMyConfigSomehow(),
	}
}
This seems convenient. Our caller doesn’t have to be aware that our Server even needs access to Config. This is all hidden from the user of our function.

However, there are some disadvantages. First of all, if we want to change the way our Config is built, we’ll have to change all the places that call the building code. Suppose, for example, our buildMyConfigSomehow function now needs an argument. Every call site would need access to that argument and would need to pass it into the building function.

Also, it gets really tricky to mock the behavior of our Config. We’ll somehow have to reach inside of our New function to monkey with the creation of Config.

Here’s the DI way to do it:

type Server struct {
	config *Config
}

func New(config *Config) *Server {
	return &Server{
		config: config,
	}
}
Now the creation of our Server is decoupled from the creation of the Config. We can use whatever logic we want to create the Config and then pass the resulting data to our New function.

Furthermore, if Config is an interface, this gives us an easy route to mocking. We can pass anything we want into New as long as it implements our interface. This makes testing our Server with mock implementations of Config simple.

The main downside is that it’s a pain to have to manually create the Config before we can create the Server. We’ve created a dependency graph here – we must create our Config first because of Server depends on it. In real applications these dependency graphs can become very large and this leads to complicated logic for building all of the components your application needs to do its job.

This is where DI frameworks can help. A DI framework generally provides two pieces of functionality:

A mechanism for “providing” new components. In a nutshell, this tells the DI framework what other components you need to build yourself (your dependencies) and how to build yourself once you have those components.
A mechanism for “retrieving” built components.
A DI framework generally builds a graph based on the “providers” you tell it about and determines how to build your objects. This is very hard to understand in the abstract, so let’s walk through a moderately-sized example.

An Example Application
We’re going to be reviewing the code for an HTTP server that delivers a JSON response when a client makes a GET request to /people. We’ll review the code piece by piece. For simplicity sake, it all lives in the same package (main). Please don’t do this in real Go applications. Full code for this example can be found here.

First, let’s look at our Person struct. It has no behavior save for some JSON tags.

type Person struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}
A Person has an Id, Name and Age. That’s it.

Next let’s look at our Config. Similar to Person, it has no dependencies. Unlike Person, we will provide a constructor.

type Config struct {
	Enabled      bool
	DatabasePath string
	Port         string
}

func NewConfig() *Config {
	return &Config{
		Enabled:      true,
		DatabasePath: "./example.db",
		Port:         "8000",
	}
}
Enabled tells us if our application should return real data. DatabasePath tells us where our database lives (we’re using sqlite). Port tells us the port on which we’ll be running our server.

Here’s the function we’ll use to open our database connection. It relies on our Config and returns a *sql.DB.

func ConnectDatabase(config *Config) (*sql.DB, error) {
	return sql.Open("sqlite3", config.DatabasePath)
}
Next we’ll look at our PersonRepository. This struct will be responsible for fetching people from our database and deserializing those database results into proper Person structs.

type PersonRepository struct {
	database *sql.DB
}

func (repository *PersonRepository) FindAll() []*Person {
	rows, _ := repository.database.Query(
		`SELECT id, name, age FROM people;`
	)
	defer rows.Close()

	people := []*Person{}

	for rows.Next() {
		var (
			id   int
			name string
			age  int
		)

		rows.Scan(&id, &name, &age)

		people = append(people, &Person{
			Id:   id,
			Name: name,
			Age:  age,
		})
	}

	return people
}

func NewPersonRepository(database *sql.DB) *PersonRepository {
	return &PersonRepository{database: database}
}
PersonRepository requires a database connection to be built. It exposes a single function called FindAll that uses our database connection to return a list of Person structs representing the data in our database.

To provide a layer between our HTTP server and the PersonRepository, we’ll create a PersonService.

type PersonService struct {
	config     *Config
	repository *PersonRepository
}

func (service *PersonService) FindAll() []*Person {
	if service.config.Enabled {
		return service.repository.FindAll()
	}

	return []*Person{}
}

func NewPersonService(config *Config, repository *PersonRepository)
*PersonService {
	return &PersonService{config: config, repository: repository}
}
Our PersonService relies on both the Config and the PersonRepository. It exposes a function called FindAll that conditionally calls the PersonRepository if the application is enabled.

Finally, we’ve got our Server. This is responsible for running an HTTP server and delegating the appropriate requests to our PersonService.

type Server struct {
	config        *Config
	personService *PersonService
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/people", s.people)

	return mux
}

func (s *Server) Run() {
	httpServer := &http.Server{
		Addr:    ":" + s.config.Port,
		Handler: s.Handler(),
	}

	httpServer.ListenAndServe()
}

func (s *Server) people(w http.ResponseWriter, r *http.Request) {
	people := s.personService.FindAll()
	bytes, _ := json.Marshal(people)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}

func NewServer(config *Config, service *PersonService) *Server {
	return &Server{
		config:        config,
		personService: service,
	}
}
The Server is dependent on the PersonService and the Config.

Ok, we know all the components of our system. Now how the hell do we actually initialize them and start our system?

The Dreaded main()
First, let’s write our main() function the old fashioned way.

func main() {
	config := NewConfig()

	db, err := ConnectDatabase(config)

	if err != nil {
		panic(err)
	}

	personRepository := NewPersonRepository(db)

	personService := NewPersonService(config, personRepository)

	server := NewServer(config, personService)

	server.Run()
}
First, we create our Config. Then, using the Config, we create our database connection. From there we can create our PersonRepository which allows us to create our PersonService. Finally, we can use this to create our Server and run it.

Phew, that was complicated. Worse, as our application becomes more complicated, our main will continue to grow in complexity. Every time we add a new dependency to any of our components, we’ll have to reflect that dependency with ordering and logic in the main function to build that component.

As you might have guessed, a Dependency Injection framework can help us solve this problem. Let’s examine how.

Building a Container
The term “container” is often used in DI frameworks to describe the thing into which you add “providers” and out of which you ask for fully-build objects. The dig library gives us the Provide function for adding providers and the Invoke function for retrieving fully-built objects out of the container.

First, we build a new container.

container := dig.New()
Now we can add new providers. To do so, we call the Provide function on the container. It takes a single argument: a function. This function can have any number of arguments (representing the dependencies of the component to be created) and one or two return values (representing the component that the function provides and optionally an error).

container.Provide(func() *Config {
	return NewConfig()
})
The above code says “I provide a Config type to the container. In order to build it, I don’t need anything else.” Now that we’ve shown the container how to build a Config type, we can use this to build other types.

container.Provide(func(config *Config) (*sql.DB, error) {
	return ConnectDatabase(config)
})
This code says “I provide a *sql.DB type to the container. In order to build it, I need a Config. I may also optionally return an error.”

In both of these cases, we’re being more verbose than necessary. Because we already have NewConfig and ConnectDatabase functions defined, we can use them directly as providers for the container.

container.Provide(NewConfig)
container.Provide(ConnectDatabase)
Now, we can ask the container to give us a fully-built component for any of the types we’ve provided. We do so using the Invoke function. The Invoke function takes a single argument – a function with any number of arguments. The arguments to the function are the types we’d like the container to build for us.

container.Invoke(func(database *sql.DB) {
	// sql.DB is ready to use here
})
The container does some really smart stuff. Here’s what happens:

The container recognizes that we’re asking for a *sql.DB
It determines that our function ConnectDatabase provides that type
It next determines that our ConnectDatabase function has a dependency of
Config
It finds the provider for Config, the NewConfig function
NewConfig doesn’t have any dependencies, so it is called
The result of NewConfig is a Config that is passed to ConnectDatabase
The result of ConnectionDatabase is a *sql.DB that is passed back to the
caller of Invoke
That’s a lot of work the container is doing for us. In fact, it’s doing even more. The container is smart enough to build one, and only one, instance of each type provided. That means we’ll never accidentally create a second database connection if we’re using it in multiple places (say multiple repositories).

A Better main()
Now that we know how the dig container works, let’s use it to build a better main.

func BuildContainer() *dig.Container {
	container := dig.New()

	container.Provide(NewConfig)
	container.Provide(ConnectDatabase)
	container.Provide(NewPersonRepository)
	container.Provide(NewPersonService)
	container.Provide(NewServer)

	return container
}

func main() {
	container := BuildContainer()

	err := container.Invoke(func(server *Server) {
		server.Run()
	})

	if err != nil {
		panic(err)
	}
}
The only thing we haven’t seen before here is the error return value from Invoke. If any provider used by Invoke returns an error, our call to Invoke will halt and that error will be returned.

Even though this example is small, it should be easy to see some of the benefits of this approach over our “standard” main. These benefits become even more obvious as our application grows larger.

One of the most important benefits is the decoupling of the creation of our components from the creation of their dependencies. Say, for example, that our PersonRepository now needs access to the Config. All we have to do is change our NewPersonRepository constructor to include the Config as an argument. Nothing else in our code changes.

Other large benefits are lack of global state, lack of calls to init (dependencies are created lazily when needed and only created once, obviating the need for error-prone init setup) and ease of testing for individual components. Imagine creating your container in your tests and asking for a fully-build object to test. Or, create an object with mock implementations of all dependencies. All of these are much easier with the DI approach.

An Idea Worth Spreading
I believe Dependency Injection helps build more robust and testable applications. This is especially true as these applications grow in size. Go is well suited to building large applications and has a great DI tool in dig. I believe the Go community should embrace DI and use it in far more applications.

Update
Google recently released their own DI container called wire. It avoids runtime reflection by building the container using code generation. I would recommend using it rather than dig. For more details, see my new post.
