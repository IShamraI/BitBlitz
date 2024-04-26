package wallet

import (
	"log"
)

type Wallet struct {
	privateKey    BitcoinPrivateKey
	publicKey     BitcoinPublicKey
	publicAddress string
	balance       int
}

func New(privateKey BitcoinPrivateKey, publicKey BitcoinPublicKey, publicAddress string) *Wallet {
	return &Wallet{
		privateKey:    privateKey,
		publicKey:     publicKey,
		publicAddress: publicAddress,
	}
}

func (w *Wallet) PrivateKey() BitcoinPrivateKey {
	return w.privateKey
}

func (w *Wallet) PrivateKeyWIF() string {
	return Bitcoin_Prikey2WIF(w.privateKey)
}

func (w *Wallet) PublicKey() BitcoinPublicKey {
	return w.publicKey
}

func (w *Wallet) PublicAddress() string {
	return w.publicAddress
}

func (w *Wallet) Balance() int {
	return w.balance
}

func GenWallet() (*Wallet, error) {
	// privateKey, publicKey, err := crypto.GenerateKeyAndAddress()
	// if err != nil {
	// 	return nil, err
	// }

	/* Generate a new ECDSA keypair */
	privateKey, publicKey, err := Bitcoin_GenerateKeypair()
	if err != nil {
		log.Fatalf("%s\n", err)
	}

	/* Convert the public key to a bitcoin network address */
	publicAddress := Bitcoin_Pubkey2Address(publicKey, 0x00)
	return New(privateKey, publicKey, publicAddress), nil
}
