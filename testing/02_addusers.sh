# Curl tests for chirpy server

BASEURL="http://localhost:8080/"


POSTURL=$BASEURL"api/validate_chirp"

# MUST reset before adding users in upcoming stanza
RESET_URL=$BASEURL"admin/reset"
echo -e "\n"
echo "#### Reset server Test ####"
curl -v -X POST -H "Content-Type: application/json" "$RESET_URL"

USER_URL=$BASEURL"api/users"
USER_TESTS=(
# Create the first user
'{"email":"mloneusk@example.co"}'
# Create the second user
'{"email":"dackjorsey@example.co"}'
)

# iterate through all of the elements of the USER_TESTS array
for item in "${USER_TESTS[@]}"; do
    echo -e "\n"
    echo "#### New USER Test ####"
    curl -v -X POST -H "Content-Type: application/json" -d "$item" "$USER_URL"
done
echo -e "\n"

# testing the ability to chirp
echo "#### New Test: CHIRP ####"
echo "  #### Reset database"
url="http://localhost:8080/admin/reset"
curl -v -X POST -H "Content-Type: application/json" "$url"
echo -e "\n"

echo "  #### Create Saul's account"
item='{"email":"saul@bettercall.com"}'
url="http://localhost:8080/api/users"
response=$(curl -v -sS -X POST -H "Content-Type: application/json" -d "$item" "$url")
userID01=$(printf '%s' "$response" | jq -r '.id')
echo -e "\n"

echo "  #### Create a chirp with author ID"
item="{\"body\":\"If you're committed enough, you can make any story work.\",\"user_id\":\"${userID01}\"}"
url="http://localhost:8080/api/chirps"
curl -v -X POST -H "Content-Type: application/json" -d "$item" "$url"
echo -e "\n"

# testing the ability to return (GET) ALL of the chirps
echo "#### New Test: GET CHIRPS ####"

echo "  #### Reset database"
url="http://localhost:8080/admin/reset"
curl -v -X POST -H "Content-Type: application/json" "$url"
echo -e "\n"

echo "  #### Create Saul's account"
item='{"email":"saul@bettercall.com"}'
url="http://localhost:8080/api/users"
response=$(curl -v -sS -X POST -H "Content-Type: application/json" -d "$item" "$url")
userID01=$(printf '%s' "$response" | jq -r '.id')
echo -e "\n"

# Now make the chirps...
url="http://localhost:8080/api/chirps"

echo "  #### Create chirp 1"
item="{\"body\":\"If you're committed enough, you can make any story work.\",\"user_id\":\"${userID01}\"}"
curl -v -sS -X POST -H "Content-Type: application/json" -d "$item" "$url"
echo -e "\n"

echo "  #### Create chirp 2"
item="{\"body\":\"I once told a woman I was Kevin Costner, and it worked because I believed it.\",\"user_id\":\"${userID01}\"}"
curl -v -sS -X POST -H "Content-Type: application/json" -d "$item" "$url"
echo -e "\n"

# Pull the chirp list
echo "  #### List chirps in creation order"
curl -v -X GET -H "Content-Type: application/json" "$url"
echo -e "\n"

# testing the ability to return (GET) one chirp via chirpID
echo "#### New Test: GET CHIRP BY CHIRP ID ####"

echo "  #### Reset database"
url="http://localhost:8080/admin/reset"
curl -v -X POST -H "Content-Type: application/json" "$url"
echo -e "\n"

## Create Saul's account
echo "  #### Create Saul's account"
item='{"email":"saul@bettercall.com"}'
url="http://localhost:8080/api/users"
response=$(curl -v -sS -X POST -H "Content-Type: application/json" -d "$item" "$url")
userID01=$(printf '%s' "$response" | jq -r '.id')
echo -e "\n"

## Create Saul's chirp
echo "  #### Create Saul's chirp"
item="{\"body\":\"I'm gonna be a damn good developer, and people are gonna know about it.\",\"user_id\":\"${userID01}\"}"
response=$(curl -v -sS -X POST -H "Content-Type: application/json" -d "$item" "$url")
chirpID01=$(printf '%s' "$response" | jq -r '.id')
echo -e "\n"

## Get the chirp by its UUID
echo "  #### GET Saul's chirp by its UUID"
url="http://localhost:8080/api/chirps"
curl -v -sS -X GET -H "Content-Type: application/json" "$url/$chirpID01"
echo -e "\n"

## Try getting an unknown chirp
echo "  #### Try getting an unkown chirp"
BadChirpID="123e4567-e89b-12d3-a456-426614174059"
curl -v -sS -X GET -H "Content-Type: application/json" "$url/$BadChirpID"

