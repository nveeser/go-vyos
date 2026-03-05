package vyos

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
)

func NewPostServer(apiKey string, h vyosHandler) (string, func()) {
	s := httptest.NewServer(postHandler{
		apiKey:  apiKey,
		handler: h,
	})
	return s.URL, s.Close
}

func NewTestGetServer(apiKey string, h func(*http.Request) response) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := h(r)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
}

type vyosHandler func(path string, data map[string]interface{}) response

type postHandler struct {
	apiKey  string
	handler vyosHandler
}

func (h postHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(32 << 10); err != nil {
		err := r.ParseForm()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "error: %s", err.Error())
			return
		}
	}

	if r.FormValue("key") != h.apiKey {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response{Success: false, Error: "Invalid API Key"})
		return
	}

	data := map[string]any{}
	err := json.Unmarshal([]byte(r.FormValue("data")), &data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error: %s", err.Error())
	}
	resp := h.handler(r.URL.Path, data)
	w.Header().Set("Content-Type", "application/json")

	if err = json.NewEncoder(w).Encode(resp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error: %s", err.Error())
	}
}
