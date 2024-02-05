package ftp

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Picassette/FtpFolderToZip/common"
	"github.com/jlaffaye/ftp"
)

type FtpClient struct {
	Session           *ftp.ServerConn
	SaveFolderPath    string
	SaveFolderPathLen int
	TmpFolder         string
}

/*
*	Download datas from FTP
 */
func (client *FtpClient) Download() error {
	// Get all files inside current folder
	entries, err := client.Session.List(client.SaveFolderPath)
	if err != nil {
		return fmt.Errorf("error listing entries in FTP")
	}

	// Process entry (folder/file)
	err = client.routeByEntryType(entries, client.SaveFolderPath)
	if err != nil {
		return fmt.Errorf("error processing FTP requests : %s", err.Error())
	}

	common.PrintMsg("download complet", "info")
	return nil
}

/*
*	Disconnect from current FTP Session
 */
func (client *FtpClient) DisconnectClient(ftpClient *ftp.ServerConn) error {
	if ftpClient == nil {
		return nil
	}

	// Stop FTP session
	err := ftpClient.Quit()
	if err != nil {
		return fmt.Errorf("error during logout : %s", err.Error())
	}
	common.PrintMsg("session disconnected", "info")
	return nil
}

/*
*	Create new Ftp Client
*	Have a TimeOut of 10s
 */
func newFtpClient(address *string, port *int64, username *string, password *string) (*ftp.ServerConn, error) {
	common.PrintMsg("creating new session", "info")
	fullAddress := fmt.Sprintf("%s:%d", *address, *port)

	// Connect to FTP with timeout
	ftpClient, err := ftp.Dial(fullAddress, ftp.DialWithTimeout(10*time.Second))
	if err != nil {
		return nil, fmt.Errorf("error during ftp client dial : %s", err.Error())
	}

	// Login to FTP
	common.PrintMsg("trying to login", "info")
	err = ftpClient.Login(*username, *password)
	if err != nil {
		return nil, fmt.Errorf("error during ftp login : %s", err.Error())
	}
	common.PrintMsg("successfuly logged", "info")
	return ftpClient, nil
}

/*
*	Download File from FTP
 */
func (client *FtpClient) downloadFile(fileName string) error {
	// Create path from ftp path
	ftpFilePath := client.savePathToFtpPath(fileName)

	// Create absolute path for local file
	fullLocalPath := filepath.FromSlash(fmt.Sprintf("%s/%s", client.TmpFolder, ftpFilePath))
	common.PrintMsg(fmt.Sprintf("trying to download file : %s to %s", fileName, fullLocalPath), "info")

	// Get file from FTP
	remoteFile, err := client.Session.Retr(fileName)
	if err != nil {
		return fmt.Errorf("error retreving FTP file : %s", err.Error())
	}
	defer remoteFile.Close()

	// Create local file
	localFile, err := os.Create(fullLocalPath)
	if err != nil {
		return fmt.Errorf("error creating local file : %s", err.Error())
	}
	defer localFile.Close()

	// Copy file datas from FTP to Local file
	_, err = io.Copy(localFile, remoteFile)
	if err != nil {
		return fmt.Errorf("error coping local file : %s", err.Error())
	}

	common.PrintMsg(fmt.Sprintf("downloaded file : %s to %s", fileName, fullLocalPath), "info")
	return nil
}

/*
*	Download folder form FTP
 */
func (client *FtpClient) downloadFolder(folderName string) error {
	// Create path from ftp path
	cleanForLocalFolderName := client.savePathToFtpPath(folderName)

	// Create path for local folder
	localFolderName := filepath.FromSlash(fmt.Sprintf("%s/%s", client.TmpFolder, cleanForLocalFolderName))
	common.PrintMsg(fmt.Sprintf("trying to create folder : %s", localFolderName), "info")

	// Create local folder
	err := os.MkdirAll(localFolderName, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error creating local folder : %s", err.Error())
	}
	common.PrintMsg(fmt.Sprintf("folder created : %s", localFolderName), "info")

	// Change working directory on the FTP server
	common.PrintMsg(fmt.Sprintf("moving to : %s", folderName), "info")
	err = client.Session.ChangeDir(folderName)
	if err != nil {
		return fmt.Errorf("error moving in FTP : %s", err.Error())
	}

	// List files and folders in the current directory
	entries, err := client.Session.List(".")
	if err != nil {
		return fmt.Errorf("error getting current directoy file list : %s", err.Error())
	}
	err = client.routeByEntryType(entries, folderName)
	if err != nil {
		return fmt.Errorf("error processing FTP files : %s", err.Error())
	}

	// Move back to the parent directory on the FTP server
	err = client.Session.ChangeDirToParent()
	if err != nil {
		return fmt.Errorf("error moving back to parent directory")
	}

	common.PrintMsg(fmt.Sprintf("downloaded folder : %s", folderName), "info")
	return nil
}

/*
*	Route entries by type (file/folder)
 */
func (client *FtpClient) routeByEntryType(entries []*ftp.Entry, initPath string) error {

	// For each entry
	for _, entry := range entries {
		// We don't want to reparse current dir or previous dir
		if entry.Name == "." || entry.Name == ".." {
			continue
		}

		// Download file or folder
		fullEntryPath := fmt.Sprintf("%s/%s", initPath, entry.Name)
		if entry.Type == ftp.EntryTypeFile {
			err := client.downloadFile(fullEntryPath)
			if err != nil {
				return fmt.Errorf("error in download file process : %s", err.Error())
			}
		} else if entry.Type == ftp.EntryTypeFolder {
			err := client.downloadFolder(fullEntryPath)
			if err != nil {
				return fmt.Errorf("error in download folder process : %s", err.Error())
			}
		}
	}
	return nil
}

/*
*	Convert Path to FTP Path
 */
func (client *FtpClient) savePathToFtpPath(input string) string {
	index := strings.Index(input, client.SaveFolderPath)
	if index == -1 {
		return input
	}
	return input[index+client.SaveFolderPathLen:]
}

/*
*	Init FTP client
 */
func InitClient(address *string, port *int64, username *string, password *string, saveFolderPath *string, tmpFolder string) (*FtpClient, error) {
	session, err := newFtpClient(address, port, username, password)
	if err != nil {
		return nil, err
	}
	client := FtpClient{
		Session:           session,
		SaveFolderPath:    *saveFolderPath,
		SaveFolderPathLen: len(*saveFolderPath),
		TmpFolder:         tmpFolder,
	}

	return &client, nil
}
