package data

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base32"
	"time"

	"movie.api/internal/validator"
)

const (
	ScopeActivation = "activation"
)

// * Data Model.

// Token
type Token struct {
	Plaintext string
	Hash      []byte
	UserID    int64
	Expiry    time.Time
	Scope     string
}

// generateToken
func generateToken(userID int64, ttl time.Duration, scope string) (*Token, error) {
	token := &Token{
		UserID: userID,
		Expiry: time.Now().Add(ttl),
		Scope:  scope,
	}
	randomBytes := make([]byte, 16)

	// fills the byte slice with random bytes from operating system's CSPRNG.
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	// omits padding
	token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)

	hash := sha256.Sum256([]byte(token.Plaintext))
	token.Hash = hash[:]
	return token, nil
}

// It’s important to point out that the plaintext token strings here are like
// Y3QMGX3PJ3WLRL2YRTQGQ6KRHU are not 16 characters long — but rather they have an underlying entropy of 16 bytes of randomness.

// The length of the plaintext token string itself depends on how those 16 random bytes are encoded to create a string.
// In our case we encode the random bytes to a base-32 string, which results in a string with 26 characters.
// In contrast, if we encoded the random bytes using hexadecimal (base-16) the string would be 32 characters long instead.

// * Validate.

// ValidateTokenPlaintext
func ValidateTokenPlaintext(v *validator.Validator, tokenPlaintext string) {
	v.Check(tokenPlaintext != "", "token", "must be provided")
	v.Check(len(tokenPlaintext) == 26, "token", "must be 26 bytes long")
}

// * CRUD.

// Define the TokenModel type.
type TokenModel struct {
	DB *sql.DB
}

// New
func (m TokenModel) New(userID int64, ttl time.Duration, scope string) (*Token, error) {
	token, err := generateToken(userID, ttl, scope)
	if err != nil {
		return nil, err
	}
	err = m.Insert(token)
	return token, err
}

// Insert adds the data for a specific token to the tokens table.
func (m TokenModel) Insert(token *Token) error {
	query := `
	INSERT INTO tokens (hash, user_id, expiry, scope)
	VALUES ($1, $2, $3, $4)`
	args := []any{token.Hash, token.UserID, token.Expiry, token.Scope}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, query, args...)
	return err
}

// DeleteAllForUser deletes all tokens for a specific user and scope.
func (m TokenModel) DeleteAllForUser(scope string, userID int64) error {
	query := `
	DELETE FROM tokens
	WHERE scope = $1 AND user_id = $2`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, query, scope, userID)
	return err
}
