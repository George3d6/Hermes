package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
	"compress/gzip"
	"strings"

	"git.cerebralab.com/george/logo"

	"github.com/ulikunitz/xz"
)

var globalFileList FileList

func initHandlers() {
	globalFileList = CreateFileList()
	initializeAdmin([]byte("the salt"), "admin", "admin")
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
		authCookie := http.Cookie{Path: "/", Name: "auth", Value: string(identifier + "#|#" + sessionId), Expires: expiration, MaxAge: 3600*24}
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
	r.ParseMultipartForm(32 << 20)

	file, _, err := r.FormFile("file")
	if logo.RuntimeError(err) {
		fmt.Fprintf(w, `Could not upload file, file can't be read`)
		return
	}

	defer file.Close()
	name := r.FormValue("name")
	compression := r.FormValue("compression")
	logo.LogDebug(compression + "\n\n\n\n")

    ttlString := r.FormValue("ttl")
	ttl, err := strconv.Atoi(ttlString)
	if logo.RuntimeError(err) {
		fmt.Fprintf(w, `Could not upload file, can't parse time to live(ttl)`)
		return
	}

	newFileModel := FileModel{Path: configuration.FilePath + name + "." + compression, Name: name, TTL: int64(ttl), Birth: time.Now(),
		Compression: compression, Size: GetFileSizeInBytes(file)}

	//Doing the authentication
	cookie, err := r.Cookie("auth")
	if err != nil {
		fmt.Fprintf(w, `Not authenticated`)
		return
	}
	values := strings.Split(cookie.Value, "#|#")
	if len(values) < 2 {
		fmt.Fprintf(w, `Authentication cookie malformed`)
	}
	valid, token := ValidateSession(values[0], values[1])
	if !valid {
		fmt.Fprintf(w, `Session Id invalid`)
		return
	}
	/* @TODO fix this, currently the comparison seems not to work
	logo.LogDebug(token.UploadSize < newFileModel.Size)
	if token.UploadSize < newFileModel.Size {
		fmt.Fprintf(w, `Trying to upload a file that is too big, your maximum upload size is: ` + strconv.Itoa(int(newFileModel.Size)))
		return
	}
	*/
	logo.LogDebug("><\n")
	if token.UploadNumber < 1 {
		fmt.Fprint(w, "Reached maximum upload limit for this token/user")
		return
	}
	logo.LogDebug("><\n")
	token.UploadNumber = token.UploadNumber - 1
	token.OwnedFiles = append(token.OwnedFiles, newFileModel.Name)
	if !ModifyToken(token) {
		fmt.Fprint(w, "Internal server error, please try again and inform the administrator about this")
		return
	}
	logo.LogDebug("><\n")
	//Uploading
	permanentFile, err := os.OpenFile(newFileModel.Path , os.O_WRONLY | os.O_CREATE, 0666)
	if logo.RuntimeError(err) {
		fmt.Fprintf(w, `Could not upload file, there was an internal file system error, try again in a few seconds`)
		return
	}
	defer permanentFile.Close()
	logo.LogDebug("><\n")
	switch compression := newFileModel.Compression; compression {
	case "gz":
		gzipped, err := gzip.NewWriterLevel(permanentFile, 6)
		if logo.RuntimeError(err) {
			fmt.Fprintf(w, `Could not upload file, there was an internal file system error, try again in a few seconds`)
			return
		}
		defer gzipped.Close()
		io.Copy(gzipped, file)
	case "xz":
		xzw, err := xz.NewWriter(permanentFile)
		if logo.RuntimeError(err) {
			fmt.Fprintf(w, `Could not upload file, there was an internal file system error, try again in a few seconds`)
			return
		}
		defer xzw.Close()
		io.Copy(xzw, file)
	case "plain":
		io.Copy(permanentFile, file)
	default:
		fmt.Fprintf(w, "Wrong compression format")
	}
	if !globalFileList.AddFile(newFileModel) {
		fmt.Fprintf(w, `Could not upload file, a similarly named file already exists`)
		return
	}
	fmt.Fprintf(w, `Successfully uploaded file`)
}
