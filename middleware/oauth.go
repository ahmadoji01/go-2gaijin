package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/twinj/uuid"
	"gitlab.com/kitalabs/go-2gaijin/config"
	"gitlab.com/kitalabs/go-2gaijin/models"
	"gitlab.com/kitalabs/go-2gaijin/responses"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// Credentials which stores google ids.
type Credentials struct {
	Cid     string `json:"cid"`
	Csecret string `json:"csecret"`
}

type oAuthCallback struct {
	AccessToken string `json:"access_token"`
}

type googleUserInfo struct {
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Profile       string `json:"profile"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Gender        string `json:"gender"`
}

type facebookUserInfo struct {
	Sub       string `json:"sub"`
	Name      string `json:"name"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Picture   string `json:"picture"`
	Email     string `json:"email"`
}

var cred Credentials
var conf *oauth2.Config
var facebookConf *oauth2.Config
var state string

const oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v3/userinfo?access_token="
const oauthFacebookUrlAPI = "https://graph.facebook.com/me?access_token="

func init() {
	conf = &oauth2.Config{
		RedirectURL:  "http://127.0.0.1:8080/auth/google/callback",
		ClientID:     "880692175404-smp8q2u85pehekh59lk2pj2n4t39u7ha.apps.googleusercontent.com",
		ClientSecret: "P5XHOsNTWhAr4KWC4AiA6wbb",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}

	facebookConf = &oauth2.Config{
		ClientID:     "936813033337153",
		ClientSecret: "d9a728bc5f435f41efd315948a45bd42",
	}
}

func OAuthGoogleCallback(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", config.CORS)
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Authorization")
	c.Writer.Header().Set("Content-Type", "application/json")
	var oAuthInfo oAuthCallback
	body, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(body, &oAuthInfo)
	if err != nil {
		log.Fatal(err)
	}

	var resp responses.GenericResponse

	data, err := getUserDataFromGoogle(oAuthInfo.AccessToken)
	if err != nil {
		resp.Status = "Error"
		resp.Message = err.Error()
		json.NewEncoder(c.Writer).Encode(resp)
		return
	}

	var result googleUserInfo
	if err := json.Unmarshal(data, &result); err != nil {
		resp.Status = "Error"
		resp.Message = err.Error()
		json.NewEncoder(c.Writer).Encode(resp)
		return
	}

	var userData responses.UserData
	var user models.User

	user, err = registerOrLoginGoogle(result)
	if err != nil {
		resp.Status = "Error"
		resp.Message = err.Error()
		json.NewEncoder(c.Writer).Encode(resp)
		return
	}

	userData.User = user
	resp.Status = "Success"
	resp.Message = "Login Successful"
	resp.Data = userData
	json.NewEncoder(c.Writer).Encode(resp)
}

func getUserDataFromGoogle(accessToken string) ([]byte, error) {
	response, err := http.Get(oauthGoogleUrlAPI + accessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed read response: %s", err.Error())
	}
	return contents, nil
}

func registerOrLoginGoogle(userInfo googleUserInfo) (models.User, error) {
	var collection = DB.Collection("users")

	var result models.User
	err := collection.FindOne(context.TODO(), bson.D{{"email", userInfo.Email}}).Decode(&result)

	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			result, err = registerGoogleUser(userInfo)
			if err != nil {
				return models.User{}, err
			}
			return result, nil
		}
		return models.User{}, err
	}
	tokenString, err := GenerateNewToken(result)
	if err != nil {
		return models.User{}, err
	}

	result.Token = tokenString.AuthToken
	result.RefreshToken = tokenString.RefreshToken
	result.AuthTokenExpiry = tokenString.AuthTokenExpiry
	result.RefreshTokenExpiry = tokenString.RefreshTokenExpiry
	result.Password = ""
	return result, nil
}

