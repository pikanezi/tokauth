package tokauth

// Client is a simple struct with a RefreshToken field.
// This struct has to be embedded by the struct populating the Client Collection set by the user.
type Client struct {
	RefreshToken string `json:"refreshToken,omitempty" bson:"refreshToken"`
	AccessToken  string `json:"accessToken,omitempty" bson:"-"`
}
