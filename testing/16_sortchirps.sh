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

# Log in as Walt
echo "  ## Log in as Walt"
loginResponseWalt=$(curl -sS -X POST -H "Content-Type: application/json" -d "$item01" "${BASEURL}api/login")
jwtokenWalt=$(printf '%s' "$loginResponseWalt" | jq -r '.token')
rtokenWalt=$(printf '%s' "$loginResponseWalt" | jq -r '.refresh_token')
WaltID=$(printf '%s' "$loginResponseWalt" | jq -r '.id')
echo -e "  * Walt's login response: $loginResponseWalt"
echo -e "  * Walt's access token: $jwtokenWalt"
echo -e "  * Walt's refresh token: $rtokenWalt"
echo -e "  * Walt's UUID: $WaltID"
echo -e "\n"

# Create Walt's first chirp
echo "  ## Create Walt's first chirp"
chirpWalt="{\"body\":\"I'm the one who knocks!\"}"
responseWalt=$(curl -sS -H "Content-Type: application/json" -H "Authorization: Bearer ${jwtokenWalt}" -d "$chirpWalt" "${BASEURL}api/chirps")
chirp_Walt=$(printf '%s' "$responseWalt" | jq -r '.body')
chirpID_Walt=$(printf '%s' "$responseWalt" | jq -r '.id')
echo -e "Walt's chirp & response: $responseWalt"
echo -e "\n"

# Create Walt's second chirp
echo "  ## Create Walt's second chirp"
chirpWalt="{\"body\":\"Gale!\"}"
responseWalt=$(curl -sS -H "Content-Type: application/json" -H "Authorization: Bearer ${jwtokenWalt}" -d "$chirpWalt" "${BASEURL}api/chirps")
chirp_Walt=$(printf '%s' "$responseWalt" | jq -r '.body')
chirpID_Walt=$(printf '%s' "$responseWalt" | jq -r '.id')
echo -e "Walt's chirp & response: $responseWalt"
echo -e "\n"

# Create Walt's third chirp
echo "  ## Create Walt's third chirp"
chirpWalt="{\"body\":\"Cmon Pinkman\"}"
responseWalt=$(curl -sS -H "Content-Type: application/json" -H "Authorization: Bearer ${jwtokenWalt}" -d "$chirpWalt" "${BASEURL}api/chirps")
chirp_Walt=$(printf '%s' "$responseWalt" | jq -r '.body')
chirpID_Walt=$(printf '%s' "$responseWalt" | jq -r '.id')
echo -e "Walt's chirp & response: $responseWalt"
echo -e "\n"

# Create Walt's fourth chirp
echo "  ## Create Walt's fourth chirp"
chirpWalt="{\"body\":\"Darn that fly, I just wanna cook\"}"
responseWalt=$(curl -sS -H "Content-Type: application/json" -H "Authorization: Bearer ${jwtokenWalt}" -d "$chirpWalt" "${BASEURL}api/chirps")
chirp_Walt=$(printf '%s' "$responseWalt" | jq -r '.body')
chirpID_Walt=$(printf '%s' "$responseWalt" | jq -r '.id')
echo -e "Walt's chirp & response: $responseWalt"
echo -e "\n"

# List chirps from newest to oldest 
echo "  ## List chirps from newest to oldest"
targetURL="${BASEURL}api/chirps?sort=desc"
responseWalt=$(curl -sS -X GET -H "Content-Type: application/json" $targetURL)
chirps_Walt=$(printf '%s' "$responseWalt" | jq -r '.body')
echo -e "  ** Walt's chirps in descending sort order: $chirps_Walt"
echo -e "\n"

# List chirps from oldest to newest
echo "  ## List chirps from oldest to newest"
targetURL="${BASEURL}api/chirps?sort=asc"
responseWalt=$(curl -sS -X GET -H "Content-Type: application/json" $targetURL)
chirps_Walt=$(printf '%s' "$responseWalt" | jq -r '.body')
echo -e "  ** Walt's chirps in ascending sort order: $chirps_Walt"
echo -e "\n"

