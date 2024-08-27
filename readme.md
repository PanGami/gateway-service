# Technical Test DOT - [gateway-service](https://github.com/PanGami/gateway-service)

## Project Structure

```
technical-test-DOT/
│
├── docker-compose.yml
├── readme.md
│
├── gateway-service/
│   ├── Dockerfile
│   ├── .env
│   ├── go.mod
│   ├── readme.md  <-- This File Should be here
│   └── (other source files)
│
└── user-service/
    ├── Dockerfile
    ├── .env
    ├── go.mod
    ├── readme.md
    └── (other source files)
```

## How to Run

To start the services, navigate to the project directory and run the following command:

```bash
cd path/to/technical-test-DOT
docker-compose up --build -d
```

To run the services without rebuilding:

```bash
docker-compose up -d
```

To rebuild the services:

```bash
docker-compose up --build
```

## Project Overview

This project utilizes **Golang** as its main programming language, structured around a microservices architecture using **gRPC**, which is 7-10 times faster than traditional REST APIs. The entire setup runs within a Docker environment.

### Key Features

1. **API Endpoints**: The project supports various HTTP methods including GET, POST, PUT, PATCH, and DELETE.

2. **Database Interaction**: The application uses **GORM** for ORM, with functionality to save data from APIs into a database. There is one relational connection between two tables, and the data-saving process includes transactions for consistency.

3. **Centralized Error Handling**: Custom error types and middleware are implemented for consistent error responses.

4. **Caching**: Implemented using **Redis** to enhance performance by caching user details.

5. **API Testing**: A Postman collection and environment are available for testing purposes.

## API Endpoints

### GET Method

- List Users: `{{host}}/user/list?page=1&page_size=5`
- User Details: `{{host}}/user/detail?id=1`
- User Activity: `{{host}}/user/activity?id=1`

### POST Method

- Create User: `{{host}}/user/create`
  
  ```json
  {
      "username": "panca",
      "full_name": "panca",
      "password": "password"
  }
  ```

### PUT Method

- Update User: `{{host}}/user/update?id=1`
  
  ```json
  {
      "username": "pan",
      "full_name": "pancaUpdated",
      "password": "password"
  }
  ```

### PATCH Method

- Patch User: `{{host}}/user/patch?id=111`
  
  ```json
  {
      "username": "pan",
      "full_name": "pancaUpdated",
      "password": "password"
  }
  ```

### DELETE Method

- Delete User: `{{host}}/user/delete?id=4`

## Detailed Implementation

### Database Operations

- **Transactional Handling**: The `Handler` method in `CreateUser` action uses a GORM transaction to save both `User` and `NewTable` entities. If any operation within the transaction fails, the entire transaction is rolled back, ensuring data consistency.

  ```go
  func (a *CreateUser) Handler(ctx context.Context, user *entity.User, newTable *entity.NewTable) error {
      return a.repoUser.DB.Transaction(func(tx *gorm.DB) error {
          if err := tx.Create(user).Error; err != nil {
              return err
          }
          newTable.UserID = user.ID
          if err := tx.Create(newTable).Error; err != nil {
              return err
          }
          return nil
      })
  }
  ```

- **Relational Data**: The `NewTable` entity has a foreign key (`UserID`) that references the `User` entity, establishing a one-to-many relationship.

### Error Handling (Centered Error Handling)

- **Custom Error Types**: Defined in `gateway-service/util/errors/error.go`

  ```go
  func Wrap(err error, code codes.Code, message string) *AppError {
      return &AppError{
          Code:    code,
          Message: message,
          Err:     err,
      }
  }
  
  func ErrNotFound(model string) error {
      return &AppError{
          Code:    codes.NotFound,
          Message: fmt.Sprintf("%s not found", model),
      }
  }
  
  func ErrBadRequest(msg string) error {
      return &AppError{
          Code:    codes.InvalidArgument,
          Message: msg,
      }
  }
  ```

- **Middleware for Error Handling**: Implemented in `gateway-service/middleware/error_handler.go`

  ```go
  func ErrorHandlerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
      return func(c echo.Context) error {
          err := next(c)
          if err != nil {
              appErr, ok := err.(*errors.AppError)
              if ok {
                  return c.JSON(int(appErr.Code), map[string]string{"error": appErr.Message})
              }
              return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
          }
          return nil
      }
  }
  ```

### Caching with Redis

Example usage in `user-service/repo/user.go`:

```go
func SaveUserToCache(ctx context.Context, redisClient *redis.Client, userDetail *entity.User) error {
    cacheKey := "user_detail_" + string(userDetail.ID)
    userData, err := json.Marshal(userDetail)
    if err != nil {
        return err
    }

    err = redisClient.Set(ctx, cacheKey, userData, 5*time.Minute).Err()
    if err != nil {
        return err
    }

    return nil
}
```

## Testing

- **API Testing**: You can use the provided Postman collection and environment files located in the project directory to test the API endpoints.

---