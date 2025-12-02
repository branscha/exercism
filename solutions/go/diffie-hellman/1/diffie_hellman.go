package diffiehellman

import (
    "math/big"
    "crypto/rand"
)

// Diffie-Hellman-Merkle key exchange
// Private keys should be generated randomly.

func PrivateKey(p *big.Int) *big.Int {
	// Generate random in range: [0, p-3)
    // Do not modify p !!!
    var max big.Int
    max.Add(p, big.NewInt(-3))
	n, err := rand.Int(rand.Reader, &max)
	if err != nil {
		return nil
	}
    // Result in range [2, p-1]
    return n.Add(n, big.NewInt(2))
}

func PublicKey(private, p *big.Int, g int64) *big.Int {
    var pk big.Int
    return pk.Exp(big.NewInt(g), private, p)
}

func NewPair(p *big.Int, g int64) (*big.Int, *big.Int) {
    priv := PrivateKey(p)
	pub := PublicKey(priv, p, g)
    return priv, pub
}

func SecretKey(private1, public2, p *big.Int) *big.Int {
	var s big.Int
    return s.Exp(public2, private1, p)
}
