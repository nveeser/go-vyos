package vyos

import (
	"context"
	"net/http"
	"testing"

	diffcmp "github.com/google/go-cmp/cmp"
)

func TestInfo(t *testing.T) {
	s := NewTestGetServer("key", func(r *http.Request) response {
		if r.URL.Path != "/info" {
			return response{Error: "wrong path"}
		}
		if r.URL.Query().Get("version") != "1" {
			return response{Error: "missing version query"}
		}
		return response{
			Success: true,
			Data: &InfoResponse{
				Version: "1.3.0",
			},
		}
	})
	defer s.Close()

	ctx := context.Background()
	c, err := NewClient(s.URL, Token("key"))
	if err != nil {
		t.Fatal(err)
	}
	resp, err := c.OpMode().Info(ctx, InfoRequest{Version: true})
	if err != nil {
		t.Errorf("Info error: %s", err)
	}
	if resp.Version != "1.3.0" {
		t.Errorf("got %s want 1.3.0", resp.Version)
	}
}

func TestShow(t *testing.T) {
	h := func(path string, data map[string]interface{}) response {
		if path != "/show" {
			t.Errorf("got %s want /show", path)
			return response{Error: "error"}
		}

		want := map[string]interface{}{
			"op":   "show",
			"path": []any{"system", "image"},
		}
		if diff := diffcmp.Diff(want, data); diff != "" {
			t.Errorf("got unexpected diff (-want +got):\n%s", diff)
			return response{Error: "error"}
		}
		return response{Success: true, Data: "show output"}
	}

	url, done := NewPostServer("key", h)
	defer done()

	ctx := context.Background()
	c, err := NewClient(url, Token("key"))
	if err != nil {
		t.Fatal(err)
	}
	msg, err := c.OpMode().Show(ctx, "system image")
	if err != nil {
		t.Errorf("Show error: %s", err)
	}
	if msg != "show output" {
		t.Errorf("got %s want 'show output'", msg)
	}
}

func TestGenerate(t *testing.T) {
	h := func(path string, data map[string]interface{}) response {
		if path != "/generate" {
			t.Errorf("got %s want /generate", path)
			return response{Error: "error"}
		}

		want := map[string]interface{}{
			"op":   "generate",
			"path": []any{"pki", "wireguard", "key-pair"},
		}
		if diff := diffcmp.Diff(want, data); diff != "" {
			t.Errorf("got unexpected diff (-want +got):\n%s", diff)
			return response{Error: "error"}
		}
		return response{Success: true, Data: "generated"}
	}

	url, done := NewPostServer("key", h)
	defer done()

	ctx := context.Background()
	c, err := NewClient(url, Token("key"))
	if err != nil {
		t.Fatal(err)
	}
	msg, err := c.OpMode().Generate(ctx, "pki wireguard key-pair")
	if err != nil {
		t.Errorf("Generate error: %s", err)
	}
	if msg != "generated" {
		t.Errorf("got %s want 'generated'", msg)
	}
}

func TestAddImage(t *testing.T) {
	h := func(path string, data map[string]interface{}) response {
		if path != "/image" {
			t.Errorf("got %s want /image", path)
			return response{Error: "error"}
		}

		want := map[string]interface{}{
			"op":  "add",
			"url": "http://example.com/vyos.iso",
		}
		if diff := diffcmp.Diff(want, data); diff != "" {
			t.Errorf("got unexpected diff (-want +got):\n%s", diff)
			return response{Error: "error"}
		}
		return response{Success: true, Data: "image added"}
	}

	url, done := NewPostServer("key", h)
	defer done()

	ctx := context.Background()
	c, err := NewClient(url, Token("key"))
	if err != nil {
		t.Fatal(err)
	}
	msg, err := c.OpMode().AddImage(ctx, "http://example.com/vyos.iso")
	if err != nil {
		t.Errorf("AddImage error: %s", err)
	}
	if msg != "image added" {
		t.Errorf("got %s want 'image added'", msg)
	}
}

func TestDeleteImage(t *testing.T) {
	h := func(path string, data map[string]interface{}) response {
		if path != "/image" {
			t.Errorf("got %s want /image", path)
			return response{Error: "error"}
		}

		want := map[string]interface{}{
			"op":   "delete",
			"name": "vyos-1.3.0",
		}
		if diff := diffcmp.Diff(want, data); diff != "" {
			t.Errorf("got unexpected diff (-want +got):\n%s", diff)
			return response{Error: "error"}
		}
		return response{Success: true, Data: "image deleted"}
	}

	url, done := NewPostServer("key", h)
	defer done()

	ctx := context.Background()
	c, err := NewClient(url, Token("key"))
	if err != nil {
		t.Fatal(err)
	}
	msg, err := c.OpMode().DeleteImage(ctx, "vyos-1.3.0")
	if err != nil {
		t.Errorf("DeleteImage error: %s", err)
	}
	if msg != "image deleted" {
		t.Errorf("got %s want 'image deleted'", msg)
	}
}

func TestReset(t *testing.T) {
	h := func(path string, data map[string]interface{}) response {
		if path != "/reset" {
			t.Errorf("got %s want /reset", path)
			return response{Error: "error"}
		}

		want := map[string]interface{}{
			"op":   "reset",
			"path": []any{"vpn", "ipsec"},
		}
		if diff := diffcmp.Diff(want, data); diff != "" {
			t.Errorf("got unexpected diff (-want +got):\n%s", diff)
			return response{Error: "error"}
		}
		return response{Success: true}
	}

	url, done := NewPostServer("key", h)
	defer done()

	ctx := context.Background()
	c, err := NewClient(url, Token("key"))
	if err != nil {
		t.Fatal(err)
	}
	err = c.OpMode().Reset(ctx, "vpn ipsec")
	if err != nil {
		t.Errorf("Reset error: %s", err)
	}
}

func TestPowerOff(t *testing.T) {
	h := func(path string, data map[string]interface{}) response {
		if path != "/poweroff" {
			t.Errorf("got %s want /poweroff", path)
			return response{Error: "error"}
		}

		want := map[string]interface{}{
			"op":   "poweroff",
			"path": []any{"now"},
		}
		if diff := diffcmp.Diff(want, data); diff != "" {
			t.Errorf("got unexpected diff (-want +got):\n%s", diff)
			return response{Error: "error"}
		}
		return response{Success: true}
	}

	url, done := NewPostServer("key", h)
	defer done()

	ctx := context.Background()
	c, err := NewClient(url, Token("key"))
	if err != nil {
		t.Fatal(err)
	}
	err = c.OpMode().PowerOff(ctx, "now")
	if err != nil {
		t.Errorf("PowerOff error: %s", err)
	}
}

func TestReboot(t *testing.T) {
	h := func(path string, data map[string]interface{}) response {
		if path != "/reboot" {
			t.Errorf("got %s want /reboot", path)
			return response{Error: "error"}
		}

		want := map[string]interface{}{
			"op":   "reboot",
			"path": []any{"now"},
		}
		if diff := diffcmp.Diff(want, data); diff != "" {
			t.Errorf("got unexpected diff (-want +got):\n%s", diff)
			return response{Error: "error"}
		}
		return response{Success: true}
	}

	url, done := NewPostServer("key", h)
	defer done()

	ctx := context.Background()
	c, err := NewClient(url, Token("key"))
	if err != nil {
		t.Fatal(err)
	}
	err = c.OpMode().Reboot(ctx, "now")
	if err != nil {
		t.Errorf("Reboot error: %s", err)
	}
}
