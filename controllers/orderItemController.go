package controllers


// gin is a framework use to create HTTP server in go frame
import (
	"golang-restaurent-management/models"

	"github.com/gin-gonic/gin"
)

type OrderItemPack struct{
 table_id *string
 Order_items []models.OrderItem

}
var orderItemCollection *mongo.Collection = database.OpenCollection(database.Client, "orderItem")

// this function will get all orderitem and use return func request where orderitem get response from database like juice,burger and so on
func GetOrderItems() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		result, err := orderItemCollection.Find(context.TODO(), bson.M{})

		defer cancel()
		if err != nil {
			msg := "error occurred while listing orders"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		var allOrderItems []bson.M
		if err := result.All(ctx, &allOrderItems); err != nil {
			log.Fatal(err)
		}

		c.JSON(http.StatusOK, allOrderItems)
	}
}

// this func will get one order only by order item by certain order id like humburge name and its price and quantity
func GetOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		orderItemId := c.Param("order_item_id")
		var orderItem models.OrderItem
	
		err := orderItemCollection.FindOne(ctx, bson.M{"order_item_id": orderItemId}).Decode(&orderItem)
		defer cancel()
		if err != nil {
			msg := "error occurred while fetching order"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		c.JSON(http.StatusOK, orderItem)
	}
}

// // this func will get order that included other orderitems by certain order id for example he system utilizes the GetOrderItemsByOrder function. 
// This function retrieves all the order items associated with the selected order ID from the database.
func GetOrderItemsByOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		orderId := c.Param("order_id")

		allOrderItems, err := ItemsByOrder(orderId)

		if err != nil {
			msg := "error occurred while listing orders"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		c.JSON(http.StatusOK, allOrderItems)
	}
}


	func ItemsByOrder(id string) ([]primitive.M, error) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
	
		matchStage := bson.D{{"$match", bson.D{{"order_id", id}}}}
		lookupStage := bson.D{{"$lookup", bson.D{{"from", "food"}, {"localField", "food_id"}, {"foreignField", "food_id"}, {"as", "food"}}}}
		unwindStage := bson.D{{"$unwind", bson.D{{"path", "$food"}, {"preserveNullAndEmptyArrays", true}}}}
		lookupOrderstage := bson.D{{"$lookup", bson.D{{"from", "order"}, {"localField", "order_id"}, {"foreignField", "order_id"}, {"as", "order"}}}}
		unwindOrderStage := bson.D{{"$unwind", bson.D{{"path", "$order"}, {"preserveNullAndEmptyArrays", true}}}}
		lookupTabletage := bson.D{{"$lookup", bson.D{{"from", "table"}, {"localField", "order.table_id"}, {"foreignField", "table_id"}, {"as", "table"}}}}
		unwindTableStage := bson.D{{"$unwind", bson.D{{"path", "$table"}, {"preserveNullAndEmptyArrays", true}}}}
	
		projectStage := bson.D{
			{"$project", bson.D{
				{"id", 0},
				{"amount", "$food.price"},
				{"total_count", 1},
				{"food_name", "$food.name"},
				{"food_image", "$food.food_image"},
				{"table_number", "$table.table_number"},
				{"table_id", "$table.table_id"},
				{"order_id", "$order.order_id"},
				{"price", "$food.price"},
				{"quantity", 1},
			}},
		}
		
		groupStage := bson.D{{"$group", bson.D{
			{"_id", bson.D{
				{"order_id", "$order_id"},
				{"table_id", "$table_id"},
				{"table_number", "$table_number"},
				{"payment_due", bson.D{{"$sum", "$amount"}}},
				{"total_count", bson.D{{"$sum", 1}}},
				{"order_items", bson.D{{"$sum", 1}}},
			}},
		}}}
	

		projectStage2 := bson.D{{"$project", bson.D{
			{"id", 0},
			{"payment_due", 1},
			{"total_count", 1},
			{"table_number", "$_id.table_number"},
			{"order_items", 1},
		}}}
		result,err := orderItemCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage,
			lookupStage,
			unwindStage,
			lookupOrderstage,
			unwindOrderStage,
			lookupTabletage,
			unwindTableStage,
			projectStage,
			groupStage,
		})
		if err != nil {

			panic(err)
		 }
	
		 if err = result.All(ctx,$OrderItems); err !=nil{
			panic(err)
		 }
		 defer cancel()

		 return OrderItems,err
		
	}








// this func will create new order item like add fried chicken item
func CreateOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var orderItemPack OrderItemPack
		if err := c.BindJSON(&orderItemPack); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var order models.Order
		order.Order_date = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		order.Table_id = orderItemPack.Table_id
		order_id := OrderItemOrderCreator(order)

		var orderItemToBeInserted []interface{}
		for _, orderItem := range orderItemPack.Order_items {
			orderItem.Order_id = order_id

			validationError := validate.Struct(orderItem)
			if validationError != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": validationError.Error()})
				return
			}

			orderItem.ID = primitive.NewObjectID()
			orderItem.Created_at = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			orderItem.Updated_at = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			orderItem.Order_item_id = orderItem.ID.Hex()

			num := toFixed(*orderItem.Unit_price, 2)
			orderItem.Unit_price = &num

			orderItemToBeInserted = append(orderItemToBeInserted, orderItem)
		}

		InsertedOrderItems, err := orderItemCollection.InsertMany(ctx, orderItemToBeInserted)
		if err != nil {
			msg := "error occurred while inserting order items"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		c.JSON(http.StatusOK, InsertedOrderItems)
	}
}

// this func will update existed order item like update fried chicken price and quantity
func UpdateOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var orderItem models.OrderItem
		orderItemId := c.Param("order_item_id")

		filter := bson.M{"order_item_id": orderItemId}
		var updateObj primitive.D

		if err := c.BindJSON(&orderItem); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if orderItem.Unit_price != nil {
			updateObj = append(updateObj, bson.E{"unit_price", *orderItem.Unit_price})
		}
		if orderItem.Quantity != nil {
			updateObj = append(updateObj, bson.E{"quantity", *orderItem.Quantity})
		}
		if orderItem.Food_id != nil {
			updateObj = append(updateObj, bson.E{"food_id", *orderItem.Food_id})
		}
		orderItem.Updated_at =  time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{"updated_at", orderItem.Updated_at})

		upsert := true
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		result, err := orderItemCollection.UpdateOne(ctx, filter, bson.D{{"$set", updateObj}}, &opt)
		if err != nil {
			msg := "order update failed"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		c.JSON(http.StatusOK, result)
	}
}

	



	
	

}