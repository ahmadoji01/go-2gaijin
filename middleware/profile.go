package middleware

import (
	"context"
	"encoding/json"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/twinj/uuid"
	"gitlab.com/kitalabs/go-2gaijin/config"
	"gitlab.com/kitalabs/go-2gaijin/models"
	"gitlab.com/kitalabs/go-2gaijin/responses"
	"go.mongodb.org/mongo-driver/bson"
)

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
