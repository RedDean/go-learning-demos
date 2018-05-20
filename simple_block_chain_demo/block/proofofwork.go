package block

import (
	"math/big"
	"bytes"
	"BC_eg1/util"
	"crypto/sha256"
	"math"
	"fmt"
)

const TargetBits = 24 //挖矿条件，即哈希符合条件，算出的哈希值前24位为零

const maxnouce = math.MaxInt64 // 计算器限制，防止计算次数过大导致溢出

type ProofOfWork struct {
	block *Block
	target *big.Int
}

func NewProofOfWork (b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256 - TargetBits)) // 左移操作，实现前24位为零
	return &ProofOfWork{b, target}
}

func (p *ProofOfWork)PrepareData(nouce int) []byte {
	data := bytes.Join([][]byte{
			 p.block.Data,
			 p.block.PrevBlockHash,
			 util.IntToHex(p.block.Timestamp),
			 util.IntToHex(int64(TargetBits)),
			 util.IntToHex(int64(nouce)),
	},[]byte{})

	return data
}

func (p *ProofOfWork)Run( ) (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	noucenum := 0

	fmt.Printf("Start mining the block containing data : %x\n", p.block.Data)
	for noucenum < maxnouce  {
		originHash := p.PrepareData(noucenum)
		hash = sha256.Sum256(originHash)
		hashInt.SetBytes(hash[:])
		if hashInt.Cmp(p.target) == -1 {
			fmt.Printf("/r hash is %x", hash)
			break
		} else {
			noucenum++
		}
	}
	fmt.Println()
	return noucenum, hash[:]
}

func (p *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := p.PrepareData(p.block.Nouce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	isValid := hashInt.Cmp(p.target) == -1
	return isValid
}
