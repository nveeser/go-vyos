package vyos

import (
	"context"
	"testing"

	diffcmp "github.com/google/go-cmp/cmp"
)

func TestShowConfig(t *testing.T) {
	h := func(path string, data map[string]interface{}) response {
		if path != "/retrieve" {
			t.Errorf("got %s want /retrieve", path)
			return response{Error: "error"}
		}

		want := map[string]interface{}{
			"op":   "showConfig",
			"path": []any{"firewall"},
		}
		if diff := diffcmp.Diff(want, data); diff != "" {
			t.Errorf("got unexpected diff (-want +got):\n%s", diff)
			return response{Error: "error"}
		}
		return response{Success: true, Data: map[string]any{"rules": []any{"10"}}}
	}

	url, done := NewPostServer("key", h)
	defer done()

	ctx := context.Background()
	c, err := NewClient(url, Token("key"))
	if err != nil {
		t.Fatal(err)
	}
	res, err := c.ConfigMode().Show(ctx, "firewall")
	if err != nil {
		t.Errorf("Show error: %s", err)
	}
	if res["rules"].([]any)[0] != "10" {
		t.Errorf("got unexpected data: %v", res)
	}
}

func TestShowValues(t *testing.T) {
	h := func(path string, data map[string]interface{}) response {
		want := map[string]interface{}{
			"op":   "returnValues",
			"path": []any{"interfaces", "ethernet", "eth1", "address"},
		}
		if diff := diffcmp.Diff(want, data); diff != "" {
			t.Errorf("got unexpected diff (-want +got):\n%s", diff)
			return response{Error: "error"}
		}
		return response{Success: true, Data: []any{"10.0.0.1/24"}}
	}

	url, done := NewPostServer("key", h)
	defer done()

	ctx := context.Background()
	c, err := NewClient(url, Token("key"))
	if err != nil {
		t.Fatal(err)
	}
	res, err := c.ConfigMode().ShowValues(ctx, "interfaces ethernet eth1 address")
	if err != nil {
		t.Errorf("ShowValues error: %s", err)
	}
	if res[0] != "10.0.0.1/24" {
		t.Errorf("got unexpected data: %v", res)
	}
}

func TestExists(t *testing.T) {
	h := func(path string, data map[string]interface{}) response {
		want := map[string]interface{}{
			"op":   "exists",
			"path": []any{"interfaces", "ethernet", "eth1"},
		}
		if diff := diffcmp.Diff(want, data); diff != "" {
			t.Errorf("got unexpected diff (-want +got):\n%s", diff)
			return response{Error: "error"}
		}
		return response{Success: true, Data: true}
	}

	url, done := NewPostServer("key", h)
	defer done()

	ctx := context.Background()
	c, err := NewClient(url, Token("key"))
	if err != nil {
		t.Fatal(err)
	}
	res, err := c.ConfigMode().Exists(ctx, "interfaces ethernet eth1")
	if err != nil {
		t.Errorf("Exists error: %s", err)
	}
	if !res {
		t.Error("got false want true")
	}
}

func TestSet(t *testing.T) {
	h := func(path string, data map[string]interface{}) response {
		if path != "/configure" {
			t.Errorf("got %s want /configure", path)
			return response{Error: "error"}
		}

		want := map[string]interface{}{
			"op":   "set",
			"path": []any{"system", "host-name", "vyos"},
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
	err = c.ConfigMode().Set(ctx, "system host-name vyos")
	if err != nil {
		t.Errorf("Set error: %s", err)
	}
}

func TestDelete(t *testing.T) {
	h := func(path string, data map[string]interface{}) response {
		if path != "/configure" {
			t.Errorf("got %s want /configure", path)
			return response{Error: "error"}
		}

		want := map[string]interface{}{
			"op":   "delete",
			"path": []any{"system", "host-name"},
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
	err = c.ConfigMode().Delete(ctx, "system host-name")
	if err != nil {
		t.Errorf("Delete error: %s", err)
	}
}

func TestComment(t *testing.T) {
	h := func(path string, data map[string]interface{}) response {
		if path != "/configure" {
			t.Errorf("got %s want /configure", path)
			return response{Error: "error"}
		}

		want := map[string]interface{}{
			"op":      "comment",
			"path":    []any{"system", "host-name"},
			"comment": "Main Router",
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
	err = c.ConfigMode().Comment(ctx, "Main Router", "system host-name")
	if err != nil {
		t.Errorf("Comment error: %s", err)
	}
}

func TestSave(t *testing.T) {
	h := func(path string, data map[string]interface{}) response {
		if path != "/config-file" {
			t.Errorf("got %s want /config-file", path)
			return response{Error: "error"}
		}

		want := map[string]interface{}{
			"op":   "save",
			"fila": "/config/config.boot",
		}
		if diff := diffcmp.Diff(want, data); diff != "" {
			t.Errorf("got unexpected diff (-want +got):\n%s", diff)
			return response{Error: "error"}
		}
		return response{Success: true, Data: "saved"}
	}

	url, done := NewPostServer("key", h)
	defer done()

	ctx := context.Background()
	c, err := NewClient(url, Token("key"))
	if err != nil {
		t.Fatal(err)
	}
	msg, err := c.ConfigMode().Save(ctx, "/config/config.boot")
	if err != nil {
		t.Errorf("Save error: %s", err)
	}
	if msg != "saved" {
		t.Errorf("got %s want 'saved'", msg)
	}
}

func TestLoad(t *testing.T) {
	h := func(path string, data map[string]interface{}) response {
		if path != "/config-file" {
			t.Errorf("got %s want /config-file", path)
			return response{Error: "error"}
		}

		want := map[string]interface{}{
			"op":   "load",
			"fila": "/config/config.boot",
		}
		if diff := diffcmp.Diff(want, data); diff != "" {
			t.Errorf("got unexpected diff (-want +got):\n%s", diff)
			return response{Error: "error"}
		}
		return response{Success: true, Data: "loaded"}
	}

	url, done := NewPostServer("key", h)
	defer done()

	ctx := context.Background()
	c, err := NewClient(url, Token("key"))
	if err != nil {
		t.Fatal(err)
	}
	msg, err := c.ConfigMode().Load(ctx, "/config/config.boot")
	if err != nil {
		t.Errorf("Load error: %s", err)
	}
	if msg != "loaded" {
		t.Errorf("got %s want 'loaded'", msg)
	}
}

func TestCommitConfirm(t *testing.T) {
	h := func(path string, data map[string]interface{}) response {
		if path != "/config-file" {
			t.Errorf("got %s want /config-file", path)
			return response{Error: "error"}
		}

		want := map[string]interface{}{
			"op": "confirm",
		}
		if diff := diffcmp.Diff(want, data); diff != "" {
			t.Errorf("got unexpected diff (-want +got):\n%s", diff)
			return response{Error: "error"}
		}
		return response{Success: true, Data: "confirmed"}
	}

	url, done := NewPostServer("key", h)
	defer done()

	ctx := context.Background()
	c, err := NewClient(url, Token("key"))
	if err != nil {
		t.Fatal(err)
	}
	msg, err := c.ConfigMode().CommitConfirm(ctx)
	if err != nil {
		t.Errorf("CommitConfirm error: %s", err)
	}
	if msg != "confirmed" {
		t.Errorf("got %s want 'confirmed'", msg)
	}
}
