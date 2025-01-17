package node

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
	chain := NewChain(NewMemoryBlockStore(), NewMemoryTxStore())
	require.Equal(t, chain.Height(), 0)

	_, err := chain.GetBlockByHeight(0)
	require.Nil(t, err)
}

func TestChainHeight(t *testing.T) {
	chain := NewChain(NewMemoryBlockStore(), NewMemoryTxStore())
	for i := 0; i < 100; i++ {
		b := randomBlock(t, chain)
		require.Nil(t, chain.AddBlock(b))
		require.Equal(t, chain.Height(), i+1)
	}
}

func TestAddBlock(t *testing.T) {
	chain := NewChain(NewMemoryBlockStore(), NewMemoryTxStore())

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

func TestAddBlockWithTxInsufficientFunds(t *testing.T) {
	var (
		chain     = NewChain(NewMemoryBlockStore(), NewMemoryTxStore())
		block     = randomBlock(t, chain)
		privKey   = crypto.GenerateNewPrivateKeyFromSeedStr(godSeed)
		recepient = crypto.GenerateNewPrivateKey().Public().Address().Bytes()
	)

	prevTx, err := chain.txStore.Get("a13f7bd36acf0e20c43f770acdde37461a092d0390bdd84d9349ff1cf70a794c")
	assert.Nil(t, err)

	inputs := []*proto.TxInput{
		{
			PrevOutIndex: 0,
			PrevTxHash:   types.HashTransaction(prevTx),
			PublicKey:    privKey.Public().Bytes(),
		},
	}
	outputs := []*proto.TxOutput{
		{Amount: 1001, Address: recepient},
	}
	tx := &proto.Transaction{Version: 1, Inputs: inputs, Outputs: outputs}

	sig := types.SignTransaction(privKey, tx)
	tx.Inputs[0].Signature = sig.Bytes()

	block.Transactions = append(block.Transactions, tx)

	require.NotNil(t, chain.AddBlock(block))
}

func TestAddBlockWithTx(t *testing.T) {
	var (
		chain     = NewChain(NewMemoryBlockStore(), NewMemoryTxStore())
		block     = randomBlock(t, chain)
		privKey   = crypto.GenerateNewPrivateKeyFromSeedStr(godSeed)
		recepient = crypto.GenerateNewPrivateKey().Public().Address().Bytes()
	)

	prevTx, err := chain.txStore.Get("a13f7bd36acf0e20c43f770acdde37461a092d0390bdd84d9349ff1cf70a794c")
	assert.Nil(t, err)

	inputs := []*proto.TxInput{
		{
			PrevOutIndex: 0,
			PrevTxHash:   types.HashTransaction(prevTx),
			PublicKey:    privKey.Public().Bytes(),
		},
	}
	outputs := []*proto.TxOutput{
		{Amount: 100, Address: recepient},
		{Amount: 900, Address: privKey.Public().Address().Bytes()},
	}
	tx := &proto.Transaction{Version: 1, Inputs: inputs, Outputs: outputs}

	sig := types.SignTransaction(privKey, tx)
	tx.Inputs[0].Signature = sig.Bytes()

	block.Transactions = append(block.Transactions, tx)
	types.SignBlock(privKey, block)
	require.Nil(t, chain.AddBlock(block))
}
