# envgaurd
EnvGuard is a Go library for validating environment variables with type safety. It supports string, int, and bool types, provides default values, and offers clear error handling.

# EnvGuard - Go Environment Variable Validation Library

EnvGuard is a lightweight Go library designed to simplify the management and validation of environment variables in your Go applications. It provides an easy way to load and validate environment variables with type safety, ensuring your application is properly configured before starting.

## Features

- **Type-Safe Validation**: Automatically validates environment variables based on expected types (`string`, `int`, `bool`).
- **Custom Default Values**: Supports setting default values for environment variables if not found or if invalid.
- **Error Handling**: Provides clear error messages if an environment variable is missing or invalid.
- **Extensible**: Easily extendable to support other types or custom validation logic.

## Installation

To install `envguard`, use the following Go command:

```bash
go get github.com/Shobhitdimri01/envgaurd
