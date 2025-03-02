package servertests

import (
	"net/http"
	"os"
	"testing"

	"github.com/sohWenMing/finance_project/internal/server"
)

func TestMain(m *testing.M) {

	server, _ := server.InitServer()
	defer server.Close()
	code := m.Run()
	os.Exit(code)
}

func TestInitServer(t *testing.T) {
	/*
		test should attempt to init the server , attempt to ping it,
		and then close the server
	*/
	client := http.DefaultClient

	req, err := http.NewRequest("GET",
		"http://localhost:8080/ping", nil)
	if err != nil {
		t.Errorf("error creating request in test\n")
	}
	res, err := client.Do(req)
	if err != nil {
		t.Error(err)
	}
	if res.StatusCode != http.StatusOK {
		t.Errorf("got %d, want %d", res.StatusCode, http.StatusOK)
	}
}
