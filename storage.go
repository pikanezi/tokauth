package tokauth

import (
	mgo "gopkg.in/mgo.v2"
	"time"
)

const ()

var (
	clientCollection *mgo.Collection
	accessCollection *mgo.Collection

	tokenExpiration = time.Minute * 10
	IDFieldName = "_id"
)

// SetTokenExpiration set a new token duration (by default it is set to 10 minutes).
// This function must be called before SetClientCollection to have any effect.
func SetTokenExpiration(expiration time.Duration) { tokenExpiration = expiration }

// SetClientCollection sets the collection used to query the clientID (clientID must be the "_id" of the collection).
func SetClientCollection(collection *mgo.Collection) { clientCollection = collection }

// SetAccessCollection sets the collection used to store the AccessData.
// By default, the AccessToken expiration is set to 10 minutes, to changes it, call SetTokenExpiration before this function.
func SetAccessCollection(collection *mgo.Collection) {
	accessCollection = collection
	accessCollection.DropIndex("expiresAt")
	index := mgo.Index{
		Key:         []string{"expiresAt"},
		ExpireAfter: 1 * time.Second,
	}
	if err := accessCollection.EnsureIndex(index); err != nil {
		panic(err)
	}
}
