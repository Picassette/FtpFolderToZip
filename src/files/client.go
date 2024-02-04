package files

import (
	"fmt"
	"path/filepath"
)

type FilesClient struct {
	Zip *ZipClient
	Tmp *TmpFolderClient
}

/*
*	Create Files Client
*	Contain ZIP and TMP clients
 */
func NewClient(localSaveFolderPath *string) (*FilesClient, error) {
	// Convert in current filesystem format
	*localSaveFolderPath = filepath.FromSlash(*localSaveFolderPath)

	// Init Tmp Client
	tmpClient, err := newTmpClient()
	if err != nil {
		return nil, fmt.Errorf("error init Files client : %s", err.Error())
	}

	// Init Zip Client
	zipClient := newZipClient(localSaveFolderPath, tmpClient.Path)

	client := FilesClient{
		Zip: zipClient,
		Tmp: tmpClient,
	}
	return &client, nil
}
