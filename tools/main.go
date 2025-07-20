package main

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"stravamcp/pkg/client"
	"stravamcp/repo"
	"strings"
)

func main() {
	fmt.Println("🚴 Strava OAuth CLI Tool")
	fmt.Println("========================")

	clientID := getInput("Enter your Strava Client ID: ")
	if clientID == "" {
		fmt.Println("❌ Client ID is required")
		os.Exit(1)
	}

	clientSecret := getInput("Enter your Strava Client Secret: ")
	if clientSecret == "" {
		fmt.Println("❌ Client Secret is required")
		os.Exit(1)
	}

	authURL := generateAuthURL(clientID)
	fmt.Println("\n📋 Step 1: Authorization")
	fmt.Println("Click the following URL to authorize the application:")
	fmt.Printf("🔗 %s\n\n", authURL)

	fmt.Println("After clicking the URL:")
	fmt.Println("1. You'll be redirected to Strava's authorization page")
	fmt.Println("2. Click 'Authorize' to grant permissions")
	fmt.Println("3. You'll be redirected to localhost (this will show an error page - that's normal)")
	fmt.Println("4. Copy the 'code' parameter from the URL in your browser")
	fmt.Println("   Example: http://localhost/exchange_token?code=YOUR_CODE_HERE")
	fmt.Println("")

	code := getInput("📥 Enter the authorization code from the URL: ")
	if code == "" {
		fmt.Println("❌ Authorization code is required")
		os.Exit(1)
	}

	stravaClient := client.NewStravaClient("https://www.strava.com")

	fmt.Println("\n🔄 Exchanging authorization code for access token...")
	token, err := stravaClient.GetTokenFromAuthCode(clientID, clientSecret, code)
	if err != nil {
		fmt.Printf("❌ Error getting token: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n✅ Success! Here are your tokens:")
	fmt.Println("===============================")
	fmt.Printf("🔑 Access Token: %s\n", token.AccessToken)
	fmt.Printf("🔄 Refresh Token: %s\n", token.RefreshToken)
	fmt.Printf("⏰ Expires At: %d\n", token.ExpiresAt)
	fmt.Printf("👤 Athlete: %s %s (%s)\n", token.Athlete.Firstname, token.Athlete.Lastname, token.Athlete.Username)
	err = repo.Save(token, "refresh_token.json")
	if err != nil {
		fmt.Printf("❌ Error saving token: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("\n💾 This token is now saved to {}")
}

func generateAuthURL(clientID string) string {
	baseURL := "https://www.strava.com/oauth/authorize"
	params := url.Values{}
	params.Add("client_id", clientID)
	params.Add("response_type", "code")
	params.Add("redirect_uri", "http://localhost/exchange_token")
	params.Add("approval_prompt", "force")
	params.Add("scope", "read,activity:read_all")
	return fmt.Sprintf("%s?%s", baseURL, params.Encode())
}

func getInput(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}
