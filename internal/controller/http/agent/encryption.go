package agent

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	config "github.com/gynshu-one/go-metric-collector/internal/config/agent"
	"github.com/rs/zerolog/log"
	"io"
	"os"
)

var publicKey *rsa.PublicKey

func init() {
	loadPublicKey(config.GetConfig().CryptoKey)
}

func encryptWithPublicKey(body []byte) []byte {
	if config.GetConfig().CryptoKey == "" {
		return body
	}
	if publicKey == nil {
		return body
	}
	// Generate a new AES key
	aesKey := make([]byte, 32) // 256 bits
	if _, err := io.ReadFull(rand.Reader, aesKey); err != nil {
		log.Error().Err(err).Msg("Error generating AES key")
		return body
	}

	// Encrypt the AES key with the RSA public key
	encryptedAESKey, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, aesKey, nil)
	if err != nil {
		log.Error().Err(err).Msg("Error encrypting AES key")
		return body
	}

	// Create a new AES cipher block
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		log.Error().Err(err).Msg("Error creating AES cipher block")
		return body
	}

	// Create a new GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Error().Err(err).Msg("Error creating GCM")
		return body
	}

	// Create a new nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		log.Error().Err(err).Msg("Error creating nonce")
		return body
	}

	// Encrypt the data using AES-GCM
	ciphertext := gcm.Seal(nonce, nonce, body, nil)

	// Return the RSA-encrypted AES key and the AES-encrypted data, both base64-encoded and separated by a colon
	return []byte(base64.StdEncoding.EncodeToString(encryptedAESKey) + ":" + base64.StdEncoding.EncodeToString(ciphertext))

}

func loadPublicKey(path string) {
	publicKeyData, err := os.ReadFile(path)
	if err != nil {
		log.Error().Err(err).Msg("Error reading public key")
		return
	}

	// Parse the public key
	block, _ := pem.Decode(publicKeyData)
	if block == nil {
		log.Error().Err(err).Msg("Error decoding public key")
		return
	}

	rsaPublicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		log.Error().Err(err).Msg("Error parsing public key")
		return
	}

	var ok bool
	publicKey, ok = rsaPublicKey.(*rsa.PublicKey)
	if !ok {
		log.Error().Err(err).Msg("Error casting public key")
		return
	}
}
