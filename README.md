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
│   │   └── open-api-spec                            # Open API Spec definition of the service's REST API
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
│       ├── config                                    # Config files for the system, in this case the db queries and migrations
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
    └── examples                                      # Example config file and postman collection

```

## API Spec

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
   
4. You can execute the binary from there pointing to the test configuration with: 

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

### Manually testing

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
2x1, I can develop way faster and I save myself from that oblivious task (HTTP Handlers I am looking at you!)

### What is your testing approach? 

Ideally I think that only libraries, entities and use cases need unit testing since they are more logic oriented. Adapters and infra are truly tested with 
integration tests. 

If we talk about e2e testing, then I feel like the best approach is to write code separated from the app being tested, I want to avoid any coupling in order
to better simulate a real usage. 

### How did you manage Dependency Injection? 

Through Wire! That DI framework is like _magic_, love it! You can take a look at it in **/app/internal/graph** folder. 


## Improvements



