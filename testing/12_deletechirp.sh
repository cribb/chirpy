# Curl tests for chirpy server

BASEURL="http://localhost:8080/"

# Reset database before performing tests 
./00_reset.sh

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

# Create Walt's chirp
echo "  ## Create Walt's chirp"
chirpWalt="{\"body\":\"I did it for me. I liked it. I was good at it. And I was really... I was alive.\"}"
responseWalt=$(curl -sS -H "Content-Type: application/json" -H "Authorization: Bearer ${jwtokenWalt}" -d "$chirpWalt" "${BASEURL}api/chirps")
chirp_Walt=$(printf '%s' "$responseWalt" | jq -r '.body')
chirpID_Walt=$(printf '%s' "$responseWalt" | jq -r '.id')
echo -e "Walt's chirp & response: $responseWalt"
echo -e "\n"

# Check if Walt's chirp exists
echo "  ## Check if Walt's chirp exists"
responseChirpExist=$(curl -sS -X GET -H "Content-Type: application/json" "${BASEURL}api/chirps/${chirpID_Walt}")
echo -e "\n"

# Reject deleting chirp without authentication
echo "  ## Reject deleting chirp without authentication"
responseChirpReject=$(curl -sS -X DELETE -H "Content-Type: application/json" "${BASEURL}api/chirps/${chirpID_Walt}")
echo -e "\n"

# Create Saul's account with a password
echo "  ## Create Saul's account with a password"
item02="{\"email\":\"saul@bettercall.com\",\"password\":\"123456\"}"
response02=$(curl -sS -X POST -H "Content-Type: application/json" -d "$item02" "${BASEURL}api/users")
userID02=$(printf '%s' "$response02" | jq -r '.id')
echo -e "  Created Saul's account, got back as UUID: ${userID02}\n"

# Log in as Saul
echo "  ## Log in as Saul"
loginResponseSaul=$(curl -sS -X POST -H "Content-Type: application/json" -d "$item02" "${BASEURL}api/login")
jwtokenSaul=$(printf '%s' "$loginResponseSaul" | jq -r '.token')
rtokenSaul=$(printf '%s' "$loginResponseSaul" | jq -r '.refresh_token')
echo -e "Saul's access token: $jwtokenSaul"
echo -e "Saul's refresh token: $rtokenSaul"
echo -e "\n"

# Reject deletion by another user
echo "  ## Reject deletion by another user"
responseChirpReject=$(curl -v -sS -X DELETE -H "Content-Type: application/json" -H "Authorization: Bearer ${jwtokenSaul}" "${BASEURL}api/chirps/${chirpID_Walt}")
echo -e "\n"

# Delete the chirp as its author
echo "  ## Delete the chirp as its author"
responseChirpDelete=$(curl -v -sS -X DELETE -H "Content-Type: application/json" -H "Authorization: Bearer ${jwtokenWalt}" "${BASEURL}api/chirps/${chirpID_Walt}")
echo -e "\n"

# Check that the deleted chirp is not found
echo "  ## Check that the deleted chirp is not found"
responseChirpExist=$(curl -v -sS -X GET -H "Content-Type: application/json" "${BASEURL}api/chirps/${chirpID_Walt}")
echo -e "\n"
