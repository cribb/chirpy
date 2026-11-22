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

# Create Skyler's account with a password (sans Chirpy Red)
echo "## Create Skyler's account with a password (sans Chirpy Red)"
item02="{\"email\":\"skyler@breakingbad.com\",\"password\":\"000111\"}"
response02=$(curl -sS -X POST -H "Content-Type: application/json" -d "$item02" "${BASEURL}api/users")
userID02=$(printf '%s' "$response02" | jq -r '.id')
echo -e "  * item02: ${item02}"
echo -e "  * response02: ${response02}"
echo -e "  * Created Skyler's account, got back as UUID: ${userID02}\n"

# Log in as Skyler
echo "  ## Log in as Skyler"
loginResponseSkyler=$(curl -sS -X POST -H "Content-Type: application/json" -d "$item02" "${BASEURL}api/login")
jwtokenSkyler=$(printf '%s' "$loginResponseSkyler" | jq -r '.token')
rtokenSkyler=$(printf '%s' "$loginResponseSkyler" | jq -r '.refresh_token')
SkylerID=$(printf '%s' "$loginResponseSkyler" | jq -r '.id')
echo -e "  * Skyler's login response: $loginResponseSkyler"
echo -e "  * Skyler's access token: $jwtokenSkyler"
echo -e "  * Skyler's refresh token: $rtokenSkyler"
echo -e "  * Skylar's UUID: $SkylerID"
echo -e "\n"

# Create Skylers's chirp
echo "  ## Create Skylers's chirp"
chirpSkyler="{\"body\":\"Mr President....\"}"
responseSkyler=$(curl -sS -H "Content-Type: application/json" -H "Authorization: Bearer ${jwtokenSkyler}" -d "$chirpSkyler" "${BASEURL}api/chirps")
chirp_Skyler=$(printf '%s' "$responseSkyler" | jq -r '.body')
chirpID_Skyler=$(printf '%s' "$responseSkyler" | jq -r '.id')
echo -e "  * Skyler's chirp & response: $responseSkyler"
echo -e "\n"

# List only Walt chirps
echo ie "## List only Walt chirps"
targetURL="${BASEURL}api/chirps?author_id=${WaltID}"
WaltsChirps=$(curl -sS -X GET -H "Content-Type: application/json" $targetURL)
echo -e "  * Target URL: $targetURL"
echo -e "  * Walt's chirps: $WaltsChirps"
echo -e "\n"

# List only Skyler chirps
echo -e "## List only Skyler chirp(s)"
targetURL="${BASEURL}api/chirps?author_id=${SkylerID}"
SkylersChirps=$(curl -sS -X GET -H "Content-Type: application/json" $targetURL)
echo -e "  * Target URL: $targetURL"
echo -e "  * Skyler's chirps: $SkylersChirps"
echo -e "\n"

# Get all chirps
echo -e "## List all chirp(s)"
targetURL="${BASEURL}api/chirps?author_id="
AllChirps=$(curl -sS -X GET -H "Content-Type: application/json" $targetURL)
echo -e "  * All chirps: $AllChirps"
echo -e "\n"

