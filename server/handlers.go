package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"git.cerebralab.com/george/logo"

	"github.com/ulikunitz/xz"
)

var globalFileList FileList

func initHandlers(adminName string, adminPassword string) {
	fmt.Println(Configuration.StatePath + "file_list.json")
	flcontent, err := ioutil.ReadFile(Configuration.StatePath + "file_list.json")
	logo.RuntimeFatal(err)
	tmcontent, err := ioutil.ReadFile(Configuration.StatePath + "token_map.json")
	logo.RuntimeFatal(err)
	globalFileList = DeserializeFileList(flcontent)
	DeserializeTokenMap(tmcontent)
	InitializeAdmin([]byte("the salt"), adminName, adminPassword)
	var oldFlSerialization string
	var oldTmSerialization string
	//Loop to preserve state
	go func() {
		for {
			//sleep
			time.Sleep(1 * time.Second)
			newFlSerialization := string(globalFileList.Serialize())
			newTmSerialization := string(SerializeTokenMap())

			if (oldTmSerialization != newTmSerialization) {
				oldTmSerialization = newTmSerialization

				//save token map
				tmSerializationFile, err := os.Create(Configuration.StatePath + "token_map.json")
				logo.RuntimeError(err)
				defer tmSerializationFile.Close()
				_, err = io.Copy(tmSerializationFile, strings.NewReader(oldTmSerialization))
				logo.RuntimeError(err)
			}

			if (oldFlSerialization != newFlSerialization) {
				oldFlSerialization = newFlSerialization

				//save file list
				fileListSerializationFile, err := os.Create(Configuration.StatePath + "file_list.json")
				logo.RuntimeError(err)
				defer fileListSerializationFile.Close()
				_, err = io.Copy(fileListSerializationFile, strings.NewReader(oldFlSerialization))
				logo.RuntimeError(err)
			}
		}
	}()
}

