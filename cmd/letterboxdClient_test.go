package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCreateRequestSignature(t *testing.T) {
	method := "POST"
	url := "http://api.letterboxd.com/api/v0/auth/token"
	body := "abc=123"
	secret := "hello"

	// Did this one on the go playground.
	truth := "d78520de57bd3ffccf6a00a48e69e4d567589e6e20b11938cf1184564d64d026"

	answer, err := createRequestSignature(method, url, body, secret)
	if err != nil {
		t.Errorf("Encountered error: %v", err)
	}
	if !cmp.Equal(truth, answer) {
		t.Errorf("Expected %v, got %v", truth, answer)
	}
}

func TestPrepareUrl(t *testing.T) {
	requestUrl := "http://api.letterboxd.com/api/v0/auth/token"
	key := "abc-123"
	url := `http://api\.letterboxd\.com/api/v0/auth/token`
	apikey := "apikey=abc-123"
	uuid := `nonce=[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`
	timestamp := `timestamp=\d+`
	truthRegex := regexp.MustCompile(
		fmt.Sprintf("%v/\\?%v\\&%v\\&%v", url, apikey, uuid, timestamp),
	)

	answer := prepareUrl(requestUrl, key)
	if !truthRegex.MatchString(answer) {
		t.Errorf("%v did not match the url regex", answer)
	}

}

func TestLoad(t *testing.T) {

	client := LetterboxdClient{}
	if err := client.Load("./test_token.json"); err != nil {
		t.Errorf("Encountered error: %v", err)
	}

	truth := LetterboxdLoginToken{
		AccessToken:  "abc",
		TokenType:    "bearer",
		RefreshToken: "def",
		ExpiresIn:    10,
		CreatedAt:    0,
	}
	if !cmp.Equal(truth, client.token) {
		t.Errorf("Expected \n%v, got \n%v", truth, client.token)
	}
}

func TestSave(t *testing.T) {
	token := LetterboxdLoginToken{
		AccessToken:  "abc",
		TokenType:    "bearer",
		RefreshToken: "def",
		ExpiresIn:    10,
		CreatedAt:    0,
	}

	client := LetterboxdClient{
		token: token,
	}

	f, err := os.CreateTemp("", "temp_creds.json")
	if err != nil {
		t.Errorf("Encountered error: %v", err)
	}
	defer os.Remove(f.Name())

	if err := client.Save(f.Name()); err != nil {
		t.Errorf("Encountered error: %v", err)
	}

	data, err := os.ReadFile(f.Name())
	if err != nil {
		t.Errorf("Encountered error: %v", err)
	}
	var readToken LetterboxdLoginToken
	if err = json.Unmarshal(data, &readToken); err != nil {
		t.Errorf("Encountered error: %v", err)
	}

	if !cmp.Equal(token, readToken) {
		t.Errorf("Expected \n%v, got \n%v", token, readToken)
	}
}
