package files

import (
	"archive/zip"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Picassette/FtpFolderToZip/common"
)

type ZipClient struct {
	SaveFolderName   string
	TmpFolder        string
	TmpFolderPathLen int
	ZipName          *string
	ZipWriter        *zip.Writer
}

/*
*	Create Zip entry point
 */
func (client *ZipClient) CreateZip() error {
	// Generate path compatible with current filesystem
	path := filepath.FromSlash(fmt.Sprintf("%s/%s", client.SaveFolderName, generateZipName()))
	client.ZipName = &path

	// Create Zip File
	zipFile, err := os.Create(*client.ZipName)
	if err != nil {
		return fmt.Errorf("error creating archive : %s", err.Error())
	}
	defer zipFile.Close()

	// Create Zip Writer
	client.ZipWriter = zip.NewWriter(zipFile)
	defer client.ZipWriter.Close()

	// Process Tmp files into Zip File
	err = filepath.WalkDir(client.TmpFolder, client.walker)
	if err != nil {
		panic(err)
	}

	return nil
}

/*
*	Convert Tmp Folder path to path inside zip
 */
func (client *ZipClient) pathToInsideZipPath(absolute string) (string, error) {
	// Get temp folder index in absolute path
	index := strings.Index(absolute, client.TmpFolder)
	if index == -1 {
		return absolute, nil
	}

	// Remove Tmp folder path and convert to the right slash format
	ret := filepath.ToSlash(absolute[index+client.TmpFolderPathLen:])

	// If first char is '/', removing it
	if ret[0] == '/' {
		ret = ret[1:]
	}

	return ret, nil
}

/*
*	File to Zip processing
 */
func (client *ZipClient) walker(path string, dir fs.DirEntry, err error) error {
	if err != nil {
		return fmt.Errorf("error in walker : %s", err.Error())
	}
	if dir.IsDir() {
		return nil
	}

	relativePath, err := client.pathToInsideZipPath(path)
	if err != nil {
		return fmt.Errorf("error in walker : %s", err.Error())
	}
	relativePath = filepath.ToSlash(relativePath)

	common.PrintMsg(fmt.Sprintf("add %s file to %s path inside zip", path, relativePath), "debug")

	// We copy file content
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("error oppening : %s", err.Error())
	}
	defer file.Close()

	// We create file
	writer, err := client.ZipWriter.Create(relativePath)
	if err != nil {
		return err
	}

	// We copy file datas
	_, err = io.Copy(writer, file)
	if err != nil {
		return err
	}

	return nil
}

/*
*	Generate Zip Name
*	File name contain full date
 */
func generateZipName() string {
	timeNow := strings.ReplaceAll(time.Now().Format("2006-01-02T15:04:05.00000"), ":", "-")
	timeNow = strings.ReplaceAll(timeNow, ".", "-")
	return fmt.Sprintf("save_%s.zip", timeNow)
}

/*
*	Create Zip client
 */
func newZipClient(folderName *string, tmpFolder string) *ZipClient {
	zipClient := ZipClient{
		SaveFolderName:   *folderName,
		TmpFolder:        tmpFolder,
		TmpFolderPathLen: len(tmpFolder),
		ZipWriter:        nil,
		ZipName:          nil,
	}
	return &zipClient
}
