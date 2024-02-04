package main

import (
	"fmt"

	"github.com/Picassette/FtpFolderToZip/common"
	"github.com/Picassette/FtpFolderToZip/files"
	"github.com/Picassette/FtpFolderToZip/ftp"
)

type CoreProcess struct {
	FilesClient *files.FilesClient
	FtpClient   *ftp.FtpClient
}

/*
*	Start Process core
*	Basicaly all start here
 */
func (core *CoreProcess) StartProcess() {
	// Get Environement vars
	common.PrintMsg("Getting env vars", "debug")
	username, err := common.GetEnvVarAsStr("USERNAME", true)
	if err != nil {
		common.PrintMsg(err.Error(), "critical")
		return
	}
	password, err := common.GetEnvVarAsStr("PASSWORD", true)
	if err != nil {
		common.PrintMsg(err.Error(), "critical")
		return
	}
	address, err := common.GetEnvVarAsStr("ADDRESS", true)
	if err != nil {
		common.PrintMsg(err.Error(), "critical")
		return
	}
	port, err := common.GetEnvVarAsNumber("PORT", true)
	if err != nil {
		common.PrintMsg(err.Error(), "critical")
		return
	}
	saveFolderPath, err := common.GetEnvVarAsStr("FTP_SAVE_FOLDER_PATH", true)
	if err != nil {
		common.PrintMsg(err.Error(), "critical")
		return
	}
	localSaveFolderPath, err := common.GetEnvVarAsStrPath("LOCAL_SAVE_FOLDER_PATH", true)
	if err != nil {
		common.PrintMsg(err.Error(), "critical")
		return
	}

	// Files Client
	common.PrintMsg("Init Files Client", "debug")
	filesClient, err := files.NewClient(localSaveFolderPath)
	if err != nil {
		common.PrintMsg(err.Error(), "critical")
		return
	}
	common.PrintMsg(fmt.Sprintf("Used TMP folder : %s", filesClient.Tmp.Path), "debug")

	// Create FTP Client
	common.PrintMsg("Creating FTP Client", "info")
	ftpClient, err := ftp.InitClient(address, port, username, password, saveFolderPath, filesClient.Tmp.Path)
	if err != nil {
		common.PrintMsg(err.Error(), "critical")
		return
	}

	// Process FTP Download
	common.PrintMsg("Starting Download process", "info")
	err = ftpClient.Download()
	if err != nil {
		common.PrintMsg(err.Error(), "critical")
		return
	}

	// Stop Ftp Client Session
	common.PrintMsg("Disconnect FTP", "info")
	ftpClient.DisconnectClient(ftpClient.Session)

	// Process archive
	common.PrintMsg("Create .ZIP", "info")
	err = filesClient.Zip.CreateZip()
	if err != nil {
		common.PrintMsg(err.Error(), "critical")
		return
	}

	// Script finished, clean
	common.PrintMsg("Copy terminated", "info")
}

/*
*	Clean TMP folder and disconnect from FTP
*	Does only necessary actions
 */
func cleanStop(core *CoreProcess) {
	if core.FilesClient != nil && core.FilesClient.Tmp != nil {
		common.PrintMsg("Remove TMP Folder", "debug")
		core.FilesClient.Tmp.CleanTmpFolder()
	}
	if core.FtpClient != nil && core.FtpClient.Session != nil {
		common.PrintMsg("Disconnect FTP Client", "info")
		core.FtpClient.DisconnectClient(core.FtpClient.Session)
	}
}

// Entry point
func main() {
	common.PrintMsg("Starting PalWorld Save Script", "info")
	// Create a new Core Process
	core := CoreProcess{}
	core.StartProcess()

	// Clean TMP and stop FTP (if present/connected)
	common.PrintMsg("Clean and Stop", "debug")
	cleanStop(&core)
}
