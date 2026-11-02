# Curl tests for chirpy server

URL="http://localhost:8080/api/validate_chirp"
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
    echo "#### New TEST ####"
    curl -v -X POST -H "Content-Type: application/json" -d "$item" "$URL"
done
echo -e "\n"
