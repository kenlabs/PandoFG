package main

import (
	"fmt"
	peerHelper "github.com/kenlabs/pando/pkg/util/peer"
	"os"
	"os/signal"
	"syscall"
	"time"

	pandoSdk "github.com/kenlabs/pando/sdk/pkg/provider"
)

const (
	privateKeyStr = "CAESQHWlReUYxW7FDvTAAqG+kNH2U7khW+iv0r+070+zKmFn9t80v5e30/NsBx5XzBLCE4uH/h3d3tpXlwCuO4YGN+w="
	pandoAddr     = "/ip4/127.0.0.1/tcp/9002"
	pandoPeerID   = "12D3KooWJCBFLWqSfdbABukxSBTN7KoYmFx96Td8QqCyRURbYx3J"
)

func main() {
	peerID, err := peerHelper.GetPeerIDFromPrivateKeyStr(privateKeyStr)
	if err != nil {
		panic(err)
	}
	fmt.Printf("provider peerID: %v\n", peerID.String())

	provider, err := pandoSdk.NewMetaProvider(privateKeyStr, 10*time.Second, 10*time.Minute)
	if err != nil {
		panic(err)
	}

	jsonData := "{\n\"date\": \"2022-04-21T12:44:12.332Z\",\n\"epoch\": 1740550,\n\"minerLocations\": [\n{\n\"miner\": \"f01012\",\n\"region\": \"AS-CN-EAST-ZJ\",\n\"long\": 120.219376,\n\"lat\": 30.259245,\n\"numLocations\": 2,\n\"country\": \"CN\",\n\"city\": \"Hangzhou\"\n},\n{\n\"miner\": \"f01012\",\n\"region\": \"AS-HK\",\n\"long\": 114.1657,\n\"lat\": 22.2578,\n\"numLocations\": 2,\n\"country\": \"HK\"\n},\n{\n\"miner\": \"f01152\",\n\"region\": \"AS-CN-NORTH-BJ\",\n\"long\": 116.40388,\n\"lat\": 39.91489,\n\"numLocations\": 1,\n\"country\": \"CN\",\n\"city\": \"Beijing\"\n}]}"

	err = provider.ConnectPando(pandoAddr, pandoPeerID)
	if err != nil {
		panic(err)
	}

	fmt.Println("pushing data to Pando...")
	metadata1, err := provider.NewMetadata([]byte(jsonData))
	if err != nil {
		panic(err)
	}
	metadata1Cid, err := provider.Push(metadata1)
	if err != nil {
		panic(err)
	}
	metadata2, err := provider.AppendMetadata(metadata1, []byte("kitty"))
	if err != nil {
		panic(err)
	}
	metadata2Cid, err := provider.Push(metadata2)
	if err != nil {
		panic(err)
	}
	fmt.Printf("pushed 2 nodes: \n\t%s\n\t%s\n", metadata1Cid.String(), metadata2Cid.String())

	//time.Sleep(20 * time.Second)

	// test for redundant push
	//_, _ = provider.Push(metadata1)

	fmt.Println("press ctrl+c to exit.")

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("Shutting down provider...")
	err = provider.Close()
	if err != nil {
		panic(err)
	}
	fmt.Println("Bye! ")
}
