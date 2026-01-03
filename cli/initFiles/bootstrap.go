package initfiles

import (
	"context"
	"path"

	"github.com/Dishank-Sen/Blockchain-Scratch-layer1/cli/utils"
)

func init(){
	InitFile(CreateBootstrap)
}

func CreateBootstrap(ctx context.Context, cancel context.CancelFunc, reinit bool) error{
	filePath := path.Join(".bloc", "bootstrap.json")
	return utils.CreateBootstrap(ctx, cancel, filePath, reinit)
}