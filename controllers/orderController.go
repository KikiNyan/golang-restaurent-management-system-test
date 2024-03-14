package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/k-kaddal/golang-projects/go-restaurant/database"
	"github.com/k-kaddal/golang-projects/go-restaurant/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

var orderCollection *mongo.Collection = database.OpenCollection(database.Client, "order")

func GetOrders() gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		result, err := orderCollection.Find(context.TODO(), bson.M{})
		
		defer cancel()
		if err != nil{
			msg := "error occured while listing orders"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		
		var allOrders []bson.M
		if err = result.All(ctx, &allOrders); err !=nil{
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allOrders)
	}
}
func GetOrder() gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		orderId := c.Param("order_id")
		var order models.

		err := orderCollection.FindOne(ctx, bson.M{"order_id": orderId}).Decode(&order)
		defer cancel()
		if err != nil{
			msg := "error occured while fetch order"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		c.JSON(http.StatusOK, order)
	}
}
func CreateOrder() gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var order models.Order
		var table models.Table

		if err := c.BindJSON(&order); err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationError := validate.Struct(order)

		if validationError != nil{
			c.JSON(http.StatusBadRequest, gin.H{"error": validationError.Error()})
			return
		}

		if order.Table_id != nil{
			err := tableCollection.FindOne(ctx, bson.M{"table_id": order.Table_id}).Decode(&table)
			defer cancel()
			if err!=nil{
				msg := fmt.Sprintf("message: table was not found")
				c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
				return
			}
		}

		order.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		order.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		order.ID = primitive.NewObjectID()
		order.Order_Id = order.ID.Hex()

		result, insertErr := orderCollection.InsertOne(ctx, order)
		if insertErr != nil{
			msg := fmt.Sprintf("order item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, result)
	}
}
func UpdateOrder() gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var table models.Table
		var order models.Order

		orderId := c.Param("order_id")
		if err := c.BindJSON(&order); err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var updateObj primitive.D

		if order.Table_id != nil{
			err := tableCollection.FindOne(ctx, bson.M{"table_id", order.table_id}).Decode(&table)
			defer cancel()
			if err != nil{
				msg := fmt.Sprintf("message: Table was not found")
				c.JSON(https.StatusInternalServerError, gin.H{"error":msg})
				return
			}
			updateObj = append(updateObj, bson.E{"menu": order.Table_id})
		}

		order.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{"updated_at": order.Updated_at})

		upser := true

		filter := bson.M{"order_id": orderId}
		opt := options.UpdateOptions{
			Upsert:&upser
		}

		result, err := orderCollection.UpdateOne(ctx, filter, bson.D{
			{"$set", updateObj}
			},
			&opt,
		)

		if err != nil{
			msg	:= fmt.Sprintf(("order item update failed"))
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg} )
			return
		}

		defer cancel()
		c.JSON(http.StatusOK, result)
	}
}

func OrderItemOrderCreate(order models.Order) string{
	order.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.ID = primitive.NewObjectID()
	order.Order_Id = order.ID.Hex()

	orderCollection.InsertOne(ctx, order)
	defer cancel()

	return order.Order_Id
}