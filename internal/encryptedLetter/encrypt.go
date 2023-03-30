package encryptedLetter

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
	"log"
	"os"

	"github.com/joho/godotenv"

	_ "github.com/Destinyxus/botLetterToFuture/pkg/config"
)

type EncrypterDecrypter interface {
	Encrypt(letter string, key string) (string, error)
	Decrypt(letter string) (string, error)
}

type encrypter struct {
}

func NewEncrypter() EncrypterDecrypter {
	return &encrypter{}
}
func encodeBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func decodeBase64(s string) []byte {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return data
}

func (e *encrypter) Encrypt(letter string, key string) (string, error) {

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	cipherText := make([]byte, aes.BlockSize+len(letter))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(cipherText[aes.BlockSize:], []byte(letter))
	return encodeBase64(cipherText), nil
}

func (e *encrypter) Decrypt(letter string) (string, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	token := os.Getenv("HASH_KEY")
	block, err := aes.NewCipher([]byte(token))
	if err != nil {
		return "", err
	}
	cipherText := decodeBase64(letter)
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(cipherText, cipherText)
	return string(cipherText), nil
}
