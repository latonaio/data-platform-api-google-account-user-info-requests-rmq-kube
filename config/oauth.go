package config

type OAuth struct {
	UserInfoURL string
}

func newOAuth() *OAuth {
	return &OAuth{
		UserInfoURL: "https://www.googleapis.com/oauth2/v2/userinfo",
	}
}
