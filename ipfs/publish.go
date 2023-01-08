package ipfs

import (
	"context"
	"strings"

	"fmt"
	"os"

	"path"

	"telegram-ipfsacer/video"

	"github.com/thedmdim/go-ipfs-api"
)


type Client struct {
	Sh *shell.Shell
	Storage string
	KeyName string
}

// init new IPFS shell and create MFS dir to store video
func NewIPFSClient(url string, storage string, keyPath string) (*Client, error) {
	
	c := new(Client)
	c.Sh = shell.NewShell(url)
	c.Storage = "/" + storage

	if keyPath != "" {
		keyFile, err := os.Open(keyPath)
		if err != nil{
			return nil, fmt.Errorf("cannot find %s: %w", keyPath, err)
		}
		keyName := path.Base(keyPath)
		err = c.Sh.KeyImport(context.Background(), keyName, keyFile)
		if err != nil && !strings.Contains(err.Error(), "already exists") {
			return nil, fmt.Errorf("cannot import key %s: %w", keyPath, err)
		}
		c.KeyName = keyName
		defer keyFile.Close()
	}
	
	return c, nil
}

func (c *Client) PublishCurrent(ctx context.Context) (*shell.PublishResponse, error) {
	file, err := c.Sh.FilesStat(ctx, c.Storage)
	if err != nil {
        return nil, fmt.Errorf("cannot get cid of %s: %s", c.Storage, err)
	}
	r, err := c.Sh.PublishWithDetails(file.Hash, c.KeyName, 0, 0, false)
	if err != nil {
		return nil, err
	}
	return r, nil
	

}

// create file /<Client.Storage>/<Video.Filename> and return files hash
func (c *Client) AddVideo(ctx context.Context, v *video.Video) (string, error) {
	
	filename := path.Join(c.Storage, v.Filename)


	err := c.Sh.FilesWrite(ctx, filename, *v.Stream, shell.FilesWrite.Create(true), shell.FilesWrite.Parents(true))
	if err != nil {
		fmt.Println(filename)
        return "", fmt.Errorf("cannot add: %s", err)
	}

	added, err := c.Sh.FilesStat(ctx, filename)
	if err != nil {
        return "", fmt.Errorf("cannot get cid of %s: %s", filename, err)
	}

	return added.Hash, nil

}
