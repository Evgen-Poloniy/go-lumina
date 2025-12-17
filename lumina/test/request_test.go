package lumina_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/nexula-rg/go-lumina/lumina"
)

func TestPing(t *testing.T) {
	client, _ := lumina.NewClient()
	err := client.Ping()

	if err != nil {
		t.Fatalf("error: %v", err)
	}
}

func TestCheckHealth(t *testing.T) {
	expected := 200

	client, _ := lumina.NewClient()
	resp, err := client.CheckHealth()

	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if resp.StatusCode != expected {
		t.Fatalf("expected %d, got %d", expected, resp.StatusCode)
	}
}

func TestMakeRequest(t *testing.T) {
	question := "What is that?"
	expected := question

	client, _ := lumina.NewClient()
	resp, err := client.MakeRequest(question)

	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if resp.Answer != expected {
		t.Fatalf("expected %s, got %s", expected, resp.Answer)
	}
}

func TestMakeRequestCtx(t *testing.T) {
	question := "What is that?"
	expected := question
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*250)
	defer cancel()

	client, _ := lumina.NewClient()
	resp, err := client.MakeRequestCtx(ctx, question)

	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if resp.Answer != expected {
		t.Fatalf("expected %s, got %s", expected, resp.Answer)
	}
}

func TestMakeAsyncRequest(t *testing.T) {
	answer_uuid := "uuid"
	question := "What is that?"
	expected := http.StatusServiceUnavailable

	client, _ := lumina.NewClient()
	_, err := client.MakeAsyncRequest(answer_uuid, question)

	if err != nil {
		if errResp, ok := err.(*lumina.ErrorResponse); ok {
			if errResp.StatusCode != expected {
				t.Fatalf("expected %d, got %d", expected, errResp.StatusCode)
			}
		} else {
			t.Fatalf("unexpected error")
		}
	} else {
		t.Fatalf("error: expected error, got no error")
	}
}
