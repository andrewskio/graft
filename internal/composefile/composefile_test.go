package composefile

import (
	"testing"

	"graft/internal/server"
)

func TestServerComposeRoundTrip(t *testing.T) {
	want := server.Server{
		ServerName:       "test-server",
		ServerType:       "paper",
		MinecraftVersion: "1.21.1",
		JavaVersion:      "21",
		Memory:           "2G",
	}

	data, err := FromServer(want).Marshal()
	if err != nil {
		t.Fatalf("marshal Compose: %v", err)
	}
	parsed, err := Parse(data)
	if err != nil {
		t.Fatalf("parse Compose: %v", err)
	}
	got, err := parsed.ToServer()
	if err != nil {
		t.Fatalf("convert Compose to server: %v", err)
	}
	if got != want {
		t.Fatalf("round trip mismatch: got %#v, want %#v", got, want)
	}
}

func TestToServerRejectsMultipleServices(t *testing.T) {
	compose := Compose{Services: map[string]Service{
		"one": {},
		"two": {},
	}}
	if _, err := compose.ToServer(); err == nil {
		t.Fatal("expected multiple services to be rejected")
	}
}
