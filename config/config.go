package config

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	envvars "github.com/sohWenMing/finance_server/env_vars"
	"github.com/sohWenMing/finance_server/internal/auth"
	"github.com/sohWenMing/finance_server/internal/database/sqlc_generated"
	usermapping "github.com/sohWenMing/finance_server/mapping/user_mapping"
)

/*
	the config holds is the main holder of all information of the application

	this includes:
	* Qeuries which are sql generated
	* Handlers, which are passed into to server in the main function
*/

type Config struct {
	Queries   *sqlc_generated.Queries
	JwtSecret []byte
}

func (c *Config) RegisterJwtSecret(envPath string) error {
	envVarErr := envvars.LoadEnv(envPath)
	if envVarErr != nil {
		return envVarErr
	}
	secret := os.Getenv("JWT_SECRET")
	c.JwtSecret = []byte(secret)
	return nil
}

func (c *Config) RegisterQueries(db *sql.DB) {
	// at this point, the database should already be loaded, so we should be passing the db type into this function
	c.Queries = sqlc_generated.New(db)
}
func (c *Config) PingHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Header().Set("content-type", "text/plain")
	w.Write([]byte("OK"))
}

func (c *Config) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	bodyJson, shouldReturn := processUserBodyToJSON(r, w)
	if shouldReturn {
		return
	}
	// after decoding, then try to encrypt the password and return hashed password using BCrypt

	hashedPassword, err := auth.GenerateHashedPassword(bodyJson.Password)
	if err != nil {
		writerInteralError(w)
		return
	}

	// if hashing fails, then write 500 and return

	params := sqlc_generated.CreateUserParams{
		ID:             uuid.New(),
		Email:          bodyJson.Email,
		HashedPassword: hashedPassword,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	user, checkErr := c.Queries.CreateUser(context.Background(), params)
	if checkErr != nil {
		fmt.Printf("error in checkErr: %v\n", checkErr)
		writerInteralError(w)
		return
	}
	// if there is a problem with craating the user in the database, then 500 and early return
	createdUserResJson := usermapping.CreatedUserResponse{
		IsSuccess: true,
		UserId:    user.ID.String(),
	}

	bodyString, err := json.Marshal(createdUserResJson)
	if err != nil {
		writerInteralError(w)
		return
	}

	w.WriteHeader(200)
	w.Header().Set("content-type", "application/json")
	w.Write([]byte(bodyString))

}

func (c *Config) LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	bodyJson, shouldReturn := processUserBodyToJSON(r, w)
	if shouldReturn {
		return
	}
	user, dbErr := c.Queries.GetUserByEmail(context.Background(), bodyJson.Email)
	if dbErr != nil {
		writerInteralError(w)
		return
	}
	compareHashErr := auth.CheckHashedPassword(bodyJson.Password, user.HashedPassword)
	if compareHashErr != nil {
		writeUnauthorizedError(w, "email and password do not match")
		return
	}
	w.WriteHeader(200)
	w.Header().Set("content-type", "application/json")

	responseJson := usermapping.LoginResponse{
		IsSuccess:   true,
		AccessToken: "placeholder_access_token",
	}
	resBytes, marshalErr := json.Marshal(responseJson)
	if marshalErr != nil {
		writerInteralError(w)
	}
	w.Write([]byte(resBytes))

}

func (c *Config) FileServerMiddleWare(fileServerHandler http.Handler) http.Handler {
	return fsMiddleWareGenerator(fileServerHandler)
}

func fsMiddleWareGenerator(next http.Handler) http.Handler {
	/*
		http.HandlerFunc takes in a function with a signature of func(http.ResponseWriter http.Request) which
		has a ServeHTTP method, satisfying the http.Handler type. Allows for passing of other handlers into the
		function, which allows the function to act as a middleware layer
	*/
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		next.ServeHTTP(w, r)
	})
}

func writerInteralError(w http.ResponseWriter) {

	w.WriteHeader(500)
	w.Header().Set("content-type", "text/plain")
	w.Write([]byte("500: Internal Error"))
}

func writeUnauthorizedError(w http.ResponseWriter, response string) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Header().Set("content-type", "text/plain")
	w.Write([]byte(response))
}

func processUserBodyToJSON(r *http.Request, w http.ResponseWriter) (usermapping.UserJSON, bool) {
	var bodyJson usermapping.UserJSON
	decoder := json.NewDecoder(r.Body)
	jsonDecodeErr := decoder.Decode(&bodyJson)
	if jsonDecodeErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		return usermapping.UserJSON{}, true
	}
	return bodyJson, false
}
