package main

import (
	"fmt"
	"strconv"
	"BC_eg1/block"
)

type BlockChain struct {
	blockChain []*block.Block
}

func (bc *BlockChain)addBlock(data string)  {
	prevBlock := bc.blockChain[len(bc.blockChain)-1]
	newBlock := block.NewBlock(data, prevBlock.Hash)
	bc.blockChain = append(bc.blockChain, newBlock)
}

func NewGenesisBlock() *block.Block {
	return block.NewBlock("源块", []byte{})
}

func NewBlockchain() *BlockChain {
	return &BlockChain{[]*block.Block{NewGenesisBlock()}}
}

func main() {
	bc := NewBlockchain()
	bc.addBlock("dean大佬的雷霆队总冠军")
	bc.addBlock("dean大佬的基金涨到飞起")
	bc.addBlock("世界和平，没有战争")
	for _, b := range bc.blockChain {
		fmt.Printf("Prev.Hash: %x\n", b.PrevBlockHash)
		fmt.Printf("data: %s\n", b.Data)
		fmt.Printf("Hash: %x\n", b.Hash)
		pow := block.NewProofOfWork(b)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
	}
}
