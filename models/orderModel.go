package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	ID         primitive.ObjectID `bson:"_id"`
	Order_date time.Time          `json:"order_date"validate:required`
	Created_at time.Time          `json:"created_id"`
	Updated_at time.Time          `json:"updated_id"`
	Order_Id   string             `json:"order_id"`
	Table_id   *string            `json:"table_id" validate:"required"`
}
