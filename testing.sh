# Curl tests for chirpy server

BASEURL="http://localhost:8080/"


POSTURL=$BASEURL"api/validate_chirp"

POST_TESTS=(
# Empty body
'{"body":""}' 
# Keep a chirp without profanity unchanged
'{"body":"I had something interesting for breakfast"}'
# Reject a chirp over the length limit 
'{"body":"lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."}'
# Filter profanity and ignore extra fields
'{"body":"I hear Mastodon is better than Chirpy. sharbert I need to migrate","extra":"this should be ignored"}'
# Filter multiple profane words regardless of case 
'{"body":"I really need a kerfuffle to go to bed sooner, Fornax !"}'
)

# iterate through all of the elements of the POST_TESTS array
for item in "${POST_TESTS[@]}"; do
    echo -e "\n"
    echo "#### New POST Test ####"
    curl -v -X POST -H "Content-Type: application/json" -d "$item" "$POSTURL"
done
echo -e "\n"

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

