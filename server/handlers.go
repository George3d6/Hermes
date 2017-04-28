package main

import (
	"fmt"
	"hash/fnv"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
	"compress/gzip"

	"git.cerebralab.com/george/logo"

	"github.com/ulikunitz/xz"
)

var globalFileList FileList

func initHandlers() {
	globalFileList = CreateFileList()
}

//Server the index.html file
func serveHome(w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie("auth")
	if !logo.RuntimeError(err) {
		logo.LogDebug(cookie.Value, false)
	}

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
	user := r.URL.Query().Get("user")
	password := r.URL.Query().Get("password")
	token := r.URL.Query().Get("token")

	hasher := fnv.New64()
	if token != "" {
		hasher.Write([]byte(token))
	} else {
		hasher.Write([]byte(user + "@" + password))

	}
	hash := hasher.Sum64()

	authCookie := http.Cookie{Name: "auth", Value: strconv.FormatUint(hash, 10), MaxAge: 3600 * 24}
	http.SetCookie(w, &authCookie)
}

//parses multipart form and saves an image called "image" to a file
func uploadFile(w http.ResponseWriter, r *http.Request) {

	r.ParseMultipartForm(32 << 20)

	file, _, err := r.FormFile("file")
	if logo.RuntimeError(err) {
		fmt.Fprintf(w, `Could not upload file, file can't be read`)
		return
	}

	defer file.Close()
	name := r.FormValue("name")
	compression := r.FormValue("compression")

    ttlString := r.FormValue("ttl")
	ttl, err := strconv.Atoi(ttlString)
	if logo.RuntimeError(err) {
		fmt.Fprintf(w, `Could not upload file, can't parse time to live(ttl)`)
		return
	}

	newFileModel := FileModel{Path: configuration.FilePath + name + "." + compression, Name: name, TTL: int64(ttl), Birth: time.Now(),
		Compression: compression}

	if !globalFileList.AddFile(newFileModel) {
		fmt.Fprintf(w, `Could not upload file, a similarly named file already exists`)
		return
	}

	//@TODO check authentication here

	permanentFile, err := os.OpenFile(newFileModel.Path , os.O_WRONLY | os.O_CREATE, 0666)
	if logo.RuntimeError(err) {
		fmt.Fprintf(w, `Could not upload file, there was an internal file system error, try again in a few seconds`)
		globalFileList.RemoveFile(name)
		return
	}
	defer permanentFile.Close()

	switch compression := newFileModel.Compression; compression {
	case "gz":
		gzipped, err := gzip.NewWriterLevel(permanentFile, 6)
		if logo.RuntimeError(err) {
			fmt.Fprintf(w, `Could not upload file, there was an internal file system error, try again in a few seconds`)
			globalFileList.RemoveFile(name)
			return
		}
		defer gzipped.Close()
		fmt.Fprintf(w, `Successfully uploaded file`)
		io.Copy(gzipped, file)

	case "xz":
		xzw, err := xz.NewWriter(permanentFile)
		if logo.RuntimeError(err) {
			fmt.Fprintf(w, `Could not upload file, there was an internal file system error, try again in a few seconds`)
			globalFileList.RemoveFile(name)
			return
		}
		defer xzw.Close()
		fmt.Fprintf(w, `Successfully uploaded file`)
		io.Copy(xzw, file)
		
	case "plain":
		fmt.Fprintf(w, `Successfully uploaded file`)
		io.Copy(permanentFile, file)
	}


}
