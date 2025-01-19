# Gojot: A Simple Go Auth Middleware

Gojot is a lightweight library for generating JSON Web Tokens (JWTs) and authenticating in Go.  It's designed for ease of use and integration into existing Go projects with the [Chi](https://go-chi.io/#/) router.  Perfect for adding authentication to your applications quickly. It also provides a simple rate limiter implementation in Go, designed for microservices and APIs.

## Overview

* **JWT Generation:** Gojot leverages the [golang-jwt (v4)](github.com/golang-jw) library for encoding the JWT payload and for signing the token.
* **Token Bucket Algorithm:**  The rate limiter uses the token bucket algorithm for smooth and predictable rate limiting.
* **Testing:**  The project includes comprehensive unit tests using the Go testing framework (`go test`). They are included in `gojot_test.go`, `token_test.go`, and `rate_limiter_test.go` to ensure code correctness and reliability.
* **Error Handling:** Gojot employs explicit error handling to provide clear and informative error messages, aiding in debugging and maintenance.

## Project Structure

```
├── token.go                // Core JWT generation logic
├── gojot.go                // Auth Middleware logic
├── rate_limiter.go         // Auth Middleware logic
├── token_test.go           // Unit tests for the token generation
├── gojot_test.go           // Unit tests for the core logic
├── rate_limiter_test.go    // Unit tests for the rate limiter
├── go.mod                  // Go module definition
├── go.sum                  // Go module checksums
├── LICENSE                 // Project license
├── README.md               // This file
└── _example/
    └── main.go             // Example application demonstrating auth middleware and rate limiter usage
```
The `_example` directory contains a simple example application demonstrating how to use the auth middleware and the rate limiter. Execute the example: `go run _example/main.go`

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

## License

This project is licensed under the MIT License - see the [LICENSE](./LICENSE) file for details.

