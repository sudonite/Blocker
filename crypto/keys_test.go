package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateNewPrivateKey(t *testing.T) {
	privKey := GenerateNewPrivateKey()
	assert.Equal(t, len(privKey.Bytes()), PrivKeyLen)

	pubKey := privKey.Public()
	assert.Equal(t, len(pubKey.Bytes()), PubKeyLen)
}

func TestGenerateNewPrivateKeyFromString(t *testing.T) {
	var (
		seed       = "285e7df3919251c3434e6b18a5b79eb5bb39866df2f4f2dac938911934ae30fd"
		addressStr = "26efee8eaed9212727f9e86d8abf0545af184f52"
		privKey    = GeneratePrivateKeyFromString(seed)
		address    = privKey.Public().Address()
	)

	assert.Equal(t, len(privKey.Bytes()), PrivKeyLen)
	assert.Equal(t, address.String(), addressStr)
}

func TestPrivateKeySign(t *testing.T) {
	privKey := GenerateNewPrivateKey()
	pubKey := privKey.Public()
	msg := []byte("foo bar baz")

	sig := privKey.Sign(msg)

	assert.True(t, sig.Verify(*pubKey, msg))
	assert.False(t, sig.Verify(*pubKey, []byte("foo")))

	invalidPrivKey := GenerateNewPrivateKey()
	invalidPubKey := invalidPrivKey.Public()

	assert.False(t, sig.Verify(*invalidPubKey, msg))
}

func TestPublicKeyToAddress(t *testing.T) {
	privKey := GenerateNewPrivateKey()
	pubKey := privKey.Public()
	address := pubKey.Address()

	assert.Equal(t, len(address.Byte()), AddressLen)
	//fmt.Println(address)
}
