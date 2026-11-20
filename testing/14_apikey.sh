# Curl tests for chirpy server

BASEURL="http://localhost:8080/"

# Reset database before performing tests 
./00_reset.sh

# Create Walt's account with a password (sans Chirpy Red)
echo "## Create Walt's account with a password (sans Chirpy Red)"
item01="{\"email\":\"walt@breakingbad.com\",\"password\":\"123456\"}"
response01=$(curl -sS -X POST -H "Content-Type: application/json" -d "$item01" "${BASEURL}api/users")
userID01=$(printf '%s' "$response01" | jq -r '.id')
echo -e "  * item01: ${item01}"
echo -e "  * response01: ${response01}"
echo -e "  Created Walt's account, got back as UUID: ${userID01}\n"

# Reject a webhook without an API key
echo "## Reject a webhook without an API key"
badwebhook="{\"data\":{\"user_id\": \"${userID01}\"},\"event\":\"user.payment_failed\"}"
badResponse=$(curl -v -sS -X POST -H "Content-Type: application/json" -d "$badwebhook" "${BASEURL}api/polka/webhooks")
echo -e "  * badwebhook=${badwebhook}"
echo -e "  * badresponse= ${badresponse}"
echo -e "\n"

# Accept an upgrade webhook with the API key 
echo "## Accept an upgrade webhook with the API key "
KEY="f271c81ff7084ee5b99a5091b42d486e"
goodwebhook="{\"event\": \"user.upgraded\",\"data\":{\"user_id\":\"${userID01}\"}}"
goodResponse=$(curl -sS -X POST -H "Content-Type: application/json" -H "Authorization: ApiKey $KEY" -d "$goodwebhook" "${BASEURL}api/polka/webhooks")
echo -e "  * goodwebhook=${goodwebhook}"
echo -e "  * goodResponse=${goodResponse}\n"

# Check if the user now has Chirpy Red
echo "## Check the user now has Chirpy Red"
loginResponseWalt=$(curl -sS -X POST -H "Content-Type: application/json" -d "$item01" "${BASEURL}api/login")
jwtokenWalt=$(printf '%s' "$loginResponseWalt" | jq -r '.token')
rtokenWalt=$(printf '%s' "$loginResponseWalt" | jq -r '.refresh_token')
redTF=$(printf '%s' "$loginResponseWalt" | jq -r '.is_chirpy_red')
echo -e "  * Walt's login response: $loginResponseWalt"
# echo -e "  * Walt's access token: $jwtokenWalt"
# echo -e "  * Walt's refresh token: $rtokenWalt"
echo -e "  * Walt's red chirpy status: $redTF\n"
