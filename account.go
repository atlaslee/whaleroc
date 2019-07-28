/*
MIT License

Copyright (c) 2019 Atlas Lee, 4859345@qq.com

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package dmt

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"golang.org/x/crypto/ripemd160"
	"math/big"
)

const (
	SIZEOF_ACCOUNT = SIZEOF_ADDRESS + SIZEOF_PRIVATEKEY + SIZEOF_PUBLICKEY
)

var (
	GLOBALRAND = make([]byte, 2*SIZEOF_BIGINT)
	_, _       = rand.Read(GLOBALRAND)
)

// 账号
type Account struct {
	address [SIZEOF_ADDRESS]byte
	priKey  [SIZEOF_PRIVATEKEY]byte
	pubKey  [SIZEOF_PUBLICKEY]byte
}

func b58Encode(b []byte) (s string) {

	/* Convert big endian bytes to big int */
	x := new(big.Int).SetBytes(b)

	/* Initialize */
	r := new(big.Int)
	m := big.NewInt(58)
	zero := big.NewInt(0)
	s = ""

	/* Convert big int to string */
	for x.Cmp(zero) > 0 {
		/* x, r = (x / 58, x % 58) */
		x.QuoRem(x, m, r)
		/* Prepend ASCII character */
		s = string(BITCOIN_BASE58_TABLE[r.Int64()]) + s
	}

	return s
}

func sha256Encode(b []byte) []byte {
	hash := sha256.New()
	hash.Reset()
	hash.Write(b)
	return hash.Sum(nil)
}

func ripemd160Encode(b []byte) []byte {
	hash := ripemd160.New()
	hash.Reset()
	hash.Write(b)
	return hash.Sum(nil)

}

func (this *Account) Address() [SIZEOF_ADDRESS]byte {
	return this.address
}

func (this *Account) AddressString() string {
	return b58Encode(this.address[:])
}

func (this *Account) PublicKey() *ecdsa.PublicKey {
	return &ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     big.NewInt(0).SetBytes(this.pubKey[:SIZEOF_BIGINT]),
		Y:     big.NewInt(0).SetBytes(this.pubKey[SIZEOF_BIGINT:])}
}

func (this *Account) PrivateKey() *ecdsa.PrivateKey {
	return &ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: elliptic.P256(),
			X:     big.NewInt(0).SetBytes(this.pubKey[:SIZEOF_BIGINT]),
			Y:     big.NewInt(0).SetBytes(this.pubKey[SIZEOF_BIGINT:])},
		D: big.NewInt(0).SetBytes(this.priKey[:])}
}

func (this *Account) PrivateKeyString() string {
	return base64.URLEncoding.EncodeToString(this.priKey[:])
}

func (this *Account) SetPrivateKey(priKey []byte) {
	copy(this.priKey[:], priKey)

	x, y := elliptic.P256().ScalarBaseMult(priKey)

	copy(this.pubKey[:SIZEOF_BIGINT], x.Bytes())
	copy(this.pubKey[SIZEOF_BIGINT:], y.Bytes())
	copy(this.address[:], ripemd160Encode(sha256Encode(this.pubKey[:])))
}

func (this *Account) SetPrivateKeyString(s string) (err error) {
	b, err := base64.URLEncoding.DecodeString(s)
	if err != nil {
		return
	}

	this.SetPrivateKey(b)
	return
}

func (this *Account) Sign(data []byte) (sign []byte) {
	h := sha256.New()
	h.Reset()
	h.Write(data)
	hash := h.Sum(nil)

	reader := bytes.NewReader(GLOBALRAND)
	sign = make([]byte, SIZEOF_SIGNATURE)

	r, s, _ := ecdsa.Sign(reader, this.PrivateKey(), hash)

	copy(sign[:SIZEOF_PUBLICKEY], this.pubKey[:])
	copy(sign[SIZEOF_PUBLICKEY:SIZEOF_PUBLICKEY+SIZEOF_BIGINT], r.Bytes())
	copy(sign[SIZEOF_PUBLICKEY+SIZEOF_BIGINT:], s.Bytes())
	return
}

func SignatureVerify(data []byte, sign []byte, t byte) bool {
	h := sha256.New()
	h.Reset()
	h.Write(data)
	hash := h.Sum(nil)

	x := big.NewInt(0).SetBytes(sign[:SIZEOF_BIGINT])
	y := big.NewInt(0).SetBytes(sign[SIZEOF_BIGINT : 2*SIZEOF_BIGINT])
	r := big.NewInt(0).SetBytes(sign[SIZEOF_PUBLICKEY : 3*SIZEOF_BIGINT])
	s := big.NewInt(0).SetBytes(sign[3*SIZEOF_BIGINT:])

	return ecdsa.Verify(&ecdsa.PublicKey{elliptic.P256(), x, y}, hash, r, s)
}

func NewAccount() *Account {
	return &Account{}
}

func GenAccount() (account *Account) {
	account = NewAccount()

	b := make([]byte, SIZEOF_RANDOM)
	rand.Read(b)
	r := bytes.NewReader(b)

	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), r)
	copy(account.priKey[:], privateKey.D.Bytes())
	copy(account.pubKey[:SIZEOF_PRIVATEKEY], privateKey.X.Bytes())
	copy(account.pubKey[SIZEOF_PRIVATEKEY:], privateKey.Y.Bytes())
	copy(account.address[:], ripemd160Encode(sha256Encode(account.pubKey[:])))
	return
}

func AccountLoad(priKey []byte) (account *Account) {
	account = NewAccount()
	account.SetPrivateKey(priKey)
	return
}
