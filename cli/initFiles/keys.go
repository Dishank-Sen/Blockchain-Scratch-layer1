package initfiles

import (
	"context"
	"path"

	"github.com/Dishank-Sen/Blockchain-Scratch-layer1/cli/utils"
)

func init(){
	InitFile(CreateKeys)
}

func CreateKeys(ctx context.Context, cancel context.CancelFunc, reinit bool) error{
	dirPath := path.Join(".bloc", "identity")
	return utils.CreateKeys(ctx, cancel, dirPath, reinit)
}