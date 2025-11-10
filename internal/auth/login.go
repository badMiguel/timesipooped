package auth

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	_ "modernc.org/sqlite"
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
	var isSecure bool
	if os.Getenv("IS_SECURE") == "true" {
		isSecure = true
	}
	return &http.Cookie{
		Name:     name,
		Value:    val,
		HttpOnly: true,
		Secure:   isSecure,
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

		clientUrl := "https://poop.badmiguel.com"
		if os.Getenv("IS_DEVELOPMENT") == "true" {
			clientUrl = "http://localhost:8000"
		}
		http.Redirect(w, r, clientUrl, http.StatusSeeOther)

		err = database.IsNewUser(db, &info, refreshToken, accessToken)
		if err != nil {
			log.Printf("Error checking if user is new: %v\n", err)
			http.Error(w, "Failed to check user status", http.StatusInternalServerError)
			return
		}
	}
}

func RefreshAccessToken(w http.ResponseWriter, userId string, authConf *OAuthConfig, db *sql.DB) error {
	log.Println("Initiating access token refresh...")
	var refreshToken string
	for i := 0; i < 10; i++ {
		err := db.QueryRow(
			`SELECT refreshToken FROM users WHERE userId = ?;`, userId,
		).Scan(&refreshToken)
		if err != nil {
			log.Printf("Error when querying user <%s>... Trying again(%d)", userId, i)
			if i == 9 {
				return fmt.Errorf("Failed to query user <%s>: %v\n", userId, err)
			}
			time.Sleep(time.Duration(i/2) * time.Second)
			continue
		}
		break
	}

	resp, err := http.PostForm("https://oauth2.googleapis.com/token", url.Values{
		"client_id":     {authConf.ClientId},
		"client_secret": {authConf.ClientSecret},
		"refresh_token": {refreshToken},
		"grant_type":    {"refresh_token"},
	})
	if err != nil || resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Error refreshing access token: %v\n", err)
	}

	var tokenData map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&tokenData); err != nil {
		return fmt.Errorf("Error decoding response: %v\n", err)
	}

	accessToken, ok := tokenData["access_token"].(string)
	if !ok {
		return fmt.Errorf("Decoded access token is invalid\n")
	}

	for i := 0; i < 10; i++ {
		_, err := db.Exec(
			`UPDATE users SET accessToken = ? WHERE userId = ?`, accessToken, userId,
		)
		if err != nil {
			log.Printf("Error updating new access token... Trying again(%d)\n", i)
			if i == 9 {
				return fmt.Errorf(
					"Error updating new access token of user <%s>\n", userId,
				)
			}
			time.Sleep(time.Duration(i/2) * time.Second)
			continue
		}
		break
	}

	http.SetCookie(w, generateCookie("access_token", accessToken))
	log.Println("Successfully refreshed access token")

	return nil
}

func logoutHelper(w http.ResponseWriter) {
	var isSecure bool
	if os.Getenv("IS_SECURE") == "true" {
		isSecure = true
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Secure:   isSecure,
		HttpOnly: true,
		Expires:  time.Unix(0, 0),
		Path:     "/",
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "user_id",
		Value:    "",
		Secure:   isSecure,
		HttpOnly: true,
		Expires:  time.Unix(0, 0),
		Path:     "/",
	})
}

func VerifyToken(r *http.Request) (*http.Response, string, error) {
	user_id, err := r.Cookie("user_id")
	if err != nil {
		return nil, "", fmt.Errorf("Error getting user_id cookie: %v\n", err)
	}

	accessTokenCookie, err := r.Cookie("access_token")
	if err != nil {
		return nil, "", fmt.Errorf("Error getting access_token of user <%v>: %v\n", user_id, err)
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v3/tokeninfo?access_token=" + accessTokenCookie.Value)
	if err != nil {
		return nil, "", fmt.Errorf("Error verifying access token of user <%v>: %v\n", user_id, err)
	}
	defer resp.Body.Close()

	return resp, user_id.Value, nil
}

func HandleStatus(authConf *OAuthConfig, db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
			return
		}

		resp, user_id, err := VerifyToken(r)
		if err != nil {
			log.Printf("Error verifying token: %v\n", err)

			logoutHelper(w)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if resp.StatusCode != http.StatusOK {
			log.Printf("Invalid access token of user <%v>\n", user_id)
			log.Println("Will attempt to refresh access token")

			err := RefreshAccessToken(w, user_id, authConf, db)
			if err != nil {
				logoutHelper(w)
				log.Println(err)
				w.WriteHeader(http.StatusForbidden)

			}
		}
	}
}

func HandleLogout(w http.ResponseWriter, r *http.Request) {
	logoutHelper(w)
	w.WriteHeader(http.StatusOK)
}
