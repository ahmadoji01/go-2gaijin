package middleware

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/twinj/uuid"
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

var cred Credentials
var conf *oauth2.Config
var state string

const oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v3/userinfo?access_token="

func init() {
	conf = &oauth2.Config{
		RedirectURL:  "http://127.0.0.1:8080/auth/google/callback",
		ClientID:     "880692175404-smp8q2u85pehekh59lk2pj2n4t39u7ha.apps.googleusercontent.com",
		ClientSecret: "P5XHOsNTWhAr4KWC4AiA6wbb",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}
}

func OAuthGoogleLogin(c *gin.Context) {
	w := c.Writer
	r := c.Request

	// Create oauthState cookie
	oauthState := generateStateOauthCookie(w)
	u := conf.AuthCodeURL(oauthState)
	fmt.Println(u)
	http.Redirect(w, r, u, http.StatusTemporaryRedirect)
}

func generateStateOauthCookie(w http.ResponseWriter) string {
	var expiration = time.Now().Add(365 * 24 * time.Hour)

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}
	http.SetCookie(w, &cookie)

	return state
}

func OAuthGoogleCallback(c *gin.Context) {
	w := c.Writer
	r := c.Request

	var resp responses.GenericResponse

	// Read oauthState from Cookie
	oauthState, _ := r.Cookie("oauthstate")

	if r.FormValue("state") != oauthState.Value {
		log.Println("invalid oauth google state")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	data, err := getUserDataFromGoogle(r.FormValue("code"))
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

func getUserDataFromGoogle(code string) ([]byte, error) {
	// Use code to get token and get user info from Google.

	token, err := conf.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("code exchange wrong: %s", err.Error())
	}
	response, err := http.Get(oauthGoogleUrlAPI + token.AccessToken)
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
