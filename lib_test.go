package ginko

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegisterAssets(t *testing.T) {
	mux := http.NewServeMux()
	registerAssets(mux)

	cases := []struct {
		path        string
		contentType string
		want        []byte
	}{
		{w3CSSPath, "text/css; charset=utf-8", w3CSS},
		{htmxJSPath, "application/javascript; charset=utf-8", htmxJS},
	}

	for _, c := range cases {
		req := httptest.NewRequest(http.MethodGet, c.path, nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("%s: got status %d, want %d", c.path, rec.Code, http.StatusOK)
		}
		if got := rec.Header().Get("Content-Type"); got != c.contentType {
			t.Errorf("%s: got Content-Type %q, want %q", c.path, got, c.contentType)
		}
		if rec.Body.String() != string(c.want) {
			t.Errorf("%s: served body does not match embedded asset", c.path)
		}
	}
}

func TestCsrfGuard(t *testing.T) {
	const port = "8090"
	selfOrigin := "http://127.0.0.1:" + port

	called := false
	action := ActionMap{
		"/action": func(w Response, req Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		},
	}

	cases := []struct {
		name       string
		method     string
		origin     string
		referer    string
		wantStatus int
		wantCalled bool
	}{
		{
			name:       "GET is always allowed regardless of headers",
			method:     http.MethodGet,
			origin:     "http://evil.example",
			wantStatus: http.StatusOK,
			wantCalled: true,
		},
		{
			name:       "PUT with matching Origin is allowed",
			method:     http.MethodPut,
			origin:     selfOrigin,
			wantStatus: http.StatusOK,
			wantCalled: true,
		},
		{
			name:       "PUT with mismatched Origin is forbidden",
			method:     http.MethodPut,
			origin:     "http://evil.example",
			wantStatus: http.StatusForbidden,
			wantCalled: false,
		},
		{
			name:       "PUT with matching Referer and no Origin is allowed",
			method:     http.MethodPut,
			referer:    selfOrigin + "/",
			wantStatus: http.StatusOK,
			wantCalled: true,
		},
		{
			name:       "PUT with mismatched Referer is forbidden",
			method:     http.MethodPut,
			referer:    "http://evil.example/",
			wantStatus: http.StatusForbidden,
			wantCalled: false,
		},
		{
			name:       "PUT with neither Origin nor Referer is allowed",
			method:     http.MethodPut,
			wantStatus: http.StatusOK,
			wantCalled: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			called = false
			mux := csrfGuard(port, action)

			req := httptest.NewRequest(c.method, "/action", nil)
			if c.origin != "" {
				req.Header.Set("Origin", c.origin)
			}
			if c.referer != "" {
				req.Header.Set("Referer", c.referer)
			}
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, req)

			if rec.Code != c.wantStatus {
				t.Errorf("got status %d, want %d", rec.Code, c.wantStatus)
			}
			if called != c.wantCalled {
				t.Errorf("handler called = %v, want %v", called, c.wantCalled)
			}
		})
	}
}
