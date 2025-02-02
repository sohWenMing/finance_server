package integrationtesting

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/sohWenMing/finance_server/internal/auth"
	tokenmapping "github.com/sohWenMing/finance_server/mapping"
	usermapping "github.com/sohWenMing/finance_server/mapping/user_mapping"
	testhelpers "github.com/sohWenMing/finance_server/test_helpers"
)

func TestPing(t *testing.T) {

	path := fmt.Sprintf("%s/ping", basePath)

	req, err := http.NewRequest(http.MethodGet, path, nil)
	testhelpers.AssertNoError(t, err)

	res, err := client.Do(req)
	testhelpers.AssertNoError(t, err)

	testhelpers.AssertIntVals(t, res.StatusCode, 200)

}

func TestFileServer(t *testing.T) {

	path := fmt.Sprintf("%s/app/", basePath)

	req, err := http.NewRequest(http.MethodGet, path, nil)
	testhelpers.AssertNoError(t, err)

	res, err := client.Do(req)
	testhelpers.AssertNoError(t, err)

	resBody, err := io.ReadAll(res.Body)
	testhelpers.AssertNoError(t, err)
	if !strings.Contains(string(resBody), "let's get this") {
		t.Errorf("%s", string(resBody))
	}
}

func TestCreateUserHandler(t *testing.T) {
	TestConfig.SetJWTValidDuration(10 * time.Minute)
	createUserReqBody, responseJson, shouldReturn := createAndRetrieveUser(t, createUserPath)
	if shouldReturn {
		return
	}
	userIdUUID, err := uuid.Parse(responseJson.UserId)

	testhelpers.AssertNoError(t, err)
	if err != nil {
		return
	}
	defer TestConfig.Queries.DeleteUserById(context.Background(), userIdUUID)
	// at this point if there is no error and the test function has not returned, user has already been created. defer the deleting of this user

	retrievedUser, err := TestConfig.Queries.GetUserById(context.Background(), userIdUUID)
	testhelpers.AssertNoError(t, err)
	if err != nil {
		return
	}
	testhelpers.AssertStringVals(t, retrievedUser.Email, createUserReqBody.Email)
	//retrieve the user from the database, assert that the email of the user is as per what was entered

	/*
		now run another test to ensure that another user cannot be created with the same email. check function createAndRetrieveUser below, to
		see the definition of the initial struct
	*/
	userWithDupEmail := usermapping.UserJSON{
		Email:    "wenming.soh@gmail.com",
		Password: "password",
	}

	bodyString, err := json.Marshal(userWithDupEmail)
	testhelpers.AssertNoError(t, err)
	if err != nil {
		return
	}
	req, err := http.NewRequest(http.MethodPost, createUserPath, bytes.NewReader(bodyString))
	testhelpers.AssertNoError(t, err)
	if err != nil {
		return
	}
	res, err := client.Do(req)
	if err != nil {
		return
	}
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}
	testhelpers.AssertStringVals(t, string(bodyBytes), "email wenming.soh@gmail.com is already being used")
}

func TestLoginUserHandler(t *testing.T) {
	TestConfig.SetJWTValidDuration(10 * time.Minute)

	type testStruct struct {
		name         string
		isExpectPass bool
	}

	tests := []testStruct{
		{
			name:         "test login user should pass",
			isExpectPass: true,
		},
		{
			name:         "test login user should fail",
			isExpectPass: false,
		},
	}
	for _, test := range tests {
		runLoginTest(t, createUserPath, test.isExpectPass)
	}
}

