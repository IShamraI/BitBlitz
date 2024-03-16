package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"

	"github.com/IShamraI/BitBlitz/internal/bech32"
	"golang.org/x/crypto/ripemd160"
)

// publicKeyToAddress calculates the Bitcoin address corresponding to an ECDSA
// public key. The resulting address will be a bech32-encoded string in the
// 'bc' prefix.
//
// The public key is first SHA256-hashed, and then the resulting hash is
// RIPEMD160-hashed. The resulting hash is then run through the bech32
// conversion algorithm to produce the final Bitcoin address.
//
// Refer to BIP13 for more information on the bech32 encoding format.
func PublicKeyToAddress(publicKey ecdsa.PublicKey) (string, error) {
	// SHA256-hash the public key
	sha256Hash := sha256.New()
	sha256Hash.Write(elliptic.Marshal(publicKey.Curve, publicKey.X, publicKey.Y))
	pubKeySHA256 := sha256Hash.Sum(nil)

	// RIPEMD160-hash the SHA256-hashed public key
	ripemd160Hash := ripemd160.New()
	ripemd160Hash.Write(pubKeySHA256)
	pubKeyHash := ripemd160Hash.Sum(nil)

	// Convert the RIPEMD160-hashed public key to bech32 using the 'bc' prefix
	converted, err := bech32.ConvertBits(pubKeyHash, 8, 5, true)
	if err != nil {
		return "", err
	}
	bech32Addr, err := bech32.Encode("bc", converted)
	if err != nil {
		return "", err
	}

	return bech32Addr, nil
}

// Sha256Checksum calculates the double SHA256 hash of the input and returns the
// first four bytes of the result. This is used to create the checksum for a
// Bitcoin transaction's input scripts and outputs.
//
// The resulting slice should be exactly four bytes long.
func Sha256Checksum(input []byte) []byte {
	hash := sha256.Sum256(input)  // Calculate the first SHA256 hash
	hash = sha256.Sum256(hash[:]) // Calculate the second SHA256 hash
	return hash[:4]               // Return the first four bytes of the second hash
}

// GenerateKeyAndAddress generates a new private key and corresponding Bitcoin
// address using secp256k1 elliptic curve cryptography. The private key is
// encoded in hex and the address is in bech32 format.
func GenerateKeyAndAddress() (privKeyHex string, address string, err error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return
	}

	publicKey := privateKey.PublicKey
	address, err = PublicKeyToAddress(publicKey)
	if err != nil {
		return
	}

	privKeyHex = hex.EncodeToString(privateKey.D.Bytes())
	return
}
