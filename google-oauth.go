package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"io"
	"net/http"
)

var (
	googleOauthConfig *oauth2.Config
	oauthStateString  = uuid.New().String()
)

// Zarigatongy youtube channel
// Demo https://8gwifi.org

func init() {
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:9090/auth/callback",
		ClientID:     "818188531951-adsgdat3ps5o9rds4o5n1hfo7amkovq1.apps.googleusercontent.com",
		ClientSecret: "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",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/o/oauth2/auth",
			TokenURL: "https://accounts.google.com/o/oauth2/token",
		},
	}

}

func main() {
	r := gin.Default()

	r.GET("/login", handleLogin)
	r.GET("/auth/callback", handleCallback)
	r.GET("/profile", handleProfile)

	r.Run(":9090")
}

func handleLogin(c *gin.Context) {
	url := googleOauthConfig.AuthCodeURL(oauthStateString)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func handleCallback(c *gin.Context) {
	state := c.Query("state")
	if state != oauthStateString {
		fmt.Println("invalid oauth state")
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	code := c.Query("code")
	fmt.Println(code)
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		fmt.Println("code exchange failed")
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	// Use the token to get user details
	client := googleOauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		fmt.Println("failed to get userinfo", err)
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}
	defer resp.Body.Close()

	var userInfo map[string]interface{}
	if err := decodeJSON(resp.Body, &userInfo); err != nil {
		fmt.Println("failed to decode JSON", err)
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, "/profile")
}

func handleProfile(c *gin.Context) {
	// Implement logic to fetch user profile using the obtained token
	// You can make requests to the Google API using the token
	// and retrieve user details.

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile endpoint",
	})
}

func decodeJSON(reader io.Reader, v interface{}) error {
	return json.NewDecoder(reader).Decode(v)
}
