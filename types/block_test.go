package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/sudonite/blocker/crypto"
	"github.com/sudonite/blocker/proto"
	"github.com/sudonite/blocker/util"
)

func TestCalculateRootHash(t *testing.T) {
	var (
		privKey = crypto.GenerateNewPrivateKey()
		block   = util.RandomBlock()
		tx      = &proto.Transaction{
			Version: 1,
		}
	)

	block.Transactions = append(block.Transactions, tx)
	SignBlock(privKey, block)

	assert.True(t, VerifyRootHash(block))
	assert.Equal(t, 32, len(block.Header.RootHash))
}

func TestSignVerifyBlock(t *testing.T) {
	var (
		block   = util.RandomBlock()
		privKey = crypto.GenerateNewPrivateKey()
		pubKey  = privKey.Public()
	)

	sig := SignBlock(privKey, block)
	assert.Equal(t, 64, len(sig.Bytes()))
	assert.True(t, sig.Verify(*pubKey, HashBlock(block)))

	assert.Equal(t, pubKey.Bytes(), block.PublicKey)
	assert.Equal(t, sig.Bytes(), block.Signature)
	assert.True(t, VerifyBlock(block))

	invalidPrivKey := crypto.GenerateNewPrivateKey()
	block.PublicKey = invalidPrivKey.Public().Bytes()

	assert.False(t, VerifyBlock(block))
}

func TestHashBlock(t *testing.T) {
	block := util.RandomBlock()
	hash := HashBlock(block)

	assert.Equal(t, 32, len(hash))
}
