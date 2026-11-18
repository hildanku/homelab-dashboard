package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/hildanku/homelab-dashboard/internal/shared"
)

func newTestApp(t *testing.T) (*fiber.App, *Service) {
	t.Helper()

	svc, err := New("admin", "secret", "test-secret", time.Hour, false)
	if err != nil {
		t.Fatalf("new service: %v", err)
	}

	app := fiber.New()
	app.Get("/api/health", func(c *fiber.Ctx) error {
		return shared.AppResponse(c, fiber.StatusOK, "ok", nil)
	})
	app.Post("/api/login", svc.Login)
	app.Use("/api", svc.Middleware())
	app.Get("/api/auth/me", svc.Me)
	app.Post("/api/logout", svc.Logout)
	app.Get("/api/system", func(c *fiber.Ctx) error {
		return shared.AppResponse(c, fiber.StatusOK, "success", nil)
	})

	return app, svc
}

func doRequest(t *testing.T, app *fiber.App, method, path, body string, cookie *http.Cookie) *http.Response {
	t.Helper()

	req := httptest.NewRequest(method, path, strings.NewReader(body))
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	if cookie != nil {
		req.AddCookie(cookie)
	}

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("%s %s: %v", method, path, err)
	}
	return resp
}

func sessionCookie(resp *http.Response) *http.Cookie {
	for _, cookie := range resp.Cookies() {
		if cookie.Name == CookieName {
			return cookie
		}
	}
	return nil
}

func responseMessage(t *testing.T, resp *http.Response) string {
	t.Helper()

	var body struct {
		Message string `json:"message"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	return body.Message
}

func TestPublicHealth(t *testing.T) {
	app, _ := newTestApp(t)

	resp := doRequest(t, app, http.MethodGet, "/api/health", "", nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("health status = %d, want 200", resp.StatusCode)
	}
}

func TestLoginProtectsAPI(t *testing.T) {
	app, _ := newTestApp(t)

	resp := doRequest(t, app, http.MethodPost, "/api/login", `{"username":"admin","password":"wrong"}`, nil)
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("wrong password status = %d, want 401", resp.StatusCode)
	}
	if got := responseMessage(t, resp); got != "invalid username or password" {
		t.Fatalf("wrong password message = %q", got)
	}

	resp = doRequest(t, app, http.MethodPost, "/api/login", `{"username":"nobody","password":"invalid-password"}`, nil)
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("wrong username status = %d, want 401", resp.StatusCode)
	}

	resp = doRequest(t, app, http.MethodPost, "/api/login", `{"username":"admin","password":"secret"}`, nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("login status = %d, want 200", resp.StatusCode)
	}
	cookie := sessionCookie(resp)
	if cookie == nil || cookie.Value == "" {
		t.Fatal("login did not set session cookie")
	}
	if !cookie.HttpOnly || cookie.SameSite != http.SameSiteLaxMode {
		t.Fatalf("cookie flags = HttpOnly:%v SameSite:%v", cookie.HttpOnly, cookie.SameSite)
	}

	resp = doRequest(t, app, http.MethodGet, "/api/system", "", nil)
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("unauthenticated API status = %d, want 401", resp.StatusCode)
	}

	resp = doRequest(t, app, http.MethodGet, "/api/system", "", cookie)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("authenticated API status = %d, want 200", resp.StatusCode)
	}

	resp = doRequest(t, app, http.MethodGet, "/api/auth/me", "", cookie)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("me status = %d, want 200", resp.StatusCode)
	}

	tampered := *cookie
	last := len(tampered.Value) - 1
	if tampered.Value[last] == 'A' {
		tampered.Value = tampered.Value[:last] + "B"
	} else {
		tampered.Value = tampered.Value[:last] + "A"
	}
	resp = doRequest(t, app, http.MethodGet, "/api/system", "", &tampered)
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("tampered cookie status = %d, want 401", resp.StatusCode)
	}

	resp = doRequest(t, app, http.MethodPost, "/api/logout", "", cookie)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("logout status = %d, want 200", resp.StatusCode)
	}
	cleared := sessionCookie(resp)
	if cleared == nil || cleared.Value != "" || !cleared.Expires.Before(time.Now()) {
		t.Fatalf("logout did not clear cookie: %+v", cleared)
	}
}

func TestLoadFromEnvRequiresCredentials(t *testing.T) {
	t.Setenv("ADMIN_USERNAME", "")
	t.Setenv("ADMIN_PASSWORD", "")

	if _, err := LoadFromEnv(); err == nil {
		t.Fatal("LoadFromEnv should fail without credentials")
	}
}
