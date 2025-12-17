package lumina_test

import (
	"testing"

	"github.com/nexula-rg/go-lumina/lumina"
)

func TestNewClient(t *testing.T) {
	_, err := lumina.NewClient()

	if err != nil {
		t.Fatalf("error: %v", err)
	}
}
