package keys

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

// Hasher abstracts PHC operations for admin & gateway.
type Hasher interface {
	// Hash generates a PHC string ($argon2id$...) for the secret using a random salt.
	Hash(secret []byte) (phc string, err error)
	// Verify checks a PHC string against the given secret.
	Verify(phc string, secret []byte) (bool, error)
}

// Argon2IDHasher implements Hasher using argon2id.
type Argon2IDHasher struct {
	Time    uint32 // iterations
	Memory  uint32 // KiB
	Threads uint8
	KeyLen  uint32 // bytes
}

func NewArgon2IDHasher(time uint32, memory uint32, threads uint8, keyLen uint32) *Argon2IDHasher {
	return &Argon2IDHasher{
		Time:    time,
		Memory:  memory,
		Threads: threads,
		KeyLen:  keyLen,
	}
}

func (h *Argon2IDHasher) Hash(secret []byte) (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	return h.hashWithSalt(salt, secret), nil
}

// hashWithSalt is handy if you ever want deterministic salts in tests.
func (h *Argon2IDHasher) hashWithSalt(salt, secret []byte) string {
	key := argon2.IDKey(secret, salt, h.Time, h.Memory, h.Threads, h.KeyLen)
	return fmt.Sprintf("$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		h.Memory, h.Time, h.Threads,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(key),
	)
}

func (h *Argon2IDHasher) Verify(phc string, secret []byte) (bool, error) {
	// Parse PHC: $argon2id$v=19$m=...,t=...,p=...$<saltb64>$<hashb64>
	if !strings.HasPrefix(phc, "$argon2id$") {
		return false, errors.New("unsupported hash")
	}
	parts := strings.Split(phc, "$")
	if len(parts) != 6 {
		return false, errors.New("bad phc format")
	}
	params := parts[3]
	var m uint32
	var t uint32
	var p uint8
	for _, kv := range strings.Split(params, ",") {
		switch {
		case strings.HasPrefix(kv, "m="):
			fmt.Sscanf(kv, "m=%d", &m)
		case strings.HasPrefix(kv, "t="):
			fmt.Sscanf(kv, "t=%d", &t)
		case strings.HasPrefix(kv, "p="):
			var pp uint32
			fmt.Sscanf(kv, "p=%d", &pp)
			p = uint8(pp)
		}
	}
	salt, err := base64.RawStdEncoding.Strict().DecodeString(parts[4])
	if err != nil {
		return false, err
	}
	want, err := base64.RawStdEncoding.Strict().DecodeString(parts[5])
	if err != nil {
		return false, err
	}

	got := argon2.IDKey(secret, salt, t, m, p, uint32(len(want)))
	return subtle.ConstantTimeCompare(got, want) == 1, nil
}
