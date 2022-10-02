package cmd

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type LetterboxdLoginToken struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	CreatedAt    int    `json:"created_at"`
}

type LetterboxdClient struct {
	rootUrl    string
	httpClient http.Client
	key        string
	secret     string
	username   string
	password   string
	token      LetterboxdLoginToken
}

type LetterboxdOAuthError struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

func createRequestSignature(
	method string, url string, body string, secret string,
) (string, error) {
	signature := hmac.New(sha256.New, []byte(secret))
	_, err := signature.Write(
		[]byte(strings.Join([]string{method, url, body}, "\u0000")),
	)
	if err != nil {
		return "", fmt.Errorf("error creating signature: %v", err)
	}
	return hex.EncodeToString(signature.Sum(nil)), nil
}

func prepareUrl(requestUrl string, key string) string {
	params := url.Values{}
	params.Set("apikey", key)
	params.Set("nonce", uuid.New().String())
	params.Set("timestamp", strconv.Itoa(int(time.Now().Unix())))
	return fmt.Sprintf("%v/?%v", requestUrl, params.Encode())
}

func (c *LetterboxdClient) Load(credFile string) error {
	data, err := os.ReadFile(credFile)
	if err != nil {
		return fmt.Errorf("error reading creds file %v: %v", credFile, err)
	}
	var token LetterboxdLoginToken
	if err = json.Unmarshal(data, &token); err != nil {
		return fmt.Errorf("error unmarshalling token %v", err)
	}
	c.token = token
	return nil
}

func (c *LetterboxdClient) Save(credFile string) error {
	tokenBytes, err := json.Marshal(&c.token)
	if err != nil {
		return fmt.Errorf("error creating token for file: %v", err)
	}
	if err = os.WriteFile(credFile, tokenBytes, 0644); err != nil {
		return fmt.Errorf("error saving token to file %v: %v", credFile, err)
	}
	return nil
}

// How to use a form-url-encoded
// https://stackoverflow.com/questions/19253469/make-a-url-encoded-post-request-using-http-newrequest
func (c *LetterboxdClient) Authenticate() error {
	baseUrl := prepareUrl(fmt.Sprintf("%v/auth/token", c.rootUrl), c.key)
	body := url.Values{}
	if int(time.Now().Unix()) > (c.token.CreatedAt + c.token.ExpiresIn) {
		body.Set("grant_type", "password")
		body.Set("username", c.username)
		body.Set("password", c.password)
	} else {
		body.Set("grant_type", "refresh_token")
		body.Set("refresh_token", c.token.RefreshToken)
	}
	signature, err := createRequestSignature(
		http.MethodPost, baseUrl, body.Encode(), c.secret,
	)
	if err != nil {
		return fmt.Errorf("error creating signature: %v", err)
	}
	params := url.Values{}
	params.Set("signature", signature)
	url := fmt.Sprintf("%v&%v", baseUrl, params.Encode())
	request, err := http.NewRequest(
		http.MethodPost, url, strings.NewReader(body.Encode()),
	)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Accept", "application/json")
	response, err := c.httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("error performing request: %v", err)
	}
	defer response.Body.Close()
	if response.StatusCode == 200 {
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return fmt.Errorf("error reading response body %v", err)
		}
		var token LetterboxdLoginToken
		if err = json.Unmarshal(body, &token); err != nil {
			return fmt.Errorf("error unmarshalling token %v", err)
		}
		token.CreatedAt = int(time.Now().Unix())
		c.token = token
	} else if response.StatusCode == 400 {
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return fmt.Errorf("error reading response body %v", err)
		}
		var oauthError LetterboxdOAuthError
		if err = json.Unmarshal(body, &oauthError); err != nil {
			return fmt.Errorf("error unmarshalling OAuth error %v", err)
		}
		return fmt.Errorf(
			"oauth error: %v - %v",
			oauthError.Error,
			oauthError.ErrorDescription,
		)
	} else {
		return fmt.Errorf("error getting token: %v", response.StatusCode)
	}

	return nil
}
