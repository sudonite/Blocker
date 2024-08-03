package main

import (
	"context"
	"log"
	"time"

	"github.com/sudonite/blocker/crypto"
	"github.com/sudonite/blocker/node"
	"github.com/sudonite/blocker/proto"
	"github.com/sudonite/blocker/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	makeNode("127.0.0.1:3000", []string{}, true)
	time.Sleep(time.Second)
	makeNode("127.0.0.1:4000", []string{"127.0.0.1:3000"}, false)
	time.Sleep(time.Second)
	makeNode("127.0.0.1:5000", []string{"127.0.0.1:4000"}, false)

	for {
		time.Sleep(time.Second * 3)
		makeTransaction()
	}
}

func makeNode(listenAddr string, bootstrapNodes []string, isValidator bool) *node.Node {
	cfg := node.ServerConfig{
		Version:    "blocker-0.1",
		ListenAddr: listenAddr,
	}

	if isValidator {
		cfg.PrivateKey = crypto.GenerateNewPrivateKey()
	}

	n := node.NewNode(cfg)

	go n.Start(listenAddr, bootstrapNodes)

	return n
}

func makeTransaction() {
	client, err := grpc.NewClient("127.0.0.1:3000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	c := proto.NewNodeClient(client)
	privKey := crypto.GenerateNewPrivateKey()

	tx := &proto.Transaction{
		Version: 1,
		Inputs: []*proto.TxInput{
			{
				PrevTxHash:   util.RandomHash(),
				PrevOutIndex: 0,
				PublicKey:    privKey.Public().Bytes(),
			},
		},
		Outputs: []*proto.TxOutput{
			{
				Amount:  99,
				Address: privKey.Public().Address().Byte(),
			},
		},
	}

	_, err = c.HandleTransaction(context.TODO(), tx)
	if err != nil {
		log.Fatal(err)
	}
}
