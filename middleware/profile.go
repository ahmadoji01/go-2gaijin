package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/twinj/uuid"
	"gitlab.com/kitalabs/go-2gaijin/config"
	"gitlab.com/kitalabs/go-2gaijin/models"
	"gitlab.com/kitalabs/go-2gaijin/responses"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetProfileInfo(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", config.CORS)
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Authorization")
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
	var tmpUser models.User
	var res responses.ResponseMessage
	var profileData responses.ProfileInfoData

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		id, err := primitive.ObjectIDFromHex(claims["_id"].(string))

		err = DB.Collection("users").FindOne(context.Background(), bson.M{"_id": id}).Decode(&tmpUser)
		if err != nil {
			res.Status = "Error"
			res.Message = "Something went wrong. Please try again"
			json.NewEncoder(c.Writer).Encode(res)
			return
		}

		AtUUID := claims["at_uuid"].(string)
		var collection = DB.Collection("tokens")
		var tokenDetail models.Token
		err = collection.FindOne(context.Background(), bson.M{"auth_token_uuid": AtUUID}).Decode(&tokenDetail)
		if err != nil {
			res.Status = "Error"
			res.Message = "Unauthorized"
			json.NewEncoder(c.Writer).Encode(res)
			return
		}

		AtExpiry, _ := time.Parse(time.RFC3339, claims["at_expiry"].(string))
		if time.Now().After(AtExpiry) {
			res.Status = "Error"
			res.Message = "Token Expired"
			json.NewEncoder(c.Writer).Encode(res)
			return
		}

		var wg sync.WaitGroup

		// Search Gold Trust Coins
		wg.Add(1)
		go func() {
			filter := bson.D{bson.E{"receiver_id", id}, bson.E{"type", "gold"}}
			result.GoldCoin, err = DB.Collection("trust_coins").CountDocuments(context.Background(), filter)
			wg.Done()
		}()

		// Search Silver Trust Coins
		wg.Add(1)
		go func() {
			filter := bson.D{bson.E{"receiver_id", id}, bson.E{"type", "silver"}}
			result.SilverCoin, err = DB.Collection("trust_coins").CountDocuments(context.Background(), filter)
			wg.Done()
		}()
		wg.Wait()

		result.ID = id
		result.Email = tmpUser.Email
		result.Phone = tmpUser.Phone
		result.FirstName = tmpUser.FirstName
		result.LastName = tmpUser.LastName
		result.AvatarURL = ""
		if tmpUser.AvatarURL != "" {
			if !strings.HasPrefix(tmpUser.AvatarURL, "https://") {
				result.AvatarURL = AvatarURLPrefix + claims["_id"].(string) + "/" + tmpUser.AvatarURL
			} else {
				result.AvatarURL = tmpUser.AvatarURL
			}
		}
		result.Role = tmpUser.Role
		result.DateOfBirth = tmpUser.DateOfBirth
		result.ShortBio = tmpUser.ShortBio
		result.EmailConfirmed = tmpUser.EmailConfirmed
		result.PhoneConfirmed = tmpUser.PhoneConfirmed

		var options = &options.FindOptions{}
		projection := bson.D{{"_id", 1}, {"name", 1}, {"price", 1}, {"img_url", 1}, {"user_id", 1}, {"seller_name", 1}, {"latitude", 1}, {"longitude", 1}, {"status_cd", 1}}
		sort := bson.D{{"created_at", -1}}
		options.SetProjection(projection)
		options.SetSort(sort)

		profileData.Profile = result

		var resp responses.GenericResponse
		resp.Status = "Success"
		resp.Message = "Profile Successfully Retrieved"
		resp.Data = profileData

		json.NewEncoder(c.Writer).Encode(resp)
		return
	}
	res.Status = "Error"
	res.Message = err.Error()
	json.NewEncoder(c.Writer).Encode(res)
	return

}

