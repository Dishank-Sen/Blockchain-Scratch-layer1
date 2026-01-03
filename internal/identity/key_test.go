package identity_test

import (
	"crypto/x509"
	"encoding/pem"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/Dishank-Sen/Blockchain-Scratch-layer1/internal/identity"
)

// Test that GenerateKeyPairPEM returns PEM blocks that parse as RSA keys.
func TestGenerateKeyPairPEM_ParsesAsRSA(t *testing.T) {
    privPEM, pubPEM, err := identity.GenerateKeyPairPEM()
    if err != nil {
        t.Fatalf("GenerateKeyPairPEM error: %v", err)
    }

    // parse private
    block, _ := pem.Decode(privPEM)
    if block == nil || block.Type != "RSA PRIVATE KEY" {
        t.Fatalf("private PEM decode failed or wrong type: %v", block)
    }
    priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
    if err != nil {
        t.Fatalf("parse private key failed: %v", err)
    }

    // parse public
    blockPub, _ := pem.Decode(pubPEM)
    if blockPub == nil || blockPub.Type != "RSA PUBLIC KEY" {
        t.Fatalf("public PEM decode failed or wrong type: %v", blockPub)
    }
    pub, err := x509.ParsePKCS1PublicKey(blockPub.Bytes)
    if err != nil {
        t.Fatalf("parse public key failed: %v", err)
    }

    // check modulus / exponent match
    if priv.PublicKey.N.Cmp(pub.N) != 0 || priv.PublicKey.E != pub.E {
        t.Fatalf("public key does not match private key")
    }
}

// Test SaveKeyPair writes files with correct permissions and content.
func TestSaveKeyPair_WritesFilesAndPerms(t *testing.T) {
    privPEM, pubPEM, err := identity.GenerateKeyPairPEM()
    if err != nil { t.Fatalf("gen: %v", err) }

    dir := t.TempDir()
    if err := identity.SaveKeyPair(dir, privPEM, pubPEM); err != nil {
        t.Fatalf("SaveKeyPair failed: %v", err)
    }

    // check existence and permissions
    p := filepath.Join(dir, "private.key")
    info, err := os.Stat(p)
    if err != nil { t.Fatalf("private.key missing: %v", err) }
    // check file mode bits (owner read/write)
    if info.Mode().Perm() & 0o777 != 0o600 {
        t.Fatalf("private.key perms incorrect: %v", info.Mode().Perm())
    }

    q := filepath.Join(dir, "public.key")
    info2, err := os.Stat(q)
    if err != nil { t.Fatalf("public.key missing: %v", err) }
    if info2.Mode().Perm() & 0o777 != 0o644 {
        t.Fatalf("public.key perms incorrect: %v", info2.Mode().Perm())
    }

    // confirm content decodeable
    privBytes, _ := os.ReadFile(p)
    if block, _ := pem.Decode(privBytes); block == nil {
        t.Fatalf("private.key not valid PEM")
    }
    pubBytes, _ := os.ReadFile(q)
    if block, _ := pem.Decode(pubBytes); block == nil {
        t.Fatalf("public.key not valid PEM")
    }
}

// Optional: smoke test that SaveKeyPair + Parse roundtrip keeps pair consistent.
func TestRoundtrip_PrivatePublicMatch(t *testing.T) {
    privPEM, pubPEM, err := identity.GenerateKeyPairPEM()
    if err != nil { t.Fatalf("gen: %v", err) }

    dir := t.TempDir()
    if err := identity.SaveKeyPair(dir, privPEM, pubPEM); err != nil { t.Fatalf("save: %v", err) }

    // parse back
    privBytes, _ := os.ReadFile(filepath.Join(dir, "private.key"))
    block, _ := pem.Decode(privBytes)
    priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
    if err != nil { t.Fatalf("parse priv: %v", err) }

    pubBytes, _ := os.ReadFile(filepath.Join(dir, "public.key"))
    blockPub, _ := pem.Decode(pubBytes)
    pub, err := x509.ParsePKCS1PublicKey(blockPub.Bytes)
    if err != nil { t.Fatalf("parse pub: %v", err) }

    if !reflect.DeepEqual(priv.PublicKey.N, pub.N) || priv.PublicKey.E != pub.E {
        t.Fatalf("keys mismatch after roundtrip")
    }
}
