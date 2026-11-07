# chirpy
boot.dev course project for http-server construction

## Notes

- Locally-built http/web server runs on port 8080
- output accesible by browser
- Example url `http://localhost:8080/app/`

```bash
go build -o out && ./out  # build and run in foreground
# ctrl-c to break

go build -o out && ./out & # build and run in background

jobs           # find jobid for bkgd server
kill %<jobid>  # where <jobid> is a number from output of "jobs", likely "1"
kill %1        # literal example

```


## Postgres

### Connection string

    "postgres://postgres:postgres@localhost:5432/chirpy"

### Packages

    go get <package>

#### Package Links

    github.com/alexedwards/argon2id
	github.com/google/uuid v1.6.0
	github.com/joho/godotenv v1.5.1 
	github.com/lib/pq v1.12.3 
	golang.org/x/crypto v0.14.0
	golang.org/x/sys v0.13.0


