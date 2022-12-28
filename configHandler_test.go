package configHandler

import "testing"

func TestReadConfigFromLocalFileEncrypted(t *testing.T) {
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
	configFileText, err := encrypt(preparedKey, []byte(plaintext))
	if err != nil {
		t.Fatalf("got error %s ", err)
	}

	writeFile(tempConfigFile, configFileText)

	config := Configuration{}
	config.Test = false

	err = GetConfig("localFile", tempConfigFile, &config, salt)
	if err != nil {
		t.Fatalf("Error reading Testfile: %v", err)
	}

	if config.Test != true {
		t.Fatalf("got wrong value from Configfile: %t", config.Test)
	}
}

func TestReadConfigFromLocalFileUnencrypted(t *testing.T) {
	plaintext := "{\"Test\":true}"
	tempFilePath := t.TempDir()
	tempConfigFile := tempFilePath + "/config.json"
	writeFile(tempConfigFile, []byte(plaintext))

	configFileWritten, err := readConfigFile(tempConfigFile)
	if err != nil {
		t.Fatalf("got error %s ", err)
	}

	if checkIfTextIsEncrypted(configFileWritten) {
		t.Fatalf("Config file is not written unencrypted")
	}

	salt := "test"
	newKey, err := generateNewKey()
	if err != nil {
		t.Fatalf("got error %s ", err)
	}
	tempKeyFile := tempFilePath + "/config.key"

	writeKeyToFile(tempKeyFile, newKey)

	config := Configuration{}
	config.Test = false

	err = GetConfig("localFile", tempConfigFile, &config, salt)
	if err != nil {
		t.Fatalf("Error reading Testfile: %v", err)
	}

	configFileWritten, err = readConfigFile(tempConfigFile)
	if err != nil {
		t.Fatalf("got error %s ", err)
	}

	if !checkIfTextIsEncrypted(configFileWritten) {
		t.Fatalf("Unencrypted config file was not encrypted")
	}
}

func TestReadConfigFromLocalFileUnencyptedWithoutKey(t *testing.T) {
	plaintext := "{\"Test\":true}"
	tempFilePath := t.TempDir()
	tempConfigFile := tempFilePath + "/config.json"
	writeFile(tempConfigFile, []byte(plaintext))

	configFileWritten, err := readConfigFile(tempConfigFile)
	if err != nil {
		t.Fatalf("got error %s ", err)
	}

	if checkIfTextIsEncrypted(configFileWritten) {
		t.Fatalf("Config file is not written unencrypted")
	}

	salt := "test"
	tempKeyFile := tempFilePath + "/config.key"

	config := Configuration{}
	config.Test = false

	err = GetConfig("localFile", tempConfigFile, &config, salt)
	if err != nil {
		t.Fatalf("Error reading Testfile: %v", err)
	}

	configFileWritten, err = readConfigFile(tempConfigFile)
	if err != nil {
		t.Fatalf("got error %s ", err)
	}

	if !checkIfTextIsEncrypted(configFileWritten) {
		t.Fatalf("Unencrypted config file was not encrypted")
	}

	key, err := readKeyFile(tempKeyFile, salt)
	if err != nil {
		t.Fatalf("got error %s ", err)
	}

	if len(key) == 0 {
		t.Fatalf("No new keyfile was written")
	}
}
