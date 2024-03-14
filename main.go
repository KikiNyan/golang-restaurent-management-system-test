package main

import{
     "os"
	 "golang-restaurent-management/database"
	 "golang-restaurent-management/routes"
	 "golang-restaurent-management/middleware"
	 "go.mongodb.org/mongo-driver/mongo"

}

var foodcollection *mongo.Collection=database.OpenCollection(database.client, "food")

func main(){

port=os.Getenv("PORT")

if port == ""{

	port="8000"
}
 router :=gin.new()
 router.Use(gin.Logger())
 routes.userRoutes(router)
 router.Use(middleware.Authentication())
 routes.foodRoutes(router)
 routes.menuRoutes(router)
 routes.tableRoutes(router)
 routes.orderRoutes(router)
 routes.invoiceRoutes(router)

 router.Run(":" + port)


 


}