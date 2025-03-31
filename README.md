# usersvc

Small microservice that manages Users written in Go (v1.24.1).

## Contents
- [Project folders architecture](#project-folders-architecture)
- [API Spec](#api-spec)
- [How to run the service](#how-to-run-the-service)
- [How to test](#how-to-test)
- [Design choices](#design-choices)
- [Improvements](#improvements)

## Project folders architecture

The service is written following hexagonal architecture:

``` ascii
├── app                                            # Complete code of the service
│   ├── bin
│   ├── cmd                                          # Service entry point's code 
│   │   ├── config                                   # Static configuration definition
│   │   └── flags                                    # Command flags definitions
│   ├── docs                                         # Service documentation
│   │   ├── config                                   # Configuration to lint the docs
│   │   └── open-api-spec               # Open API Spec definition of the service's REST API
│   │       ├── components
│   │       │   ├── parameters
│   │       │   │   ├── path
│   │       │   │   └── query
│   │       │   ├── responses
│   │       │   └── schemas
│   │       │       ├── common
│   │       │       ├── health
│   │       │       └── user
│   │       └── paths
│   │           ├── health
│   │           └── user
│   └── internal                       # Private service code 
│       ├── adapters                   # Adapter layer, here we have both input and output adapters
│       │   ├── input
│       │   │   └── http
│       │   └── output
│       │       ├── queue
│       │       │   ├── healthqueue
│       │       │   └── userqueue
│       │       └── storage
│       │           └── postgresql
│       │               ├── healthstorage
│       │               └── userstorage
│       ├── config                         # Config files for the system, in this case the db queries and migrations
│       │   ├── mappers
│       │   │   └── postgresql
│       │   └── migrations
│       │       └── postgres
│       ├── core                                      # Business logic layer
│       │   ├── entities
│       │   ├── events
│       │   └── usecases
│       │       ├── health
│       │       └── user
│       ├── graph                                     # The package that handles all dependency injection using wire
│       ├── infra                                     # Infrastructure layer of the service
│       │   ├── http
│       │   │   ├── generatedhttpinfra
│       │   │   └── middlewares
│       │   ├── queue
│       │   │   └── rabbitmq
│       │   └── storage
│       │       └── postgresql
│       │           └── userstorageinfra
│       ├── pkg                                       # Definition of the shared libraries used across layers
│       │   ├── commons
│       │   │   ├── logger
│       │   │   ├── time
│       │   │   └── validators
│       │   ├── errors
│       │   └── persistence
│       │       └── sql
│       │           ├── init
│       │           └── postgres
│       └── test                                  # Testing related mocks/helpers
│           └── e2e                                # Code of the consumer client for e2e testing
│           ├── helpers                           # Unit testing helpers that use mocks to easy testing
│           └── generatedmocks
└── test                                              # Testing folder with entry point to locally start the app
    ├── deploy
    │       └── postgresql      
    └── examples                              # Example config file and postman collection

```

## API Spec

### Requirements

- [swagger-cli](https://www.npmjs.com/package/swagger-cli)
- [swagger-ui-watcher](https://www.npmjs.com/package/swagger-ui-watcher)

You can serve the API Spec from ``/app/docs`` by executing:

```bash
make serve_api
```

## How to run the service

### How to run the required services (broker & database)

#### Prerequisites 

 - **Docker** installed

1. Navigate to the project's ``test/deploy`` directory:

   ```bash
   cd test/deploy
   ```

2. Execute the following command which will start **Postgres** and **RabbitMQ** containers and run the Postgres migrations:

   ```bash
   make
   ```
   
> **_NOTE:_**  You can see the default configurations for RabbitMQ and Postgres inside **/test/examples/service-config.yaml**.

### How to build the binary

To install the service follow this steps:

1. Navigate to the project's ``app`` directory:

   ```bash
   cd app
   ```

2. Build the binary:

   ```bash
   make build
   ```

3. Change directory to ``bin`` where the built binary will be stored:

   ```bash
   cd bin  
   ```
   
4. You can execute the binary from there pointing to the test configuration(``test/examples/service-config.yaml``) with: 

   ```bash
   ./usersvc
   ```
> **_NOTE:_**  By default the apps runs on localhost:8080, you can change it with the **address** and **port** flags

## How to test

1. Navigate to the ``app`` directory. 

2. Execute the unit test with: 

   ```bash
   make unit_tests
   ```
   
3. To execute e2e tests while your service is running execute: 

   ```bash
   make e2e
   ```

### Manual testing

Inside ``test/examples`` there's a postman collection that covers all the API endpoints.


## Design choices

I will structure this part with questions that I think one might ask themselves when looking at my code: 

### Why hexagonal architecture? This is just a small microservice!

Well you are right but hexagonal architecture once you get the grasp of it is not that slow to bootstrap and helps keep the code detached from implementation details. 
The main benefits that I love as a dev are:

1. **Testability**: Very easy to test, plug and play (literally!) which also highly impacts its maintenance. 
2. **Adaptability and future of the app**: Right now I only support Postgres and RabbitMQ but just by writing new adapters and infra and extending a bit the static configuration I could support
any database/broker I needed. This can help keep the microservice fresh if new technologies drop!
3. **Separation of concerns**: Don't you hate when a piece of code does 29309302930 things? Well I do!

### Postgres can generate UUIDs for you! Didn't you read the docs?

I did! But don't you think that coupling business logic to your database technology can damage the future of your app? All business logic in my opinion has to be written inside the 
defined use cases! What if you use another database that doesn't generate UUIDs?

### Why didn't you use an ORM as your database framework?

I just like to be in control of my mappers :P (Tho nothing bad with ORM, in fact way easier to maintain!)

### Why JSON logging? Are you a robot? 

Who doesn't love querying logs through traceable attributes? Tho I agree they are uglier to look at, it just feels worse to debug without them am I right?

### Why did you rely on code generation? Are you really that lazy?

Well relying on code generators is something that I personally enjoy since it saves me from writing boilerplate code and test it.

An example of how convenient they are is that I was able to do changes in the API Spec during development (even if you do API first, you know sometimes you need to make adjustments),
and just by running a command I had the http infra aligned with the new changes. Sweet!

### What is your testing approach? 

Ideally I think that only libraries, entities and use cases need unit testing since they are more logic oriented. Adapters and infra are truly tested with 
integration tests. 

If we talk about e2e testing, then I feel like the best approach is to write code separated from the app being tested, I want to avoid any coupling in order
to better simulate a real usage. The e2e I wrote for this project are just an http client of the API with a RabbitMQ consumer for the notify messages. (Lot of room to improve them obv!)

### How did you manage Dependency Injection? 

Through Wire! That DI framework is like _magic_, love it! You can take a look at it in **/app/internal/graph** folder.


**TL;DR:** The core principle of my design approach is to stay as decoupled as possible

## Improvements

There are a lot of things that I was not able to do as well as things I left "incomplete", feel free to take a look at the TODOs I left! (I promise they are not as many)

Meanwhile, here's a list that I wrote of possible improvements while writing this service: 

* Create integration tests.
* Make the users' password encryption pluggable through an interface making the service more flexible to changes to the encryption approach.
* Extend the binary functionality to be able to migrate the db itself since right now the database migration is tied to a script. 
* Right now the whole dependency injection graph is coupled to the RabbitMQ and Postgres implementations of the adapters/infra even though the app code is decoupled from them, a great improvement would be to decouple the dependency injection too
* Stacktrace for the custom errors' library.
* Extend the health check in the kubernetes way, separating liveness and readiness.
* Add an error translator in the RabbitMQ infra layer.
* Configurable min log level.
* Write more unit tests, the whole health use case was left without tests.

Obviously way more things left to do, but I hope this gives a nice overview, thanks! 



