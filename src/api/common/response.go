package common

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type MetaData struct {
	CurrentPage   uint64  `json:"currentPage"`
	PerPage       uint64  `json:"perPage"`
	TotalPages    uint64  `json:"totalPages"`
	TotalElements uint64  `json:"totalElements"`
	NextPage      *uint64 `json:"nextPage"`
	PrevPage      *uint64 `json:"prevPage"`
}

type ResponseDTO struct {
	Message string      `json:"message" bson:"message"`
	Status  string      `json:"status" bson:"status"`
	Data    interface{} `json:"data" bson:"data"`
}

type PaginateDataDTO struct {
	Content  interface{} `json:"content" bson:"content"`
	MetaData *MetaData   `json:"metadata" bson:"metadata"`
}

type PaginationResponseDTO struct {
	Message string          `json:"message" bson:"message"`
	Status  string          `json:"status" bson:"status"`
	Data    PaginateDataDTO `json:"data" bson:"data"`
}

type ResponseOption struct {
	HttpCode int
	MetaData *MetaData
}

func mergeOptions(opts ...*ResponseOption) *ResponseOption {
	successOpt := &ResponseOption{}
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if opt.HttpCode != 0 {
			successOpt.HttpCode = opt.HttpCode
		}
		if opt.MetaData != nil {
			successOpt.MetaData = opt.MetaData
		}
	}
	return successOpt
}

// Response methods

func GenerateSuccessResponse(c echo.Context, data interface{}, message string, options ...*ResponseOption) error {
	_httpCode := http.StatusOK
	_statusText := "success"

	option := mergeOptions(options...)

	if option.HttpCode != 0 {
		_httpCode = option.HttpCode
	}
	if option.MetaData != nil {
		return c.JSON(_httpCode, &PaginationResponseDTO{
			Data: PaginateDataDTO{
				MetaData: option.MetaData,
				Content:  data,
			},
			Message: message,
			Status:  _statusText,
		})
	}

	return c.JSON(_httpCode, ResponseDTO{
		Data:    data,
		Message: message,
		Status:  _statusText,
	})
}

func GenerateErrorResponse(c echo.Context, data interface{}, message string, options ...*ResponseOption) error {
	_httpCode := http.StatusBadRequest
	_statusText := "error"

	option := mergeOptions(options...)

	if option.HttpCode != 0 {
		_httpCode = option.HttpCode
	}

	return c.JSON(_httpCode, ResponseDTO{
		Data:    data,
		Message: message,
		Status:  _statusText,
	})
}
