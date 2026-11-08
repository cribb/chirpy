# Curl tests for chirpy server

BASEURL="http://localhost:8080/"

# Reset database before performing tests 
./00_reset.sh

# create user with password 
echo -e "#### New Test: Add USER with PASSWORD ####\n"

echo "  ## Create Saul's account with a password"
item="{\"email\":\"saul@bettercall.com\",\"password\":\"123456\"}"
response=$(curl -sS -H "Content-Type: application/json" -d "$item" "${BASEURL}api/users")
userID01=$(printf '%s' "$response" | jq -r '.id')
echo -e "  Created Saul's account, got back as UUID: ${userID01}\n"

# log in with password
echo "  ## Log in with Saul's email and password"
curl -H "Content-Type: application/json" -d "$item" "${BASEURL}api/login"
echo -e "\n"

