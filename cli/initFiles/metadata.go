package initfiles

import (
	"context"
	"path"
	"github.com/Dishank-Sen/Blockchain-Scratch-layer1/cli/utils"
)

func init(){
	InitFile(CreateMetadata)
}

func CreateMetadata(ctx context.Context, cancel context.CancelFunc, reinit bool) error{
	filepath := path.Join(".bloc", "identity", "metadata.json")
	hashPath := path.Join(".bloc", "identity", "public.key")
	return utils.CreateMetadata(ctx, cancel, filepath, hashPath, reinit)
}