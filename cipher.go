/*
 * @Author: tr3e
 * @Date: 2019-11-26 20:50:43
 * @Last Modified by: tr3e
 * @Last Modified time: 2019-11-26 20:56:38
 */

package main

import "errors"

type Cipher interface {
	Encrypt([]byte) []byte
	Decrypt([]byte) []byte
	Copy() Cipher
	Reset()
}

var cipherMethod = map[string]func(string) (Cipher, error){
	"plain": NewPlainCipher,
	"xor":   NewXorCipher,
}

func NewCipher(method, password string) (Cipher, error) {
	cipher, ok := cipherMethod[method]
	if !ok {
		return nil, errors.New("Unsupported cipher method")
	}
	return cipher(password)
}

type PlainCipher struct {
}

func NewPlainCipher(string) (Cipher, error) {
	return &PlainCipher{}, nil
}

func (c *PlainCipher) Encrypt(p []byte) []byte {
	return p
}

func (c *PlainCipher) Decrypt(p []byte) []byte {
	return p
}

func (c *PlainCipher) Copy() Cipher {
	return c
}

func (c *PlainCipher) Reset() {

}

type XorCipher struct {
	encInd int
	decInd int
	secret string
}

func NewXorCipher(secret string) (Cipher, error) {
	return &XorCipher{secret: secret}, nil
}

func (c *XorCipher) Encrypt(p []byte) []byte {
	for i := 0; i < len(p); i++ {
		c.encInd %= len(c.secret)
		p[i] ^= c.secret[c.encInd]
		c.encInd++
	}
	return p
}

func (c *XorCipher) Decrypt(p []byte) []byte {
	for i := 0; i < len(p); i++ {
		c.decInd %= len(c.secret)
		p[i] ^= c.secret[c.decInd]
		c.decInd++
	}
	return p
}

func (c *XorCipher) Copy() Cipher {
	return &XorCipher{secret: c.secret}
}

func (c *XorCipher) Reset() {
	c.encInd = 0
	c.decInd = 0
}
