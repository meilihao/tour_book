package token

import (
	"errors"
	"log"
	"time"

	"golang.org/x/crypto/ed25519"
	"golang.org/x/crypto/ssh"
	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

var (
	ErrTokenExpired = errors.New("token is expired (exp)")
)

var ed25519PrivateKey ed25519.PrivateKey
var ed25519PublicKey ed25519.PublicKey

func init() {
	ed25519PrivateKeyRaw := []byte(`
-----BEGIN OPENSSH PRIVATE KEY-----
b3BXXX...XXX=
-----END OPENSSH PRIVATE KEY-----
`)

	pri, err := ssh.ParseRawPrivateKey(ed25519PrivateKeyRaw)
	if err != nil {
		log.Fatalln(err)
	}

	ed25519PrivateKey = *(pri.(*ed25519.PrivateKey))
	ed25519PublicKey = ed25519PrivateKey.Public().(ed25519.PublicKey)
}

type Token struct {
	UserId int64 `json:"user_id"`
	RootId int64 `json:"root_id,omitempty"`
	Role   byte  `json:"role,omitempty"`
}

func GenerateBaseToken(userId, rootId int64, role byte) (string, error) {
	opts := &jose.SignerOptions{}
	opts.WithType("JWT")

	siger, err := jose.NewSigner(jose.SigningKey{
		Algorithm: jose.EdDSA,
		Key:       ed25519PrivateKey,
	}, opts)
	if err != nil {
		return "", err
	}

	cl := jwt.Claims{
		Expiry: jwt.NumericDate(time.Now().AddDate(0, 0, 30).Unix()),
	}

	t := Token{
		UserId: userId,
		RootId: rootId,
		Role:   role,
	}

	raw, err := jwt.Signed(siger).Claims(cl).Claims(t).CompactSerialize()
	if err != nil {
		return "", err
	}

	return raw, nil
}

func ParseToken(raw string) (*Token, error) {
	tok, err := jwt.ParseSigned(raw)
	if err != nil {
		return nil, err
	}

	cl := jwt.Claims{}
	t := new(Token)
	if err = tok.Claims(ed25519PublicKey, &cl, t); err != nil {
		return nil, err
	}

	if time.Now().Unix() > int64(cl.Expiry) {
		return nil, ErrTokenExpired
	}

	return t, nil
}
