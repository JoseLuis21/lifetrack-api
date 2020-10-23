# Domain-Driven Design Simplified, CQRS and Hexagonal Architecture for Modern APIs

Created: Oct 22, 2020 8:23 PM

# Context

When we talk about Domain-Driven Design, we often refer to complex business domains and even more with a CQRS implementation. Whatsoever, sometimes one can just simplify some concepts from these amazing tools so you may feel comfortable even with simple domains. Another reason to use the theoretical approach I'm about to expose, apart from having a simple domain as mentioned before, is using programming languages such as Go or Python that tends to be simple by their own nature. The main benefits of DDD and CQRS such as loosely coupled layers and readable software using Ubiquitous Language will be our guidelines to keep consistency and a robust standard in our ecosystem.

In the other hand, webservices in form of APIs (using HTTP protocols + REST mostly) are being widely used nowadays. Here is where I'd like to mention the integration of Hexagonal Architecture, which gives the developer a way to decouple his/her applications even more using an input/output principle.

It is completely true these approaches were made mainly for complex applications/ecosystems, yet one can take some benefits from them and at the same time maintaining the complexity level lower. 

After this little introduction, we shall begin.

# Modern ecosystem principles

Before explaining any kind of application design or architecture, I'd like to give you a brief introduction about the actual state of modern applications/platforms. If you are not interested, then you shall skip this section.

Nowadays applications must comply a trivial standard to keep them on a level a final user expects. This standard, which considers the Internet as a high-entropy environment, can be simplified in the following principles:

- A software ecosystem as a whole must be scalable.
- A software ecosystem as a whole must be resilient to most common of failures.
- A software ecosystem as a whole must be secure to most common security attacks.

Based on these principles, It is recommended to use different techniques and architectures to comply with the principles. 

The main architectural styles used in nowadays applications should be distributed and loosely coupled to comply with the 1st principle. Either serverless or microservices are a good choice (depending on the needs).

To comply with the 2nd principle, one must instrument its ecosystem with resiliency patterns such as circuit breaker, load balancer, and retry strategy. Most of these patterns and many more can be simplified outside our code by using tools as service meshes (Istio), container orchestrators (Kubernetes) and sidecars (Envoy).

Just to conclude, to comply with the 2nd and 3rd principles, one should instrument its ecosystem using observability tools which can be divided in just three categories: **logging, tracing and monitoring**. There are many tools for each category such as Prometheus (monitoring), Fluentd (logging), Jaeger (tracing) and OpenCensus (tracing and monitoring). This is extremely useful for any member of the development team since one can know what's actually happening inside our whole ecosystem and anomalies can be detected very easily using tools like Grafana or even Kibana for logging if we were using Elasticsearch as log container. Just to be clear, to fully comply with the last principle, one should deploy its ecosystem with different strategies and techniques such as OSI-Layer 7 firewalls, TLS encryption, rate limiters, correct cloud-governance, fine-grained ACL rules and correct NAT configurations.

# Domain-Driven Design Simplified

## Domain

### Value Object

1. If a value object is shared between entities and aggregates, then it should have an string fieldName variable and every factory should accept a fieldName variable and set default values to it if null/empty string
2. A shared value object should throw custom domain errors using the fieldName variable, if not then it should throw common domain errors
3. A value object should contain the following methods:

    *Nomenclature: Method - Short description - return or throw value(s)*

    - IsValid() - Validate the current value object → domain error
    - FromPrimitive() - Accept non-validated primitive values for marshaling purposes → n/a
    - String(), Time(), ... - Retrieve current value in primitive data → primitive value
    - *setFieldName() - Un-exported, Sets and sanitizes value object's field name (if no field name specified, set a default value) → n/a

*setFieldName() is only required when a value object is shared

### Entity

1. An entity should be only used when an aggregate requires more than one entity
2. An entity should contain an IsValid() method
3. An entity should have exported all its fields (value objects)

### Aggregate

1. An aggregate could be using shared value objects, if that is the case, then define field names for each value object shared if applicable
2. An aggregate should have its own domain rules such as required fields
3. An aggregate should interact with the application layer using only value objects
4. An aggregate should contain business actions as methods, this is the only way to interact with domain states (value objects are an exception since they are required to interact with aggregates)
5. An aggregate should have a primitive-only model replica (w/o business logic)
6. An aggregate should have a domain event factory, these events shall be used only in aggregates or in the application layer (Commands mostly)
7. An aggregate should throw domain errors
8. An aggregate root should handle other aggregates just if the aggregate does not have its own life-cycle (e.g. shopping cart - item)
9. An aggregate should contain the following metadata fields:
    - Create Time - Time aggregate was created
    - Update Time - Last time aggregate state was edited
    - Active - Boolean (flag) state of the aggregate, this enables restoring and soft-removal. All queries should retrieve active-only aggregates by default
