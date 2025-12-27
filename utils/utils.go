package utils

import "os"

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