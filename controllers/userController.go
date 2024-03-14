package controllers

import (
	"context"
	"fmt"
	"golang-restaurent-management/database"
	helper "golang-restaurent-management/helpers"
	"golang-restaurent-management/models"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"honnef.co/go/tools/structlayout"
)



var tableCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

// this func will display all users who registered in system
func GetUsers() gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		
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

		projectStage := bson.D{
			{
				"$project", bson.D{
					{"_id",0},
					{"total_count", 1},
					{"user_items", bson.D{{"slice", []interface{}{"$data", startIndex, recordPerPage}}}},
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

		var allUsers []bson.M
		if err = result.All(ctx, &allUsers); err != nil{
			log.Fatal(err)
		}

		c.JSON(http.StatusOK, allUsers[0])





	}
}

// only display specific user
func GetUser() gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		userId := c.Param("user_id")
		var user models.User

		err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)

		defer cancel()
		if err != nil{
			msg := "error occured while listing orders"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		
		c.JSON(http.StatusOK, user)



	}
}

func SignUp() gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User

    //   convert the json data coming from POSTMAN that golang understand
	if err := c.BindJSON(&user); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//   validate the data based on user struct in models
	
	validationErr := validate.Struct(user)
	if validationErr != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		return
	}
	//   you all check if email has been used by another user
	count, err := userCollection.CountDocuments(ctx,bson.M("email":user.Email))
	if err != nil {
		log.Panic(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking email"})
		return
	}

	//   hash password
	password := HashPassword(*user.Password)
	user.password= &password

	//   you will check if phone no already has been used by other user
	count, err := userCollection.CountDocuments(ctx,bson.M("phone":user.Phone))
	if err != nil {
		log.Panic(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking phone"})
		return
	}

	if count > 0{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while exsiting email or phone number"})
		return

	}

	//   create some extra details for user objects some details created at updated at
	user.Created_at = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.Updated_at = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.ID = primitive.NewObjectID()
	user.User_id= user.ID.Hex()


	// generate all tokes and refresh all token(its comes from token helper)
	token, refreshToken, _ := helper.GenerateAlltokens(*user.Email,*user.First_name,*user.Last_name, user.User_id)
	user.Token =&token
	user.Refresh_Token = &refreshToken


	// if all ok then u insert this new user into user collection in database
	resultInsertationNumber, err := userCollection.Insertone(ctx, user)
	if err != nil {
		msg := fmt.Sprint("user item was not created")
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}
	defer.cancel()
	// return status ok and send the status back
	c.JSON(http.StatusOK, resultInsertationNumber)
   








	}
}

func Login() gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User
		var founduser models.User

		//   convert the login data coming from POSTMAN is in json that golang readable format
		if err := c.BindJSON(&user); err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// find a user with the email and see if that user even exist
		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&founduser)
		defer.cancel()
		if err != nil {
			msg := fmt.Sprint("user item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		
		// then you will verify the password

		passwordIsValid,msg := VerifyPassword(*user.Password,*founduser.Password)
		defer.cancel()
		if passwordIsValid != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
        // if all goes well then u generate token 
		token, refreshToken, _ := helper.GenerateAlltokens(*founduser.Email,*founduser.First_name,*founduser.Last_name, founduser.User_id)
	    user.Token =&token
	    user.Refresh_Token = &refreshToken

		// update token and refresh token()
        helper.UpdateAlltokens(token, refreshToken,founduser.User_id)

		// return status ok
		c.JSON(http.StatusOK, founduser)


    

	}

}

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPassword(userPassword string, providePassword string)(bool,string){
	 err := bcrypt.CompareHashAndPassword([]byte(providedpassword), []byte(userpassword))
	 check := true
	 msg := ""

	 if err != nil {
		msg := fmt.Sprint("login and username password is incorrect")
		check =false
		
	}
	return check, msg





}

