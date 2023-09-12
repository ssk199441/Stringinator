# API Design Document: Stringinator

## API Overview

- **API Name:** Stringinator
- **API Version:** 1.0
- **Authors:** Satheesh Kumar
- **Date Created:** 12 September 2023
- **Last Updated:** 12 September 2023
- **Status:** Stable
- **Change Log:** N/A

## Introduction

The Stringinator API is designed to provide various string manipulation and statistics gathering functionalities. It allows users to analyze, transform, and retrieve statistics about strings. This document outlines the architecture, endpoints, and usage guidelines for the API.

## Use Cases

The primary use cases for the Stringinator API include:

1. Analyzing input strings and extracting information such as length, most frequent character, and character count.
2. Transforming input strings into different formats (uppercase, lowercase, title case).
3. Retrieving statistics about all strings processed by the server, including the longest and most popular strings.

## API Architecture

The Stringinator API is built using the Echo framework for handling HTTP requests. It utilizes the BoltDB database for storing statistics data. The API includes the following components:

- **Endpoints:** API endpoints for string analysis, transformation, and statistics retrieval.
- **Validation:** Input data validation using the Go Playground Validator library.
- **Middleware:** Middleware functions for request validation, response validation, and storing statistics.
- **Data Models:** Data structures for input and response data.
- **BoltDB:** A key-value store for persisting statistics data.

## API Endpoints

### 1. Analyze String

- **Endpoint URL:** `/stringinate`
- **HTTP Method:** `POST`
- **Description:** Analyze a string and retrieve information about it.
- **Request Parameters:**
  - `input` (query): The input string to be analyzed.
- **Request Headers:** N/A
- **Request Body:** N/A
- **Response:**
  - Status Code: `200 OK`
  - Response Body: JSON containing analyzed string data (StringData).

### 2. Transform String

- **Endpoint URL:** `/transform`
- **HTTP Method:** `POST`
- **Description:** Transform a string into a specified format (uppercase, lowercase, title case).
- **Request Body:**
  - `text` (JSON): The input string to be transformed.
  - `transformation` (JSON): The desired transformation type.
- **Response:**
  - Status Code: `200 OK`
  - Response Body: JSON containing original and transformed text (OriginalText, TransformedText).

### 3. Reset Statistics

- **Endpoint URL:** `/reset-stats`
- **HTTP Method:** `GET`
- **Description:** Reset all statistics data, both in-memory and in the database.
- **Response:**
  - Status Code: `200 OK`
  - Response Body: JSON with a success message (message).

### 4. Get Statistics

- **Endpoint URL:** `/stats`
- **HTTP Method:** `GET`
- **Description:** Retrieve statistics about all processed strings, including the most popular and longest strings.
- **Response:**
  - Status Code: `200 OK`
  - Response Body: JSON containing statistics data (StatsData).

## Authentication

The Stringinator  API does not require authentication for access. It is open for public use.

## Error Handling

Errors in the API are handled by returning appropriate HTTP status codes and error messages in the response body. Common status codes include:

- `400 Bad Request`: Invalid input or request.
- `500 Internal Server Error`: Server-side errors.

## Data Models

### StringData

```json
{
  "input": "banana",
  "length": 6,
  "most_frequent": "a",
  "frequent_count": 3
}
```

### StatsData

```json
{
  "inputs": {
        "The quick brown fox jumps over a lazy dog": 1,
        "The quick brown fox jumps over the lazy dog. Pack my box with five dozen liquor jugs": 1,
        "banana": 1
},
  "most_popular": "The quick brown fox jumps over the lazy dog. Pack my box with five dozen liquor jugs",
  "longest_input": "The quick brown fox jumps over the lazy dog. Pack my box with five dozen liquor jugs",
  "longest_input_len": 84
}
```
### Transform

```json
{
  "original_text": "hello world",
  "transformed_text": "HELLO WORLD"
}

```

## Versioning

The Stringinator  API follows a simple versioning strategy, starting with version 1.0. Changes to the API will be managed by incrementing the version number and documenting changes in the change log.

## Security

Security measures are implemented to protect the API and user data:

- Input validation to protect against malformed requests.
- Proper error handling to avoid exposing sensitive information.
- Regular database backups to prevent data loss.

## Testing

The API is thoroughly tested with unit tests, integration tests, and sample requests. Testing ensures the reliability and correctness of the API's functionality.


## Usage Guidelines

- Use appropriate endpoints for string analysis and transformation.
- Respect rate limits to avoid API throttling.
- Handle errors gracefully and display user-friendly error messages.
- Follow best practices for string manipulation.

## References

- [Echo Framework Documentation](https://echo.labstack.com/)
- [BoltDB Documentation](https://pkg.go.dev/github.com/boltdb/bolt)
- [Go Playground Validator Documentation](https://pkg.go.dev/github.com/go-playground/validator/v10)

## Conclusion

Developers can utilize this API to build applications that require string manipulation and analysis functionalities. This design document serves as a reference for understanding and using the API effectively.

