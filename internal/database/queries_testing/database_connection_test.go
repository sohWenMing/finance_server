package queries_testing

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	errorutils "github.com/sohWenMing/finance_server/error_utils"
	database "github.com/sohWenMing/finance_server/internal/database/connection"
	"github.com/sohWenMing/finance_server/internal/database/sqlc_generated"
	testhelpers "github.com/sohWenMing/finance_server/test_helpers"
)

var registeredQueries *sqlc_generated.Queries
var createdUserUuids []uuid.UUID

func TestMain(m *testing.M) {
	db, err := database.ConnectToDB("../../../.env")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	createdUserUuids = make([]uuid.UUID, 2)
	// create a space in memory for the create user uuids, start, with space for 2
	registeredQueries = sqlc_generated.New(db)

	code := m.Run()
	clearCreatedUsers(createdUserUuids)
	os.Exit(code)

}

func TestCreateUsers(t *testing.T) {
	firstParams := sqlc_generated.CreateUserParams{
		ID:             uuid.New(),
		IsAdmin:        false,
		Email:          "wenming.soh@gmail.com",
		HashedPassword: "test_password",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	createdUser, err := registeredQueries.CreateUser(context.Background(), firstParams)
	testhelpers.AssertNoError(t, err)
	if err != nil {
		return
	}
	createdUserUuids = append(createdUserUuids, createdUser.ID)

	secondParams := sqlc_generated.CreateUserParams{
		ID:             createdUser.ID,
		IsAdmin:        false,
		Email:          "wenming.soh@temus,com",
		HashedPassword: "test_password",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	_, duplicateIdErr := registeredQueries.CreateUser(context.Background(), secondParams)
	testhelpers.AssertHasError(t, duplicateIdErr)
	if duplicateIdErr == nil {
		return
	}
	isUniqueViolation, pqErr, _ := errorutils.CheckIsUniqueConstraintPqError(duplicateIdErr)
	fmt.Printf("string returned from pqErr: %s", pqErr.Message)
	testhelpers.AssertBool(t, isUniqueViolation, true)

}

func clearCreatedUsers(idList []uuid.UUID) {
	for _, uuid := range idList {
		registeredQueries.DeleteUserById(context.Background(), uuid)
	}
}
