# Strava Oauth CLI Tool

A command-line tool for authenticating with the Strava API using OAuth 2.0. This tool guides you through the OAuth flow and saves your access tokens for future API calls.

## Prerequisites

Before you can use this tool, you need to create a Strava application to get your Client ID and Client Secret.

### Creating a Strava Application

1. **Go to Strava Developers**: Visit [https://www.strava.com/settings/api](https://www.strava.com/settings/api)

2. **Create Your Application**:
    - Click "Create & Manage Your App"
    - Fill in the application details:
        - **Application Name**: Choose any name (e.g., "My Strava CLI Tool")
        - **Category**: Choose the most appropriate category
        - **Club**: Leave blank unless you're creating for a specific club
        - **Website**: You can use `http://localhost` for testing
        - **Authorization Callback Domain**: Set this to `localhost`
    - Agree to the API Agreement
    - Click "Create"

3. **Get Your Credentials**:
    - After creation, you'll see your application details
    - Note down your **Client ID** (publicly visible)
    - Click "Show" next to **Client Secret** and note it down (keep this private!)

## Installation

### Step 1: Run the Tool

**Clone or create the project**:
```bash
git clone https://github.com/lyledean1/strava-mcp
```

```bash
go run tools/main.go
```

The tool will display:
```
🚴 Strava OAuth CLI Tool
========================
```

### Step 2: Enter Your Credentials

When prompted, enter your Strava application credentials:

```
Enter your Strava Client ID: [YOUR_CLIENT_ID]
Enter your Strava Client Secret: [YOUR_CLIENT_SECRET]
```

### Step 3: Authorize the Application

The tool will generate an authorization URL:

```
📋 Step 1: Authorization
Click the following URL to authorize the application:
🔗 https://www.strava.com/oauth/authorize?client_id=...
```

1. **Click the URL** or copy-paste it into your browser
2. **Login to Strava** if you're not already logged in
3. **Review the permissions** the app is requesting:
    - Read your profile information
    - Read all your activities
4. **Click "Authorize"** to grant permissions

### Step 4: Get the Authorization Code

After clicking "Authorize":

1. **You'll be redirected** to `http://localhost/exchange_token?code=...`
2. **The page will show an error** - this is normal! The important part is the URL
3. **Copy the code parameter** from the URL

Example URL:
```
http://localhost/exchange_token?code=abc123def456ghi789
```
Copy: `abc123def456ghi789`

### Step 5: Complete Authentication

Paste the authorization code when prompted:

```
📥 Enter the authorization code from the URL: abc123def456ghi789
```

### Step 6: Success!

The tool will exchange your code for tokens and display:

```
✅ Success! Here are your tokens:
===============================
🔑 Access Token: your_access_token_here
🔄 Refresh Token: your_refresh_token_here
⏰ Expires At: 1234567890
👤 Athlete: John Doe (johndoe)

💾 Save these tokens securely - you'll need them to make API calls!
```

The tokens are also automatically saved to `refresh_token.json` in the current directory.

## What You Get

After successful authentication, you'll receive:

- **Access Token**: Use this to make API calls (expires in ~6 hours)
- **Refresh Token**: Use this to get new access tokens when they expire
- **Expires At**: Unix timestamp when the access token expires
- **Athlete Info**: Basic information about the authenticated user

## Permissions (Scopes)

This tool requests the following permissions:

- `read`: Access to read your profile information
- `activity:read_all`: Access to read all your activities (public and private)

## Security Notes

- **Keep your Client Secret private** - never share it or commit it to version control
- **Store your tokens securely** - they provide access to your Strava data
- **Refresh tokens don't expire** - treat them like passwords
- **Access tokens expire in ~6 hours** - use refresh tokens to get new ones

## API Documentation

For more information about the Strava API:
- **API Reference**: [https://developers.strava.com/docs/reference/](https://developers.strava.com/docs/reference/)
- **Getting Started**: [https://developers.strava.com/docs/getting-started/](https://developers.strava.com/docs/getting-started/)
- **Rate Limits**: [https://developers.strava.com/docs/rate-limits/](https://developers.strava.com/docs/rate-limits/)