package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509/pkix"
	"encoding/asn1"
	"errors"
	"fmt"
	"math/big"

	mozpkcs7 "go.mozilla.org/pkcs7"
)

type pkcs7ContentInfo struct {
	ContentType asn1.ObjectIdentifier
	Content     asn1.RawValue `asn1:"explicit,optional,tag:0"`
}

type pkcs7EnvelopedData struct {
	Version              int
	RecipientInfos       []pkcs7RecipientInfo `asn1:"set"`
	EncryptedContentInfo pkcs7EncryptedContentInfo
}

type pkcs7RecipientInfo struct {
	Version                int
	IssuerAndSerialNumber  pkcs7IssuerAndSerial
	KeyEncryptionAlgorithm pkix.AlgorithmIdentifier
	EncryptedKey           []byte
}

type pkcs7IssuerAndSerial struct {
	IssuerName   asn1.RawValue
	SerialNumber *big.Int
}

type pkcs7EncryptedContentInfo struct {
	ContentType                asn1.ObjectIdentifier
	ContentEncryptionAlgorithm pkix.AlgorithmIdentifier
	EncryptedContent           asn1.RawValue `asn1:"tag:0,optional"`
}

type pkcs7AESGCMParameters struct {
	Nonce  []byte `asn1:"tag:4"`
	ICVLen int
}

func decryptPKCS7WithKey(encryptedData []byte, key *rsa.PrivateKey) ([]byte, error) {
	p7, err := mozpkcs7.Parse(encryptedData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse PKCS7 data: %w", err)
	}

	for _, cert := range p7.Certificates {
		if cert == nil {
			continue
		}
		plain, err := p7.Decrypt(cert, key)
		if err == nil {
			return plain, nil
		}
	}

	plain, err := decryptEnvelopedWithRSA(encryptedData, key)
	if err == nil {
		return plain, nil
	}
	return nil, fmt.Errorf("failed to decrypt data: %w", err)
}

func decryptEnvelopedWithRSA(encryptedData []byte, key *rsa.PrivateKey) ([]byte, error) {
	var info pkcs7ContentInfo
	rest, err := asn1.Unmarshal(encryptedData, &info)
	if err != nil || len(rest) > 0 {
		return nil, fmt.Errorf("failed to unmarshal PKCS7 content info: %w", err)
	}
	if !info.ContentType.Equal(mozpkcs7.OIDEnvelopedData) {
		return nil, errors.New("pkcs7: content is not enveloped data")
	}

	var env pkcs7EnvelopedData
	if _, err := asn1.Unmarshal(info.Content.Bytes, &env); err != nil {
		return nil, fmt.Errorf("failed to unmarshal enveloped data: %w", err)
	}
	if len(env.RecipientInfos) == 0 {
		return nil, errors.New("pkcs7: no recipients")
	}

	var lastErr error
	for _, recp := range env.RecipientInfos {
		if len(recp.EncryptedKey) == 0 {
			continue
		}
		contentKey, err := rsa.DecryptPKCS1v15(rand.Reader, key, recp.EncryptedKey)
		if err != nil {
			lastErr = err
			continue
		}
		plain, err := decryptEncryptedContent(env.EncryptedContentInfo, contentKey)
		if err != nil {
			lastErr = err
			continue
		}
		return plain, nil
	}
	if lastErr == nil {
		lastErr = errors.New("pkcs7: no enveloped recipient for provided key")
	}
	return nil, lastErr
}

func decryptEncryptedContent(eci pkcs7EncryptedContentInfo, key []byte) ([]byte, error) {
	alg := eci.ContentEncryptionAlgorithm.Algorithm
	supported := alg.Equal(mozpkcs7.OIDEncryptionAlgorithmDESCBC) ||
		alg.Equal(mozpkcs7.OIDEncryptionAlgorithmDESEDE3CBC) ||
		alg.Equal(mozpkcs7.OIDEncryptionAlgorithmAES256CBC) ||
		alg.Equal(mozpkcs7.OIDEncryptionAlgorithmAES128CBC) ||
		alg.Equal(mozpkcs7.OIDEncryptionAlgorithmAES128GCM) ||
		alg.Equal(mozpkcs7.OIDEncryptionAlgorithmAES256GCM)
	if !supported {
		return nil, fmt.Errorf("pkcs7: unsupported content encryption algorithm %s", alg)
	}

	var cyphertext []byte
	if eci.EncryptedContent.IsCompound {
		var buf bytes.Buffer
		cypherbytes := eci.EncryptedContent.Bytes
		for {
			var part []byte
			cypherbytes, _ = asn1.Unmarshal(cypherbytes, &part)
			buf.Write(part)
			if cypherbytes == nil {
				break
			}
		}
		cyphertext = buf.Bytes()
	} else {
		cyphertext = eci.EncryptedContent.Bytes
	}

	var block cipher.Block
	var err error
	switch {
	case alg.Equal(mozpkcs7.OIDEncryptionAlgorithmDESCBC):
		block, err = des.NewCipher(key)
	case alg.Equal(mozpkcs7.OIDEncryptionAlgorithmDESEDE3CBC):
		block, err = des.NewTripleDESCipher(key)
	default:
		block, err = aes.NewCipher(key)
	}
	if err != nil {
		return nil, err
	}

	if alg.Equal(mozpkcs7.OIDEncryptionAlgorithmAES128GCM) || alg.Equal(mozpkcs7.OIDEncryptionAlgorithmAES256GCM) {
		var params pkcs7AESGCMParameters
		if _, err := asn1.Unmarshal(eci.ContentEncryptionAlgorithm.Parameters.Bytes, &params); err != nil {
			return nil, err
		}
		gcm, err := cipher.NewGCM(block)
		if err != nil {
			return nil, err
		}
		if len(params.Nonce) != gcm.NonceSize() || params.ICVLen != gcm.Overhead() {
			return nil, errors.New("pkcs7: encryption algorithm parameters are incorrect")
		}
		return gcm.Open(nil, params.Nonce, cyphertext, nil)
	}

	iv := eci.ContentEncryptionAlgorithm.Parameters.Bytes
	if len(iv) != block.BlockSize() {
		return nil, errors.New("pkcs7: encryption algorithm parameters are malformed")
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	plaintext := make([]byte, len(cyphertext))
	mode.CryptBlocks(plaintext, cyphertext)
	return pkcs7Unpad(plaintext, mode.BlockSize())
}

func pkcs7Unpad(data []byte, blocklen int) ([]byte, error) {
	if blocklen < 1 {
		return nil, fmt.Errorf("invalid blocklen %d", blocklen)
	}
	if len(data)%blocklen != 0 || len(data) == 0 {
		return nil, fmt.Errorf("invalid data len %d", len(data))
	}
	padlen := int(data[len(data)-1])
	if padlen < 1 || padlen > blocklen || padlen > len(data) {
		return nil, errors.New("invalid padding")
	}
	pad := data[len(data)-padlen:]
	for _, padbyte := range pad {
		if padbyte != byte(padlen) {
			return nil, errors.New("invalid padding")
		}
	}
	return data[:len(data)-padlen], nil
}
