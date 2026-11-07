# Curl tests for chirpy server

BASEURL="http://localhost:8080/"

# MUST reset before adding users in upcoming stanza
./00_reset.sh

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
    curl -v -H "Content-Type: application/json" -d "$item" "${BASEURL}api/users"
done
echo -e "\n"

