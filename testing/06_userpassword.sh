# Curl tests for chirpy server

BASEURL="http://localhost:8080/"

# Reset database before performing tests 
./00_reset.sh


# create user with password 
echo -e "\n"
echo "#### New Test: Add USER with PASSWORD ####"

echo "#### Create Saul's account with a password"
item="{\"email\":\"saul@bettercall.com\",\"password\":\"123456\"}"
response=$(curl -v -sS -H "Content-Type: application/json" -d "$item" "${BASEURL}api/users")
echo "response: ${response}"
userID01=$(printf '%s' "$response" | jq -r '.id')
echo "  Created Saul's account, got back as UUID: ${userID01}"
echo -e "\n"


echo "  #### Log in with Saul's email and password"
item="{\"body\":\"If you're committed enough, you can make any story work.\",\"user_id\":\"${userID01}\"}"
curl -v -X POST -H "Content-Type: application/json" -d "$item" "${BASEURL}api/login"
echo -e "\n"

