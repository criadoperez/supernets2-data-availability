package batch

import (
	"crypto/ecdsa"

	"github.com/0xPolygon/cdk-validium-node/jsonrpc/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// Batch represents a batch used for synchronization
type Batch struct {
	Number         types.ArgUint64 `json:"number"`
	GlobalExitRoot common.Hash     `json:"globalExitRoot"`
	Timestamp      types.ArgUint64 `json:"timestamp"`
	Coinbase       common.Address  `json:"coinbase"`
	L2Data         types.ArgBytes  `json:"batchL2Data"`
}

// HashToSign returns a hash that uniquely identifies the batch
func (b *Batch) HashToSign() []byte {
	return crypto.Keccak256(
		[]byte(b.Number.Hex()),
		b.GlobalExitRoot[:],
		[]byte(b.Timestamp.Hex()),
		b.Coinbase[:],
		b.L2Data,
	)
}

// Sign returns a signed batch by the private key
func (b *Batch) Sign(privateKey *ecdsa.PrivateKey) (*SignedBatch, error) {
	hashToSign := b.HashToSign()
	sig, err := crypto.Sign(hashToSign, privateKey)
	if err != nil {
		return nil, err
	}
	return &SignedBatch{
		Batch:     *b,
		Signature: sig,
	}, nil
}

// SignedBatch is a batch but signed
type SignedBatch struct {
	Batch     Batch          `json:"batch"`
	Signature types.ArgBytes `json:"signature"`
}

// Signer returns the address of the signer
func (s *SignedBatch) Signer() (common.Address, error) {
	pubKey, err := crypto.SigToPub(s.Batch.HashToSign(), s.Signature)
	if err != nil {
		return common.Address{}, err
	}
	return crypto.PubkeyToAddress(*pubKey), nil
}
