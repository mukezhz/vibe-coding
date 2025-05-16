# Testing Instructions for Go Clean Architecture

This document provides guidelines and examples for writing tests in the Go Clean Architecture project. We use Ginkgo, TestContainers, and various testing utilities to create comprehensive integration and unit tests.

## Testing Framework

- **Ginkgo**: BDD testing framework
- **Testify**: Assertion library
- **TestContainers**: For database integration testing
- **apitest**: HTTP testing utility

## Test Organization

Tests should be organized following the domain structure:

1. Create a test suite file for each domain (e.g., `todo_suite_test.go`)
2. Create specific test files for controllers, services, or routes (e.g., `route_test.go`)

## Creating a Test Suite

For each domain or package, create a test suite file:

```go
package todo_test // Note the _test suffix to avoid import cycles

import (
	"clean-architecture/pkg/utils"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTodo(t *testing.T) {
	utils.ChDir() // Ensures tests run from the project root
	RegisterFailHandler(Fail)
	RunSpecs(t, "Todo Suite")
}

// Define the test interface that will be used in test files
var t GinkgoTInterface
var _ = BeforeSuite(func() {
	t = GinkgoT()
})
```
## API Testing Example

Here's an example of how to test API endpoints using our architecture:

```go
package todo_test

import (
	"clean-architecture/domain/todo"
	"clean-architecture/pkg/infrastructure"
	"clean-architecture/pkg/responses"
	"clean-architecture/testutil"
	"encoding/json"
	"fmt"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/steinfletcher/apitest"
	"go.uber.org/fx"
)

var _ = Describe("Domain/Todo/Route", Ordered, func() {
	var (
		router      infrastructure.Router
		todoService *todo.Service
		todoRepo    *todo.Repository
	)

	BeforeAll(func() {
		setupDI := func() {
			err := testutil.DI(t,
				fx.Populate(&router),
				fx.Populate(&todoService),
				fx.Populate(&todoRepo),
			)
			if err != nil {
				t.Error(err)
			}
		}
		setupDI()
	})

	// Helper function for reuse
	createTodo := func(title string, description string) (string, error) {
		reqBody := fmt.Sprintf(`{"title": "%s", "description": "%s"}`, title, description)

		result := apitest.
			New().
			Handler(router).
			Post("/api/todos").
			Body(reqBody).
			Expect(t).
			Status(http.StatusCreated).
			End()

		var responseBody responses.DetailResponseType[todo.TodoResponse]
		if err := json.NewDecoder(result.Response.Body).Decode(&responseBody); err != nil {
			return "", err
		}

		return responseBody.Item.ID, nil
	}

	It("should create a new todo", func() {
		reqBody := `{"title": "Test Todo", "description": "Test Description"}`

		result := apitest.
			New().
			Handler(router).
			Post("/api/todos").
			Body(reqBody).
			Expect(t).
			Status(http.StatusCreated).
			End()

		response := result.Response
		var responseBody responses.DetailResponseType[todo.TodoResponse]

		err := json.NewDecoder(response.Body).Decode(&responseBody)
		Expect(err).To(BeNil())
		Expect(responseBody.Message).To(Equal("success"))
		Expect(responseBody.Item.Title).To(Equal("Test Todo"))
		Expect(responseBody.Item.Description).To(Equal("Test Description"))
		Expect(responseBody.Item.ID).NotTo(BeEmpty())
	})
})
```
## Testing Error Cases

It's important to test both happy paths and error cases. Here's how to test error scenarios:

```go
It("should return error when creating todo without title", func() {
    reqBody := `{"description": "Missing Title"}`

    result := apitest.
        New().
        Handler(router).
        Post("/api/todos").
        Body(reqBody).
        Expect(t).
        Status(http.StatusBadRequest).
        End()

    // Check for validation error in response
    response := result.Response
    Expect(response.StatusCode).To(Equal(http.StatusBadRequest))
})

It("should return error for non-existent todo ID", func() {
    // Use a non-existent UUID
    nonExistentID := "00000000-0000-0000-0000-000000000000"

    result := apitest.
        New().
        Handler(router).
        Get(fmt.Sprintf("/api/todos/%s", nonExistentID)).
        Expect(t).
        Status(http.StatusNotFound).
        End()

    response := result.Response
    Expect(response.StatusCode).To(Equal(http.StatusNotFound))
})

It("should return error for invalid todo ID format", func() {
    invalidID := "not-a-uuid"

    result := apitest.
        New().
        Handler(router).
        Get(fmt.Sprintf("/api/todos/%s", invalidID)).
        Expect(t).
        Status(http.StatusBadRequest).
        End()

    response := result.Response
    Expect(response.StatusCode).To(Equal(http.StatusBadRequest))
})
```

## Testing Pagination

For endpoints that support pagination, you should test different page sizes and page numbers:

