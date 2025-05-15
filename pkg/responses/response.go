package responses

import (
	"clean-architecture/pkg/framework"

	"github.com/gin-gonic/gin"
)

type DetailResponseType[T any] struct {
	Item    T      `json:"item"`
	Message string `json:"message,omitempty"`
}

type ErrorResponseType struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type PaginationResponseType struct {
	Total   int64 `json:"total"`
	HasNext bool  `json:"has_next"`
}

type ListResponseType[T any] struct {
	Items      []T                    `json:"items"`
	Message    string                 `json:"message,omitempty"`
	Pagination PaginationResponseType `json:"pagination,omitempty"`
}

// MessageOnlyResponse: json response function
func MessageOnlyResponse(ctx *gin.Context, statusCode int, message string) {
	ctx.JSON(statusCode, gin.H{
		"message": message,
	})
}

// DetailResponse: json response function
func DetailResponse[T any](ctx *gin.Context, statusCode int, response DetailResponseType[T]) {
	ctx.JSON(statusCode, response)
}

// ErrorResponse: json error response function
func ErrorResponse(ctx *gin.Context, statusCode int, response ErrorResponseType) {
	ctx.JSON(statusCode, response)
}

// ListResponse: json response function
func ListResponse[T any](ctx *gin.Context, statusCode int, response ListResponseType[T]) {
	ctx.JSON(statusCode, response)
}

// JSONWithPagination : json response function
func JSONWithPagination[T any](ctx *gin.Context, statusCode int, response ListResponseType[T]) {
	limit, _ := ctx.MustGet(framework.Limit).(int64)
	size, _ := ctx.MustGet(framework.Page).(int64)

	ctx.JSON(
		statusCode,
		ListResponseType[T]{
			Items: response.Items,
			Pagination: PaginationResponseType{
				Total:   response.Pagination.Total,
				HasNext: limit*size < response.Pagination.Total,
			},
		},
	)
}