func TestJWTValidation(t *testing.T) {
	type testStruct struct {
		name        string
		jwtValidity time.Duration
		isExpectErr bool
	}

	tests := []testStruct{
		{
			"test jwt is valid, should pass",
			10 * time.Minute,
			false,
		},
		{
			"test jwt is valid, should fail",
			0 * time.Minute,
			true,
		},
	}
	createUserReqBody, responseJson, shouldReturn := createAndRetrieveUser(t, createUserPath)
	if shouldReturn {
		return
	}
	userIdUUID, err := uuid.Parse(responseJson.UserId)
	testhelpers.AssertNoError(t, err)
	if err != nil {
		return
	}
	defer TestConfig.Queries.DeleteUserById(context.Background(), userIdUUID)
	for _, test := range tests {
		TestConfig.SetJWTValidDuration(test.jwtValidity)

		//set the validity of the token according to the test before running the test
		t.Run(test.name, func(t *testing.T) {
			res, shouldReturn := runLoginAndGetResponse(t, createUserReqBody)
			if shouldReturn {
				return
			}
			var loginResponse usermapping.LoginResponse
			decoder := json.NewDecoder(res.Body)
			jsonDecodeErr := decoder.Decode(&loginResponse)
			testhelpers.AssertNoError(t, jsonDecodeErr)
			if jsonDecodeErr != nil {
				return
			}
			req, err := http.NewRequest(http.MethodGet, testJWtAccessPath, nil)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer: %s", loginResponse.AccessToken))
			testhelpers.AssertNoError(t, err)
			if err != nil {
				return
			}
			testJWTAccessRes, err := client.Do(req)
			testhelpers.AssertNoError(t, err)
			if err != nil {
				return
			}
			testhelpers.AssertNoError(t, err)
			switch test.isExpectErr {
			case false:
				testhelpers.AssertIntVals(t, testJWTAccessRes.StatusCode, 200)
			case true:
				testhelpers.AssertIntVals(t, testJWTAccessRes.StatusCode, 401)
			}
		})
	}
}

func TestRefreshTokenGeneration(t *testing.T) {

	createUserReqBody, responseJson, shouldReturn := createAndRetrieveUser(t, createUserPath)
	if shouldReturn {
		return
	}

	userIdUUID, err := uuid.Parse(responseJson.UserId)
	testhelpers.AssertNoError(t, err)
	if err != nil {
		return
	}
	defer TestConfig.Queries.DeleteUserById(testContext, userIdUUID)
	res, shouldReturn := runLoginAndGetResponse(t, createUserReqBody)
	if shouldReturn {
		return
	}
	var loginResponse usermapping.LoginResponse
	decoder := json.NewDecoder(res.Body)
	jsonDecodeErr := decoder.Decode(&loginResponse)
	testhelpers.AssertNoError(t, jsonDecodeErr)
	if jsonDecodeErr != nil {
		fmt.Println("error occured at line 216")
		return
	}

	refreshTokenCreatedAtLogin, err := TestConfig.Queries.GetRefreshTokenInfoByToken(testContext, loginResponse.RefreshToken)
	defer TestConfig.Queries.DeleteRefreshTokenById(testContext, refreshTokenCreatedAtLogin.ID)
	testhelpers.AssertNoError(t, err)
	if err != nil {
		return
	}
	refreshRequest := tokenmapping.RefreshTokenJSON{
		RefreshToken: loginResponse.RefreshToken,
	}
	reqBody, err := json.Marshal(refreshRequest)
	testhelpers.AssertNoError(t, err)
	if err != nil {
		return
	}
	req, err := http.NewRequest(http.MethodPost, testRefreshTokenPath, bytes.NewReader(reqBody))
	testhelpers.AssertNoError(t, err)
	if err != nil {
		return
	}
	refreshResponse, refreshErr := client.Do(req)
	testhelpers.AssertNoError(t, refreshErr)
	if refreshErr != nil {
		return
	}
	// bodyBytes, err := io.ReadAll(bufio.NewReader(refreshResponse.Body))
	// if err != nil {
	// 	fmt.Println("error reading")
	// }
	// fmt.Printf("bodybytes:  %s", string(bodyBytes))

	var refreshResponseJSON usermapping.LoginResponse
	refreshDecoder := json.NewDecoder(refreshResponse.Body)
	jsonDecodeErr = refreshDecoder.Decode(&refreshResponseJSON)
	testhelpers.AssertNoError(t, jsonDecodeErr)
	if jsonDecodeErr != nil {
		fmt.Println("error occured at line 250")
		return
	}
	retrievedRefreshToken, err := TestConfig.Queries.GetRefreshTokenInfoByToken(testContext, refreshResponseJSON.RefreshToken)
	if err != nil {
		return
	}
	defer TestConfig.Queries.DeleteRefreshTokenById(testContext, retrievedRefreshToken.ID)
	testhelpers.AssertStringVals(t, responseJson.UserId, retrievedRefreshToken.UserID.String())
	if retrievedRefreshToken.Token == loginResponse.RefreshToken {
		t.Errorf("token should be different from the one that was created at login")
	}

}

