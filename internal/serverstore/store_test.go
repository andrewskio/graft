package serverstore

import (
	"os"
	"path/filepath"
	"testing"

	"graft/internal/composefile"
	"graft/internal/server"
)

func TestSaveLoadAndList(t *testing.T) {
	store := New(t.TempDir())
	want := server.Server{
		ServerName:       "survival_1",
		ServerType:       "paper",
		MinecraftVersion: "1.21.1",
		JavaVersion:      "21",
		Memory:           "2G",
	}

	composeData, err := composefile.FromServer(want).Marshal()
	if err != nil {
		t.Fatalf("marshal server Compose: %v", err)
	}
	path, err := store.SaveCompose(want.ServerName, composeData, false)
	if err != nil {
		t.Fatalf("save server: %v", err)
	}
	if path != filepath.Join(store.Root(), "servers", want.ServerName, "compose.yml") {
		t.Fatalf("unexpected compose path: %s", path)
	}

	loadedData, err := store.LoadCompose(want.ServerName)
	if err != nil {
		t.Fatalf("load Compose: %v", err)
	}
	parsed, err := composefile.Parse(loadedData)
	if err != nil {
		t.Fatalf("parse loaded Compose: %v", err)
	}
	got, err := parsed.ToServer()
	if err != nil {
		t.Fatalf("convert loaded Compose: %v", err)
	}
	if got != want {
		t.Fatalf("loaded server mismatch: got %#v, want %#v", got, want)
	}

	names, err := store.ListNames()
	if err != nil {
		t.Fatalf("list servers: %v", err)
	}
	if len(names) != 1 || names[0] != want.ServerName {
		t.Fatalf("unexpected server names: %#v", names)
	}
}

func TestSaveDoesNotOverwriteByDefault(t *testing.T) {
	store := New(t.TempDir())
	configuration := server.Server{ServerName: "survival"}
	composeData, err := composefile.FromServer(configuration).Marshal()
	if err != nil {
		t.Fatalf("marshal server Compose: %v", err)
	}
	if _, err := store.SaveCompose(configuration.ServerName, composeData, false); err != nil {
		t.Fatalf("initial save: %v", err)
	}
	if _, err := store.SaveCompose(configuration.ServerName, composeData, false); err == nil {
		t.Fatal("expected second save to reject overwrite")
	}
	if _, err := store.SaveCompose(configuration.ServerName, composeData, true); err != nil {
		t.Fatalf("save with overwrite enabled: %v", err)
	}
}

func TestStoreRejectsUnsafeNames(t *testing.T) {
	store := New(t.TempDir())
	for _, name := range []string{"../outside", "a/b", "-invalid", "x"} {
		t.Run(name, func(t *testing.T) {
			if _, err := store.ComposePath(name); err == nil {
				t.Fatalf("expected unsafe name %q to be rejected", name)
			}
		})
	}
	if _, err := os.Stat(filepath.Join(store.Root(), "outside")); !os.IsNotExist(err) {
		t.Fatalf("unsafe path unexpectedly exists, stat error: %v", err)
	}
}

func TestListNamesIgnoresInvalidDirectoryNames(t *testing.T) {
	store := New(t.TempDir())
	invalidPath := filepath.Join(store.Root(), "servers", "..escape")
	if err := os.MkdirAll(invalidPath, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(invalidPath, "compose.yml"), []byte("invalid"), 0644); err != nil {
		t.Fatal(err)
	}
	names, err := store.ListNames()
	if err != nil {
		t.Fatalf("list names: %v", err)
	}
	if len(names) != 0 {
		t.Fatalf("expected invalid directory to be ignored, got %#v", names)
	}
}

func TestListNamesForMissingStoreIsEmpty(t *testing.T) {
	store := New(filepath.Join(t.TempDir(), "missing"))
	names, err := store.ListNames()
	if err != nil {
		t.Fatalf("list missing store: %v", err)
	}
	if len(names) != 0 {
		t.Fatalf("expected no names, got %#v", names)
	}
}

func TestValidateNewNameDetectsExistingDirectory(t *testing.T) {
	store := New(t.TempDir())
	serverPath := filepath.Join(store.Root(), "servers", "already-there")
	if err := os.MkdirAll(serverPath, 0755); err != nil {
		t.Fatal(err)
	}
	if err := store.ValidateNewName("already-there"); err == nil {
		t.Fatalf("expected name conflict, got %v", err)
	}
}
