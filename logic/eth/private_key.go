package eth

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

type PrivateKey struct {
	key *ecdsa.PrivateKey
}

func NewPrivateKey(hex string) (*PrivateKey, error) {
	priKey, err := crypto.HexToECDSA(hex)
	if err != nil {
		return nil, err
	}

	return &PrivateKey{key: priKey}, nil
}

func (p *PrivateKey) Signature(data []byte) ([]byte, error) {
	hash := crypto.Keccak256Hash(data)

	return crypto.Sign(hash.Bytes(), p.key)
}

func (p *PrivateKey) PublicKey() *PublicKey {
	pub := p.key.Public()

	if publicKeyECDSA, ok := pub.(*ecdsa.PublicKey); ok {
		return &PublicKey{key: publicKeyECDSA}
	} else {
		return nil
	}
}

func (p *PrivateKey) Hex() []byte {
	return crypto.FromECDSA(p.key)
}

func (p *PrivateKey) HexString() string {
	return hexutil.Encode(p.Hex())
}

func (p *PrivateKey) Address() common.Address {
	return p.PublicKey().Address()
}
