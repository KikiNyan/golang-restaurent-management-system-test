package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var menuCollection *mongo.Collection = database.OpenCollection(database.Client, "menu")
var validate = validator.New()

func GetMenus() gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage<1{
			recordPerPage=10
		}

		page, err := strconv.Atoi(c.Query("page"))
		if err != nil || page<1{
			page=1
		}
		
		startIndex := (page-1) * recordPerPage
		startIndex, err = strconv.Atoi(c.Query("startIndex"))
		
		matchStage := bson.D{{"$match"}, bson.D{{}}}
		groupStage := bson.D{{"$group", bson.D{{"_id", bson.D{{"_id", "null"}}}, {"total_count", bson.D{{"$sum, 1"}}},{"data",bson.D{{"$push", "$$ROOT"}}} }}}
		projectStage := bson.D{
			{
				"$project", bson.D{
					{"_id",0},
					{"total_count", 1},
					{"food_items", bson.D{{"slice", []interface{}{"$data", startIndex, recordPerPage}}}},
				}
			}
		}

		result, err := menuCollection.Aggregate(ctx, mongo.Pipeline{matchStage, groupStage, projectStage})
		defer cancel()

		if err != nil{
			msg := "error occuerd while listing menus"
			c.JSON(http.StatusInternalServerError, gin.H{"error", msg})
			return
		}

		var allMenus []bson.M
		if err = result.All(ctx, &allMenus); err != nil{
			log.Fatal(err)
		}

		c.JSON(http.StatusOK, allMenus[0])

	}
}

func GetMenu() gin.HandlerFunc{
	return func(c *gin.Context){
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		menuId := c.Param("menu_id")
		var menu models.Menu
		
		err := menuCollection.FindOne(ctx, bson.M{"menu_id":menuId}).Decode(&menu)
		defer cancel()

		if err!=nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error":"error occured while fetching menu"})
		}

		c.JSON(http.StatusOK, menu)
	}
}

func CreateMenu() gin.HandlerFunc{
	return func(c *gin.Context){
		var menu models.Menu

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		if err := c.BindJSON(&menu); err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
		}

		validationErr := validate.Struct(menu)
		if validationErr != nil{
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		}

		menu.Created_at,_ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		menu.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		menu.ID = primitive.NewObjectID()
		menu.Menu_id = menu.ID.Hex()

		result,insertErr := menuCollection.InsertOne(ctx, menu)
		if insertErr!=nil{
			msg := fmt.Sprintf("Menu item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error":msg})
			return
		}

		defer cancel()
		c.JSON(http.StatusOK, result)
		defer cancel()
	}
}

func inTimeSpan(start, end, check time.Time) bool{
	return start.After(time.Now()) && end.After(start)
}

func UpdateMenu() gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var menu models.Menu

		if err := c.BindJSON(&menu); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		menuId := c.Param("menu_id")
		filter := bson.M{"menu_id":menuId}

		var updateObj primitive.D

		if menu.Start_Date !=nil && menu.End_Date !=nil{
			if !inTimeSpan(*menu.Start_Date, *menu.End_Date, time.Now()){
				msg := "Kindly retype the time"
				c.JSON(http.StatusInternalServerError, gin.H{"error":msg})
				defer cancel()
				return
			}

			updateObj = append(updateObj, bson.E{"start_date", menu.Start_Date})
			updateObj = append(updateObj, bson.E{"end_date", menu.End_Date})

			if menu.Name != ""{
				updateObj = append(updateObj, bson.E{"name", menu.Name})
			}
			if menu.Category != ""{
				updateObj = append(updateObj, bson.E{"category", menu.Category})
			}

			menu.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			updateObj = append(updateObj, bson.E{"updated_date", menu.Updated_at})

			upsert:=true

			opt := options.UpdateOptions{
				Upsert: &upsert,
			}

			result, err := menuCollection.UpdateOne(ctx,filter,bson.D{
				{"$set", updateObj}
			},
			&opt

			if err!=nil{
				msg := "Menu update failed"
				c.Json(http.StatusInternalServerError, gin.H{"error":msg})
				return
			}

			defer cancel()
			c.JSON(http.StatusOk, result)
		)

		}

	}
}
