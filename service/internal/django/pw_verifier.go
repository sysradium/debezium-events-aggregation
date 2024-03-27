package django

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/pbkdf2"
)

var (
	ErrInvalidHash = errors.New("invalid hash")
)

func VerifyDjangoPassword(hash, plaintextPassword string) (bool, error) {
	parts := strings.Split(hash, "$")
	if len(parts) != 4 {
		return false, ErrInvalidHash
	}

	iterations := 0
	fmt.Sscanf(parts[1], "%d", &iterations)
	salt := parts[2]
	storedHash := parts[3]

	encodedHash := base64.StdEncoding.EncodeToString(
		pbkdf2.Key([]byte(plaintextPassword), []byte(salt), iterations, sha256.Size, sha256.New),
	)

	return encodedHash == storedHash, nil
}
