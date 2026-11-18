package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/hildanku/homelab-dashboard/internal/shared"
	"golang.org/x/crypto/bcrypt"
)

const CookieName = "dashboard_session"

type Service struct {
	username     string
	passwordHash []byte
	dummyHash    []byte
	secret       []byte
	ttl          time.Duration
	secure       bool
}

type credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func New(username, password, secret string, ttl time.Duration, secure bool) (*Service, error) {
	return newService(username, password, secret, bcrypt.DefaultCost, ttl, secure)
}

func newService(username, password, secret string, cost int, ttl time.Duration, secure bool) (*Service, error) {
	if username == "" || password == "" {
		return nil, errors.New("ADMIN_USERNAME and ADMIN_PASSWORD must be set")
	}
	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		return nil, fmt.Errorf("bcrypt cost must be between %d and %d", bcrypt.MinCost, bcrypt.MaxCost)
	}
	if ttl <= 0 {
		ttl = 12 * time.Hour
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return nil, fmt.Errorf("hash admin password: %w", err)
	}
	dummyHash, err := bcrypt.GenerateFromPassword([]byte("invalid-password"), cost)
	if err != nil {
		return nil, fmt.Errorf("hash dummy password: %w", err)
	}

	secretBytes := []byte(secret)
	if secret == "" {
		sum := sha256.Sum256([]byte(username + "\x00" + password))
		secretBytes = sum[:]
	}

	return &Service{
		username:     username,
		passwordHash: passwordHash,
		dummyHash:    dummyHash,
		secret:       secretBytes,
		ttl:          ttl,
		secure:       secure,
	}, nil
}

func LoadFromEnv() (*Service, error) {
	username := os.Getenv("ADMIN_USERNAME")
	password := os.Getenv("ADMIN_PASSWORD")
	if username == "" || password == "" {
		return nil, errors.New("ADMIN_USERNAME and ADMIN_PASSWORD must be set")
	}

	cost := bcrypt.DefaultCost
	if raw := os.Getenv("BCRYPT_COST"); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil || parsed < bcrypt.MinCost || parsed > bcrypt.MaxCost {
			return nil, fmt.Errorf("invalid BCRYPT_COST: %q", raw)
		}
		cost = parsed
	}

	ttl := 12 * time.Hour
	if raw := os.Getenv("SESSION_TTL"); raw != "" {
		parsed, err := time.ParseDuration(raw)
		if err != nil || parsed <= 0 {
			return nil, fmt.Errorf("invalid SESSION_TTL: %q", raw)
		}
		ttl = parsed
	}

	secure := false
	if raw := os.Getenv("COOKIE_SECURE"); raw != "" {
		parsed, err := strconv.ParseBool(raw)
		if err != nil {
			return nil, fmt.Errorf("invalid COOKIE_SECURE: %q", raw)
		}
		secure = parsed
	}

	return newService(username, password, os.Getenv("SESSION_SECRET"), cost, ttl, secure)
}

func (s *Service) Login(c *fiber.Ctx) error {
	var creds credentials
	if err := c.BodyParser(&creds); err != nil {
		return shared.AppResponse(c, fiber.StatusBadRequest, "invalid request", nil)
	}

	usernameOK := subtle.ConstantTimeCompare([]byte(creds.Username), []byte(s.username)) == 1
	hash := s.passwordHash
	if !usernameOK {
		hash = s.dummyHash
	}
	if !usernameOK || bcrypt.CompareHashAndPassword(hash, []byte(creds.Password)) != nil {
		return shared.AppResponse(c, fiber.StatusUnauthorized, "invalid username or password", nil)
	}

	token, err := s.issueToken()
	if err != nil {
		return shared.AppResponse(c, fiber.StatusInternalServerError, "failed to create session", nil)
	}
	s.setCookie(c, token, s.ttl)

	return shared.AppResponse(c, fiber.StatusOK, "login success", fiber.Map{"username": s.username})
}

func (s *Service) Logout(c *fiber.Ctx) error {
	s.clearCookie(c)
	return shared.AppResponse(c, fiber.StatusOK, "logout success", nil)
}

func (s *Service) Me(c *fiber.Ctx) error {
	username, _ := c.Locals("username").(string)
	return shared.AppResponse(c, fiber.StatusOK, "success", fiber.Map{
		"authenticated": true,
		"username":      username,
	})
}

func (s *Service) Middleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		path := c.Path()
		if path == "/api/health" || path == "/api/login" {
			return c.Next()
		}

		token := c.Cookies(CookieName)
		if token == "" || !s.verifyToken(token) {
			return shared.AppResponse(c, fiber.StatusUnauthorized, "unauthorized", nil)
		}

		c.Locals("username", s.username)
		return c.Next()
	}
}

func (s *Service) issueToken() (string, error) {
	nonce := make([]byte, 16)
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}

	message := fmt.Sprintf("%d.%s", time.Now().Add(s.ttl).Unix(), hex.EncodeToString(nonce))
	mac := hmac.New(sha256.New, s.secret)
	mac.Write([]byte(message))
	return message + "." + base64.RawURLEncoding.EncodeToString(mac.Sum(nil)), nil
}

func (s *Service) verifyToken(token string) bool {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return false
	}

	exp, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil || exp <= time.Now().Unix() {
		return false
	}
	if len(parts[1]) != 32 {
		return false
	}
	if _, err := hex.DecodeString(parts[1]); err != nil {
		return false
	}

	signature, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return false
	}

	mac := hmac.New(sha256.New, s.secret)
	mac.Write([]byte(parts[0] + "." + parts[1]))
	return hmac.Equal(signature, mac.Sum(nil))
}

func (s *Service) setCookie(c *fiber.Ctx, token string, ttl time.Duration) {
	c.Cookie(&fiber.Cookie{
		Name:     CookieName,
		Value:    token,
		Path:     "/",
		HTTPOnly: true,
		Secure:   s.secure,
		SameSite: "Lax",
		Expires:  time.Now().Add(ttl),
		MaxAge:   int(ttl.Seconds()),
	})
}

func (s *Service) clearCookie(c *fiber.Ctx) {
	c.Cookie(&fiber.Cookie{
		Name:     CookieName,
		Value:    "",
		Path:     "/",
		HTTPOnly: true,
		Secure:   s.secure,
		SameSite: "Lax",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
	})
}
