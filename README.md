# Go HTTP CRUD Web Application with PostgreSQL

This is a simple CRUD (Create, Read, Update, Delete) web application built in Go that uses a PostgreSQL database to store and manage user data. You can use this as a starting point for building web applications in Go that require database interactions.

## Prerequisites

Before you get started, ensure you have the following prerequisites installed on your system:

- Go (1.20 or higher)
- PostgreSQL
- `github.com/lib/pq` Go package for PostgreSQL driver (you can install it using `go get github.com/lib/pq`)

## Configuration

In the `main.go` file, you'll find the following constants that you can configure according to your environment:

```go
const (
    DB_USER     = "postgres" // Your PostgreSQL username
    DB_PASSWORD = "password" // Your PostgreSQL password
    DB_NAME     = "godb"     // Your PostgreSQL database name
    DB_IP       = "localhost" // Your PostgreSQL server IP
)
```

Ensure that you replace these values with your PostgreSQL server configuration.

## Database Setup

This application automatically creates a `users` table in the PostgreSQL database if it doesn't exist. You can find the table creation code in the `main` function.

```go
if _, err := db.Exec(`
    CREATE TABLE IF NOT EXISTS users(
        id SERIAL PRIMARY KEY,
        name VARCHAR(50) NOT NULL,
        email VARCHAR(50) UNIQUE NOT NULL,
        sex INT NOT NULL CHECK (sex IN (1, 2)),
        interest TEXT
    );
`); err != nil {
    log.Fatal(err)
}
```

## Running the Application

To run the application, execute the following command in your terminal:

```shell
go run main.go
```

The application will start and listen on `http://localhost:9000`. You can access it using your web browser.

## Usage

### Adding a User

1. Access the application in your web browser (`http://localhost:9000`).
2. Click on the "Add User" link in the navigation bar.
3. Fill in the user details (Name, Email, Sex, and Interest) in the form.
4. Click the "Submit" button to add the user to the database.

### Viewing Users

1. Access the application in your web browser (`http://localhost:9000`).
2. You will see a list of all users in the database displayed on the main page.

### Updating a User

1. Access the application in your web browser (`http://localhost:9000`).
2. Click on the "Update User" link in the navigation bar.
3. Enter the User ID of the user you want to update.
4. Modify the user details in the form.
5. Click the "Submit" button to update the user's information.

### Deleting a User

1. Access the application in your web browser (`http://localhost:9000`).
2. Click on the "Delete User" link in the navigation bar.
3. Enter the User ID of the user you want to delete.
4. Click the "Delete" button to remove the user from the database.

## Important Notes

- This application does not implement user authentication and authorization, so it is not suitable for production use without adding security features.
- It's recommended to implement proper error handling and validation for production applications.
- Make sure you secure your database credentials and connection details in a production environment.

Feel free to modify and extend this code to suit your specific needs for web application development in Go.
