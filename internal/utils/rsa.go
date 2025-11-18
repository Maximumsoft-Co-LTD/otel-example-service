package utils

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
)

// GenerateRSAKeyPair สร้างคู่กุญแจ RSA
func GenerateRSAKeyPair(bits int) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, err
	}
	return privKey, &privKey.PublicKey, nil
}

// SavePrivateKeyPKCS1 บันทึก private key แบบ PKCS#1
func SavePrivateKeyPKCS1(path string, priv *rsa.PrivateKey) error {
	privBytes := x509.MarshalPKCS1PrivateKey(priv)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privBytes,
	}
	return os.WriteFile(path, pem.EncodeToMemory(block), 0600)
}

// SavePrivateKeyPKCS8 บันทึก private key แบบ PKCS#8
func SavePrivateKeyPKCS8(path string, priv *rsa.PrivateKey) error {
	privBytes, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return err
	}
	block := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privBytes,
	}
	return os.WriteFile(path, pem.EncodeToMemory(block), 0600)
}

// SavePublicKeyPKIX บันทึก public key แบบ PKIX (X.509 SubjectPublicKeyInfo)
func SavePublicKeyPKIX(path string, pub *rsa.PublicKey) error {
	pubBytes, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return err
	}
	block := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubBytes,
	}
	return os.WriteFile(path, pem.EncodeToMemory(block), 0644)
}

// SavePublicKeyPKCS1 บันทึก public key แบบ PKCS#1
func SavePublicKeyPKCS1(path string, pub *rsa.PublicKey) error {
	pubBytes := x509.MarshalPKCS1PublicKey(pub)
	block := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubBytes,
	}
	return os.WriteFile(path, pem.EncodeToMemory(block), 0644)
}

// LoadRSAPublicKeyFromFile รองรับทั้ง "PUBLIC KEY" (PKIX) และ "RSA PUBLIC KEY" (PKCS#1)
func LoadRSAPublicKeyFromFile(path string) (*rsa.PublicKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("invalid public key PEM format")
	}

	switch block.Type {
	case "PUBLIC KEY":
		pub, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		rsaPub, ok := pub.(*rsa.PublicKey)
		if !ok {
			return nil, errors.New("not RSA public key")
		}
		return rsaPub, nil

	case "RSA PUBLIC KEY":
		rsaPub, err := x509.ParsePKCS1PublicKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		return rsaPub, nil

	default:
		return nil, errors.New("unsupported public key type: " + block.Type)
	}
}

// LoadRSAPrivateKeyFromFile รองรับทั้ง PKCS#8 ("PRIVATE KEY") และ PKCS#1 ("RSA PRIVATE KEY")
func LoadRSAPrivateKeyFromFile(path string) (*rsa.PrivateKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("invalid private key PEM format")
	}

	switch block.Type {
	case "PRIVATE KEY":
		priv, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		rsaPriv, ok := priv.(*rsa.PrivateKey)
		if !ok {
			return nil, errors.New("not RSA private key")
		}
		return rsaPriv, nil

	case "RSA PRIVATE KEY":
		priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		return priv, nil

	default:
		return nil, errors.New("unsupported private key type: " + block.Type)
	}
}

// RSAEncrypt เข้ารหัสด้วย RSA public key
func RSAEncrypt(pub *rsa.PublicKey, plaintext []byte) ([]byte, error) {
	return rsa.EncryptPKCS1v15(rand.Reader, pub, plaintext)
}

// RSADecrypt ถอดรหัสด้วย RSA private key
func RSADecrypt(priv *rsa.PrivateKey, ciphertext []byte) ([]byte, error) {
	return rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
}

// RSAEncryptOAEP เข้ารหัสด้วย RSA-OAEP (ใช้ SHA256)
func RSAEncryptOAEP(pub *rsa.PublicKey, plaintext []byte, label []byte) ([]byte, error) {
	hash := sha256.New()
	return rsa.EncryptOAEP(hash, rand.Reader, pub, plaintext, label)
}

// RSADecryptOAEP ถอดรหัสด้วย RSA-OAEP (ใช้ SHA256)
func RSADecryptOAEP(priv *rsa.PrivateKey, ciphertext []byte, label []byte) ([]byte, error) {
	hash := sha256.New()
	return rsa.DecryptOAEP(hash, rand.Reader, priv, ciphertext, label)
}

// RSASignSHA256 เซ็นข้อความด้วย RSA + SHA256
func RSASignSHA256(priv *rsa.PrivateKey, message []byte) ([]byte, error) {
	hash := sha256.Sum256(message)
	return rsa.SignPKCS1v15(rand.Reader, priv, crypto.SHA256, hash[:])
}

// RSAVerifySHA256 ตรวจสอบลายเซ็นด้วย RSA + SHA256
func RSAVerifySHA256(pub *rsa.PublicKey, message []byte, signature []byte) error {
	hash := sha256.Sum256(message)
	return rsa.VerifyPKCS1v15(pub, crypto.SHA256, hash[:], signature)
}
