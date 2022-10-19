package main

import (
	"encoding/json"
	"github.com/edward/file-collector/logger"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	logger.Init("file_collector.log")
	jobs, err := ReadConfig("config.json")
	if err != nil {
		logger.Log(err.Error())
		return
	}
	for _, job := range jobs {
		logger.Log("Starting a job...")
		jobBytes, err := json.Marshal(job)
		if err != nil {
			logger.Log(err.Error())
			return
		}
		logger.Log("Job content: " + string(jobBytes))

		copyObjs, err := CollectCopyObjs(job.SrcRootPath, job.DesRootPath, job.DirPattern, job.ExcludedDirPattern, job.FileNames)
		if err != nil {
			logger.Log(err.Error())
			return
		}
		for _, copyObj := range copyObjs {
			_, err := CopyFile(copyObj.SrcPath, copyObj.DesPath)
			if err != nil {
				logger.Log(err.Error())
				return
			}
		}
		logger.Log("Job is executed")
	}
}

type Job struct {
	SrcRootPath        string `json:"SrcRootPath"`
	DesRootPath        string `json:"DesRootPath"`
	FileNames          string `json:"FileNames"`
	DirPattern         string `json:"DirPattern"`
	ExcludedDirPattern string `json:"ExcludedDirPattern"`
}

func ReadConfig(configPath string) (jobs []Job, err error) {
	configFile, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer configFile.Close()
	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&jobs)
	if err != nil {
		return nil, err
	}
	return jobs, nil
}

func CollectCopyObjs(srcRootPath string, desRootPath string, dirPattern string, excludedDirPattern string, fileNames string) ([]FileCopyObj, error) {
	var fileNameArr []string
	if fileNames != "" {
		fileNameArr = strings.Split(fileNames, ",")
		for key := range fileNameArr {
			fileNameArr[key] = strings.ToLower(strings.TrimSpace(fileNameArr[key]))
		}
	}

	absFilePaths, relFilePaths, err := GetFilePaths(srcRootPath, "", dirPattern, excludedDirPattern, fileNameArr)
	if err != nil {
		return nil, err
	}

	fileCopyObjs := make([]FileCopyObj, 0, len(absFilePaths))
	for key, absFilePath := range absFilePaths {
		exists := false
		for _, copyObj := range fileCopyObjs {
			if absFilePath == copyObj.SrcPath {
				exists = true
				break
			}
		}
		if exists {
			continue
		}

		fileCopyObj := FileCopyObj{
			SrcPath: absFilePath,
			DesPath: filepath.Join(desRootPath, relFilePaths[key]),
		}
		fileCopyObjs = append(fileCopyObjs, fileCopyObj)
	}

	return fileCopyObjs, nil
}
