package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Invoice struct{
	ID						primitive.ObjectID		`bson:"_id"`
	Invoice_Id				string					`json:"invoice_id"`
	Order_Id				string					`json:"order_id"`
	Payment_method			*string					`json:"payment_method" validate:"eq=CARD|eq=CASH|eq="`
	Payment_status			*string					`json:"payment_status" validate:"eq=PENDING|eq=PAID"`
	Payment_due_date		time.Time				`json:"payment_due_date"`
	Created_at				time.Time				`json:"created_at"`
	Updated_at				time.Time				`json:"updated_at"`
}
