package main

import (
	"encoding/json"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

//FileModel used to model a file and control it
type FileModel struct {
	Path        string    `json:"path"`
	Name        string    `json:"name"`
	Compression string    `json:"compression"`
	Birth       time.Time `json:"birth"`
	TTL         int64     `json:"ttl"`
	Size        int64     `json:"size"`
}

//GetDeathTime gives the unix timestamp when the file is due for deletion
func (file FileModel) GetDeathTime() time.Time {
	var birthTimestamp int64 = file.Birth.Unix()
	var deathTimestamp int64 = birthTimestamp + file.TTL
	return time.Unix(deathTimestamp, 0)
}

//Delete removes the file from the system
func (file *FileModel) Delete() bool {
	err := os.Remove(file.Path)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

//Update does maintanance work
//Currently maintanance only involves deleting old files
func (file *FileModel) Update() bool {
	if file.GetDeathTime().Unix() < time.Now().UTC().Unix() {
		return file.Delete()
	}
	return false
}

//Serialize gives a string (as a byte slice) represntation of a FileModel struct
func (file *FileModel) Serialize() []byte {
	serialization, err := json.Marshal(file)
	if err != nil {
		log.Println(err)
	}
	return serialization
}

//DeserializeFileModel takes a byte slice and create a FileModel
func DeserializeFileModel(serialization []byte) FileModel {
	var newFileModel FileModel
	if err := json.Unmarshal(serialization, &newFileModel); err != nil {
		log.Println(err)
	}
	return newFileModel
}

//FileList is meant to hold multiple file models
//FileList is meant to be used by multiple threads and as such all methods are thread safe unless otherwise specified
type FileList struct {
	fileList []FileModel
	mutex    sync.RWMutex
}

//RunOnFileList runs a function on each file in the internal file list under a mutex
func (list *FileList) RunOnFileList(task func(fileList []FileModel) interface{}) interface{} {
	list.mutex.Lock()
	result := task(list.fileList)
	list.mutex.Unlock()
	return result
}

//ReadOnFileList runs a read-only function on each file in the internal file list under a mutex
func (list *FileList) ReadOnFileList(task func(fileList []FileModel) interface{}) interface{} {
	list.mutex.RLock()
	result := task(list.fileList)
	list.mutex.RUnlock()
	return result
}

//Serialize gives a string (as a byte slice) represntation of a FileList struct
func (list *FileList) Serialize() []byte {
	list.mutex.RLock()
	defer list.mutex.RUnlock()
	var serialization string
	var stringifiedFiles []string = []string{}
	for _, file := range list.fileList {
		stringifiedFiles = append(stringifiedFiles, string(file.Serialize()))
	}

	serialization = strings.Join(stringifiedFiles, "#|#")
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
func (list *FileList) AddFile(file FileModel) bool {
	list.mutex.Lock()
	defer list.mutex.Unlock()
	exists, _ := list.FindFile(file.Name, false)
	if exists {
		return false
	}

	list.fileList = append(list.fileList, file)
	return true
}

//RemoveFile removes a file from the list if it exists and returns true
//If the file doesn't exist it returns false
func (list *FileList) RemoveFile(name string) (bool, FileModel) {
	list.mutex.Lock()
	defer list.mutex.Unlock()
	var newList []FileModel = make([]FileModel, 0, len(list.fileList)-1)
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

//DeleteFile removes a file from the list if it exists and delets it from the fs if it exists and returns true
//If either operation was unsuecesful it returns false and the file list remains unchanged
//If the file doesn't exist it returns false
func (list *FileList) DeleteFile(name string) (bool, FileModel) {
	list.mutex.Lock()
	defer list.mutex.Unlock()
	var success bool
	var fileDeleted FileModel
	var newList []FileModel = make([]FileModel, 0, len(list.fileList)-1)
	for _, file := range list.fileList {
		if file.Name != name {
			//@TODO: Look into doing this more efficiently using copy or a similar mechanism
			newList = append(newList, file)
		} else {
			success = !file.Delete()
			fileDeleted = file
			if !success {
				newList = append(newList, file)
			}
		}
	}
	list.fileList = newList
	return success, fileDeleted
}

//FindFile finds a file in the list and returns in togehter with true if the file is found
//Otherwise it return false
func (list *FileList) FindFile(name string, lock ...bool) (bool, FileModel) {
	if len(lock) > 0 {
		if lock[0] {
			list.mutex.RLock()
			defer list.mutex.RUnlock()
		}
	} else {
		list.mutex.RLock()
		defer list.mutex.RUnlock()
	}
	for _, file := range list.fileList {
		if file.Name == name {
			return true, file
		}
	}
	return false, FileModel{}
}

//CleanUp is an efficient routine for cleaning up the list of old files and deleting old files
func (list *FileList) CleanUp() {
	//Could be a rlock and than a lock but this way is probably faster in reality
	list.mutex.Lock()
	defer list.mutex.Unlock()
	var newList []FileModel = make([]FileModel, 0, len(list.fileList))
	for _, file := range list.fileList {
		if !file.Update() {
			//@TODO: Look into doing this more efficiently using copy or a similar mechanism
			newList = append(newList, file)
		}
	}
	list.fileList = newList
	return
}

//CreateFileList creates a new file list and returns it
func CreateFileList() FileList {
	return FileList{fileList: make([]FileModel, 0, 233), mutex: sync.RWMutex{}}
}
