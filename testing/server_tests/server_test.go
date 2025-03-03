package servertests

import (
	"net/http"
	"os"
	"os/exec"
	"testing"
)

func TestMain(m *testing.M) {
	startServerCmd := cmdBuilder("../..", "make", "docker-server-start")
	if err := startServerCmd.Run(); err != nil {
		os.Stderr.WriteString("Error starting server: " + err.Error() + "\n")
		os.Exit(1)
	}

	code := m.Run()
	stopServerCmd := cmdBuilder("../..", "make", "docker-server-stop")
	stopServerCmd.Dir = "../.."
	if err := stopServerCmd.Run(); err != nil {
		os.Stderr.WriteString("Error stopping server: " + err.Error() + "\n")
		os.Exit(1)
	}

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

func cmdBuilder(dir string, command string, args ...string) *exec.Cmd {
	cmd := exec.Command(command, args...)
	if dir != "" {
		cmd.Dir = dir
	}
	return cmd
}
