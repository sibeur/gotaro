package entity

type AccessToken struct {
	Token     string `json:"token"`
	IssuedAt  int64  `json:"iat"`
	ExpiredAt int64  `json:"exp"`
}

type RefreshToken struct {
	Token     string `json:"token"`
	IssuedAt  int64  `json:"iat"`
	ExpiredAt int64  `json:"exp"`
}

type Auth struct {
	AccessToken  *AccessToken  `json:"access_token"`
	RefreshToken *RefreshToken `json:"refresh_token"`
}
