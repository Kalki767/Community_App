package helper

import (
    "crypto/sha512"
    "encoding/hex"
)

func HashTokenSHA512(token string) string {
    hash := sha512.Sum512([]byte(token))
    return hex.EncodeToString(hash[:]) // 128 hex characters
}

func CompareTokenSHA512(token, storedHash string) bool {
    return HashTokenSHA512(token) == storedHash
}
