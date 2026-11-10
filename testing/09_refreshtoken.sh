# Curl tests for chirpy server

BASEURL="http://localhost:8080/"

# Reset database before performing tests 
./00_reset.sh

# testing the ability to return (GET) ALL of the chirps
echo "#### New Test: GET CHIRPS ####"


# create Saul's accound with passwords
echo -e "#### New Test: Add USER with PASSWORD ####\n"

echo "  ## Create Saul's account with a password"
item01="{\"email\":\"saul@bettercall.com\",\"password\":\"123456\"}"
response01=$(curl -sS -H "Content-Type: application/json" -d "$item01" "${BASEURL}api/users")
userID01=$(printf '%s' "$response01" | jq -r '.id')
echo -e "  Created Saul's account, got back as UUID: ${userID01}\n"

# Get Saul's access and refresh token
echo "  ## Get Saul's access and refresh token"
loginResponseSaul=$(curl -H "Content-Type: application/json" -d "$item01" "${BASEURL}api/login")
tokenSaul=$(printf '%s' "$loginResponseSaul" | jq -r '.token')
echo -e "\n"

# Reject refresh token used to create chirp
echo "  #### Reject refresh token used to create chirp"
chirpSaul="{\"body\":\"Let's just say I know a guy... who knows a guy... who knows another guy.\"}"
curl -H "Content-Type: application/json" -H "Authorization: Bearer ${tokenSaul}" -d "$chirpSaul" "${BASEURL}api/chirps"
echo -e "\n"

# Create chirp with the access token
echo "  #### Create chirp with the access token"
chirpSaul="{\"body\":\"Let's just say I know a guy... who knows a guy... who knows another guy.\"}"
curl -H "Content-Type: application/json" -H "Authorization: Bearer ${tokenSaul}" -d "$chirpSaul" "${BASEURL}api/chirps"
echo -e "\n"
