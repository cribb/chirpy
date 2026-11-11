# Curl tests for chirpy server

BASEURL="http://localhost:8080/"

# Reset database before performing tests 
./00_reset.sh

# testing the refresh tokens
echo "#### New Test: REFRESH TOKENS ####"


# create Saul's accound with passwords
echo "  ## Create Saul's account with a password"
item01="{\"email\":\"saul@bettercall.com\",\"password\":\"123456\"}"
response01=$(curl -sS -H "Content-Type: application/json" -d "$item01" "${BASEURL}api/users")
userID01=$(printf '%s' "$response01" | jq -r '.id')
echo -e "  Created Saul's account, got back as UUID: ${userID01}\n"

# Get Saul's access and refresh token
echo "  ## Get Saul's access and refresh token"
loginResponseSaul=$(curl -H "Content-Type: application/json" -d "$item01" "${BASEURL}api/login")
jwtokenSaul=$(printf '%s' "$loginResponseSaul" | jq -r '.token')
rtokenSaul=$(printf '%s' "$loginResponseSaul" | jq -r '.refresh_token')
echo -e "Saul's access token: $jwtokenSaul"
echo -e "Saul's refresh token: $rtokenSaul"
echo -e "\n"

# Refresh Saul's access token
echo "  #### Refresh Saul's access token"
refreshResponseSaul=$(curl -H "Content-Type: application/json;Authorization: Bearer ${rtokenSaul}" -d "$item01" "${BASEURL}api/refresh")
new_jwtokenSaul=$(printf '%s', "$refreshResponseSaul" | jq -r ".token")
new_rtokenSaul=$(printf '%s', "$refreshResponseSaul" | jq -r ".refresh_token")
echo -e "Saul's new access token: $new_jwtokenSaul"
echo -e "Saul's new refresh token: $new_rtokenSaul"

# Create chirp with the refreshed access token
echo "  #### Create chirp with the access token"
chirpSaul="{\"body\":\"I'm the guy who's gonna win you this case.\"}"
curl -H "Content-Type: application/json;Authorization: Bearer ${new_rtokenSaul}" -d "$chirpSaul" "${BASEURL}api/chirps"
echo -e "\n"