func UploadProfilePhoto(c *gin.Context) {
	c.Writer.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	c.Writer.Header().Set("Access-Control-Allow-Origin", config.CORS)
	c.Writer.Header().Set("Content-Type", "application/json")
	var res responses.GenericResponse

	tokenString := c.Request.Header.Get("Authorization")
	userData, isLoggedIn := LoggedInUser(tokenString)

	if isLoggedIn {
		var uploadAvatar models.UploadAvatar

		body, _ := ioutil.ReadAll(c.Request.Body)
		err := json.Unmarshal(body, &uploadAvatar)
		if err != nil {
			res.Status = "Error"
			res.Message = err.Error()
			json.NewEncoder(c.Writer).Encode(res)
			return
		}

		imgPath := uuid.NewV4().String()
		imgPath = imgPath + "/"

		imgName := uuid.NewV4().String()
		imgName = imgName + ".jpg"

		DecodeBase64ToImage(uploadAvatar.Avatar, imgName)
		UploadToGCS(AvatarImagePrefix+imgPath, imgName)

		update := bson.M{"$set": bson.D{{"avatar", AvatarURLPrefix + imgPath + imgName}}}
		_, err = DB.Collection("users").UpdateOne(context.Background(), bson.M{"_id": userData.ID}, update)
		if err != nil {
			res.Status = "Error"
			res.Message = err.Error()
			json.NewEncoder(c.Writer).Encode(res)
			return
		}

		var upResult = struct {
			AvatarURL string `json:"avatar_url"`
		}{}
		upResult.AvatarURL = AvatarURLPrefix + imgPath + imgName

		res.Status = "Success"
		res.Message = "Avatar Uploaded Successfully"
		res.Data = upResult
		json.NewEncoder(c.Writer).Encode(res)
		return
	}

	res.Status = "Error"
	res.Message = "Unauthorized"
	json.NewEncoder(c.Writer).Encode(res)
	return
}

func GetEmailConfirmationStatus(c *gin.Context) {
	c.Writer.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	c.Writer.Header().Set("Access-Control-Allow-Origin", config.CORS)
	c.Writer.Header().Set("Content-Type", "application/json")
	var res responses.GenericResponse

	tokenString := c.Request.Header.Get("Authorization")
	userData, isLoggedIn := LoggedInUser(tokenString)
	var tmpUser models.User

	if isLoggedIn {
		err := DB.Collection("users").FindOne(context.Background(), bson.M{"_id": userData.ID}).Decode(&tmpUser)

		if err != nil {
			res.Status = "Error"
			res.Message = err.Error()
			json.NewEncoder(c.Writer).Encode(res)
			return
		}

		var emailConfirm = struct {
			EmailConfirmed bool `json:"email_confirmed"`
		}{}
		emailConfirm.EmailConfirmed = tmpUser.EmailConfirmed

		res.Status = "Success"
		res.Message = "Email Confirmation Status Retrieved"
		res.Data = emailConfirm
		json.NewEncoder(c.Writer).Encode(res)
		return
	}
	res.Status = "Error"
	res.Message = "Unauthorized"
	json.NewEncoder(c.Writer).Encode(res)
	return
}

func GetPhoneConfirmationStatus(c *gin.Context) {
	c.Writer.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	c.Writer.Header().Set("Access-Control-Allow-Origin", config.CORS)
	c.Writer.Header().Set("Content-Type", "application/json")
	var res responses.GenericResponse

	tokenString := c.Request.Header.Get("Authorization")
	userData, isLoggedIn := LoggedInUser(tokenString)
	var tmpUser models.User

	if isLoggedIn {
		err := DB.Collection("users").FindOne(context.Background(), bson.M{"_id": userData.ID}).Decode(&tmpUser)

		if err != nil {
			res.Status = "Error"
			res.Message = err.Error()
			json.NewEncoder(c.Writer).Encode(res)
			return
		}

		var phoneConfirm = struct {
			PhoneConfirmed bool `json:"phone_confirmed"`
		}{}
		phoneConfirm.PhoneConfirmed = tmpUser.PhoneConfirmed

		res.Status = "Success"
		res.Message = "Phone Confirmation Status Retrieved"
		res.Data = phoneConfirm
		json.NewEncoder(c.Writer).Encode(res)
		return
	}
	res.Status = "Error"
	res.Message = "Unauthorized"
	json.NewEncoder(c.Writer).Encode(res)
	return
}
