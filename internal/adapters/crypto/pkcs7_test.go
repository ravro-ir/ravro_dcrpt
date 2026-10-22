package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"path/filepath"
	"testing"
	"time"

	mozpkcs7 "go.mozilla.org/pkcs7"
)

func TestDecryptPKCS7DoesNotPanicWithoutCert(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}

	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: "test"},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(time.Hour),
		KeyUsage:     x509.KeyUsageKeyEncipherment,
	}
	certDER, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	if err != nil {
		t.Fatal(err)
	}
	cert, err := x509.ParseCertificate(certDER)
	if err != nil {
		t.Fatal(err)
	}

	plaintext := []byte("ravro decrypt test payload")
	encrypted, err := mozpkcs7.Encrypt(plaintext, []*x509.Certificate{cert})
	if err != nil {
		t.Fatal(err)
	}

	dir := t.TempDir()
	keyPath := filepath.Join(dir, "key.pem")
	pemBlock := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)}
	if err := os.WriteFile(keyPath, pem.EncodeToMemory(pemBlock), 0600); err != nil {
		t.Fatal(err)
	}

	svc := NewPKCS7Service()
	got, err := svc.DecryptPKCS7(encrypted, keyPath)
	if err != nil {
		t.Fatalf("DecryptPKCS7: %v", err)
	}
	if string(got) != string(plaintext) {
		t.Fatalf("got %q want %q", got, plaintext)
	}
}
