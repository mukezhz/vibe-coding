package todo_test

import (
	"clean-architecture/domain/todo"
	"clean-architecture/pkg/infrastructure"
	"clean-architecture/pkg/responses"
	"clean-architecture/testutil"
	"encoding/json"
	"fmt"
	"log"
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

	It("should return empty list of todos", func() {
		expected := todo.TodoListResponse{
			Message: "success",
			Items:   []todo.TodoListItem{},
			Pagination: responses.PaginationResponseType{
				Total:   0,
				HasNext: false,
			},
		}
		expectedJSON, _ := json.MarshalIndent(expected, "", "  ")

		result := apitest.
			New().
			Handler(router).
			Get("/api/todos").
			Query("page", "1").
			Query("limit", "10").
			Expect(t).
			Body(string(expectedJSON)).
			Status(http.StatusOK).
			End()

		response := result.Response
		var responseBody todo.TodoListResponse
		if err := json.NewDecoder(response.Body).Decode(&responseBody); err != nil {
			t.Errorf("Failed to decode response: %v", err)
		}
		log.Printf("Server response: %+v\n", responseBody)
	})

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

	It("should update a todo", func() {
		todoID, err := createTodo("Update Todo Test", "Initial description")
		Expect(err).To(BeNil())

		updateReqBody := `{"title": "Updated Todo", "description": "Updated description"}`

		result := apitest.
			New().
			Handler(router).
			Put(fmt.Sprintf("/api/todos/%s", todoID)).
			Body(updateReqBody).
			Expect(t).
			Status(http.StatusOK).
			End()

		response := result.Response
		var responseBody responses.DetailResponseType[todo.TodoResponse]

		err = json.NewDecoder(response.Body).Decode(&responseBody)
		Expect(err).To(BeNil())
		Expect(responseBody.Message).To(Equal("success"))
		Expect(responseBody.Item.ID).To(Equal(todoID))
		Expect(responseBody.Item.Title).To(Equal("Updated Todo"))
		Expect(responseBody.Item.Description).To(Equal("Updated description"))
	})

	It("should partially update a todo with only title", func() {
		todoID, err := createTodo("Partial Update Test", "Original description")
		Expect(err).To(BeNil())

		updateReqBody := `{"title": "Only Title Updated"}`

		result := apitest.
			New().
			Handler(router).
			Put(fmt.Sprintf("/api/todos/%s", todoID)).
			Body(updateReqBody).
			Expect(t).
			Status(http.StatusOK).
			End()

		response := result.Response
		var responseBody responses.DetailResponseType[todo.TodoResponse]

		err = json.NewDecoder(response.Body).Decode(&responseBody)
		Expect(err).To(BeNil())
		Expect(responseBody.Message).To(Equal("success"))
		Expect(responseBody.Item.Title).To(Equal("Only Title Updated"))
		// Description should remain unchanged
		Expect(responseBody.Item.Description).To(Equal("Original description"))
	})

	It("should list todos with pagination after creating multiple todos", func() {
		for i := 1; i <= 15; i++ {
			_, err := createTodo(fmt.Sprintf("Pagination Todo %d", i), fmt.Sprintf("Description %d", i))
			Expect(err).To(BeNil())
		}

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
		// The second page should have the remaining todos (5 or more)
		Expect(responseBody2.Pagination.HasNext).To(BeFalse())
	})

	It("should handle custom pagination limits", func() {
		result := apitest.
			New().
			Handler(router).
			Get("/api/todos").
			Query("page", "1").
			Query("limit", "5").
			Expect(t).
			Status(http.StatusOK).
			End()

		var responseBody todo.TodoListResponse
		err := json.NewDecoder(result.Response.Body).Decode(&responseBody)
		Expect(err).To(BeNil())

		// Should have at most 5 items per page with this limit
		Expect(len(responseBody.Items)).To(BeNumerically("<=", 5))
	})
})
