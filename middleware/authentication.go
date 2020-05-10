package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gitlab.com/kitalabs/go-2gaijin/models"
	"gitlab.com/kitalabs/go-2gaijin/responses"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", CORS)
	c.Writer.Header().Set("Content-Type", "application/json")
	var user models.User
	body, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(body, &user)
	var res responses.ResponseMessage
	if err != nil {
		res.Status = "Error"
		res.Message = err.Error()
		json.NewEncoder(c.Writer).Encode(res)
		return
	}

	collection := DB.Collection("users")

	if err != nil {
		res.Status = "Error"
		res.Message = err.Error()
		json.NewEncoder(c.Writer).Encode(res)
		return
	}
	var result models.User
	err = collection.FindOne(context.TODO(), bson.D{{"email", user.Email}}).Decode(&result)

	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 5)

			if err != nil {
				res.Status = "Error"
				res.Message = "Error While Hashing Password, Try Again"
				json.NewEncoder(c.Writer).Encode(res)
				return
			}
			user.Password = string(hash)

			_, err = collection.InsertOne(context.TODO(), user)
			if err != nil {
				res.Status = "Error"
				res.Message = "Error While Creating User, Try Again"
				json.NewEncoder(c.Writer).Encode(res)
				return
			}

			tokenString, err := generateNewToken(user)
			update := bson.M{"$set": bson.M{"token": tokenString}}
			_, err = collection.UpdateOne(context.Background(), bson.D{{"email", user.Email}}, update)
			if err != nil {
				res.Status = "Error"
				res.Message = "Error while generating token, try again"
				json.NewEncoder(c.Writer).Encode(res)
				return
			}

			user.Token = tokenString
			user.Password = ""

			var result = struct {
				Status   string      `json:"status"`
				Message  string      `json:"message" `
				UserData models.User `json:"data"`
			}{}

			result.Status = "Success"
			result.Message = "Registration Successful"
			result.UserData = user

			json.NewEncoder(c.Writer).Encode(result)
			return
		}

		res.Status = "Error"
		res.Message = err.Error()
		json.NewEncoder(c.Writer).Encode(res)
		return
	}

	res.Status = "Error"
	res.Message = "Email already Exists!!"
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
	update := bson.M{"$set": bson.M{"token": tokenString}}
	_, err = collection.UpdateOne(context.Background(), bson.D{{"email", user.Email}}, update)
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

	c.SetCookie("jid", tokenString, 10, "/", "localhost", false, true)

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

func ResetPasswordHandler(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", CORS)
	c.Writer.Header().Set("Content-Type", "application/json")
	var user models.User
	body, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(body, &user)
	var res responses.ResponseMessage
	if err != nil {
		res.Status = "Error"
		res.Message = err.Error()
		json.NewEncoder(c.Writer).Encode(res)
		return
	}

	var collection = DB.Collection("users")

	var result models.User
	err = collection.FindOne(context.Background(), bson.D{{"email", user.Email}}).Decode(&result)
	if err != nil {
		res.Status = "Error"
		res.Message = err.Error()
		json.NewEncoder(c.Writer).Encode(res)
		return
	}

	tokenString, err := generateResetToken(result)
	if err != nil {
		res.Status = "Error"
		res.Message = err.Error()
		json.NewEncoder(c.Writer).Encode(res)
		return
	}

	tokenExpiry := primitive.NewDateTimeFromTime(time.Now().Add(time.Hour * 1))
	update := bson.M{"$set": bson.M{"reset_password_token": tokenString, "reset_token_expiry": tokenExpiry}}

	_, err = collection.UpdateOne(context.Background(), bson.D{{"email", user.Email}}, update)
	if err != nil {
		res.Status = "Error"
		res.Message = err.Error()
		json.NewEncoder(c.Writer).Encode(res)
		return
	}

	res.Status = "Success"
	res.Message = "Check your email to reset your password"

	json.NewEncoder(c.Writer).Encode(res)
	return
}

