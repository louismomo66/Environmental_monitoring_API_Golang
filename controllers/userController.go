package controllers
import (
	"log"
	"context"
	"net/http"
  "fmt"
	"time"
	"louis/go_projects/models"
	"github.com/go-playground/validator/v10"
	helper "louis/go_projects/helpers"
	"louis/go_projects/database"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bycrypt"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client,"user")
var validate = validate.New()

func HashPassword(password string)string{
	bytes,err := bycrypt.GenerateFromPassword([]byte(password),14)
	if err != nil{
		log.Panic(err)
	}
	return string(bytes)
}

func VeryfyPassword(userPassword string, providedPassword string)(bool,string){
	err := bycrypt.CompareHashAndPassword([]byte(providedPassword),[]byte(userPassword))
	check := true
	msg := ""
	if err != nil{
		msg = fmt.Sprintf("password is incorrect")
		check = false
	}
	return check,msg
}

func SignUp()gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second)
		defer cancel()
		var user models.User

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
		return
	}
	validationErr := validate.Struct(user)
	fmt.Println("This is the validation error",validationErr)
	if validationErr != nil{
		c.JSON(http.StatusBadRequest,gin.H{"error":validationErr.Error()})
		return
	}
	count,err := userCollection.CountDocuments(ctx,bson.M{"phone":user.Phone})
	// defer cancel()
	if err != nil{
		log.Panic(err)
		c.JSON(http.StatusInternalServerError,gin.H{"error":"error occured while checking for the phone number"})
	}
if count > 0{
	c.JSON(http.StatusInternalServerError,gin.H{"error":"This phone number of email already exists"})
}
user.Created_at, _ = time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))
user.Updated_at,_ = time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))
user.ID = primitive.NewObjectID()
user.User_id = user.ID.Hex()
token,refreshToken,_ := helper.GenerateAllTokens(*user.Email,*user.FirstName,*user.LastName,*user.User_type,user.User_id)
user.Token = &token
user.Refresh_token = &refreshToken

resultInsertionNumber, inserErr := userCollection.InsertOne(ctx,user)
if inserErr != nil{
	msg := fmt.Sprintf("User item was not created")
	c.JSON(http.StatusInternalServerError,gin.H{"error":msg})
	return
}
defer cancel()
c.JSON(http.StatusOK,resultInsertionNumber)

	}

}


func Login() gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second)
    defer cancel()
		var user models.User
		var foundUser models.User

		if err := c.BindJSON(&user);err != nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
			return
		}
		err := userCollection.FindOne(ctx, bson.M{"email":user.Email}).Decode(&foundUser)
		defer cancel()
		if err != nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"email or pasword is incorrect"})
			return
		}

		passwordIsValid,msg := VerifyPassword(*user.Password,*foundUser.Password)
		defer cancel()
		if !passwordIsValid {
			c.JSON(http.StatusInternalServerError,gin.H{"error":msg})
			return
		}

		if foundUser.Email == nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"user not found"})
		}
		token,refreshToken,_ := helper.GenerateAllTokens(*foundUser.Email,*foundUser.FirstName,*foundUser.LastName,*foundUser.User_type,foundUser.User_id)
		err = userCollection.FindOne(ctx,bson.M{"user_id":foundUser.User_id}).Decode(&foundUser)

		if err != nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
			return
		}
		c.JSON(http.StatusOK,foundUser)
	}


}



