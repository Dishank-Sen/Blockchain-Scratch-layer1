package utils

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/Dishank-Sen/Blockchain-Scratch-layer1/types"
)

/* -------------------- DIR -------------------- */

func CreateDir(ctx context.Context, cancel context.CancelFunc, dirPath string, reinit bool) error {
	if reinit {
		if _, err := os.Stat(dirPath); err == nil {
			return nil
		} else if !os.IsNotExist(err) {
			return err
		}
	}

	if err := os.MkdirAll(dirPath, 0700); err != nil {
		return err
	}

	if ctx.Err() != nil {
		cancel()
		return errors.New("operation canceled during directory creation")
	}
	return nil
}

/* -------------------- KEYS -------------------- */

func CreateKeys(ctx context.Context, cancel context.CancelFunc, dirPath string, reinit bool) error {
	privatePath := filepath.Join(dirPath, "private.key")
	publicPath := filepath.Join(dirPath, "public.key")

	if reinit {
		if _, err := os.Stat(privatePath); err == nil {
			if _, err := os.Stat(publicPath); err == nil {
				return nil
			}
		}
	}

	privatePEM, publicPEM, err := generateKeyPairPEM()
	if err != nil {
		return err
	}

	if err := saveKeyPair(dirPath, privatePEM, publicPEM); err != nil {
		return err
	}

	if ctx.Err() != nil {
		cancel()
		return errors.New("operation canceled during key creation")
	}
	return nil
}

func generateKeyPairPEM() ([]byte, []byte, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	privatePEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	publicPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(&privateKey.PublicKey),
	})

	return privatePEM, publicPEM, nil
}

func saveKeyPair(dir string, privatePEM, publicPEM []byte) error {
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dir, "private.key"), privatePEM, 0o600); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dir, "public.key"), publicPEM, 0o644); err != nil {
		return err
	}
	return nil
}

/* -------------------- METADATA -------------------- */

func CreateMetadata(ctx context.Context, cancel context.CancelFunc, filePath string, hashPath string, reinit bool) error {
	if reinit {
		info, err := os.Stat(filePath)
		if err == nil && info.Size() > 0 {
			return nil
		}
		if err != nil && !os.IsNotExist(err) {
			return err
		}
	}

	id, err := hashFile(hashPath)
	if err != nil {
		return err
	}

	m := types.Metadata{ID: id}
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}

	if err := os.WriteFile(filePath, data, 0755); err != nil {
		return err
	}

	if ctx.Err() != nil {
		cancel()
		return errors.New("operation canceled during metadata creation")
	}
	return nil
}

/* -------------------- HASH -------------------- */

func hashFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return "", err
	}
	if stat.Size() == 0 {
		return "", nil
	}

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}
