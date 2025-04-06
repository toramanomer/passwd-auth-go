package password

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/argon2"
)

func generateSalt() []byte {
	salt := make([]byte, 16)
	rand.Read(salt)
	return salt
}

func HashPassword(password string) string {
	var (
		salt    = generateSalt()
		time    = uint32(3)
		memory  = uint32(256 * 1024)
		threads = uint8(8)
		keyLen  = uint32(32)
	)

	hash := argon2.IDKey([]byte(password), salt, time, memory, threads, keyLen)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	return fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		memory,
		time,
		threads,
		b64Salt,
		b64Hash,
	)
}
