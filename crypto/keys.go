package crypto

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"io"
)

const (
	PrivKeyLen   = 64
	PubKeyLen    = 32
	SeedLen      = 32
	AddressLen   = 20
	SignatureLen = 64
)

type PrivateKey struct {
	key ed25519.PrivateKey
}

func GeneratePrivateKeyFromString(s string) *PrivateKey {
	b, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}

	return GeneratePrivateKeyFromSeed(b)
}

func GeneratePrivateKeyFromSeed(seed []byte) *PrivateKey {
	if len(seed) != SeedLen {
		panic("invalid seed length")
	}

	return &PrivateKey{
		key: ed25519.NewKeyFromSeed(seed),
	}
}

func GenerateNewPrivateKeyFromSeedStr(seed string) *PrivateKey {
	seedBytes, err := hex.DecodeString(seed)
	if err != nil {
		panic(err)
	}

	return GeneratePrivateKeyFromSeed(seedBytes)
}

func GenerateNewPrivateKey() *PrivateKey {
	seed := make([]byte, SeedLen)

	_, err := io.ReadFull(rand.Reader, seed)
	if err != nil {
		panic(err)
	}

	return &PrivateKey{
		key: ed25519.NewKeyFromSeed(seed),
	}
}

func (p *PrivateKey) Bytes() []byte {
	return p.key
}

func (p *PrivateKey) Sign(msg []byte) *Signature {
	return &Signature{
		value: ed25519.Sign(p.key, msg),
	}
}

func (p *PrivateKey) Public() *PublicKey {
	b := make([]byte, PubKeyLen)
	copy(b, p.key[32:])

	return &PublicKey{
		key: b,
	}
}

type PublicKey struct {
	key ed25519.PublicKey
}

func PublicKeyFromBytes(b []byte) *PublicKey {
	if len(b) != PubKeyLen {
		panic("invalid public key length")
	}

	return &PublicKey{
		key: ed25519.PublicKey(b),
	}
}

func (p *PublicKey) Bytes() []byte {
	return p.key
}

func (p *PublicKey) Address() Address {
	return Address{
		value: p.key[len(p.key)-AddressLen:],
	}
}

type Signature struct {
	value []byte
}

func SignatureFromBytes(b []byte) *Signature {
	if len(b) != SignatureLen {
		panic("invalid signature length")
	}

	return &Signature{
		value: b,
	}
}

func (s *Signature) Bytes() []byte {
	return s.value
}

func (s *Signature) Verify(pubKey PublicKey, msg []byte) bool {
	return ed25519.Verify(pubKey.key, msg, s.value)
}

type Address struct {
	value []byte
}

func AddressFromBytes(b []byte) Address {
	if len(b) != AddressLen {
		panic("length of the (address) bytes not equal to 20")
	}

	return Address{
		value: b,
	}
}

func (a Address) Bytes() []byte {
	return a.value
}

func (a Address) String() string {
	return hex.EncodeToString(a.value)
}