func registerGoogleUser(userInfo googleUserInfo) (models.User, error) {
	var user models.User
	var collection = DB.Collection("users")

	user.FirstName = userInfo.GivenName
	user.LastName = userInfo.FamilyName
	user.AvatarURL = userInfo.Picture
	user.Email = userInfo.Email
	user.GoogleSub = userInfo.Sub
	user.EmailConfirmed = true
	user.Password = uuid.NewV4().String()

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 5)
	if err != nil {
		return user, err
	}

	user.ID = primitive.NewObjectIDFromTimestamp(time.Now())
	user.Password = string(hash)
	user.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	user.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
	user.NotifRead = true
	user.MessageRead = true

	_, err = collection.InsertOne(context.TODO(), user)
	if err != nil {
		return user, err
	}

	tokenString, err := GenerateNewToken(user)
	if err != nil {
		return user, err
	}

	user.Token = tokenString.AuthToken
	user.RefreshToken = tokenString.RefreshToken
	user.AuthTokenExpiry = tokenString.AuthTokenExpiry
	user.RefreshTokenExpiry = tokenString.RefreshTokenExpiry
	user.Password = ""

	var userData responses.UserData
	userData.User = user

	return user, nil
}

func OAuthFacebookCallback(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", config.CORS)
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Authorization")
	c.Writer.Header().Set("Content-Type", "application/json")
	var oAuthInfo oAuthCallback
	body, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(body, &oAuthInfo)
	if err != nil {
		log.Fatal(err)
	}

	var resp responses.GenericResponse

	data, err := getUserDataFromFacebook(oAuthInfo.AccessToken)
	if err != nil {
		resp.Status = "Error"
		resp.Message = err.Error()
		json.NewEncoder(c.Writer).Encode(resp)
		return
	}

	var result facebookUserInfo
	if err := json.Unmarshal(data, &result); err != nil {
		resp.Status = "Error"
		resp.Message = err.Error()
		json.NewEncoder(c.Writer).Encode(resp)
		return
	}

	var userData responses.UserData
	var user models.User

	user, err = registerOrLoginFacebook(result)
	if err != nil {
		resp.Status = "Error"
		resp.Message = err.Error()
		json.NewEncoder(c.Writer).Encode(resp)
		return
	}

	userData.User = user
	resp.Status = "Success"
	resp.Message = "Login Successful"
	resp.Data = userData
	json.NewEncoder(c.Writer).Encode(resp)
}

func getUserDataFromFacebook(accessToken string) ([]byte, error) {
	response, err := http.Get(oauthFacebookUrlAPI + accessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed read response: %s", err.Error())
	}
	return contents, nil
}

func registerOrLoginFacebook(userInfo facebookUserInfo) (models.User, error) {
	var collection = DB.Collection("users")

	var result models.User
	err := collection.FindOne(context.TODO(), bson.D{{"email", userInfo.Email}}).Decode(&result)

	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			result, err = registerFacebookUser(userInfo)
			if err != nil {
				return models.User{}, err
			}
			return result, nil
		}
		return models.User{}, err
	}
	tokenString, err := GenerateNewToken(result)
	if err != nil {
		return models.User{}, err
	}

	result.Token = tokenString.AuthToken
	result.RefreshToken = tokenString.RefreshToken
	result.AuthTokenExpiry = tokenString.AuthTokenExpiry
	result.RefreshTokenExpiry = tokenString.RefreshTokenExpiry
	result.Password = ""
	return result, nil
}

func registerFacebookUser(userInfo facebookUserInfo) (models.User, error) {
	var user models.User
	var collection = DB.Collection("users")

	user.FirstName = userInfo.FirstName
	user.LastName = userInfo.LastName
	user.AvatarURL = userInfo.Picture
	user.Email = userInfo.Email
	user.EmailConfirmed = true
	user.Password = uuid.NewV4().String()

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 5)
	if err != nil {
		return user, err
	}

	user.ID = primitive.NewObjectIDFromTimestamp(time.Now())
	user.Password = string(hash)
	user.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	user.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
	user.NotifRead = true
	user.MessageRead = true

	_, err = collection.InsertOne(context.TODO(), user)
	if err != nil {
		return user, err
	}

	tokenString, err := GenerateNewToken(user)
	if err != nil {
		return user, err
	}

	user.Token = tokenString.AuthToken
	user.RefreshToken = tokenString.RefreshToken
	user.AuthTokenExpiry = tokenString.AuthTokenExpiry
	user.RefreshTokenExpiry = tokenString.RefreshTokenExpiry
	user.Password = ""

	var userData responses.UserData
	userData.User = user

	return user, nil
}
