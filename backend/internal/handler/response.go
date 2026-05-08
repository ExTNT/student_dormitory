package handler

import (
	"database/sql"
	"errors"
	"net/http"

	"dormitory/backend/internal/errs"
	"github.com/gin-gonic/gin"
)

func ok(c *gin.Context, data any) {
	c.JSON(http.StatusOK, data)
}

func createdID(c *gin.Context, id int64) {
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func noContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

func fail(c *gin.Context, err error) {
	code := http.StatusInternalServerError
	short := "internal_error"
	switch {
	case errors.Is(err, errs.ErrUnauthorized):
		code, short = http.StatusUnauthorized, "unauthorized"
	case errors.Is(err, errs.ErrForbidden):
		code, short = http.StatusForbidden, "forbidden"
	case errors.Is(err, errs.ErrNotFound), errors.Is(err, sql.ErrNoRows):
		code, short = http.StatusNotFound, "not_found"
	case errors.Is(err, errs.ErrBadRequest):
		code, short = http.StatusBadRequest, "bad_request"
	case errors.Is(err, errs.ErrConflict):
		code, short = http.StatusConflict, "conflict"
	}
	c.JSON(code, gin.H{"error": short, "message": err.Error()})
}

func bindJSON[T any](c *gin.Context) (T, bool) {
	var req T
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad_request", "message": err.Error()})
		return req, false
	}
	return req, true
}
