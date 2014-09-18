package tokauth

import (
	"code.google.com/p/go-uuid/uuid"
	"encoding/hex"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type accessData struct {
	ExpiresAt    time.Time `json:"expiresAt,omitempty" bson:"createdAt"`
	AccessToken  string    `json:"accessToken,omitempty" bson:"accessToken"`
	RefreshToken string    `json:"refreshToken,omitempty" bson:"refreshToken"`
}

// newAccessData creates new AccessData from the ID and insert it in the database.
func newAccessData(refreshToken string) (data *accessData, err error) {
	hexRefreshToken, err := hex.DecodeString(refreshToken)
	if err != nil {
		return
	}
	hexAccessToken, err := hex.DecodeString(uuid.New())
	if err != nil {
		return
	}
	data = &accessData{
		ExpiresAt:    time.Now().Add(tokenExpiration),
		AccessToken:  string(hexAccessToken),
		RefreshToken: string(hexRefreshToken),
	}
	_, err = accessCollection.UpsertId(data.AccessToken, data)
	return
}

// Register a new client. The client must exists in the collection given via the function SetClientCollection.
// The ID parameter is the ID of the client that needs to be registered.
func Register(id interface{}) (refreshToken, accessToken string, err error) {
	refreshToken = uuid.New()
	data, err := newAccessData(refreshToken)
	if err != nil {
		return
	}
	accessToken = data.AccessToken
	if err = clientCollection.Update(bson.M{"_id": id}, bson.M{"$set": bson.M{"refreshToken": refreshToken}}); err != nil {
		return
	}
	return
}

// SetRemainingTime set the new date at which the token will expires.
func SetRemainingTime(accessToken string, date time.Time) (err error) {
	return accessCollection.Update(bson.M{"accessToken": accessToken}, bson.M{"expiresAt": date})
}

// Refresh the AccessData of the user.
func Refresh(refreshToken string) (accessToken string, err error) {
	if err = clientCollection.Find(bson.M{"refreshToken": refreshToken}).One(&Client{}); err != nil {
		return
	}
	data, err := newAccessData(refreshToken)
	if err != nil {
		return
	}
	accessToken = data.AccessToken
	return
}

// Authorize check if the given token match with the AccessToken of the client.
func Authorize(token string) bool {
	var data *accessData
	if err := accessCollection.FindId(token).One(&data); err == mgo.ErrNotFound {
		return false
	}
	return true
}

// GetRefreshTokenFromAccessToken returns the refreshToken of the user associated by this accessToken.
func GetRefreshTokenFromAccessToken(accessToken string) (refreshToken string, err error) {
	client := &Client{}
	if err = accessCollection.FindId(accessToken).One(&client); err != nil {
		return
	}
	refreshToken = client.RefreshToken
	return
}

// GetAccessTokenFromRefreshToken returns the accessToken of the user associated by this refreshToken.
func GetAccessTokenFromRefreshToken(refreshToken string) (accessToken string, err error) {
	data := &accessData{}
	if err = accessCollection.Find(bson.M{"refreshToken": refreshToken}).One(&data); err != nil {
		return
	}
	accessToken = data.AccessToken
	return
}
