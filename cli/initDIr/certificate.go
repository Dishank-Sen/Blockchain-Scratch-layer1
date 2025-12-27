package initdir

import (
	"context"
	"path"

	"github.com/Dishank-Sen/Blockchain-Scratch-layer1/cli/utils"
)

func init(){
	InitDir(CreateCertificate)
}

func CreateCertificate(ctx context.Context, cancel context.CancelFunc, reinit bool) error{
	dirPath := path.Join(".bloc", "certificate")
	return utils.CreateDir(ctx, cancel, dirPath, reinit)
}