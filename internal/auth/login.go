package auth

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"timesipooped.fyi/internal/database"
)

type OAuthConfig struct {
	ClientId     string
	ClientSecret string
	RedirectURI  string
	OauthURL     string
	TokenURL     string
	UserInfoURL  string
}

func NewOAuthConfig() *OAuthConfig {
	config := &OAuthConfig{}
	config.ClientId = os.Getenv("CLIENT_ID")
	config.ClientSecret = os.Getenv("CLIENT_SECRET")
	config.RedirectURI = os.Getenv("REDIRECT_URI")
	config.OauthURL = "https://accounts.google.com/o/oauth2/auth"
	config.TokenURL = "https://accounts.google.com/o/oauth2/token"
	config.UserInfoURL = "https://www.googleapis.com/oauth2/v1/userinfo?alt=json"

	return config
}

func HandleLogin(authConf *OAuthConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authUrl := fmt.Sprintf(
			"%s?response_type=code&client_id=%s&redirect_uri=%s&scope=openid%%20email%%20https://www.googleapis.com/auth/userinfo.profile&access_type=offline&prompt=consent",
			authConf.OauthURL,
			authConf.ClientId,
			url.QueryEscape(authConf.RedirectURI),
		)
		http.Redirect(w, r, authUrl, http.StatusFound)
		log.Println("Redirected to auth url")
	}
}

func generateCookie(name string, val string) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    val,
		HttpOnly: true,
		Secure:   false, // TODO TRUE IN PROD
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	}
}

func HandleCallback(authConf *OAuthConfig, db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			http.Error(w, "Missing code parameter", http.StatusBadRequest)
			return
		}

		resp, err := http.PostForm(authConf.TokenURL, url.Values{
			"code":          {code},
			"client_id":     {authConf.ClientId},
			"client_secret": {authConf.ClientSecret},
			"redirect_uri":  {authConf.RedirectURI},
			"grant_type":    {"authorization_code"},
		})

		if err != nil {
			log.Printf("Error exchanging code for token: %v\n", err)
			http.Error(w, "Failed to exchange code", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		var tokenData map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&tokenData); err != nil {
			http.Error(w, "Failed to decode token response", http.StatusInternalServerError)
			return
		}

		accessToken, ok := tokenData["access_token"].(string)
		if !ok {
			log.Println("Error getting access token")
			http.Error(w, "Failed to get access token", http.StatusInternalServerError)
			return
		}
		log.Println("Successfully got access token")

		req, err := http.NewRequest("GET", authConf.UserInfoURL, nil)
		if err != nil {
			log.Printf("Error creating user info request: %v\n", err)
			http.Error(w, "Failed to create user info request", http.StatusInternalServerError)
			return
		}
		req.Header.Set("Authorization", "Bearer "+accessToken)
		client := &http.Client{}
		resp, err = client.Do(req)
		if err != nil {
			log.Printf("Error fetching user info: %v\n", err)
			http.Error(w, "Failed to fetch user info", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		refreshToken, ok := tokenData["refresh_token"].(string)
		if !ok {
			log.Println("Error getting access token")
			http.Error(w, "Failed to get access token", http.StatusInternalServerError)
			return
		}
		log.Println("Successfully got access token")

		var info database.UserInfo
		if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
			log.Printf("Error decoding user info: %v\n", err)
			http.Error(w, "Failed to decode user info", http.StatusInternalServerError)
			return
		}
		log.Println("User data received")

		if !info.VerifiedEmail {
			log.Printf("User <%v> email is not verified", info.Id)
			http.Error(w, "Email is not verified", http.StatusUnauthorized)
			return
		}
		log.Println("User email is verified")

		http.SetCookie(w, generateCookie("access_token", accessToken))
		http.SetCookie(w, generateCookie("user_id", info.Id))
		http.Redirect(w, r, "http://localhost:8080", http.StatusSeeOther)

		err = database.IsNewUser(db, &info, refreshToken, accessToken)
		if err != nil {
			log.Printf("Error checking if user is new: %v\n", err)
			http.Error(w, "Failed to check user status", http.StatusInternalServerError)
			return
		}
	}
}

func VerifyToken(r *http.Request) (*http.Response, string, error) {
	user_id, err := r.Cookie("user_id")
	if err != nil {
		return nil, "", fmt.Errorf("Error getting user_id cookie \nError:\n%v\n\n", err)
	}

	accessTokenCookie, err := r.Cookie("access_token")
	if err != nil {
		return nil, "", fmt.Errorf("Error getting access_token of user <%v>\nError:\n%v\n\n", user_id, err)
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v3/tokeninfo?access_token=" + accessTokenCookie.Value)
	if err != nil {
		return nil, "", fmt.Errorf("Error verifying access token of user <%v>\nError:\n%v\n\n", user_id, err)
	}
	defer resp.Body.Close()

	return resp, user_id.Value, nil
}

func HandleStatus(authConf *OAuthConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
			return
		}

		resp, user_id, err := VerifyToken(r)
		if err != nil {
			log.Printf("Error verifying token: %v\n", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if resp.StatusCode != http.StatusOK {
			log.Printf("Invalid access token of user <%v>\nError:\n%v\n\n", user_id, err)
			w.WriteHeader(http.StatusForbidden)
			return
		}
	}
}