//Server the index.html file
func serveHome(w http.ResponseWriter, r *http.Request) {
	index, err := ioutil.ReadFile("./client/index.html")
	logo.RuntimeError(err)

	if err != nil {
		logo.RuntimeError(err)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	fmt.Fprintf(w, string(index))
}

//Authentication
func engageAuthSession(w http.ResponseWriter, r *http.Request) {
	identifier := r.URL.Query().Get("identifier")
	credentials := r.URL.Query().Get("credentials")

	isValid, sessionId := ValidateToke(identifier, credentials)
	if isValid {
		expiration := time.Now().Add(24 * time.Hour)
		authCookie := http.Cookie{Path: "/", Name: "auth", Value: string(identifier + "#|#" + sessionId), Expires: expiration, MaxAge: 3600 * 24}
		http.SetCookie(w, &authCookie)
		//http.Redirect(w, r, "/", http.StatusFound)
		return
	} else {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
	}
}

//parses multipart form and saves an image called "image" to a file
func uploadFile(w http.ResponseWriter, r *http.Request) {
	//Parsing the upload arguments into the values we shall be working with
	r.ParseMultipartForm(5000000000000000)

	file, _, err := r.FormFile("file")
	if logo.RuntimeError(err) {
		fmt.Fprintf(w, `Could not upload file, file can't be read\n`)
		return
	}
	//defer file.Close()
	name := r.FormValue("name")
	compression := r.FormValue("compression")
	public := r.FormValue("public")

	ttlString := r.FormValue("ttl")
	ttl, err := strconv.Atoi(ttlString)
	if logo.RuntimeError(err) {
		fmt.Fprintf(w, `Could not upload file, can't parse time to live(ttl)\n`)
		return
	}

	extension := ""
	if compression == "gz" {
		extension = ".gz"
	}
	if compression == "xz" {
		extension = ".xz"
	}

	newFileModel := FileModel{Path: Configuration.FilePath + name + extension, Name: name, TTL: int64(ttl), Birth: time.Now(),
		Compression: compression, Size: GetFileSizeInBytes(file)}

	//Doing the authentication
	cookie, err := r.Cookie("auth")
	if err != nil {
		fmt.Fprintf(w, `Not authenticated\n`)
		return
	}
	values := strings.Split(cookie.Value, "#|#")
	if len(values) < 2 {
		fmt.Fprintf(w, `Authentication cookie malformed\n`)
		return
	}
	valid, token := ValidateSession(values[0], values[1])
	if !valid {
		fmt.Fprintf(w, `Session Id invalid\n`)
		return
	}
	/* @TODO fix this, currently the comparison seems not to work
	logo.LogDebug(token.UploadSize < newFileModel.Size)
	if token.UploadSize < newFileModel.Size {
		fmt.Fprintf(w, `Trying to upload a file that is too big, your maximum upload size is: ` + strconv.Itoa(int(newFileModel.Size)))
		return
	}
	*/

	if token.UploadNumber < 1 {
		fmt.Fprint(w, "Reached maximum upload limit for this token/user\n")
		return
	}

	token.UploadNumber = token.UploadNumber - 1
	token.OwnedFiles = append(token.OwnedFiles, newFileModel.Name)
	if !ModifyToken(token) {
		fmt.Fprint(w, "Internal server error, please try again and inform the administrator about this\n")
		return
	}

	if public == "true" {
		UpdatePublicToken(newFileModel.Name)
	}

	//Uploading
	logo.LogDebug(newFileModel.Path)
	permanentFile, err := os.OpenFile(newFileModel.Path, os.O_WRONLY|os.O_CREATE, 0666)

	if logo.RuntimeError(err) {
		fmt.Fprintf(w, `Could not upload file, there was an internal file system error, try again in a few seconds\n`)
		return
	}
	defer permanentFile.Close()

	switch compression := newFileModel.Compression; compression {
	case "gz":
		gzipped, err := gzip.NewWriterLevel(permanentFile, 6)
		if logo.RuntimeError(err) {
			fmt.Fprintf(w, `Could not upload file, there was an internal file system error, try again in a few seconds\n`)
			return
		}
		defer gzipped.Close()
		io.Copy(gzipped, file)
	case "xz":
		xzw, err := xz.NewWriter(permanentFile)
		if logo.RuntimeError(err) {
			fmt.Fprintf(w, `Could not upload file, there was an internal file system error, try again in a few seconds\n`)
			return
		}
		defer xzw.Close()
		io.Copy(xzw, file)
	case "plain":
		io.Copy(permanentFile, file)
	default:
		fmt.Fprintf(w, "Wrong compression format\n")
	}
	if !globalFileList.AddFile(newFileModel) {
		fmt.Fprintf(w, `Could not upload file, a similarly named file already exists\n`)
		return
	}
	fmt.Fprintf(w, "Successfully uploaded file\n")
	return
}

//List files
func listFiles(w http.ResponseWriter, r *http.Request) {
	//Doing the authentication
	cookie, err := r.Cookie("auth")
	if err != nil {
		fmt.Fprintf(w, `Not authenticated\n`)
		return
	}
	values := strings.Split(cookie.Value, "#|#")
	if len(values) < 2 {
		fmt.Fprintf(w, `Authentication cookie malformed\n`)
		return
	}
	valid, token := ValidateSession(values[0], values[1])
	if !valid {
		fmt.Fprintf(w, `Session Id invalid\n`)
		return
	}
	stringified := ""
	for _, val := range token.OwnedFiles {
		found, file := globalFileList.FindFile(val)
		if found {
			stringified += file.Name + "|#|" + strconv.Itoa(int(file.Size)) + "|#|" + strconv.Itoa(int(file.GetDeathTime().Unix())) + "#|#"
		}
	}
	publicToken := GetPublicToken()
	for _, val := range publicToken.OwnedFiles {
		found, file := globalFileList.FindFile(val)
		if found {
			stringified += file.Name + "|#|" + strconv.Itoa(int(file.GetDeathTime().Unix())) + "|#|" + file.Compression + "#|#"
		}
	}
	fmt.Fprintf(w, stringified)
	return
}
