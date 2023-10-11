# Dev Challenge 2023 Backend / Online Round

## How to run?

Run this command from a root of the project directory:

```sh
docker compose up --build
# or 
make container
```

After the container has been started the `Starting server...` message should be printed to the console which means you're good to *go*.

By default the server will be listening on `http://localhost:8080`.

## API Endpoints

```
[GET]   /api/v1/:sheet_id            // get an array of cells by sheet id

[GET]   /api/v1/:sheet_id/:cell_id   // get a cell by sheet and cell ids

[POST]  /api/v1/:sheet_id/:cell_id   // create/update a cell
```

## Tests

This project includes intergration and unit tests.
The integration tests cover app use cases such as creating and accessing cells etc.
The unit tests cover critical app services such as parser and evaluator.

To run all tests inside a docker container run:

```sh
docker compose -f ./docker-compose.test.yml up --build
# or
make container-tests
```

## Data persistence

The project uses `postgres` to store data, therefore I've added a `github.com/lib/pq` driver as an essential dependency.

## Thoughts about my choises

The programming language. I've chosen golang because it ideally suits for a web service. Faster than any of existing mature javascript runtimes, but not as complicated as languages with manual memory management.

The architecture. The application is divided into four main parts.

1. Front-end (routing and handlers) layer which is responsible for processing response and returning request, including input validation and json encoding.
2. Database (or repository) is a simple data accessing layer.
3. Service layer is responsible for business logic. It provides an interface to the front-end and repository layers to make the app more flexible.
4. Core is a parser which parses input formula into an abstract syntax tree and evaluator which takes an AST as an input and returns a result.

## Possible improvements

- Adjust parse and evaluation functions to use goroutines (concurrency) to make the app faster.
- With the project growth, implement better dependency injection mechanism e.g. DI Container pattern.
- Introduce support for string concatenation and other arithmetical operations like exponentiation, modulus etc.
