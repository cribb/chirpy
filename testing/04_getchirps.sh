# Curl tests for chirpy server

BASEURL="http://localhost:8080/"

# Reset database before performing tests 
./00_reset.sh

# testing the ability to return (GET) ALL of the chirps
echo "#### New Test: GET CHIRPS ####"

# Create user account
echo "#### Create Saul's account"
item='{"email":"saul@bettercall.com"}'
response=$(curl -v -sS -X POST -H "Content-Type: application/json" -d "$item" "${BASEURL}api/users")
userID01=$(printf '%s' "$response" | jq -r '.id')
echo "  Created Saul's account, got back as UUID: ${userID01}"
echo -e "\n"

# Now make the chirps...
echo "  #### Create chirp 1 (with userID)"
item="{\"body\":\"If you're committed enough, you can make any story work.\",\"user_id\":\"${userID01}\"}"
curl -v -sS -X POST -H "Content-Type: application/json" -d "$item" "${BASEURL}api/chirps"
echo -e "\n"

echo "  #### Create chirp 2 (with userID)"
item="{\"body\":\"I once told a woman I was Kevin Costner, and it worked because I believed it.\",\"user_id\":\"${userID01}\"}"
curl -v -sS -X POST -H "Content-Type: application/json" -d "$item" "${BASEURL}api/chirps"
echo -e "\n"

# Pull the chirp list
echo "  #### List chirps in creation order"
curl -v -X GET -H "Content-Type: application/json" "${BASEURL}api/chirps"
echo -e "\n"
