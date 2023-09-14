```markdown
# Go REST API for Managing Persons

This is a simple Go-based REST API that allows you to perform CRUD (Create, Read, Update, Delete) operations on a "person" resource. You can interact with this API to add, retrieve, update, or delete person records.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Getting Started](#getting-started)
  - [Installation](#installation)
  - [Running the API](#running-the-api)
- [API Endpoints](#api-endpoints)
- [Sample Requests and Responses](#sample-requests-and-responses)
- [Testing](#testing)
- [Database](#database)
- [Documentation](#documentation)

## Prerequisites

Before you begin, ensure you have met the following requirements:

- Go (Golang) is installed. You can download it [here](https://golang.org/dl/).
- Git is installed. You can download it [here](https://git-scm.com/downloads).
- A text editor or IDE of your choice.

## Getting Started

### Installation

1. Clone this repository to your local machine using:

   ```shell
   git clone https://github.com/shiewhun/hng2.git
   ```

2. Change the directory to the project folder:

   ```shell
   cd hng2
   ```

### Running the API

1. Start the API by running:

   ```shell
   go run main.go
   ```

   The API will start on port 8080 by default. You can change the port in the `main.go` file if needed.

2. Your API is now running and accessible at `http://localhost:8080`.

## API Endpoints

- **GET /api/{user_id}**: Retrieve details of a person by ID.
- **POST /api**: Create a new person record.
- **PUT /api/{user_id}** or **PATCH /api/{user_id}**: Update details of an existing person by ID.
- **DELETE /api/{user_id}**: Delete a person record by ID.

## Sample Requests and Responses

### Create a New Person (POST /api)

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

### Retrieve Person Details (GET /api/{user_id})

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

### Update Person Details (PUT /api/{user_id} or PATCH /api/{user_id})

**Request:**

```http
PUT /api/1
Content-Type: application/json

{
  "name": "Updated Name"
}
```

**Response (200 OK):**

```json
{
  "message": "Person updated successfully"
}
```

### Delete a Person (DELETE /api/{user_id})

**Request:**

```http
DELETE /api/1
```

**Response (200 OK):**

```json
{
  "message": "Person deleted successfully"
}
```

## Testing

You can test the API using tools like Postman or by writing scripts in your preferred programming language. Ensure you test all CRUD operations with valid and invalid inputs.

## Database

This API uses a CSV file (`persons.csv`) to store person records. The CSV file is created in the project directory when you run the API for the first time. Make sure the file has write permissions.

## Documentation

For detailed API documentation, please refer to [DOCUMENTATION.md](./DOCUMENTATION.md).
