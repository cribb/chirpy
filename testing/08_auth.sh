# Curl tests for chirpy server

BASEURL="http://localhost:8080/"

# Reset database before performing tests 
./00_reset.sh

# testing the ability to return (GET) ALL of the chirps
echo "#### New Test: GET CHIRPS ####"


# create users with passwords
echo -e "#### New Test: Add USER with PASSWORD ####\n"

echo "  ## Create Saul's account with a password"
item01="{\"email\":\"saul@bettercall.com\",\"password\":\"123456\"}"
response01=$(curl -sS -H "Content-Type: application/json" -d "$item01" "${BASEURL}api/users")
userID01=$(printf '%s' "$response01" | jq -r '.id')
echo -e "  Created Saul's account, got back as UUID: ${userID01}\n"

echo "  ## Create Mike's account with a password"
item02="{\"email\":\"mikel@bettercall.com\",\"password\":\"98765\"}"
response02=$(curl -sS -H "Content-Type: application/json" -d "$item02" "${BASEURL}api/users")
userID02=$(printf '%s' "$response02" | jq -r '.id')
echo -e "  Created Mike's account, got back as UUID: ${userID02}\n"

# log in with passwords
echo "  ## Log in with Saul's email and password"
loginResponseSaul=$(curl -H "Content-Type: application/json" -d "$item01" "${BASEURL}api/login")
tokenSaul=$(printf '%s' "$loginResponseSaul" | jq -r '.token')
echo -e "\n"

echo "  ## Log in with Mike's email and password"
loginResponseMike=$(curl -H "Content-Type: application/json" -d "$item02" "${BASEURL}api/login")
tokenMike=$(printf '%s' "$loginResponseMike" | jq -r '.token')
echo -e "\n"

# Create chirps
echo "  #### Saul creates a chirp"
chirpSaul="{\"body\":\"Clearly his taste in women is the same as his taste in lawyers: only the very best... with just a right amount of dirty!\"}"
curl -H "Content-Type: application/json" -H "Authorization: Bearer ${tokenSaul}" -d "$chirpSaul" "${BASEURL}api/chirps"
echo -e "\n"

echo "  #### Mike creates a chirp"
chirpMike="{\"body\":\"No more half-measures.\"}"
curl -H "Content-Type: application/json" -H "Authorization: Bearer ${tokenSaul}" -d "$chirpMike" "${BASEURL}api/chirps"
echo -e "\n"