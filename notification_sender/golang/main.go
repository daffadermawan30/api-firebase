// Web Pengirim Notifikasi FCM (HTTP v1) - versi Golang
// Hanya memakai standard library (tanpa package tambahan).
//
// Syarat: file service-account.json (dari Firebase Console) ada di folder ini.
package main

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	serviceAccountFile = "service-account.json"
	fcmScope           = "https://www.googleapis.com/auth/firebase.messaging"
)

// ServiceAccount berisi field penting dari service-account.json
type ServiceAccount struct {
	ProjectID   string `json:"project_id"`
	ClientEmail string `json:"client_email"`
	PrivateKey  string `json:"private_key"`
	TokenURI    string `json:"token_uri"`
}

// Cache access token agar tidak meminta token baru setiap kirim.
var (
	tokenMu     sync.Mutex
	cachedToken string
	tokenExpiry time.Time
)

func loadServiceAccount() (*ServiceAccount, error) {
	// Urutan pencarian: folder kerja, folder induk, lalu folder tempat file .exe berada
	candidates := []string{serviceAccountFile, filepath.Join("..", serviceAccountFile)}
	if exe, e := os.Executable(); e == nil {
		candidates = append(candidates, filepath.Join(filepath.Dir(exe), serviceAccountFile))
	}
	var raw []byte
	var err error
	for _, path := range candidates {
		if raw, err = os.ReadFile(path); err == nil {
			break
		}
	}
	if err != nil {
		return nil, errors.New("file service-account.json tidak ditemukan di folder golang/ maupun notification_sender/")
	}
	var sa ServiceAccount
	if err := json.Unmarshal(raw, &sa); err != nil {
		return nil, fmt.Errorf("service-account.json tidak valid: %w", err)
	}
	return &sa, nil
}

func b64url(b []byte) string { return base64.RawURLEncoding.EncodeToString(b) }

// getAccessToken - LANGKAH 1: buat JWT, tanda tangani (RS256),
// lalu tukar ke Google untuk mendapatkan access token OAuth 2.0.
func getAccessToken(sa *ServiceAccount) (string, error) {
	tokenMu.Lock()
	defer tokenMu.Unlock()
	if cachedToken != "" && time.Now().Before(tokenExpiry) {
		return cachedToken, nil
	}

	block, _ := pem.Decode([]byte(sa.PrivateKey))
	if block == nil {
		return "", errors.New("private_key tidak bisa dibaca (format PEM)")
	}
	parsed, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("private_key tidak valid: %w", err)
	}
	key, ok := parsed.(*rsa.PrivateKey)
	if !ok {
		return "", errors.New("private_key bukan RSA")
	}

	now := time.Now().Unix()
	header, _ := json.Marshal(map[string]string{"alg": "RS256", "typ": "JWT"})
	claims, _ := json.Marshal(map[string]any{
		"iss":   sa.ClientEmail,
		"scope": fcmScope,
		"aud":   sa.TokenURI,
		"iat":   now,
		"exp":   now + 3600,
	})
	unsigned := b64url(header) + "." + b64url(claims)
	hash := sha256.Sum256([]byte(unsigned))
	sig, err := rsa.SignPKCS1v15(nil, key, crypto.SHA256, hash[:])
	if err != nil {
		return "", err
	}
	jwt := unsigned + "." + b64url(sig)

	resp, err := http.PostForm(sa.TokenURI, url.Values{
		"grant_type": {"urn:ietf:params:oauth:grant-type:jwt-bearer"},
		"assertion":  {jwt},
	})
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var tok struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if resp.StatusCode != 200 || json.Unmarshal(body, &tok) != nil || tok.AccessToken == "" {
		return "", fmt.Errorf("gagal mendapatkan access token (%d): %s", resp.StatusCode, body)
	}
	cachedToken = tok.AccessToken
	tokenExpiry = time.Now().Add(time.Duration(tok.ExpiresIn-60) * time.Second)
	return cachedToken, nil
}

func main() {
	sa, err := loadServiceAccount()
	if err != nil {
		log.Fatal(err)
	}
	if _, err := getAccessToken(sa); err != nil {
		log.Fatal(err)
	}
	log.Printf("Autentikasi OAuth 2.0 berhasil untuk project %s", sa.ProjectID)
}
