package repo

import (
	"os"

	logging "github.com/ipfs/go-log/v2"
	"github.com/mitchellh/go-homedir"
	"golang.org/x/xerrors"
)

var log = logging.Logger("repo")

const (
	FsRepoEnv  = "VIDS_DATA_REPO"
	fsAPI      = "api"
	fsAPIToken = "token"
	fsLock     = "repo.lock"
)

var ErrRepoExists = xerrors.New("repo exists")

type FsRepo struct {
	path string
}

func NewRepo(path string) (*FsRepo, error) {
	path, err := homedir.Expand(path)
	if err != nil {
		return nil, err
	}

	return &FsRepo{
		path: path,
	}, nil
}

func (fs *FsRepo) Init() error {
	_, err := os.Stat(fs.path)
	if err == nil {
		return ErrRepoExists
	}

	log.Infof("Initializing repo at %s", fs.path)

	err = os.MkdirAll(fs.path, 0755)
	if err != nil && !os.IsExist(err) {
		return err
	}

	return nil
}

func (fs *FsRepo) Path() string {
	return fs.path
}
