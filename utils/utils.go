package utils

import (
	"encoding/json"
	"os"
	"path"

	"github.com/Dishank-Sen/Blockchain-Scratch-layer1/types"
)

func CheckDirExist(path string) bool{
	info, err := os.Stat(path)
	if err != nil{
		if os.IsNotExist(err){
			return false
		}
	}
	if info.IsDir(){
		return true
	}
	return false
}

func GetPeerID() (string, error){
	filePath := path.Join(".bloc", "identity", "metadata.json")
	byteData, err := os.ReadFile(filePath)
	if err != nil{
		return "", err
	}

	var m types.Metadata
	if err := json.Unmarshal(byteData, &m); err != nil{
		return "", err
	}

	return m.ID, nil
}