package todo_test

import (
	"clean-architecture/domain/models"
	"clean-architecture/domain/todo"
	"clean-architecture/pkg/infrastructure"
	"clean-architecture/pkg/types"
	"clean-architecture/testutil"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/fx"
)

var _ = Describe("Domain/Todo/Service", Ordered, func() {
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
		ok := errors.Is(err, todo.ErrTodoNotFound)
		log.Println("===Error:===", err, ok)
		Expect(ok).To(BeTrue(), "Expected a TodoNotFoundError")
		Expect(todo.ErrTodoNotFound.Error()).To(ContainSubstring("Todo not found"))
	})

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

	It("should handle custom pagination limits", func() {
		// Act
		todos, _, err := todoService.List(1, 5)

		// Assert
		Expect(err).To(BeNil())
		Expect(len(todos)).To(BeNumerically("<=", 5))
	})
})
