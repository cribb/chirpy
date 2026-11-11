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
loginResponseSaul=$(curl -sS -H "Content-Type: application/json" -d "$item01" "${BASEURL}api/login")
jwtokenSaul=$(printf '%s' "$loginResponseSaul" | jq -r '.token')
rtokenSaul=$(printf '%s' "$loginResponseSaul" | jq -r '.refresh_token')
echo -e "Saul's access token: $jwtokenSaul"
echo -e "Saul's refresh token: $rtokenSaul"
echo -e "\n"

# Refresh Saul's access token
echo "  #### Refresh Saul's access token"
refreshResponseSaul=$(curl -sS -H "Content-Type: application/json" -H "Authorization: Bearer ${rtokenSaul}" -d "$item01" "${BASEURL}api/refresh")
printf '%s\n' "refreshResponseSaul: $refreshResponseSaul"
new_jwtokenSaul=$(printf '%s' "$refreshResponseSaul" | jq -r ".token")
echo -e "Saul's new access token: $new_jwtokenSaul"
echo -e "\n"

# Reject a refresh token used to create a chirp
echo "  #### Reject a refresh token used to create a chirp"
chirpSaul="{\"body\": \"Let's just say I know a guy... who knows a guy... who knows another guy.\"}"
curl -sS -H "Content-Type: application/json" -H "Authorization: Bearer ${rtokenSaul}" -d "$chirpSaul" "${BASEURL}api/chirps"
echo -e "\n"

# Create a chirp with the access token 
echo "  #### Create a chirp with the access token"
curl -sS -H "Content-Type: application/json" -H "Authorization: Bearer ${jwtokenSaul}" -d "$chirpSaul" "${BASEURL}api/chirps"
echo -e "\n"

