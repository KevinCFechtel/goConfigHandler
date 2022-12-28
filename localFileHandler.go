package configHandler

import (
	"encoding/base64"
	"encoding/json"
	"os"
	"path/filepath"
)

func writeFile(filePath string, file []byte) {
	os.WriteFile(filePath, file, 0777)
}

func writeKeyToFile(keyFilePath string, key []byte) {
	keyString := base64.StdEncoding.EncodeToString(key)
	writeFile(keyFilePath, []byte(keyString))
}

func readKeyFile(keyFilePath string, salt string) ([]byte, error) {
	keyString, err := os.ReadFile(keyFilePath)
	if err != nil {
		return nil, err
	}
	key, err := base64.StdEncoding.DecodeString(string(keyString))
	if err != nil {
		return nil, err
	}
	return prepareKey([]byte(key), salt), err
}

func readConfigFile(configFilePath string) ([]byte, error) {
	file, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, err
	}
	return file, err
}

func readConfigFromLocalFile(configFilePath string, configStruct any, salt string) error {
	var keyFilePath string
	var err error
	var file []byte
	var key []byte

	path := filepath.Dir(configFilePath)
	if path == "." {
		keyFilePath = "config.key"
	} else {
		keyFilePath = path + "/config.key"
	}

	file, err = readConfigFile(configFilePath)
	if err != nil {
		return err
	}

	key, err = readKeyFile(keyFilePath, salt)
	if err != nil {
		key, err = generateNewKey()
		if err != nil {
			return err
		}
		writeKeyToFile(keyFilePath, key)
	}

	key = prepareKey(key, salt)

	if file != nil {
		if checkIfTextIsEncrypted(file) {
			file, err = decrypt(key, file)
			if err != nil {
				return err
			}
			err = json.Unmarshal([]byte(file), &configStruct)
			if err != nil {
				return err
			}
		} else {
			err = json.Unmarshal([]byte(file), &configStruct)
			if err != nil {
				return err
			}
			file, err = encrypt(key, file)
			if err != nil {
				return err
			}
			writeFile(configFilePath, file)
		}
	}
	return err
}

func exportConfigToLocalFile(configFilePath string, configStruct any, encrypted bool, salt string) error {
	var keyFilePath string
	var err error
	var file []byte

	path := filepath.Dir(configFilePath)
	if path == "." {
		keyFilePath = "config.key"
	} else {
		keyFilePath = path + "/config.key"
	}

	key, err := readKeyFile(keyFilePath, salt)
	if err != nil {
		key, err = generateNewKey()
		if err != nil {
			return err
		}
		writeKeyToFile(keyFilePath, key)
	}

	key = prepareKey(key, salt)

	if encrypted {
		file, err = json.Marshal(&configStruct)
		if err != nil {
			return err
		}
		file, err = encrypt(key, file)
		if err != nil {
			return err
		}
		writeFile(configFilePath, file)
	} else {
		file, err = json.Marshal(&configStruct)
		if err != nil {
			return err
		}
		writeFile(configFilePath, file)
	}
	return err
}
