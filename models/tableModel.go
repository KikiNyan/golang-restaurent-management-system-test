package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Table struct{
	ID				 primitive.ObjectID			`bson:"_id"`
	Number_of_guests *int			            `bson:"number_of_guests"validate:required`
	Table_number     *int			            `bson:"table_number "validate:required`
	Created_at		 time.Time					`json:"updated_id"`
	Updated_at		 time.Time					`json:"updated_id"`
	Table_id		 string						`json:"table_id"`
}
