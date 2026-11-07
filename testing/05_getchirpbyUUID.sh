# Curl tests for chirpy server

BASEURL="http://localhost:8080/"

# Reset database before performing tests 
./00_reset.sh

# testing the ability to return (GET) one chirp via chirpID
# echo -e "\n"
echo -e "\n#### New Test: GET CHIRP BY ITS UUID ####"

## Create Saul's account
echo "  #### Create Saul's account"
item='{"email":"saul@bettercall.com"}'
url="http://localhost:8080/api/users"
response=$(curl -v -sS -X POST -H "Content-Type: application/json" -d "$item" "${BASEURL}api/users")
userID01=$(printf '%s' "$response" | jq -r '.id')
echo "  Created Saul's account, got back as UUID: ${userID01}"
echo -e "\n"

## Create Saul's chirp
echo "  #### Create Saul's chirp"
item="{\"body\":\"I'm gonna be a damn good developer, and people are gonna know about it.\",\"user_id\":\"${userID01}\"}"
response=$(curl -v -sS -X POST -H "Content-Type: application/json" -d "$item" "${BASEURL}api/chirps")
chirpID01=$(printf '%s' "$response" | jq -r '.id')
echo "  Created Saul's chirp, got back as its UUID: ${chirpID01}"
echo -e "\n"

## Get the chirp by its UUID
echo "  #### GET Saul's chirp by its UUID"
curl -v -sS -X GET -H "Content-Type: application/json" "${BASEURL}api/chirps/${chirpID01}"
echo -e "\n"

## Try getting an unknown chirp
echo "  #### Try getting an unknown chirp"
BadChirpID="123e4567-e89b-12d3-a456-426614174059"
curl -v -sS -X GET -H "Content-Type: application/json" "${BASEURL}api/chirps/${BadChirpID}"

