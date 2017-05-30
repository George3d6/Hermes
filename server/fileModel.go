package main

import (
	"os"
	"time"
	"encoding/json"
	"strings"

	"git.cerebralab.com/george/logo"
)

//FileModel used to model a file and control it
type FileModel struct {
	Path        string		`json:"path"`
    Name        string		`json:"name"`
	Compression string		`json:"compression"`
	Birth       time.Time	`json:"birth"`
	TTL         int64		`json:"ttl"`
	Size		int64		`json:"size"`
}


//GetDeathTime gives the unix timestamp when the file is due for deletion
func (file FileModel) GetDeathTime() time.Time {
	var birthTimestamp int64 = file.Birth.Unix()
	var deathTimestamp int64 = birthTimestamp + file.TTL
	return time.Unix(deathTimestamp, 0)
}

//Delete removes the file from the system
func (file * FileModel) Delete() bool {
	return logo.RuntimeError(os.Remove(file.Path))
}

//Update does maintanance work
//Currently maintanance only involves deleting old files
func (file * FileModel) Update() bool {
    if(file.GetDeathTime().Unix() < time.Now().UTC().Unix()) {
        return file.Delete()
    }
    return false
}

//Serialize gives a string (as a byte slice) represntation of a FileModel struct
func (file * FileModel) Serialize() []byte {
    serialization, err := json.Marshal(file)
	logo.RuntimeError(err)
	return serialization
}

//DeserializeFileModel takes a byte slice and create a FileModel
func DeserializeFileModel(serialization []byte) FileModel {
	var newFileModel FileModel
	logo.RuntimeError(json.Unmarshal(serialization, &newFileModel))
	return newFileModel
}




//FileList is meant to hold multiple file models
//FileList is meant to be used by multiple threads and as such all methods are thread safe unless otherwise specified
type FileList struct {
    fileList []FileModel
}

//Serialize gives a string (as a byte slice) represntation of a FileList struct
func (list * FileList) Serialize() []byte {
    var serialization string
	var stringifiedFiles []string = []string{}
	for _, file := range list.fileList {
		stringifiedFiles = append(stringifiedFiles, string(file.Serialize()))
	}

	serialization = strings.Join(stringifiedFiles,"#|#")
	return []byte(serialization)
}

//DeserializeFileList takes a byte slice and create a FileList
func DeserializeFileList(serialization []byte) FileList {
	fileArr := strings.Split(string(serialization), "#|#")
	var newFileList FileList = CreateFileList()
	for _, serializedFileModel := range fileArr {
		newFileList.AddFile(DeserializeFileModel([]byte(serializedFileModel)))
	}
	return newFileList
}

//AddFile adds a file to the list if a file with the same name doesn't exist and return true
//If a file with the same name exists it does nothing and returns false
func (list * FileList) AddFile(file FileModel) bool {
    exists, _ := list.FindFile(file.Name)
    if(exists) {
        return false
    }

    list.fileList = append(list.fileList, file)
    return true
}

//RemoveFile removes a file from the list if it exists and returns true
//If the file doesn't exist it returns false
func (list * FileList) RemoveFile(name string) (bool, FileModel) {
    var newList []FileModel = make([]FileModel, 0, len(list.fileList) - 1)
    var success bool = false
    var removedFile FileModel
    for _, file := range list.fileList {
        if !(file.Name == name) {
            //@TODO: Look into doing this more efficiently using copy or a similar mechanism
            newList = append(newList, file)
        } else {
            success = true
            removedFile = file
        }
    }
    list.fileList = newList
    return success, removedFile
}

//RemoveFile removes a file from the list if it exists and delets it from the fs if it exists and returns true
//If either operation was unsuecesful it returns false and the file list remains unchanged
//If the file doesn't exist it returns false
func (list * FileList) DeleteFile(name string) (bool, FileModel) {
    var success bool
    var fileDeleted FileModel
    var newList []FileModel = make([]FileModel, 0, len(list.fileList) - 1)
    for _, file := range list.fileList {
        if !(file.Name == name) {
            //@TODO: Look into doing this more efficiently using copy or a similar mechanism
            newList = append(newList, file)
        } else {
            success = file.Delete()
            fileDeleted = file
            if(!success) {
                newList = append(newList, file)
            }
        }
    }
    list.fileList = newList
    return success, fileDeleted
}

//FindFile finds a file in the list and returns in togehter with true if the file is found
//Otherwise it return false
func (list * FileList) FindFile(name string) (bool, FileModel) {
    for _, file := range list.fileList {
        if(file.Name == name) {
            return true, file
        }
    }
    return false, FileModel{}
}

//CleanUp is an efficient routine for cleaning up the list of old files and deleting old files
func (list * FileList) CleanUp() {
    var newList []FileModel = make([]FileModel, 0, len(list.fileList))
    for _, file := range list.fileList {
        if(!file.Update()) {
            //@TODO: Look into doing this more efficiently using copy or a similar mechanism
            newList = append(newList, file)
        }
    }
    list.fileList = newList
}

//CreateFileList creates a new file list and returns it
func CreateFileList() FileList {
    return FileList{fileList: make([]FileModel, 0, 233)}
}
