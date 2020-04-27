package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gitlab.com/kitalabs/go-2gaijin/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(c *gin.Context) {

	c.Writer.Header().Set("Content-Type", "application/json")
	var user models.User
	body, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(body, &user)
	var res models.ResponseResult
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(c.Writer).Encode(res)
		return
	}

	collection := DB.Collection("users")

	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(c.Writer).Encode(res)
		return
	}
	var result models.User
	err = collection.FindOne(context.TODO(), bson.D{{"email", user.Email}}).Decode(&result)

	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 5)

			if err != nil {
				res.Error = "Error While Hashing Password, Try Again"
				json.NewEncoder(c.Writer).Encode(res)
				return
			}
			user.Password = string(hash)

			_, err = collection.InsertOne(context.TODO(), user)
			if err != nil {
				res.Error = "Error While Creating User, Try Again"
				json.NewEncoder(c.Writer).Encode(res)
				return
			}
			res.Result = "Registration Successful"
			json.NewEncoder(c.Writer).Encode(res)
			return
		}

		res.Error = err.Error()
		json.NewEncoder(c.Writer).Encode(res)
		return
	}

	res.Result = "Email already Exists!!"
	json.NewEncoder(c.Writer).Encode(res)
	return
}

func LoginHandler(c *gin.Context) {

	c.Writer.Header().Set("Content-Type", "application/json")
	var user models.User
	body, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(body, &user)
	if err != nil {
		log.Fatal(err)
	}

	collection := DB.Collection("users")

	if err != nil {
		log.Fatal(err)
	}
	var result = struct {
		ID        primitive.ObjectID `json:"_id" bson:"_id"`
		Email     string             `json:"email" bson:"email"`
		FirstName string             `json:"first_name" bson:"first_name"`
		LastName  string             `json:"last_name" bson:"last_name"`
		Avatar    string             `json:"avatar" bson:"avatar"`
		Token     string             `json:"authentication_token" bson:"token"`
		Password  string             `json:"password" bson:"password"`
	}{}
	var res models.ResponseResult
	var options = &options.FindOptions{}
	options.SetProjection(bson.M{
		"_id":        1,
		"first_name": 1,
		"last_name":  1,
		"avatar":     1,
		"token":      1,
	})

	err = collection.FindOne(context.TODO(), bson.D{{"email", user.Email}}).Decode(&result)

	if err != nil {
		res.Error = "Invalid email"
		json.NewEncoder(c.Writer).Encode(res)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(user.Password))

	if err != nil {
		res.Error = "Invalid password"
		json.NewEncoder(c.Writer).Encode(res)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":      result.Email,
		"first_name": result.FirstName,
		"last_name":  result.LastName,
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("MY_JWT_TOKEN")))

	if err != nil {
		res.Error = "Error while generating token, try again"
		json.NewEncoder(c.Writer).Encode(res)
		return
	}

	result.Token = tokenString
	result.Password = ""

	json.NewEncoder(c.Writer).Encode(result)

}

func ProfileHandler(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")
	tokenString := c.Request.Header.Get("Authorization")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method")
		}
		return []byte(os.Getenv("MY_JWT_TOKEN")), nil
	})
	var result models.User
	var res models.ResponseResult
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		result.Email = claims["email"].(string)
		result.FirstName = claims["first_name"].(string)
		result.LastName = claims["last_name"].(string)

		json.NewEncoder(c.Writer).Encode(result)
		return
	} else {
		res.Error = err.Error()
		json.NewEncoder(c.Writer).Encode(res)
		return
	}

}
