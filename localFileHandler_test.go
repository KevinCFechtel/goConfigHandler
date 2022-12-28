package configHandler

import (
	"bytes"
	"testing"
)

type Configuration struct {
	Test bool `json:"Test"`
}

func TestWriteAndReadKeyFile(t *testing.T) {
	newKey, err := generateNewKey()
	if err != nil {
		t.Fatalf("got error %s ", err)
	}
	salt := "test"
	estimatedKey := prepareKey(newKey, salt)
	tempFilePath := t.TempDir() + "/testKey.key"

	writeKeyToFile(tempFilePath, estimatedKey)
	readFileResult, err := readKeyFile(tempFilePath, salt)
	if err != nil {
		t.Fatalf("Error reading Testfile: %v", err)
	}

	if !bytes.Equal(estimatedKey, readFileResult) {
		t.Fatalf("got wrong value from Keyfile")
	}
}

func TestReadConfigFile(t *testing.T) {
	plaintext := "{\"Test\":true}"
	tempFilePath := t.TempDir() + "/testKey.key"

	writeFile(tempFilePath, []byte(plaintext))
	readFileResult, err := readConfigFile(tempFilePath)
	if err != nil {
		t.Fatalf("Error reading Testfile: %v", err)
	}

	if plaintext != string(readFileResult) {
		t.Fatalf("got wrong value from Configfile: %s", string(readFileResult))
	}
}

func TestWriteConfigFileWithKey(t *testing.T) {
	plaintext := "{\"Test\":true}"
	tempFilePath := t.TempDir()
	tempConfigFile := tempFilePath + "/config.json"
	salt := "test"
	newKey, err := generateNewKey()
	if err != nil {
		t.Fatalf("got error %s ", err)
	}
	preparedKey := prepareKey(newKey, salt)
	tempKeyFile := tempFilePath + "/config.key"

	writeKeyToFile(tempKeyFile, newKey)

	config := Configuration{}
	config.Test = true

	err = exportConfigToLocalFile(tempConfigFile, &config, true, salt)
	if err != nil {
		t.Fatalf("Error writing Testfile: %v", err)
	}

	writtenConfigFile, err := readConfigFile(tempConfigFile)
	if err != nil {
		t.Fatalf("got error %s ", err)
	}

	writtenConfigFile, err = decrypt(preparedKey, writtenConfigFile)
	if err != nil {
		t.Fatalf("got error %s ", err)
	}

	if string(writtenConfigFile) != plaintext {
		t.Fatalf("got wrong value from Configfile: %s", writtenConfigFile)
	}
}

func TestWriteConfigFileWithoutKey(t *testing.T) {
	plaintext := "{\"Test\":true}"
	tempFilePath := t.TempDir()
	tempConfigFile := tempFilePath + "/config.json"
	tempKeyFile := tempFilePath + "/config.key"
	salt := "test"

	config := Configuration{}
	config.Test = true

	err := exportConfigToLocalFile(tempConfigFile, &config, true, salt)
	if err != nil {
		t.Fatalf("Error writing Testfile: %v", err)
	}

	writtenConfigFile, err := readConfigFile(tempConfigFile)
	if err != nil {
		t.Fatalf("got error %s ", err)
	}

	writtenKey, err := readKeyFile(tempKeyFile, salt)
	if err != nil {
		t.Fatalf("got error %s ", err)
	}

	preparedKey := prepareKey(writtenKey, salt)

	writtenConfigFile, err = decrypt(preparedKey, writtenConfigFile)
	if err != nil {
		t.Fatalf("got error %s ", err)
	}

	if string(writtenConfigFile) != plaintext {
		t.Fatalf("got wrong value from Configfile: %s", writtenConfigFile)
	}
}

func TestWriteConfigFileUnencrypted(t *testing.T) {
	plaintext := "{\"Test\":true}"
	tempFilePath := t.TempDir()
	tempConfigFile := tempFilePath + "/config.json"
	salt := "test"

	config := Configuration{}
	config.Test = true

	err := exportConfigToLocalFile(tempConfigFile, &config, false, salt)
	if err != nil {
		t.Fatalf("Error writing Testfile: %v", err)
	}

	writtenConfigFile, err := readConfigFile(tempConfigFile)
	if err != nil {
		t.Fatalf("got error %s ", err)
	}

	if string(writtenConfigFile) != plaintext {
		t.Fatalf("got wrong value from Configfile: %s", writtenConfigFile)
	}
}