func UpdatePasswordHandler(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", CORS)
	c.Writer.Header().Set("Content-Type", "application/json")
	var user models.User
	body, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(body, &user)
	var res responses.ResponseMessage
	if err != nil {
		res.Status = "Error"
		res.Message = err.Error()
		json.NewEncoder(c.Writer).Encode(res)
		return
	}

	var collection = DB.Collection("users")

	var result models.User
	err = collection.FindOne(context.Background(), bson.D{{"email", user.Email}}).Decode(&result)
	if err != nil {
		res.Status = "Error"
		res.Message = err.Error()
		json.NewEncoder(c.Writer).Encode(res)
		return
	}

	if result.ResetPasswordToken != "" && user.ResetPasswordToken != "" {
		if result.ResetPasswordToken == user.ResetPasswordToken {
			expiryTime := result.ResetTokenExpiry.Time()

			if time.Now().Before(expiryTime) {
				hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 5)
				if err != nil {
					res.Status = "Error"
					res.Message = "Error While Hashing Password, Try Again"
					json.NewEncoder(c.Writer).Encode(res)
					return
				}
				update := bson.M{"$set": bson.M{"password": string(hash)}}

				_, err = collection.UpdateOne(context.Background(), bson.D{{"email", user.Email}}, update)
				if err != nil {
					res.Status = "Error"
					res.Message = err.Error()
					json.NewEncoder(c.Writer).Encode(res)
					return
				}

				update = bson.M{"$set": bson.M{"reset_password_token": "", "reset_token_expiry": primitive.NewDateTimeFromTime(time.Now())}}
				_, err = collection.UpdateOne(context.Background(), bson.D{{"email", user.Email}}, update)
				if err != nil {
					res.Status = "Error"
					res.Message = err.Error()
					json.NewEncoder(c.Writer).Encode(res)
					return
				}

				res.Status = "Success"
				res.Message = "Password Successfully Changed"
				json.NewEncoder(c.Writer).Encode(res)
				return
			} else {
				res.Status = "Error"
				res.Message = "Session has expired. Try again"
				json.NewEncoder(c.Writer).Encode(res)
				return
			}
		} else {
			res.Status = "Error"
			res.Message = "Wrong credentials. Try again"
			json.NewEncoder(c.Writer).Encode(res)
			return
		}
	} else {
		res.Status = "Error"
		res.Message = "Something went wrong. Try again"
		json.NewEncoder(c.Writer).Encode(res)
		return
	}
}

func EmailConfirmation(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", CORS)
	c.Writer.Header().Set("Content-Type", "application/json")

	var res responses.ResponseMessage
	var err error

	email := c.Request.URL.Query().Get("email")
	token := c.Request.URL.Query().Get("confirm_token")

	var collection = DB.Collection("users")

	var result models.User
	err = collection.FindOne(context.Background(), bson.D{{"email", email}}).Decode(&result)
	if err != nil {
		res.Status = "Error"
		res.Message = err.Error()
		json.NewEncoder(c.Writer).Encode(res)
		return
	}

	if result.ConfirmToken != "" && token != "" {
		if result.ConfirmToken == token {
			expiryTime := result.ConfirmTokenExpiry.Time()

			if time.Now().Before(expiryTime) {
				update := bson.M{"$set": bson.M{"email_confirmed": true, "confirm_token": "", "confirm_token_expiry": primitive.NewDateTimeFromTime(time.Now())}}
				_, err = collection.UpdateOne(context.Background(), bson.D{{"email", email}}, update)
				if err != nil {
					res.Status = "Error"
					res.Message = err.Error()
					json.NewEncoder(c.Writer).Encode(res)
					return
				}

				res.Status = "Success"
				res.Message = "Email has successfully been confirmed"
				json.NewEncoder(c.Writer).Encode(res)
				return
			} else {
				res.Status = "Error"
				res.Message = "Session has expired. Try again"
				json.NewEncoder(c.Writer).Encode(res)
				return
			}
		} else {
			res.Status = "Error"
			res.Message = "Wrong credentials. Try again"
			json.NewEncoder(c.Writer).Encode(res)
			return
		}
	} else {
		res.Status = "Error"
		res.Message = "Wrong credentials. Try again"
		json.NewEncoder(c.Writer).Encode(res)
		return
	}
}

