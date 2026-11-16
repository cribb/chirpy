# Curl tests for chirpy server

BASEURL="http://localhost:8080/"

# Reset database before performing tests 
./00_reset.sh

# testing the refresh tokens
echo "#### New Test: AUTHORIZATION ####"

# Create Walt's account with a password
echo "  ## Create Walt's account with a password"
item01="{\"email\":\"walt@breakingbad.com\",\"password\":\"123456\"}"
response01=$(curl -sS -X POST -H "Content-Type: application/json" -d "$item01" "${BASEURL}api/users")
userID01=$(printf '%s' "$response01" | jq -r '.id')
echo -e "  Created Walt's account, got back as UUID: ${userID01}\n"

# Log in as Walt
echo "  ## Log in as Walt"
loginResponseWalt=$(curl -sS -X POST -H "Content-Type: application/json" -d "$item01" "${BASEURL}api/login")
jwtokenWalt=$(printf '%s' "$loginResponseWalt" | jq -r '.token')
rtokenWalt=$(printf '%s' "$loginResponseWalt" | jq -r '.refresh_token')
echo -e "Walt's access token: $jwtokenWalt"
echo -e "Walt's refresh token: $rtokenWalt"
echo -e "\n"

# Update authenticated Walt's email and password 
echo "  ## Update the authenticated Walt's email and password"
new_WaltUnPw="{\"email\":\"walter@breakingbad.com\",\"password\":\"losPollosHermanos\"}" # the request body
loginResponseWalt=$(curl -sS -X PUT -H "Content-Type: application/json" -H "Authorization: Bearer ${jwtokenWalt}" -d "$new_WaltUnPw" "${BASEURL}api/users")
printf 'response from Walt PUT: %s\n' "$loginResponseWalt"
EmailWalt=$(printf '%s' "$loginResponseWalt" | jq -r '.email')
echo -e "Walt's new email: $EmailWalt"

# Reject a user update without authentication (empty/no header)
echo "  ## Reject a user update without authentication"
bad_WaltUnPw="{\"email\":\"walter@breakingbad.com\",\"password\":\"j3ssePinkM@nCantCook\"}" # the request body
EmptyJWtoken=""
loginResponseBad=$(curl -v -sS -X PUT -H "Content-Type: application/json" -H "Authorization: Bearer ${EmptyJWtoken}" -d "$bad_WaltUnPw" "${BASEURL}api/users")
RejectWaltEmail=$(printf '%s' "$loginResponseBad" | jq -r '.email')
RejectWaltToken=$(printf '%s' "$loginResponseBad" | jq -r '.token')
echo -e "loginResponseBad: $loginResponseBad"
echo -e "Walt's bad email: $RejectWaltEmail"
echo -e "Walt's refresh token: $RejectWaltToken"
echo -e "\n"

