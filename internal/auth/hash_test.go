package auth

import (
	"encoding/base64"
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "testPassword123"

	salt, hash, err := hashPassword(password)
	if err != nil {
		t.Fatalf("hashPassword returned error: %v", err)
	}

	if salt == "" {
		t.Error("salt should not be empty")
	}

	if hash == "" {
		t.Error("hash should not be empty")
	}

	// Verify salt is valid base64
	if _, err := base64.RawStdEncoding.DecodeString(salt); err != nil {
		t.Errorf("salt is not valid base64: %v", err)
	}

	// Verify hash is valid base64
	if _, err := base64.RawStdEncoding.DecodeString(hash); err != nil {
		t.Errorf("hash is not valid base64: %v", err)
	}
}

func TestHashPasswordUniqueness(t *testing.T) {
	password := "samePassword"

	salt1, hash1, err1 := hashPassword(password)
	salt2, hash2, err2 := hashPassword(password)

	if err1 != nil || err2 != nil {
		t.Fatalf("hashPassword returned errors: %v, %v", err1, err2)
	}

	// Salts should be different (random)
	if salt1 == salt2 {
		t.Error("salts should be unique for each hash operation")
	}

	// Hashes should be different because salts are different
	if hash1 == hash2 {
		t.Error("hashes should be different when salts are different")
	}
}

func TestVerifyPasswordCorrect(t *testing.T) {
	password := "mySecurePassword"

	salt, hash, err := hashPassword(password)
	if err != nil {
		t.Fatalf("hashPassword failed: %v", err)
	}

	if !verifyPassword(password, salt, hash) {
		t.Error("verifyPassword should return true for correct password")
	}
}

func TestVerifyPasswordIncorrect(t *testing.T) {
	password := "correctPassword"
	wrongPassword := "wrongPassword"

	salt, hash, err := hashPassword(password)
	if err != nil {
		t.Fatalf("hashPassword failed: %v", err)
	}

	if verifyPassword(wrongPassword, salt, hash) {
		t.Error("verifyPassword should return false for incorrect password")
	}
}

func TestVerifyPasswordInvalidSalt(t *testing.T) {
	password := "testPassword"

	_, hash, err := hashPassword(password)
	if err != nil {
		t.Fatalf("hashPassword failed: %v", err)
	}

	if verifyPassword(password, "invalidBase64!@#", hash) {
		t.Error("verifyPassword should return false for invalid salt")
	}
}

func TestVerifyPasswordInvalidHash(t *testing.T) {
	password := "testPassword"

	salt, _, err := hashPassword(password)
	if err != nil {
		t.Fatalf("hashPassword failed: %v", err)
	}

	if verifyPassword(password, salt, "invalidBase64!@#") {
		t.Error("verifyPassword should return false for invalid hash")
	}
}

func TestVerifyPasswordEmptyPassword(t *testing.T) {
	emptyPassword := ""

	salt, hash, err := hashPassword(emptyPassword)
	if err != nil {
		t.Fatalf("hashPassword failed for empty password: %v", err)
	}

	if !verifyPassword(emptyPassword, salt, hash) {
		t.Error("verifyPassword should work with empty passwords")
	}

	if verifyPassword("notEmpty", salt, hash) {
		t.Error("verifyPassword should return false for non-matching password")
	}
}

func TestHashPasswordConsistency(t *testing.T) {
	password := "consistentPassword"

	salt, hash, err := hashPassword(password)
	if err != nil {
		t.Fatalf("hashPassword failed: %v", err)
	}

	// Verify multiple times to ensure consistency
	for i := 0; i < 10; i++ {
		if !verifyPassword(password, salt, hash) {
			t.Errorf("verification failed on attempt %d", i+1)
		}
	}
}
