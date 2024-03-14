package helper


import(
	"context"
	"fmt"
	"log"
	"os"
	"time"
	"golang-restaurent-management/database"
	jwt "github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"



)
type SignedDetails struct{
Email string
First_name string
Last_name string
Uid string
jwt.StandardClaims


}

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var SECRET_KEY string = os.Getenv("SECRET_KEY")


// this func will display all users who registered in system
func GenerateAlltokens(email string,firstname string,lastname string, uid string)(signedToken string, refreshToken string, err error){
	claims := &SignedDetails{
 Email: email,
 First_name: firstname,
 Last_name: lastname,
 Uid: uid,
 StandardClaims: jwt.StandardClaims{
 ExpiresAt: time.Now().Local.Add(time.Hour * time.Duration (24)).Unix,

 },

 refreshClaims := &SignedDetails{
    StandardClaims: jwt.StandardClaims{
	ExpiresAt: time.Now().Local.Add(time.Hour * time.Duration (168)).Unix,
   
	},
	


	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodH256, claims).SignedString([]byte(SECRET_KEY))
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodH256,  refreshClaims).SignedString([]byte(SECRET_KEY))
    if err!= nil{

		log.panic(err)
		return
	}

	return token, refreshToken,err


}

func UpdateAlltokens(){
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var updateObj primitive.D
	updateObj = append(updateObj, bson.E{"token", signedToken})
	updateObj = append(updateObj, bson.E{"refresh_token", signedrefreshToken})
	Updated_at = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{"updated_at", Updated_at})

	
	upsert := true
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}
	filter := bson.M{"user_id": UpdateOptions}
	_, err := userCollection.UpdateOne(ctx, filter, bson.D{
		{"$set", updateObj}
		},
		&opt,
	)
	defer.cancel()
	if err!= nil{

		log.panic(err)
		return
	}
	return






	
}

func ValidateAlltokens(){
 token,err:= jwt.ParseWithClaims{
 signedToken,
 &SignedDetails{}
 func(token *jwt.Token)(Interface{}, error)
 {
	return[]byte(SECRET_KEY),nil
 },


 }



//   the token is invalid

claims, ok := token.Claims.(*SignedDetails)
if !ok {
	msg := fmt.Sprintf("the token is invalid")
	msg =err.Error()

	return
}

//the token is expired
if claims.ExpiresAt < time.Now().Local().Unix(){
	msg := fmt.Sprintf("the token is expired")
	msg =err.Error()

	return

}
return claims,msg 

	
}
