package models

import (
	"time"

	"github.com/go-playground/validator/v10"
)

var vld = validator.New()

func Validate(modelsStruct interface{}) error {
	return vld.Struct(modelsStruct)
}

type Comment struct {
	ID        int64     `json:"id" validate:"required"`
	ParentID  int64     `json:"parent_id" validate:"required"`
	Content   string    `json:"content" validate:"required"`
	CreatedAt time.Time `json:"created_at" validate:"required"`
}