func PhoneConfirmation(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", CORS)
	c.Writer.Header().Set("Content-Type", "application/json")

	var res responses.ResponseMessage
	var err error

	phone := c.Request.URL.Query().Get("phone")
	token := c.Request.URL.Query().Get("confirm_token")

	var collection = DB.Collection("users")

	var result models.User
	err = collection.FindOne(context.Background(), bson.D{{"phone", phone}}).Decode(&result)
	if err != nil {
		res.Status = "Error"
		res.Message = err.Error()
		json.NewEncoder(c.Writer).Encode(res)
		return
	}

	if result.ConfirmToken != "" && token != "" {
		if result.ConfirmToken == token {
			expiryTime := result.ConfirmTokenExpiry.Time()

			if time.Now().Before(expiryTime) {
				update := bson.M{"$set": bson.M{"phone_confirmed": true, "confirm_token": "", "confirm_token_expiry": primitive.NewDateTimeFromTime(time.Now())}}
				_, err = collection.UpdateOne(context.Background(), bson.D{{"phone", phone}}, update)
				if err != nil {
					res.Status = "Error"
					res.Message = err.Error()
					json.NewEncoder(c.Writer).Encode(res)
					return
				}

				res.Status = "Success"
				res.Message = "Phone has successfully been confirmed"
				json.NewEncoder(c.Writer).Encode(res)
				return
			} else {
				res.Status = "Error"
				res.Message = "Session has expired. Try again"
				json.NewEncoder(c.Writer).Encode(res)
				return
			}
		} else {
			res.Status = "Error"
			res.Message = "Wrong credentials. Try again"
			json.NewEncoder(c.Writer).Encode(res)
			return
		}
	} else {
		res.Status = "Error"
		res.Message = "Empty credentials. Try again"
		json.NewEncoder(c.Writer).Encode(res)
		return
	}
}

func GenerateConfirmToken(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", CORS)
	c.Writer.Header().Set("Content-Type", "application/json")
	var user models.User
	body, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(body, &user)
	var res responses.ResponseMessage
	if err != nil {
		res.Status = "Error"
		res.Message = err.Error()
		json.NewEncoder(c.Writer).Encode(res)
		return
	}

	var filter bson.D
	var message string
	if user.Email != "" {
		filter = append(filter, bson.E{"email", user.Email})
		message = "Confirmation has been sent to your email"
	}
	if user.Phone != "" {
		filter = append(filter, bson.E{"phone", user.Phone})
		message = "Confirmation has been sent to your email"
	}

	var result models.User
	var collection = DB.Collection("users")
	err = collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		res.Status = "Error"
		res.Message = err.Error()
		json.NewEncoder(c.Writer).Encode(res)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"_id":                  user.ID,
		"email":                user.Email,
		"first_name":           user.FirstName,
		"last_name":            user.LastName,
		"avatar":               user.AvatarURL,
		"confirm_token_expiry": primitive.NewDateTimeFromTime(time.Now().Add(time.Hour * 1)),
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("MY_JWT_TOKEN")))

	update := bson.M{"$set": bson.M{"confirm_token": tokenString, "confirm_token_expiry": primitive.NewDateTimeFromTime(time.Now().Add(time.Hour * 1))}}
	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		res.Status = "Error"
		res.Message = err.Error()
		json.NewEncoder(c.Writer).Encode(res)
		return
	}

	if user.Email != "" {
		SendEmailConfirmation(tokenString, user.Email)
	}
	if user.Phone != "" {
		SendPhoneConfirmation(tokenString, user.Phone)
	}

	res.Status = "Success"
	res.Message = message
	json.NewEncoder(c.Writer).Encode(res)
	return
}

func generateNewToken(user models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"_id":            user.ID,
		"email":          user.Email,
		"first_name":     user.FirstName,
		"last_name":      user.LastName,
		"avatar":         user.AvatarURL,
		"last_active_at": primitive.NewDateTimeFromTime(time.Now()),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("MY_JWT_TOKEN")))

	if err != nil {
		return "Error while generating token, try again", err
	}

	return tokenString, err
}

func generateResetToken(user models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"_id":                user.ID,
		"email":              user.Email,
		"first_name":         user.FirstName,
		"reset_token_expiry": primitive.NewDateTimeFromTime(time.Now().Add(time.Hour * 1)),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("MY_JWT_TOKEN")))

	if err != nil {
		return "Error while generating token, try again", err
	}

	return tokenString, err
}