func runLoginTest(t *testing.T, createUserPath string, isExpectPass bool) {
	createUserReqBody, responseJson, shouldReturn := createAndRetrieveUser(t, createUserPath)
	if shouldReturn {
		return
	}

	//first create the user in the DB

	userIdUUID, err := uuid.Parse(responseJson.UserId)
	testhelpers.AssertNoError(t, err)
	if err != nil {
		return
	}
	defer TestConfig.Queries.DeleteUserById(context.Background(), userIdUUID)
	//at this point, the user is already created, so defer deletion so that it will occur even if the function hits exception or error and returns early

	if !isExpectPass {
		createUserReqBody.Password = createUserReqBody.Password + "add fail"
	}
	res, shouldReturn := runLoginAndGetResponse(t, createUserReqBody)
	if shouldReturn {
		return
	}

	switch isExpectPass {
	case true:
		testhelpers.AssertIntVals(t, res.StatusCode, 200)

		var loginResponseJson usermapping.LoginResponse
		decoder := json.NewDecoder(res.Body)
		decodeErr := decoder.Decode(&loginResponseJson)
		testhelpers.AssertNoError(t, decodeErr)
		if decodeErr != nil {
			return
		}

		testhelpers.AssertBool(t, loginResponseJson.IsSuccess, true)

		token, err := auth.ValidateAndReturnClaims(TestConfig.JwtSecret, loginResponseJson.AccessToken)
		testhelpers.AssertNoError(t, err)
		if err != nil {
			fmt.Println("test failed at line 181")
			return
		}
		testhelpers.AssertStringVals(t, responseJson.UserId, token.UserId)
		testhelpers.AssertBool(t, token.IsAdmin, false)
		if len(loginResponseJson.RefreshToken) == 0 {
			t.Errorf("refresh token returned was null\n")

		}
		return
	case false:
		testhelpers.AssertIntVals(t, res.StatusCode, 401)
	}
}

func runLoginAndGetResponse(t *testing.T, createUserReqBody usermapping.UserJSON) (*http.Response, bool) {
	bodyString, err := json.Marshal(createUserReqBody)
	testhelpers.AssertNoError(t, err)
	if err != nil {
		fmt.Printf("error occured at line 233, %s\n", err.Error())
		return nil, true
	}
	req, err := http.NewRequest(http.MethodPost, loginUserPath, bytes.NewReader(bodyString))
	testhelpers.AssertNoError(t, err)
	if err != nil {
		fmt.Printf("error occured at line 239, %s\n", err.Error())
		return nil, true
	}
	res, err := client.Do(req)
	testhelpers.AssertNoError(t, err)
	if err != nil {
		fmt.Printf("error occured at line 245, %s\n", err.Error())
		return nil, true
	}
	return res, false
}

func createAndRetrieveUser(t *testing.T, path string) (usermapping.UserJSON, usermapping.CreatedUserResponse, bool) {

	createUserReqBody := usermapping.UserJSON{
		Email:    "wenming.soh@gmail.com",
		Password: "password",
	}

	bodyString, err := json.Marshal(createUserReqBody)
	testhelpers.AssertNoError(t, err)
	if err != nil {
		fmt.Printf("error occured at line 260, %s\n", err.Error())
		return usermapping.UserJSON{}, usermapping.CreatedUserResponse{}, true
	}

	req, reqErr := http.NewRequest(http.MethodPost, path, bytes.NewReader(bodyString))
	testhelpers.AssertNoError(t, reqErr)
	if reqErr != nil {
		fmt.Printf("error occured at line 266, %s\n", reqErr.Error())
		return usermapping.UserJSON{}, usermapping.CreatedUserResponse{}, true
	}
	res, resErr := client.Do(req)
	testhelpers.AssertNoError(t, resErr)
	if resErr != nil {
		fmt.Printf("error occured at line 273, %s\n", resErr.Error())
		return usermapping.UserJSON{}, usermapping.CreatedUserResponse{}, true
	}
	var responseJson usermapping.CreatedUserResponse

	decoder := json.NewDecoder(res.Body)
	jsonDecodeErr := decoder.Decode(&responseJson)

	testhelpers.AssertNoError(t, jsonDecodeErr)
	if jsonDecodeErr != nil {
		fmt.Printf("error occured at line 283, %s\n", jsonDecodeErr.Error())
		return usermapping.UserJSON{}, usermapping.CreatedUserResponse{}, true
	}
	testhelpers.AssertBool(t, responseJson.IsSuccess, true)
	return createUserReqBody, responseJson, false
}
