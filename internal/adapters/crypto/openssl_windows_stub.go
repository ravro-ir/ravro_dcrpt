//go:build windows && !opensslcgo
// +build windows,!opensslcgo

package crypto

import "ravro_dcrpt/internal/ports"

// NewOpenSSLCGOService uses the Pure Go PKCS7 implementation on Windows
// unless the opensslcgo build tag is set. This lets GUI binaries ship
// without OpenSSL DLLs and be cross-compiled from Linux.
func NewOpenSSLCGOService() ports.CryptoService {
	return NewPKCS7Service()
}
