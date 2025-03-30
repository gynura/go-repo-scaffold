# users-microservice

## 🗒 Contents
- [Installation](#installation)
- [Getting Started](#getting-started)

## Project folder architecture

The service is written following hexagonal architecture:

```
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
│       ├── core                                      # Business logic core layer
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

## How to build the binary

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


