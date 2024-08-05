package node

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/sudonite/blocker/crypto"
	"github.com/sudonite/blocker/proto"
	"github.com/sudonite/blocker/types"
	"github.com/sudonite/blocker/util"
)

func randomBlock(t *testing.T, chain *Chain) *proto.Block {
	privKey := crypto.GenerateNewPrivateKey()
	b := util.RandomBlock()

	prevBlock, err := chain.GetBlockByHeight(chain.Height())
	require.Nil(t, err)

	b.Header.PrevHash = types.HashBlock(prevBlock)

	types.SignBlock(privKey, b)

	return b
}

func TestNewChain(t *testing.T) {
	chain := NewChain(NewMemoryBlockStore())
	require.Equal(t, chain.Height(), 0)

	_, err := chain.GetBlockByHeight(0)
	require.Nil(t, err)
}

func TestChainHeight(t *testing.T) {
	chain := NewChain(NewMemoryBlockStore())
	for i := 0; i < 100; i++ {
		b := randomBlock(t, chain)
		require.Nil(t, chain.AddBlock(b))
		require.Equal(t, chain.Height(), i+1)
	}
}

func TestAddBlock(t *testing.T) {
	chain := NewChain(NewMemoryBlockStore())

	for i := 0; i < 100; i++ {
		block := randomBlock(t, chain)
		hashBlock := types.HashBlock(block)

		require.Nil(t, chain.AddBlock(block))

		fetchedBlock, err := chain.GetBlockByHash(hashBlock)
		require.Nil(t, err)
		require.Equal(t, block, fetchedBlock)

		fetchedBlockByHeight, err := chain.GetBlockByHeight(i + 1)
		require.Nil(t, err)
		require.Equal(t, block, fetchedBlockByHeight)
	}
}
