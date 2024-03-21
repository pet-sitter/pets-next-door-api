package kakaoinfra

import "net/url"

type KakaoTokenRequest struct {
	GrantType   string `json:"grant_type"`
	ClientID    string `json:"client_id"`
	RedirectURI string `json:"redirect_uri"`
	Code        string `json:"code"`
}

func NewKakaoTokenRequest(clientID, redirectURI, code string) *KakaoTokenRequest {
	return &KakaoTokenRequest{
		GrantType:   "authorization_code",
		ClientID:    clientID,
		RedirectURI: redirectURI,
		Code:        code,
	}
}

func (r KakaoTokenRequest) ToURLValues() url.Values {
	values := url.Values{}
	values.Add("grant_type", r.GrantType)
	values.Add("client_id", r.ClientID)
	values.Add("redirect_uri", r.RedirectURI)
	values.Add("code", r.Code)

	return values
}

type KakaoTokenResponse struct {
	TokenType             string `json:"token_type"`
	AccessToken           string `json:"access_token"`
	ExpiresIn             int64  `json:"expires_in"`
	RefreshToken          string `json:"refresh_token"`
	RefreshTokenExpiresIn int64  `json:"refresh_token_expires_in"`
	Scope                 string `json:"scope"`
}

type KakaoUserProfile struct {
	ID          int64  `json:"id"`
	ConnectedAt string `json:"connected_at"`
	Properties  struct {
		Nickname       string `json:"nickname"`
		ProfileImage   string `json:"profile_image"`
		ThumbnailImage string `json:"thumbnail_image"`
	} `json:"properties"`
	KakaoAccount struct {
		ProfileNeedsAgreement bool `json:"profile_needs_agreement"`
		Profile               struct {
			Nickname          string `json:"nickname"`
			ProfileImageURL   string `json:"profile_image_url"`
			ThumbnailImageURL string `json:"thumbnail_image_url"`
		} `json:"profile"`
		HasEmail               bool   `json:"has_email"`
		EmailNeedsAgreement    bool   `json:"email_needs_agreement"`
		IsEmailValid           bool   `json:"is_email_valid"`
		IsEmailVerified        bool   `json:"is_email_verified"`
		Email                  string `json:"email"`
		HasAgeRange            bool   `json:"has_age_range"`
		AgeRangeNeedsAgreement bool   `json:"age_range_needs_agreement"`
		AgeRange               string `json:"age_range"`
		HasGender              bool   `json:"has_gender"`
		GenderNeedsAgreement   bool   `json:"gender_needs_agreement"`
		Gender                 string `json:"gender"`
	} `json:"kakao_account"`
}
