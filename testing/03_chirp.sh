# Curl tests for chirpy server

BASEURL="http://localhost:8080/"

# Reset database before performing tests 
./00_reset.sh


# testing the user's ability to chirp
echo -e "\n"
echo "#### New Test: USER CHIRP ####"

echo "#### Create Saul's account"
item='{"email":"saul@bettercall.com"}'
response=$(curl -v -sS -X POST -H "Content-Type: application/json" -d "$item" "${BASEURL}api/users")
userID01=$(printf '%s' "$response" | jq -r '.id')
echo "  Created Saul's account, got back as UUID: ${userID01}"
echo -e "\n"

echo "  #### Create a chirp with author ID"
item="{\"body\":\"If you're committed enough, you can make any story work.\",\"user_id\":\"${userID01}\"}"
curl -v -X POST -H "Content-Type: application/json" -d "$item" "${BASEURL}api/chirps"
echo -e "\n"

