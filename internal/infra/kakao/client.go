package kakaoinfra

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/pet-sitter/pets-next-door-api/internal/configs"
)

func FetchAccessToken(code string) (*kakaoTokenResponse, error) {
	kakaoTokenRequest := NewKakaoTokenRequest(
		configs.KakaoRestAPIKey,
		configs.KakaoRedirectURI,
		code,
	)

	req, _ := http.NewRequest(
		"POST",
		"https://kauth.kakao.com/oauth/token",
		strings.NewReader(kakaoTokenRequest.ToURLValues().Encode()),
	)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("charset", "utf-8")
	req.Header.Add("Accept", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	kakaoTokenResponse := &kakaoTokenResponse{}
	if err = json.Unmarshal(body, kakaoTokenResponse); err != nil {
		return nil, err
	}

	return kakaoTokenResponse, nil
}

func FetchUserProfile(code string) (*kakaoUserProfile, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://kapi.kakao.com/v2/user/me", nil)
	req.Header.Add("Authorization", "Bearer "+code)
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	kakaoUserProfile := &kakaoUserProfile{}
	if err = json.Unmarshal(body, kakaoUserProfile); err != nil {
		return nil, err
	}

	return kakaoUserProfile, nil
}
