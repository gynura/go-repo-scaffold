# usersvc

Small microservice that manages Users written in Go (v1.24.1). 

## Contents
- [Project folders architecture](#project-folders-architecture)
- [How to build the binary](#how-to-build-the-binary)
- [How to run the service](#how-to-run-the-service)

## Project folders architecture

The service is written following hexagonal architecture:

``` ascii
├── app                                            # Complete code of the service
│   ├── bin
│   ├── cmd                                          # Service entry point's code 
│   │   ├── config                                   # Static configuration definition
│   │   └── flags                                    # Command flags definitions
│   ├── docs                                         # Service documentation
│   │   ├── config                                   # Configuration to serve/lint/format the docs
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
│   └── internal                                    # Private service code 
│       ├── adapters                                # Adapter layer, here we have both input and output adapters
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
│       └── test                                      # Testing related mocks/helpers
│           └── generatedmocks
└── test                                              # Testing folder with entry point to locally start the app
    ├── deploy
    │   └── postgresql
    └── examples

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