10. *An aggregate should contain getter methods returning only primitive data
11. An aggregate should contain the following methods:

    *Nomenclature: Method - Short description - return or throw value(s)*

    - IsValid() - Validates the current aggregate's state → throws domain error if anomaly is found in each value object or aggregate's own domain rules
    - RecordEvents([] event) - Stores an array of domain events to the current aggregate → n/a
    - PullEvents() - Retrieves all domain events that had happened in the current aggregate → array of events stored in aggregate
    - MarshalJSON() - Convert the current aggregate into a JSON binary→ array of bytes, throws domain error if marshaling cannot be done
    - MarshalPrimitive() - Convert the current aggregate to a primitive-only model replica, highly used by the infrastructure layer → aggregate's primitive model
    - UnmarshalPrimitive( primitiveModel{} ) - Parse a primitive-only model replica into a valid aggregate → throws domain error if a value object or aggregate's domain rule is violated

### Model

1. A primitive model should have primitive-only values
2. A primitive model must not have business logic
3. A primitive model should have JSON tags if applicable

### Factory

1. An aggregate/entity factory should be used only to create valid aggregates/entities with just required values
2. An aggregate factory must add a create-like domain event to the new aggregate
3. A value object should have two factories, one for validated values and other using non-validated primitives for marshaling purposes

### Adapter

1. An adapter should adapt aggregates only
2. An adapter should process bulk operations only since aggregates already have marshaling properties by their own
3. An adapter should be using aggregates marshaling implementations only

### Event

1. An event must contain the following methods:

    *Nomenclature: Method - Short description - return or throw value(s)*

    - MarshalBinary() - Convert the event into a JSON binary → array of bytes, throws domain error if marshaling cannot be done
    - UnmarshalBinary([]byte) - Parse the given array of bytes into a domain event → throws domain error if un-marshaling process cannot be done
2. An event should have the following structure:
    - Correlation ID: Distributed unique identifier of a domain event
    - Topic: Name of the topic the event will be published into
    - Publisher: Name of the module which published the event
    - Action: Operation that triggered the event
    - Publish Time: UNIX timestamp when event was published
    - Aggregate ID: Domain's aggregate unique identifier
    - Body: Actual message that will be published as body on the event bus, if an aggregation process has happened, it is recommended to use the primitive model as body for data eventual consistency and projections in the ecosystem

### Exception

1. A domain exception should be thrown in application, domain and infrastructure layers
2. Any exception apart from domain exceptions should be considered harmful to the whole ecosystem (like infrastructure exceptions such as network and driver problems)
3. If using a third-party library, wrap exceptions such as not found, already exists, invalid values and user-fault exceptions with domain exceptions, they should be harmless to the ecosystem since they can be handled
4. A domain exception should be identified using the following groups:
    - Not found
    - Out of range
    - Invalid format
    - Required
    - Already exists

### Repository

1. A repository must be abstracted using an interface, thus every time a persistence layer is required by the application layer, it must use the interface instead of concrete implementations
2. A repository fetch method should accept an specific criteria struct to define the fetching strategy
3. An event bus should contain the two following methods:
    - Publish([] event) - Send the given array of events to the event bus → throws domain error if marshaling cannot be done or if topic was not found, throws driver error
    - SubscribeTo(topic) - Subscribe to the given topic from the event bus → channel(event), throws domain error if topic was not found, throws driver error

### Event Bus

1. An event bus must be abstracted using an interface, thus every time an event bus layer is required by the application layer, it must use the interface instead of concrete implementations
2. An event bus should contain the two following methods:
    - Publish([] event) - Send the given array of events to the event bus → throws domain error if marshaling cannot be done or if topic was not found, throws driver error
    - SubscribeTo(topic) - Subscribe to the given topic from the event bus → channel(event), throws domain error if topic was not found, throws driver error

## Application

### Command

