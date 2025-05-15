package todo_test

import (
	"clean-architecture/domain/todo"
	"clean-architecture/pkg/infrastructure"
	"clean-architecture/pkg/responses"
	"clean-architecture/testutil"
	"encoding/json"
	"log"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	_ "github.com/onsi/gomega"
	"github.com/steinfletcher/apitest"
	"go.uber.org/fx"
)

var _ = Describe("Domain/Todo/Route", func() {
	var (
		t      GinkgoTInterface
		router infrastructure.Router
	)

	setupDI := func() {
		err := testutil.DI(t,
			fx.Populate(&router),
		)
		if err != nil {
			t.Error(err)
		}
	}

	BeforeEach(func() {
		t = GinkgoT()
	})

	It("should return empty list of todos", func() {
		setupDI()
		expected := todo.TodoListResponse{
			Message: "success",
			Items:   []todo.TodoListItem{},
			Pagination: responses.PaginationResponseType{
				Total:   0,
				HasNext: false,
			},
		}
		expectedJSON, _ := json.MarshalIndent(expected, "", "  ")
		log.Println(string(expectedJSON))

		result := apitest.
			New().
			Handler(router).
			Get("/api/todos").
			Query("page", "1").
			Query("limit", "10").
			Expect(t).
			Status(http.StatusOK).
			End()

		response := result.Response
		var responseBody todo.TodoListResponse
		if err := json.NewDecoder(response.Body).Decode(&responseBody); err != nil {
			t.Errorf("Failed to decode response: %v", err)
		}
		log.Printf("Server response: %+v\n", responseBody)
	})
})
