package files

import (
	"fmt"
	"os"
)

type TmpFolderClient struct {
	Path string
}

/*
*	Create a TMP Folder in your system
 */
func generateTmpFolder() (*string, error) {
	tmpDir, err := os.MkdirTemp("", "FtpSaveCpy")
	if err != nil {
		return nil, fmt.Errorf("error creating temporary directory: %s", err.Error())
	}
	return &tmpDir, nil
}

/*
*	Create new Tmp Client
 */
func newTmpClient() (*TmpFolderClient, error) {
	// Generate Tmp Folder
	path, err := generateTmpFolder()
	if err != nil {
		return nil, fmt.Errorf("error init TmpClient : %s", err.Error())
	}
	client := TmpFolderClient{
		Path: *path,
	}
	return &client, nil
}

/*
*	Delete Tmp Folder
 */
func (client *TmpFolderClient) CleanTmpFolder() error {
	err := os.RemoveAll(client.Path)
	if err != nil {
		return fmt.Errorf("error removing tmp folder : %s", err.Error())
	}
	return nil
}
