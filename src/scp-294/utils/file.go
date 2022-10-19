package utils

import (
	"os"
	"path/filepath"
	"strings"
)

type FileCopyObj struct {
	SrcPath string
	DesPath string
}

func GetFilePaths(absDirPath string, relDirPath string, dirPattern string, excludedDirPattern string, fileNameArr []string) (absFilePaths []string, relFilePaths []string, err error) {
	folder, err := os.Open(absDirPath)
	defer folder.Close()
	if err != nil {
		return nil, nil, err
	}
	files, err := folder.Readdir(-1)
	if err != nil {
		return nil, nil, err
	}

	absFilePaths = make([]string, 0, 100)
	relFilePaths = make([]string, 0, 100)
	for _, file := range files {
		absFilePath := filepath.Join(absDirPath, file.Name())
		relFilePath := filepath.Join(relDirPath, file.Name())
		if file.IsDir() {
			if excludedDirPattern != "" && match(excludedDirPattern, file.Name()) {
				continue
			}
			absPaths, relPaths, err := GetFilePaths(absFilePath, relFilePath, dirPattern, excludedDirPattern, fileNameArr)
			if err != nil {
				return nil, nil, err
			}
			absFilePaths = append(absFilePaths, absPaths...)
			relFilePaths = append(relFilePaths, relPaths...)
		} else {
			if dirPattern == "" || match(dirPattern, filepath.Base(folder.Name())) {
				for _, fileName := range fileNameArr {
					if fileName == strings.ToLower(file.Name()) {
						absFilePaths = append(absFilePaths, absFilePath)
						relFilePaths = append(relFilePaths, relFilePath)
					}
				}
			}
		}
	}

	return absFilePaths, relFilePaths, nil
}

func match(pattern string, value string) bool {
	noPattern, startWith, endWith, contains := 0, 1, 2, 3
	patternType := noPattern
	indexByte := strings.IndexByte(pattern, '%')
	lastIndexByte := strings.LastIndexByte(pattern, '%')
	if indexByte == lastIndexByte {
		if indexByte == 0 {
			patternType = endWith
		} else if indexByte == len(pattern)-1 {
			patternType = startWith
		}
	} else {
		if indexByte == 0 && lastIndexByte == len(pattern)-1 {
			patternType = contains
		} else if indexByte == 0 {
			patternType = startWith
		} else if lastIndexByte == len(pattern)-1 {
			patternType = endWith
		}
	}
	lowerValue := strings.ToLower(value)
	lowerPattern := strings.ToLower(pattern)
	switch patternType {
	case noPattern:
		return lowerValue == lowerPattern
	case startWith:
		lowerPattern = lowerPattern[0 : len(lowerPattern)-1]
		return strings.HasPrefix(lowerValue, lowerPattern)
	case endWith:
		lowerPattern = lowerPattern[1:]
		return strings.HasSuffix(lowerValue, lowerPattern)
	case contains:
		lowerPattern = lowerPattern[1 : len(lowerPattern)-1]
		return strings.Contains(lowerValue, lowerPattern)
	default:
		return lowerValue == lowerPattern
	}
}
