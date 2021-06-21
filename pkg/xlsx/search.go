package xlsx

import (
	"archive/zip"
	"errors"
	"io"
	"io/ioutil"
	"regexp"
)

func Contains(path, term string) (bool, int, error) {
	data, err := ReadDocxFile(path)
	if err != nil {
		return false, 0, err
	}

	var (
		validID = regexp.MustCompile(`(?i)` + term)
		matches = validID.FindAllString(data, -1)
	)

	if occurrences := len(matches); occurrences == 0 {
		return false, occurrences, nil
	} else {
		return true, occurrences, nil
	}
}

func ReadDocxFile(path string) (string, error) {
	reader, err := zip.OpenReader(path)
	if err != nil {
		return "", err
	}

	zipData := ZipFile{data: reader}

	return ReadDocx(zipData)
}

func ReadDocx(reader ZipData) (string, error) {
	content, err := readText(reader.files())
	if err != nil {
		return "", err
	}

	return content, nil
}

type ZipFile struct {
	data *zip.ReadCloser
}

func (d ZipFile) files() []*zip.File {
	return d.data.File
}

func (d ZipFile) close() error {
	return d.data.Close()
}

type ZipData interface {
	files() []*zip.File
	close() error
}

func readText(files []*zip.File) (text string, err error) {
	var documentFile *zip.File

	documentFile, err = retrieveWordDoc(files)
	if err != nil {
		return text, err
	}

	var documentReader io.ReadCloser

	documentReader, err = documentFile.Open()
	if err != nil {
		return text, err
	}

	text, err = wordDocToString(documentReader)

	return text, err
}

func retrieveWordDoc(files []*zip.File) (file *zip.File, err error) {
	for _, f := range files {
		if f.Name == "xl/sharedStrings.xml" {
			file = f
		}
	}

	if file == nil {
		err = errors.New("sharedStrings.xml file not found")
	}

	return file, err
}

func wordDocToString(reader io.Reader) (string, error) {
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
