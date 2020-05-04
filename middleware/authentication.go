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
	"gitlab.com/kitalabs/go-2gaijin/responses"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

func generateNewToken(user models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"_id":        user.ID,
		"email":      user.Email,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"avatar":     user.AvatarURL,
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("MY_JWT_TOKEN")))

	if err != nil {
		return "Error while generating token, try again", err
	}

	return tokenString, err
}

func RegisterHandler(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", CORS)
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

			tokenString, err := generateNewToken(user)
			if err != nil {
				res.Error = "Error while generating token, try again"
				json.NewEncoder(c.Writer).Encode(res)
				return
			}

			user.Token = tokenString
			user.Password = ""

			var result = struct {
				Message  string      `json:"message" bson:"message"`
				UserData models.User `json:"data"`
			}{}

			result.Message = "Registration Successful"

			json.NewEncoder(c.Writer).Encode(result)
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
	var result models.User
	var res responses.ResponseMessage
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
		res.Status = "Error"
		res.Message = "Invalid email"
		json.NewEncoder(c.Writer).Encode(res)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(user.Password))

	if err != nil {
		res.Status = "Error"
		res.Message = "Invalid password"
		json.NewEncoder(c.Writer).Encode(res)
		return
	}

	tokenString, err := generateNewToken(result)

	if err != nil {
		res.Status = "Error"
		res.Message = "Error while generating token, try again"
		json.NewEncoder(c.Writer).Encode(res)
		return
	}

	result.Token = tokenString
	result.Password = ""

	var results = struct {
		Status   string      `json:"status" bson:"message"`
		Message  string      `json:"message" bson:"message"`
		UserData models.User `json:"data"`
	}{}
	results.Status = "Success"
	results.Message = "Login Success"
	results.UserData = result

	json.NewEncoder(c.Writer).Encode(results)
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
	var res responses.ResponseMessage

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		id, err := primitive.ObjectIDFromHex(claims["_id"].(string))

		if err != nil {
			res.Status = "Error"
			res.Message = "Something went wrong. Please try again"
			json.NewEncoder(c.Writer).Encode(res)
			return
		}

		result.ID = id
		result.Email = claims["email"].(string)
		result.FirstName = claims["first_name"].(string)
		result.LastName = claims["last_name"].(string)
		result.AvatarURL = claims["avatar"].(string)

		var resp = struct {
			Status   string      `json:"status" bson:"message"`
			Message  string      `json:"message" bson:"message"`
			UserData models.User `json:"data"`
		}{}
		resp.Status = "Success"
		resp.Message = "Profile Successfully Retrieved"
		resp.UserData = result

		json.NewEncoder(c.Writer).Encode(resp)
		return
	} else {
		res.Status = "Error"
		res.Message = err.Error()
		json.NewEncoder(c.Writer).Encode(res)
		return
	}

}

func LoggedInUser(tokenString string) (models.User, bool) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method")
		}
		return []byte(os.Getenv("MY_JWT_TOKEN")), nil
	})

	var result models.User

	if err != nil {
		return result, false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		id, err := primitive.ObjectIDFromHex(claims["_id"].(string))
		if err != nil {
			return result, false
		}

		result.ID = id
		result.Email = claims["email"].(string)
		result.FirstName = claims["first_name"].(string)
		result.LastName = claims["last_name"].(string)
		result.AvatarURL = claims["avatar"].(string)
	} else {
		return result, false
	}

	return result, true
}
