include .env

oauthstep1:
	open https://raindrop.io/oauth/authorize?client_id=$(CLIENT_ID)&redirect_uri=https://oauthdebugger.com/debug

oauthstep2:
	curl -X POST https://raindrop.io/oauth/access_token \
		-H "Content-Type: application/json" \
		-d '{
			"client_id": "$(CLIENT_ID)",
			"client_secret": "$(CLIENT_SECRET)",
			"code": "$(OAUTH_STEP1_CODE)",
			"grant_type":	"authorization_code",
			"redirect_uri": "https://oauthdebugger.com/debug"
		}'