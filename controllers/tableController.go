package controllers

import (
	"golang-restaurent-management/models"
	"golang-restaurent-management/database"

	"github.com/gin-gonic/gin"
)

var tableCollection *mongo.Collection = database.OpenCollection(database.Client, "table")
// this used to retrieve a list of all tables. It is used in the web application to display a list of all tables.like table 1 table 2 table 3





func GetTables() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		result, err := tableCollection.Find(ctx, bson.M{})
		if err != nil {
			msg := "error occurred while listing tables"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		var allTables []bson.M
		if err = result.All(ctx, &allTables); err != nil {
			msg := "error occurred while decoding tables"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		c.JSON(http.StatusOK, allTables)
	}
}

// The GetTable function is used to retrieve the details of a specific table. It is utilized in the web application to display the details of a particular table.
// for example when waiter check details of table no 5 to see if table is already occuiped by customer and their order
func GetTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		tableId := c.Param("table_id")
		var table models.Table

		err := tableCollection.FindOne(ctx, bson.M{"table_id": tableId}).Decode(&table)
		if err != nil {
			msg := "error occurred while fetching table"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		c.JSON(http.StatusOK, table)
	}
}

func CreateTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var table models.Table

		if err := c.BindJSON(&table); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationError := validate.Struct(table)
		if validationError != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationError.Error()})
			return
		}

		table.Created_at = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		table.Updated_at = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		table.ID = primitive.NewObjectID()
		result, err := tableCollection.InsertOne(ctx, table)

		if err != nil {
			msg := "error occurred while inserting table"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

func UpdateTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var table models.Table
		tableId := c.Param("table_id")

		if err := c.BindJSON(&table); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var updateObj primitive.D

		if table.Number_of_guests != nil {
			updateObj = append(updateObj, bson.E{"number_of_guests", *table.Number_of_guests})
		}

		if table.Table_number != nil {
			updateObj = append(updateObj, bson.E{"table_number", *table.Table_number})
		}

		table.Updated_at = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		upsert := true
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}
		filter := bson.M{"table_id": tableId}
		result, err := tableCollection.UpdateOne(ctx, filter, bson.D{{"$set", updateObj}}, &opt)

		if err != nil {
			msg := "error occurred while updating table items"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}
