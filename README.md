# Gojot: A Simple Go Auth Middleware

Gojot is a lightweight library for generating JSON Web Tokens (JWTs) and authenticating in Go.  It's designed for ease of use and integration into existing Go projects with the [Chi](https://go-chi.io/#/) router.  Perfect for adding authentication to your applications quickly.

## Overview

* **JWT Generation:** Gojot leverages the standard library's `encoding/json` package for encoding the JWT payload and `crypto/rsa` for signing the token.  This ensures robust security and avoids external dependencies for core functionality.  You can learn more about JWTs [here](https://jwt.io/introduction/).  The RSA algorithm used is documented [here](https://developer.mozilla.org/en-US/docs/Web/API/SubtleCrypto/encrypt).
* **Testing:**  The project includes comprehensive unit tests using the Go testing framework (`go test`).  These tests cover the JWT generation process, ensuring reliability. See `token_test.go` and `gojot_test.go` for examples.
* **Error Handling:** Gojot employs explicit error handling to provide clear and informative error messages, aiding in debugging and maintenance.

This project uses only Go's standard library.  This simplifies dependencies and improves the portability of your application.

## Project Structure

```
├── token.go         // Core JWT generation logic
├── gojot.go         // Auth Middleware logic
├── go.mod           // Go module definition
├── LICENSE          // Project license
├── gojot_test.go    // Unit tests for the core logic
├── go.sum           // Go module checksums
├── README.md        // This file
├── token_test.go    // Unit tests for the token generation
└── _example/
    └── main.go     // Example usage of Gojot
```
The `_example` directory contains a simple `main.go` file that creates a simple server with public routes and protected routes with gojot middleware. Execute the example: `go run _example/main.go`

## License

This project is licensed under the MIT License - see the [LICENSE](./LICENSE) file for details.

