package crypto

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

//Zip 压缩文件/文件夹
func Zip(srcFileOrDir, destZipFile string) error {
	zipFile, err := os.Create(destZipFile)
	if err != nil {
		return err
	}
	defer func() {
		_ = zipFile.Close()
	}()

	archive := zip.NewWriter(zipFile)
	defer func() {
		_ = archive.Close()
	}()

	err = filepath.Walk(srcFileOrDir, func(path string, info os.FileInfo, err error) error {
		var header *zip.FileHeader
		var writer io.Writer
		var file *os.File

		if err != nil {
			return err
		}

		header, err = zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		header.Name = path
		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err = archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return err
		}

		file, err = os.Open(path)
		if err != nil {
			return err
		}

		defer func() {
			_ = file.Close()
		}()

		_, err = io.Copy(writer, file)

		return err
	})

	return err
}

//Unzip 解压缩文件到文件夹
func Unzip(zipFile, destDir string) error {
	zipReader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer func() {
		_ = zipReader.Close()
	}()

	for _, file := range zipReader.File {
		filePath := filepath.Join(destDir, file.FileInfo().Name())
		if file.FileInfo().IsDir() {
			continue
		}
		var inFile io.ReadCloser
		var outFile *os.File

		if err = os.MkdirAll(filepath.Dir(filePath), 0644); err != nil {
			return err
		}
		inFile, err = file.Open()
		if err != nil {
			return err
		}

		outFile, err = os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			_ = inFile.Close()
			return err
		}

		_, err = io.Copy(outFile, inFile)

		_ = inFile.Close()
		_ = outFile.Close()

		if err != nil {
			return err
		}
	}
	return err
}
