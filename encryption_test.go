package configHandler

import (
	"bytes"
	"testing"
)

func TestCheckIfTextIsEncrypted(t *testing.T) {
	plaintext := "{\"Test\":true}"
	encryptedText := "Encrypted"

	testPlaintext := checkIfTextIsEncrypted([]byte(plaintext))
	testEncryptedText := checkIfTextIsEncrypted([]byte(encryptedText))

	if testPlaintext {
		t.Fatalf("got %t for Plaintext %s", testPlaintext, plaintext)
	}

	if !testEncryptedText {
		t.Fatalf("got %t for encrypted Text \"%s\"", testEncryptedText, encryptedText)
	}
}

func TestGenerateNewKey(t *testing.T) {
	newKey, err := generateNewKey()
	if err != nil {
		t.Fatalf("got error %s ", err)
	}

	if len(string(newKey)) != 32 {
		t.Fatalf("got a Key with length of %d ", len(string(newKey)))
	}
}

func TestPrepareKey(t *testing.T) {
	newKey, err := generateNewKey()
	if err != nil {
		t.Fatalf("got error %s ", err)
	}
	newKeyString := string(newKey)
	salt := "Test"
	estimatedKey := []byte(salt + newKeyString[len(salt):])
	preaperdKey := prepareKey(newKey, salt)

	if !bytes.Equal(preaperdKey, estimatedKey) {
		t.Fatalf("Keys are not equal")
	}
}

func TestEncryptionAndDecryption(t *testing.T) {
	plaintext := "{\"Test\":true}"
	newKey, err := generateNewKey()
	if err != nil {
		t.Fatalf("got error %s ", err)
	}
	salt := "Test"
	preaperdKey := prepareKey(newKey, salt)
	encryptedText, err := encrypt(preaperdKey, []byte(plaintext))
	if err != nil {
		t.Fatalf("got error %s ", err)
	}
	decryptedText, err := decrypt(preaperdKey, encryptedText)
	if err != nil {
		t.Fatalf("got error %s ", err)
	}

	if string(decryptedText) != plaintext {
		t.Fatalf("got wrong value for decrypted Text \"%s\"", decryptedText)
	}
}
