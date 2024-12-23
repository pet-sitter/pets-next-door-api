package kakaoinfra

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/pet-sitter/pets-next-door-api/internal/configs"
)

type KakaoClient interface {
	FetchAccessToken(ctx context.Context, code string) (*KakaoTokenResponse, error)
	FetchUserProfile(ctx context.Context, code string) (*KakaoUserProfile, error)
}

type KakaoDefaultClient struct{}

func NewKakaoDefaultClient() *KakaoDefaultClient {
	return &KakaoDefaultClient{}
}

func (kakaoClient *KakaoDefaultClient) FetchAccessToken(
	ctx context.Context,
	code string,
) (*KakaoTokenResponse, error) {
	kakaoTokenRequest := NewKakaoTokenRequest(
		configs.KakaoRestAPIKey,
		configs.KakaoRedirectURI,
		code,
	)

	req, _ := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
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

	kakaoTokenResponse := &KakaoTokenResponse{}
	if err = json.Unmarshal(body, kakaoTokenResponse); err != nil {
		return nil, err
	}

	return kakaoTokenResponse, nil
}

func (kakaoClient *KakaoDefaultClient) FetchUserProfile(
	ctx context.Context,
	code string,
) (*KakaoUserProfile, error) {
	client := &http.Client{}
	req, _ := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"https://kapi.kakao.com/v2/user/me",
		nil,
	)
	req.Header.Add("Authorization", "Bearer "+code)
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch user profile from Kakao server")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	kakaoUserProfile := &KakaoUserProfile{}
	if err = json.Unmarshal(body, kakaoUserProfile); err != nil {
		return nil, err
	}

	return kakaoUserProfile, nil
}
