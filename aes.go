package public

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

type AesAuth struct {
	KEY []byte
	IV 	[]byte
}

type cfb8 struct {
	b         cipher.Block
	blockSize int
	in        []byte
	out       []byte

	decrypt bool
}

func (x *cfb8) XORKeyStream(dst, src []byte) {
	for i := range src {
		x.b.Encrypt(x.out, x.in)
		copy(x.in[:x.blockSize-1], x.in[1:])
		if x.decrypt {
			x.in[x.blockSize-1] = src[i]
		}
		dst[i] = src[i] ^ x.out[0]
		if !x.decrypt {
			x.in[x.blockSize-1] = dst[i]
		}
	}
}

// NewCFB8Encrypter returns a Stream which encrypts with cipher feedback mode
// (segment size = 8), using the given Block. The iv must be the same length as
// the Block's block size.
func newCFB8Encrypter(block cipher.Block, iv []byte) cipher.Stream {
	return newCFB8(block, iv, false)
}

// NewCFB8Decrypter returns a Stream which decrypts with cipher feedback mode
// (segment size = 8), using the given Block. The iv must be the same length as
// the Block's block size.
func newCFB8Decrypter(block cipher.Block, iv []byte) cipher.Stream {
	return newCFB8(block, iv, true)
}

func newCFB8(block cipher.Block, iv []byte, decrypt bool) cipher.Stream {
	blockSize := block.BlockSize()
	if len(iv) != blockSize {
		// stack trace will indicate whether it was de or encryption
		panic("cipher.newCFB: IV length must equal block size")
	}
	x := &cfb8{
		b:         block,
		blockSize: blockSize,
		out:       make([]byte, blockSize),
		in:        make([]byte, blockSize),
		decrypt:   decrypt,
	}
	copy(x.in, iv)

	return x
}

func genKey(text string) []byte {
	padNum := len(text) % 16
	if padNum != 0 {
		for i := 0; i < 16-padNum; i++ {
			text += "\x00" // change to what you want
		}
	}
	return []byte(text)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func NewAesAuth() *AesAuth {
	key := genKey("servyou_")
	iv := RandStringBytes(16)
	return &AesAuth{
		KEY: key,
		IV:  iv,
	}

}



func RandStringBytes(n int) []byte {
	b := make([]byte, n)
	for i := range b {
		b[i] = LETTER_BYTES[rand.Intn(len(LETTER_BYTES))]
	}
	return b
}

func (this *AesAuth) AesEncryptCFB(origData []byte) (encrypted string) {
	block, err := aes.NewCipher(this.KEY)
	if err != nil {
		panic(err)
	}
	encrypt := make([]byte,len(origData))
	stream := newCFB8Encrypter(block, this.IV)
	stream.XORKeyStream(encrypt,origData)
	ivs := hex.EncodeToString(this.IV)
	enc := hex.EncodeToString(encrypt)
	tm := ivs + enc
	return tm
}
func (this *AesAuth) AesDecryptCFB(ciphertext string) (decrypted string,err error) {
	bs,err :=hex.DecodeString(ciphertext)
	if err!=nil{
		fmt.Printf("可能是明文,不支持的解码:%v",err)
		return
	}
	iv := bs[:aes.BlockSize]
	block, _ := aes.NewCipher(this.KEY)
	if len(iv) < aes.BlockSize {
		err = errors.New("密钥太短,无效")
		return
	}
	encrypted := bs[aes.BlockSize:]
	stream := newCFB8Decrypter(block,iv)
	stream.XORKeyStream(encrypted, encrypted)
	decrypted = string(encrypted)
	return
}







