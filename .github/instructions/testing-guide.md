# Testing Guide for Go Clean Architecture

This comprehensive guide provides instructions for writing tests in our Go Clean Architecture project, covering both service and API endpoint testing.

## Testing Framework and Tools

- **Ginkgo**: BDD testing framework for Go
- **Testify**: Assertion library
- **TestContainers**: For database integration testing
- **apitest**: HTTP testing utility

## Test Organization

Tests should follow the domain structure:

1. **Test Suite File**: Create a test suite file for each domain (e.g., `todo_suite_test.go`)
2. **Component Tests**: Create specific test files for components:
   - `service_test.go`: For testing business logic
   - `route_test.go`: For testing API endpoints
   - `repository_test.go`: For testing data access (optional)

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

## Service Testing

Service tests focus on business logic in isolation from HTTP concerns:

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

	It("should create a new todo", func() {
		title := "Test Todo"
		description := "Test Description"
		todo, err := createTestTodo(title, description)

		Expect(err).To(BeNil())
		Expect(todo.ID).NotTo(BeEmpty())
		Expect(todo.Title).To(Equal(title))
		Expect(todo.Description).To(Equal(description))
	})

	It("should get a todo by ID", func() {
		title := "Get Todo Test"
		todo, err := createTestTodo(title, "Description")
		Expect(err).To(BeNil())

		// Test the GetByID method
		result, err := todoService.GetByID(todo.ID)
		Expect(err).To(BeNil())
		Expect(result.ID).To(Equal(todo.ID))
		Expect(result.Title).To(Equal(title))
	})

	It("should return error for non-existent todo", func() {
		nonExistentID := types.ParseUUID(uuid.New().String())
		_, err := todoService.GetByID(nonExistentID)
		Expect(err).NotTo(BeNil())
		// Check if error is the expected type
	})

	It("should update todo fields", func() {
		todo, err := createTestTodo("Original Title", "Original Description")
		Expect(err).To(BeNil())

		newTitle := "Updated Title"
		updated, err := todoService.Update(todo.ID, func(t *models.Todo) error {
			t.Title = newTitle
			return nil
		})

		Expect(err).To(BeNil())
		Expect(updated.Title).To(Equal(newTitle))
		Expect(updated.Description).To(Equal(todo.Description)) // Unchanged
	})
})
```

## API Endpoint Testing

For testing API endpoints, use the `apitest` library:

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
	)

	BeforeAll(func() {
		setupDI := func() {
			err := testutil.DI(t,
				fx.Populate(&router),
				fx.Populate(&todoService),
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

		var response map[string]interface{}
		if err := json.Unmarshal(result.Response.Body.Bytes(), &response); err != nil {
			return "", err
		}

		item, ok := response["item"].(map[string]interface{})
		if !ok {
			return "", fmt.Errorf("item not found in response")
		}

		id, ok := item["id"].(string)
		if !ok {
			return "", fmt.Errorf("id not found in item")
		}

		return id, nil
	}

	It("should create a todo via API", func() {
		title := "API Test Todo"
		description := "Created via API test"
		reqBody := fmt.Sprintf(`{"title": "%s", "description": "%s"}`, title, description)

		apitest.
			New().
			Handler(router).
			Post("/api/todos").
			Body(reqBody).
			Expect(t).
			Status(http.StatusCreated).
			Assert(func(res *http.Response, req *http.Request) error {
				var response map[string]interface{}
				if err := json.Unmarshal(res.Body.Bytes(), &response); err != nil {
					return err
				}

				item, ok := response["item"].(map[string]interface{})
				if !ok {
					return fmt.Errorf("item not found in response")
				}

				if item["title"] != title {
					return fmt.Errorf("expected title %s, got %s", title, item["title"])
				}

				if item["description"] != description {
					return fmt.Errorf("expected description %s, got %s", description, item["description"])
				}

				return nil
			}).
			End()
	})

	It("should get a todo by ID via API", func() {
		title := "Get API Test Todo"
		description := "For testing Get endpoint"
		id, err := createTodo(title, description)
		Expect(err).To(BeNil())

		apitest.
			New().
			Handler(router).
			Get(fmt.Sprintf("/api/todos/%s", id)).
			Expect(t).
			Status(http.StatusOK).
			Assert(func(res *http.Response, req *http.Request) error {
				var response map[string]interface{}
				if err := json.Unmarshal(res.Body.Bytes(), &response); err != nil {
					return err
				}

				item, ok := response["item"].(map[string]interface{})
				if !ok {
					return fmt.Errorf("item not found in response")
				}

				if item["id"] != id {
					return fmt.Errorf("expected id %s, got %s", id, item["id"])
				}

				if item["title"] != title {
					return fmt.Errorf("expected title %s, got %s", title, item["title"])
				}

				return nil
			}).
			End()
	})

	It("should return 404 for non-existent todo", func() {
		nonExistentID := "00000000-0000-0000-0000-000000000000"
		apitest.
			New().
			Handler(router).
			Get(fmt.Sprintf("/api/todos/%s", nonExistentID)).
			Expect(t).
			Status(http.StatusNotFound).
			End()
	})

	It("should update a todo via API", func() {
		id, err := createTodo("Original API Title", "Original API Description")
		Expect(err).To(BeNil())

		updatedTitle := "Updated API Title"
		reqBody := fmt.Sprintf(`{"title": "%s"}`, updatedTitle)

		apitest.
			New().
			Handler(router).
			Put(fmt.Sprintf("/api/todos/%s", id)).
			Body(reqBody).
			Expect(t).
			Status(http.StatusOK).
			Assert(func(res *http.Response, req *http.Request) error {
				var response map[string]interface{}
				if err := json.Unmarshal(res.Body.Bytes(), &response); err != nil {
					return err
				}

				item, ok := response["item"].(map[string]interface{})
				if !ok {
					return fmt.Errorf("item not found in response")
				}

				if item["title"] != updatedTitle {
					return fmt.Errorf("expected title %s, got %s", updatedTitle, item["title"])
				}

				return nil
			}).
			End()
	})

	It("should list todos with pagination", func() {
		// Create a few todos first
		for i := 0; i < 3; i++ {
			_, err := createTodo(fmt.Sprintf("List Test Todo %d", i), "For testing list endpoint")
			Expect(err).To(BeNil())
		}

		apitest.
			New().
			Handler(router).
			Get("/api/todos?page=1&limit=10").
			Expect(t).
			Status(http.StatusOK).
			Assert(func(res *http.Response, req *http.Request) error {
				var response map[string]interface{}
				if err := json.Unmarshal(res.Body.Bytes(), &response); err != nil {
					return err
				}

				items, ok := response["items"].([]interface{})
				if !ok {
					return fmt.Errorf("items not found in response")
				}

				if len(items) < 3 {
					return fmt.Errorf("expected at least 3 items, got %d", len(items))
				}

				page, ok := response["page"].(map[string]interface{})
				if !ok {
					return fmt.Errorf("page not found in response")
				}

				if _, ok := page["total"]; !ok {
					return fmt.Errorf("total not found in page")
				}

				if _, ok := page["has_next"]; !ok {
					return fmt.Errorf("has_next not found in page")
				}

				return nil
			}).
			End()
	})
})
```

## Best Practices

1. **Use Descriptive Test Names**: Make test names descriptive of what they're testing
2. **Test Each Layer Independently**: Test service layer business logic separately from API endpoints
3. **Dependency Injection**: Use DI for test setup to get real instances of components
4. **Test Data Management**: Create helper functions for test data creation and cleanup
5. **Test Common Flows**: Test both happy paths and error cases
6. **Isolation**: Each test should be independent and not depend on the state from other tests
7. **Cleanup**: Clean up test data after tests where appropriate
