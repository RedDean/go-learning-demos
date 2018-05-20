package block

import (
	"time"
)

type Block struct {
	Timestamp 		int64
	Data	  		[]byte
	Hash      		[]byte
	PrevBlockHash 	[]byte
	Nouce 			int
}


func NewBlock (data string, prevBlockHash []byte,) *Block {
	block := &Block{time.Now().Unix(), []byte(data), []byte{}, prevBlockHash, 0}
	pow := NewProofOfWork(block)
	nouce,hash := pow.Run()
	block.Hash = hash
	block.Nouce = nouce

	return block
}