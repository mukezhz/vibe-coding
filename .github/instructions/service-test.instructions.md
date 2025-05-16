# Service Testing Instructions for Go Clean Architecture

This document provides guidelines and examples for writing service tests in the Go Clean Architecture project. Service tests focus on testing the business logic layer in isolation from HTTP concerns.

## Testing Framework

- **Ginkgo**: BDD testing framework
- **Testify**: Assertion library
- **TestContainers**: For database integration testing

## Test Organization

Service tests should focus on the business logic and interaction with the repository layer:

1. Create a test file named `service_test.go` in the domain directory
2. Use a separate test package (`package todo_test` instead of `package todo`) to avoid circular imports
3. Tests should cover all service methods and business logic cases

## Setting Up Service Tests

Here's how to set up a service test correctly:

```go
package todo_test

import (
	"clean-architecture/domain/models"
	"clean-architecture/domain/todo"
	"clean-architecture/pkg/types"
	"clean-architecture/testutil"
	"errors"
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/fx"
)

var _ = Describe("Domain/Todo/Service", Ordered, func() {
	var (
		todoService *todo.Service
		todoRepo    *todo.Repository
	)

	BeforeAll(func() {
		setupDI := func() {
			err := testutil.DI(t,
				fx.Populate(&todoService),
				fx.Populate(&todoRepo),
			)
			if err != nil {
				t.Error(err)
			}
		}
		setupDI()
	})

	// Helper function to create test data
	createTestTodo := func(title string, description string) (*models.Todo, error) {
		newTodo := &models.Todo{
			ID:          types.ParseUUID(uuid.New().String()),
			Title:       title,
			Description: description,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		err := todoService.Create(newTodo)
		return newTodo, err
	}

	// Test cases follow...
})
```

## Testing Service Methods

### Create Method

```go
It("should create a new todo", func() {
	// Arrange
	title := "Test Todo"
	description := "Test Description"

	// Act
	todo, err := createTestTodo(title, description)

	// Assert
	Expect(err).To(BeNil())
	Expect(todo.ID).NotTo(BeEmpty())
	Expect(todo.Title).To(Equal(title))
	Expect(todo.Description).To(Equal(description))
})
```

### GetByID Method

```go
It("should get a todo by ID", func() {
	// Arrange
	title := "Get Todo Test"
	description := "For testing get by ID"
	createdTodo, err := createTestTodo(title, description)
	Expect(err).To(BeNil())

	// Act
	todo, err := todoService.GetByID(createdTodo.ID)

	// Assert
	Expect(err).To(BeNil())
	Expect(todo.ID.String()).To(Equal(createdTodo.ID.String()))
	Expect(todo.Title).To(Equal(title))
	Expect(todo.Description).To(Equal(description))
})

It("should return error for non-existent todo ID", func() {
	// Arrange
	nonExistentID := types.ParseUUID(uuid.New().String())

	// Act
	_, err := todoService.GetByID(nonExistentID)

	// Assert
	Expect(err).NotTo(BeNil())
	Expect(errors.Is(err, todo.ErrTodoNotFound)).To(BeTrue(), "Expected a TodoNotFoundError")
})
```

### Update Method

```go
It("should update a todo", func() {
	// Arrange
	createdTodo, err := createTestTodo("Update Todo Test", "Initial description")
	Expect(err).To(BeNil())

	// Act
	createdTodo.Title = "Updated Todo"
	createdTodo.Description = "Updated description"
	err = todoService.Update(createdTodo)

	// Assert
	Expect(err).To(BeNil())

	// Verify the update by fetching the todo again
	updatedTodo, err := todoService.GetByID(createdTodo.ID)
	Expect(err).To(BeNil())
	Expect(updatedTodo.Title).To(Equal("Updated Todo"))
	Expect(updatedTodo.Description).To(Equal("Updated description"))
})

It("should partially update a todo with only title", func() {
	// Arrange
	originalDescription := "Original description"
	createdTodo, err := createTestTodo("Partial Update Test", originalDescription)
	Expect(err).To(BeNil())

	// Act
	createdTodo.Title = "Only Title Updated"
	// Keep the description the same
	err = todoService.Update(createdTodo)

	// Assert
	Expect(err).To(BeNil())

	// Verify the update by fetching the todo again
	updatedTodo, err := todoService.GetByID(createdTodo.ID)
	Expect(err).To(BeNil())
	Expect(updatedTodo.Title).To(Equal("Only Title Updated"))
	// Description should remain unchanged
	Expect(updatedTodo.Description).To(Equal(originalDescription))
})
```

### Pagination and List Method

```go
It("should list todos with pagination", func() {
	// Arrange - Create multiple todos
	for i := 1; i <= 15; i++ {
		_, err := createTestTodo(
			"Pagination Todo "+GinkgoT().Name(), 
			"Description for pagination test",
		)
		Expect(err).To(BeNil())
	}

	// Act - Get first page with 10 items
	todos, total, err := todoService.List(1, 10)

	// Assert
	Expect(err).To(BeNil())
	Expect(len(todos)).To(Equal(10))
	Expect(total).To(BeNumerically(">=", 15))

	// Act - Get second page with remaining items
	todosPage2, totalPage2, err := todoService.List(2, 10)

	// Assert
	Expect(err).To(BeNil())
	Expect(len(todosPage2)).To(BeNumerically(">=", 5))
	Expect(totalPage2).To(Equal(total))
})
```

## Best Practices for Service Tests

1. **Test Independence**: Each test should be independent and not rely on the state from previous tests
2. **Error Handling**: Test both success and error cases for all service methods
3. **Validation**: Test all validation logic and business rules in the service
4. **DB Interaction**: Use real database connections with transactions that rollback after each test
5. **Arrange-Act-Assert**: Structure tests with clear arrangement, action, and assertion phases
6. **Descriptive Names**: Use descriptive test names that explain the behavior being tested

## Testing Edge Cases

Remember to test edge cases such as:

1. Empty or invalid input data
2. Non-existent records
3. Boundary conditions in pagination
4. Authorization and permission checks
5. Field validation rules

## Running Service Tests

To run just the service tests:

```bash
ginkgo -v -focus="Domain/Todo/Service" ./domain/todo/
```

For running with coverage:

```bash
ginkgo -v --cover -focus="Domain/Todo/Service" ./domain/todo/
```
