package client_test

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/anexia-it/go-anxcloud/pkg/client"
)

func TestTokenClient(t *testing.T) {
	dummyToken := "ie7dois8Ooquoo1ieB9kae8Od9ooshee3nejuach4inae3gai0Re0Shaipeihail" //nolint:gosec // Not a real token.
	c := client.NewTokenClient(dummyToken, &http.Client{})
	expectedAuthorizationHeader := fmt.Sprintf("Token %s", dummyToken)
	echoHandler := echoTestHandler(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != expectedAuthorizationHeader {
			t.Fatalf("token client did not add expected header Authorization. Was %s", r.Header.Get("Authorization"))
		}
		echoHandler.ServeHTTP(w, r)
	})

	cw, server := client.NewTestClient(c, handler)
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), client.DefaultRequestTimeout)
	defer cancel()
	if err := client.Echo(ctx, cw); err != nil {
		t.Fatalf("echo test failed: %v", err)
	}
}

func TestTokenClientIntegration(t *testing.T) {
	var set bool
	if _, set = os.LookupEnv(client.IntegrationTestEnvName); !set {
		t.Skip("integration tests disabled")
	}
	c, err := client.NewAnyClientFromEnvs(false, nil)
	if err != nil {
		t.Fatalf("could not create client: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), client.DefaultRequestTimeout)
	defer cancel()
	if err := client.Echo(ctx, c); err != nil {
		t.Fatalf("[%s] echo test failed: %v", time.Now(), err)
	}
}
