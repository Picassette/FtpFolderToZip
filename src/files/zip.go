package files

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
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
	err = filepath.Walk(client.TmpFolder, client.walker)
	if err != nil {
		panic(err)
	}

	return nil
}

/*
*	Convert Tmp Folder path from Absolute to Relative
 */
func (client *ZipClient) absoluteToRelativePath(absolute string) (string, error) {
	index := strings.Index(absolute, client.TmpFolder)
	if index == -1 {
		return absolute, nil
	}
	ret := absolute[index+client.TmpFolderPathLen:]
	if len(ret) == 0 {
		// Unless, we have an empty string
		ret = filepath.FromSlash("/")
	}
	return ret, nil
}

/*
*	File to Zip processing
 */
func (client *ZipClient) walker(path string, info os.FileInfo, err error) error {
	if err != nil {
		return fmt.Errorf("error in walker : %s", err.Error())
	}
	relativePath, err := client.absoluteToRelativePath(path)
	if err != nil {
		return fmt.Errorf("error in walker : %s", err.Error())
	}
	if info.IsDir() {
		return nil
	}

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
