package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/sudonite/blocker/crypto"
	"github.com/sudonite/blocker/util"
)

func TestSignBlock(t *testing.T) {
	var (
		block   = util.RandomBlock()
		privKey = crypto.GenerateNewPrivateKey()
		pubKey  = privKey.Public()
	)

	sig := SignBlock(privKey, block)
	assert.Equal(t, 64, len(sig.Bytes()))
	assert.True(t, sig.Verify(*pubKey, HashBlock(block)))
}

func TestHashBlock(t *testing.T) {
	block := util.RandomBlock()
	hash := HashBlock(block)

	assert.Equal(t, 32, len(hash))
}
