package utils

import (
	"AScan/common/utils/gologger"
	"bufio"
	"io"
	"os"
)

// FileExists checks if a file exists and is not a directory
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// FolderExists checks if a folder exists
func FolderExists(folderpath string) bool {
	_, err := os.Stat(folderpath)
	return !os.IsNotExist(err)
}

// HasStdin determines if the user has piped input
func HasStdin() bool {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return false
	}

	mode := stat.Mode()

	isPipedFromChrDev := (mode & os.ModeCharDevice) == 0
	isPipedFromFIFO := (mode & os.ModeNamedPipe) != 0

	return isPipedFromChrDev || isPipedFromFIFO
}

func ReadImf(input io.Reader) []string {
	var imf []string
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		imf = append(imf, scanner.Text())
	}
	return imf
}

// 文件读取，将读取内容按行保存到数组返回
func FileReadByline(fileName string) []string{
	datas := []string{}
	fileData := FileRead(fileName)
	for fileData.Scan(){
		datas = append(datas, fileData.Text())
	}
	return datas
}

//文件读取，返回文件流
func FileRead(fileName string) *bufio.Scanner{
	f, err := os.OpenFile(fileName, os.O_RDONLY, 0)
	if err != nil {
		gologger.Fatalf("file:%v read err!",fileName)
		os.Exit(0)
	}
	datas := bufio.NewScanner(f)
	return datas
}