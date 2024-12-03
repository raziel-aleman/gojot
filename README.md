# Gojot: A Simple Go-based JWT middleware

Gojot is simple JWT package for use in personal projects. It provides a basic JWT authentication middleware. This release focuses on core functionality and serves as a foundation for future feature additions.


## Project Structure

```
├── token.go         // Handles tokenization (if applicable, details in the specific file's documentation)
├── gojot.go         // Main middleware logic.
├── go.mod           // Go module definition file.
├── go.sum           // Go checksums.
├── LICENSE          // License information.
├── README.md        // This file.
├── _example         // Example usage of the Gojot package.
│   └── main.go      // Entry point for the example application.
```

The `_example` directory contains a simple example demonstrating how to use the `gojot` package. This is a great starting point to understand how to integrate `gojot` into your own applications.


## How to Use

1. **Clone the repository:** `git clone https://github.com/raziel-aleman/gojot.git`
2. **Navigate to the directory:** `cd gojot`
3. **Run the example:** `go run _example/main.go`


This will start a basic Go application in localhost:8080. Note that this is a minimal example and the features may be limited at this stage. Further development will add more advanced features.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.


## License

This project is licensed under the [MIT License](LICENSE).