```go
It("should list todos with pagination after creating multiple todos", func() {
	// Create test data
	for i := 1; i <= 15; i++ {
		_, err := createTodo(fmt.Sprintf("Pagination Todo %d", i), fmt.Sprintf("Description %d", i))
		Expect(err).To(BeNil())
	}

	// Test first page
	result := apitest.
		New().
		Handler(router).
		Get("/api/todos").
		Query("page", "1").
		Query("limit", "10").
		Expect(t).
		Status(http.StatusOK).
		End()

	var responseBody1 todo.TodoListResponse
	err := json.NewDecoder(result.Response.Body).Decode(&responseBody1)
	Expect(err).To(BeNil())
	Expect(responseBody1.Message).To(Equal("success"))
	Expect(len(responseBody1.Items)).To(Equal(10))
	Expect(responseBody1.Pagination.HasNext).To(BeTrue())

	// Test second page
	result2 := apitest.
		New().
		Handler(router).
		Get("/api/todos").
		Query("page", "2").
		Query("limit", "10").
		Expect(t).
		Status(http.StatusOK).
		End()

	var responseBody2 todo.TodoListResponse
	err = json.NewDecoder(result2.Response.Body).Decode(&responseBody2)
	Expect(err).To(BeNil())
	Expect(responseBody2.Message).To(Equal("success"))
	Expect(len(responseBody2.Items)).To(BeNumerically(">", 0))
	Expect(responseBody2.Pagination.HasNext).To(BeFalse())
})

It("should get a todo by ID", func() {
	todoID, err := createTodo("Get Todo Test", "For testing get by ID")
	Expect(err).To(BeNil())

	result := apitest.
		New().
		Handler(router).
		Get(fmt.Sprintf("/api/todos/%s", todoID)).
		Expect(t).
		Status(http.StatusOK).
		End()

	response := result.Response
	var responseBody responses.DetailResponseType[todo.TodoResponse]

	err = json.NewDecoder(response.Body).Decode(&responseBody)
	Expect(err).To(BeNil())
	Expect(responseBody.Message).To(Equal("success"))
	Expect(responseBody.Item.ID).To(Equal(todoID))
	Expect(responseBody.Item.Title).To(Equal("Get Todo Test"))
	Expect(responseBody.Item.Description).To(Equal("For testing get by ID"))
})
```

## Testing Database Operations

Our application uses TestContainers to create an isolated MySQL container for integration testing. The test utilities handle:

1. Starting a MySQL container
2. Setting up the database
3. Running migrations
4. Connecting to the database
5. Cleaning up resources after tests

You don't need to manage these directly - the `testutil.DI()` function sets up everything.

## Running Tests

To run all tests in the project:

```bash
go test ./... -v
```

For running tests with coverage:

```bash
ginkgo -v --cover -r ./domain/... ./pkg/...
```

To generate a coverage report:

```bash
go test ./... -v -coverprofile=cover.txt -coverpkg=./...
go tool cover -html=cover.txt -o coverage.html
```

## Best Practices

1. **Test Independence**: Each test should be independent and not rely on the state from previous tests
2. **Use Helper Functions**: Create helper functions for common operations
3. **Clean Test Data**: Always clean up test data after tests to prevent state leakage
4. **Use Descriptive Names**: Use descriptive names for your test cases with the `It()` function
5. **Organize with Describe/Context**: Use `Describe()` and `Context()` to organize tests logically
6. **Assertions**: Use Gomega's expressive assertions like `Expect(x).To(Equal(y))` for readable tests

## Common Testing Patterns

### 1. Testing API Endpoints

For API endpoints, use the `apitest` library to make requests to your handlers:

```go
result := apitest.
    New().
    Handler(router).
    Get("/api/todos").
    Query("page", "1").
    Query("limit", "10").
    Expect(t).
    Status(http.StatusOK).
    End()
```

### 2. Testing Database Operations

For database operations, inject the repository and test it directly:

```go
var _ = Describe("Todo Repository", func() {
    var repository todo.Repository
    
    BeforeSuite(func() {
        t = GinkgoT()
        setupDI := func() {
            err := testutil.DI(t,
                fx.Populate(&repository),
            )
            if err != nil {
                t.Error(err)
            }
        }
        setupDI()
    })
    
    It("should create a todo", func() {
        newTodo := &models.Todo{
            Title:       "Test Todo",
            Description: "Test Description",
        }
        
        createdTodo, err := repository.Create(newTodo)
        Expect(err).To(BeNil())
        Expect(createdTodo.ID).NotTo(BeEmpty())
        Expect(createdTodo.Title).To(Equal("Test Todo"))
    })
})
```

## Conclusion

Following these patterns will help maintain consistency across all tests in the project. For more examples, refer to the existing tests in the `domain/todo` package.
