# Curl tests for chirpy server

BASEURL="http://localhost:8080/"

# Reset database
echo -e "\n"
echo "#### Reset database"
curl -X POST -H "Content-Type: application/json" "http://localhost:8080/admin/reset"



