package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bobinette/papernet/errors"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`

	IsAdmin bool `json:"isAdmin"`

	Owns      []int `json:"owns"`
	CanSee    []int `json:"canSee"`
	CanEdit   []int `json:"canEdit"`
	Bookmarks []int `json:"bookmarks"`
}

type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

type Client struct {
	baseURL string
	client  HTTPClient
}

func NewClient(c HTTPClient, baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		client:  c,
	}
}

func (c *Client) User(id int) (User, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/auth/v2/users/%d", c.baseURL, id), nil)
	if err != nil {
		return User{}, err
	}

	res, err := c.client.Do(req)
	if err != nil {
		return User{}, err
	}
	defer res.Body.Close()

	var user User
	err = json.NewDecoder(res.Body).Decode(&user)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (c *Client) Token(id int) (string, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/auth/v2/users/%d/token", c.baseURL, id), nil)
	if err != nil {
		return "", err
	}

	res, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	var token struct {
		AccessToken string `json:"access_token"`
	}
	err = json.NewDecoder(res.Body).Decode(&token)
	if err != nil {
		return "", err
	}

	return token.AccessToken, nil
}

func (c *Client) CreatePaper(userID, paperID int) error {
	body := bytes.Buffer{}
	_ = json.NewEncoder(&body).Encode(map[string]int{"paperId": paperID}) // Cannot fail
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/auth/v2/users/%d/papers", c.baseURL, userID), &body)
	if err != nil {
		return err
	}

	res, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		var callErr struct {
			Message string `json:"message"`
		}
		err := json.NewDecoder(res.Body).Decode(&callErr)
		if err != nil {
			return err
		}

		return errors.New(fmt.Sprintf("error in call: %v", err), errors.WithCode(res.StatusCode))
	}

	return nil
}