1. A command should receive primitive-only data
2. A command should receive a context for context propagation purposes (used by distributed tracing)
3. A command handler must have exported only an Invoke() method
4. A command handler invoker must return either a simple domain/driver exception or the aggregate id which was modified/created and the exception if applicable
5. A command handler should get injected a repository and an event bus concrete implementation
6. A command handler should parse given primitive data into value object
7. A command handler should attach any domain event depending on its operations, if used factory to create an aggregate, then skip (factories must add the create domain event)
8. A command handler should implement a doer strategy if business domain is complex to avoid code smells
9. A command handler must fetch data directly from the repository (do not use query structs/classes)
10. A command handler should separate each operation with isolated (non-exported) functions (e.g. store & fetch data, doer, parsing, etc ...)
11. A command handler should be abstracted with an interface
12. A command handler should be wrapped with observability tools such as logging and monitoring using a chain of responsibility pattern (middleware) for just business cases

### Query

1. A query must receive primitive only data
2. A query should receive a custom criteria struct replica using JSON tags if applicable (inside application layer) to eventually map it into the domain criteria struct
3. A query struct/class must get injected a repository concrete implementation (this can be different from the command repository implementation if the ecosystem had heterogenous databases for a service/function like PostgreSQL - Write / Elasticsearch - Read)
4. A query struct/class should get injected an event bus concrete implementation
5. A query should be abstracted with an interface
6. A query should be wrapped with observability tools such as logging and monitoring using a chain of responsibility pattern (middleware) for just business cases
7. *If applicable, every time an aggregate is fetched correctly and organically, a query should either throw an event or write to the time-series monitoring database through a middleware pattern for data analytics services/purposes
8. *If applicable, a query should increment a total view counter for each aggregate fetched correctly and organically (metadata)

## Infrastructure

### Exception

To make a resilient and fault-tolerant ecosystem, it is recommended to wrap certain infrastructure exceptions such as network errors. This will help circuit breakers and retry strategies to identify exceptions that should be actually handled

### Persistence (Repository)

1. A repository should be using Read-Write mutability locks (mutex)
2. A repository should be instrumented with logging, monitoring, tracing and resiliency, it is recommended to wrap the interface using a chain of responsibility pattern (middleware)
3. A repository criteria should contain the following fields:
    - ID - Aggregate ID
    - Limit - Page size, it is recommended to set default values if given value violates domain rules
    - Token - Next page token, sometimes it could be the next aggregate ID or an stringified JSON
4. A repository should implement the following fetching strategies:
    - Fetch All - Uses limit and token fields, fetch all data with pagination techniques
    - Fetch By ID - Uses ID field, fetch one single aggregate with the matching unique identifier

### Event Bus

1. An event bus should be using write mutability locks (mutex)
2. An event bus should be using a distributed queue system, every group of nodes (e.g. a microservice or a serverless function) should have its own queue for each topic working as a consumer group (like Apache Kafka)
3. An event bus should be instrumented with logging, monitoring, tracing and resiliency, it is recommended to wrap the interface using a chain of responsibility pattern (middleware)
4. If the given context happens inside a distributed transaction and a domain exception was thrown then do not retry, log the error at a warning level and acknowledge the message to the distributed queue system

### Configuration

1. A configuration should read its field values from environment variables or a configuration file such as YAML, JSON or Dotenv
2. A configuration file should not contain secrets, for that kind of purposes use services such as AWS Secrets Manager or Hashicorp Vault. If still want to proceed, then do not push any secrets to a public repository
3. An environment variable used by application configuration should have a unique prefix such as the first letters of the platform name to avoid collisions with OS-required variables
4. A configuration should be a simple class/struct which contains infrastructure configuration,  development stages and application/service versioning preferably using Semantic Versioning, for any other purpose use global variables inside each layer
5. A configuration must have default values for each field
6. If a configuration field value has changed at runtime, write changes to each source
7. A configuration file should be listening to changes for each source

## View

### Transport

1. A server should implement one of the following transport strategies:
    - Handler - HTTP
    - Schema - GraphQL
    - Action - GRPC, Thrift (TCP/UDP with binary)
    - Consumer - PubSub (TCP/UDP with binary, AMQP, etc)
2. A server should implement observability techniques such as logging, tracing, monitoring and resiliency (if not consumer strategy, then rate-limit only) for each strategy
3. If working with a microservice and if applicable, a server should contain an strategy abstraction using interfaces with the following methods:
    - SetRoutes(router) - Inject transport routes into the given router (calls mapRoutes() method) → n/a
    - Name() - Retrieves the transport strategy name → transport strategy name
4. If applicable, a transport strategy implementation should contain a mapRoutes() method which allows the implementation to attach handlers with transport routing
5. If working with a serverless environment, each strategy handler must be exported
6. If a distributed process is required (through events), consider using services such as AWS Step Functions (serverless state machine)
7. If using multiple transport strategies, consider using a facade pattern that wraps all transport strategies
8. A transport strategy handler should isolate each request with concurrency techniques