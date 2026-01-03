package identity

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
)

func GenerateKeyPairPEM() (privatePEM, publicPEM []byte, err error) {
    privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
    if err != nil { return nil, nil, err }

    privatePEM = pem.EncodeToMemory(&pem.Block{
        Type:  "RSA PRIVATE KEY",
        Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
    })

    publicPEM = pem.EncodeToMemory(&pem.Block{
        Type:  "RSA PUBLIC KEY",
        Bytes: x509.MarshalPKCS1PublicKey(&privateKey.PublicKey),
    })

    return privatePEM, publicPEM, nil
}

func SaveKeyPair(dir string, privatePEM, publicPEM []byte) error {
    if err := os.MkdirAll(dir, 0700); err != nil { return err }
    if err := os.WriteFile(dir+"/private.key", privatePEM, 0o600); err != nil { return err }
    if err := os.WriteFile(dir+"/public.key", publicPEM, 0o644); err != nil { return err }
    return nil
}