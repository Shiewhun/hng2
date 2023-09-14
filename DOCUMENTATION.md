```markdown
### API Documentation

This document provides detailed documentation for the Go REST API for managing persons.

### Table of Contents

- [Standard Request and Response Formats](#standard-request-and-response-formats)
- [Sample API Usage](#sample-api-usage)
- [Known Limitations and Assumptions](#known-limitations-and-assumptions)
- [Setup and Deployment Instructions](#setup-and-deployment-instructions)

## Standard Request and Response Formats

### Request Format

#### Create a New Person (POST /api)

- Method: POST
- Endpoint: `/api`
- **Request Body:**

  ```json
  {
    "name": "John Doe"
  }
  ```

### Response Format

#### Successful Response (201 Created)

- Status Code: 201 Created
- Response Body:

  ```json
  {
    "id": 1,
    "name": "John Doe"
  }
  ```

#### Error Response (400 Bad Request)

- Status Code: 400 Bad Request
- Response Body:

  ```json
  {
    "error": "Invalid JSON payload"
  }
  ```

### Request Format

#### Retrieve Person Details (GET /api/{user_id})

- Method: GET
- Endpoint: `/api/{user_id}`

### Response Format

#### Successful Response (200 OK)

- Status Code: 200 OK
- Response Body:

  ```json
  {
    "id": 1,
    "name": "John Doe"
  }
  ```

#### Error Response (404 Not Found)

- Status Code: 404 Not Found
- Response Body:

  ```json
  {
    "error": "Person not found"
  }
  ```

## Sample API Usage

### Creating a New Person

**Request:**

```http
POST /api
Content-Type: application/json

{
  "name": "John Doe"
}
```

**Response (201 Created):**

```json
{
  "id": 1,
  "name": "John Doe"
}
```

### Retrieving Person Details

**Request:**

```http
GET /api/1
```

**Response (200 OK):**

```json
{
  "id": 1,
  "name": "John Doe"
}
```

## Known Limitations and Assumptions

- The API uses a CSV file (`persons.csv`) to store person records. Make sure the file has write permissions.
- Validations:
  - Name should contain only valid characters (letters and spaces).
- The API does not support bulk operations.

## Setup and Deployment Instructions

1. Clone the repository to your local machine:

   ```shell
   git clone https://github.com/shiewhun/hng2.git
   ```

2. Navigate to the project directory:

   ```shell
   cd hng2
   ```

3. Run the API:

   ```shell
   go run main.go
   ```

   By default, the API runs on port 8080. You can change the port in the `main.go` file.

4. Access the API at `http://localhost:8080`.

You can also deploy this API on a server following your server provider's deployment instructions. Make sure to update the server configuration and firewall rules as needed.

For more information, refer to the [README.md](./README.md) file.
```
