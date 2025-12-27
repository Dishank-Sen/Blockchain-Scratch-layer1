package utils

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"github.com/Dishank-Sen/Blockchain-Scratch-layer1/types"
)


func CreateDir(ctx context.Context, cancel context.CancelFunc, dirPath string, reinit bool) error{
    if reinit{
        if _, err := os.Stat(dirPath); err == nil {
            // directory already exists → skip
            return nil
        } else if !os.IsNotExist(err) {
            // unexpected error
            return err
        }

        // directory does not exist → create it
        return CreateDir(ctx, cancel, dirPath, false)
    }else{
        if err := os.MkdirAll(dirPath, 0700); err != nil{
            return err
        }
    }

    if ctx.Err() != nil {
        cancel()
        return errors.New("operation canceled during directory creation")
    }
    return nil
}

func CreateBootstrap(ctx context.Context, cancel context.CancelFunc, filePath string, reinit bool) error{
    if reinit{
        if info, err := os.Stat(filePath); err == nil{
            // checking if file is empty
            if info.Size() == 0{
                CreateBootstrap(ctx, cancel, filePath, false)
            }
            return nil
        }else if !os.IsNotExist(err) {
            return err
        }

        return CreateBootstrap(ctx, cancel, filePath, false)
    }else{
        u := types.UsersIdentity{
            Peers: []string{},
        }
        data, err := json.Marshal(u)
        if err != nil{
            return err
        }
        return os.WriteFile(filePath, data, 0700)
    }
}