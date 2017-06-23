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
			//wait
			time.Sleep(1 * time.Second)
			newFlSerialization := string(globalFileList.Serialize())
			newTmSerialization := string(SerializeTokenMap())

			if oldTmSerialization != newTmSerialization {
				oldTmSerialization = newTmSerialization

				//save token map
				tmSerializationFile, err := os.Create(Configuration.StatePath + "token_map.json")
				logo.RuntimeError(err)
				_, err = io.Copy(tmSerializationFile, strings.NewReader(oldTmSerialization))
				logo.RuntimeError(err)
				err = tmSerializationFile.Close()
				logo.RuntimeError(err)
			}

			if oldFlSerialization != newFlSerialization {
				oldFlSerialization = newFlSerialization

				//save file list
				fileListSerializationFile, err := os.Create(Configuration.StatePath + "file_list.json")
				logo.RuntimeError(err)
				_, err = io.Copy(fileListSerializationFile, strings.NewReader(oldFlSerialization))
				logo.RuntimeError(err)
				err = fileListSerializationFile.Close()
				logo.RuntimeError(err)
			}
		}
	}()
	//Loop to update token's file permissions
	go func() {
		for {
			//wait
			time.Sleep(10 * time.Second)
			RunUnderAuthWMutex(func(authMap *map[string]Token) interface{} {
				//Synchronize permissions between readers and equals
				for _, token := range (*authMap) {
					for _, equalId := range token.Equals {
						equal := (*authMap)[equalId]
						if (len(equal.OwnedFiles) != len(token.OwnedFiles)) {
							//@TODO see how the hell the copy function works
							//Extend capacity of first array than copy ~!?
							//Extend length or capacity ? Does it matter ??
							//Pressumably length ? How do I extend length ?
							for _, ele := range token.OwnedFiles {
								equal.OwnedFiles = append(equal.OwnedFiles, ele)
							}
							tokenMap[equal.Identifier] = equal
						}
					}
					for _, readerId := range token.Equals {
						reader := (*authMap)[readerId]
						if (len(reader.ReadPermission) != len(token.ReadPermission)) {
							for _, ele := range token.ReadPermission {
								reader.ReadPermission = append(reader.ReadPermission, ele)
							}
							tokenMap[reader.Identifier] = reader
						}
					}
				}
				//Remove duplicates
				for _, token := range (*authMap) {
					var alreadyPresentR map[string]bool = make(map[string]bool)
					var newRFileList []string
					for _, fileName := range token.ReadPermission {
						_, exists := alreadyPresentR[fileName]
						if (!exists) {
							alreadyPresentR[fileName] = true
							newRFileList = append(newRFileList, fileName)
						}
					}
					token.ReadPermission = newRFileList

					var alreadyPresentW map[string]bool = make(map[string]bool)
					var newWFileList []string
					for _, fileName := range token.OwnedFiles {
						_, exists := alreadyPresentW[fileName]
						if (!exists) {
							alreadyPresentW[fileName] = true
							newWFileList = append(newWFileList, fileName)
						}
					}
					token.OwnedFiles = newWFileList
					(*authMap)[token.Identifier] = token
				}
				return true
			})
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
	} else {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
		return
	}
	http.Redirect(w, r, "/", http.StatusOK)
	return
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
	public := r.FormValue("ispublis")
	logo.LogDebug("\n\nThe value of public is " + public + " \n\n")
	logo.LogDebug("\n\nThe value of compression is " + compression + " \n\n")

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
	fmt.Println("1")
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
	fmt.Println("2")
	RunUnderAuthWMutex(func(tokenMap *map[string]Token) interface{} {
		valid, token := ValidateSession(values[0], values[1])
		if !valid {
			fmt.Fprintf(w, `Session Id invalid\n`)
			return false
		}
		/*
			@TODO fix this, currently the comparison seems not to work
			logo.LogDebug(token.UploadSize < newFileModel.Size)
			if token.UploadSize < newFileModel.Size {
				fmt.Fprintf(w, `Trying to upload a file that is too big, your maximum upload size is: ` + strconv.Itoa(int(newFileModel.Size)))
				return
			}
		*/

		if token.UploadNumber < 1 {
			fmt.Fprint(w, "Reached maximum upload limit for this token/user\n")
			return false
		}

		token.UploadNumber = token.UploadNumber - 1
		token.OwnedFiles = append(token.OwnedFiles, newFileModel.Name)
		(*tokenMap)[token.Identifier] = token
		return true

		if public == "true" || public == "on" {
			var publicToken Token = (*tokenMap)["public"]
			publicToken.ReadPermission = append(publicToken.ReadPermission, name)
			(*tokenMap)["public"] = publicToken
		}
		return true
	})
	fmt.Println("3")
	//Uploading
	logo.LogDebug(newFileModel.Path)
	permanentFile, err := os.OpenFile(newFileModel.Path, os.O_WRONLY|os.O_CREATE, 0666)
	fmt.Println("4")
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
	fmt.Println("5")
	fmt.Fprintf(w, "Successfully uploaded file\n")
	return
}

//List files
func listFiles(w http.ResponseWriter, r *http.Request) {
	RunUnderAuthRMutex(func(tokenMap *map[string]Token) interface{} {
		stringified := ""
		var oldNames map[string]bool = make(map[string]bool)
		publicToken := (*tokenMap)["public"]
		for _, val := range publicToken.ReadPermission {
			found, file := globalFileList.FindFile(val)
			if found {
				oldNames[file.Name] = true
				stringified += file.Name + "|#|" + strconv.Itoa(int(file.GetDeathTime().Unix())) + "|#|" + file.Compression + "#|#"
			}
		}

		//Doing the authentication
		cookie, err := r.Cookie("auth")
		if err != nil {
			fmt.Fprintf(w, stringified)
			return false
		}
		values := strings.Split(cookie.Value, "#|#")
		if len(values) < 2 {
			fmt.Fprintf(w, stringified)
			return false
		}
		valid, token := ValidateSession(values[0], values[1])
		if !valid {
			fmt.Fprintf(w, stringified)
			return false
		}

		for _, val := range token.ReadPermission {
			found, file := globalFileList.FindFile(val)
			if found {
				_, ok := oldNames[file.Name]
				if !ok {
					oldNames[file.Name] = true
					stringified += file.Name + "|#|" + strconv.Itoa(int(file.Size)) + "|#|" + strconv.Itoa(int(file.GetDeathTime().Unix())) + "#|#"
				}
			}
		}
		stringified = strings.TrimRight(stringified, "#|#")
		fmt.Fprintf(w, stringified)
		return true
	})
	return
}

//getFile
func getFile(w http.ResponseWriter, r *http.Request) {
	RunUnderAuthWMutex(func(tokenMap *map[string]Token) interface{} {
		fileName := r.URL.Query().Get("file")
		publicToken := (*tokenMap)["public"]

		//There's actually not a build in find '( ...
		//also the list should probably be sroted and I should implement a basic search function
		for _, val := range publicToken.OwnedFiles {
			if val == fileName {
				found, file := globalFileList.FindFile(val)
				if found {
					w.Header().Set("Content-Disposition", "attachment; filename="+file.Name)
					http.ServeFile(w, r, file.Path)
					return true
				}
			}
		}

		//Doing the authentication
		cookie, err := r.Cookie("auth")
		if err != nil {
			fmt.Fprintf(w, "File is not public and you are not authenticated")
			return false
		}
		values := strings.Split(cookie.Value, "#|#")
		if len(values) < 2 {
			fmt.Fprintf(w, "File is not public and you are not authenticated")
			return false
		}
		valid, token := ValidateSession(values[0], values[1])
		if !valid {
			fmt.Fprintf(w, "File is not public and you are not authenticated")
			return false
		}

		for _, val := range token.OwnedFiles {
			if val == fileName {
				found, file := globalFileList.FindFile(val)
				if found {
					w.Header().Set("Content-Disposition", "attachment; filename="+file.Name)
					http.ServeFile(w, r, file.Path)
					return true
				}
			}
		}

		fmt.Fprintf(w, "File not found")
		return true
	})
	return
}

//removeFile
func removeFile(w http.ResponseWriter, r *http.Request) {
	//@TODO FIND WAYS TO REMOVE FROM EVERY TOKEN.... ARGH -_-
	fileName := r.URL.Query().Get("file")

	//Doing the authentication
	cookie, err := r.Cookie("auth")
	if err != nil {
		fmt.Fprintf(w, "File is not public and you are not authenticated")
		return
	}
	values := strings.Split(cookie.Value, "#|#")
	if len(values) < 2 {
		fmt.Fprintf(w, "File is not public and you are not authenticated")
		return
	}
	valid, token := ValidateSession(values[0], values[1])
	if !valid {
		fmt.Fprintf(w, "File is not public and you are not authenticated")
		return
	}

	for _, val := range token.OwnedFiles {
		if val == fileName {
			succ, _ := globalFileList.DeleteFile(fileName)
			if !succ {
				logo.LogDebug("Error deleting file")
				fmt.Fprintf(w, "Error deleting file")
				return
			}
			fmt.Fprintf(w, "Successfully removed file")
			return
		}
	}
	fmt.Fprintf(w, "File not found")
}

//creatToken
func createToken(w http.ResponseWriter, r *http.Request) {
	RunUnderAuthWMutex(func(tokenMap *map[string]Token) interface{} {
		//Doing the authentication
		cookie, err := r.Cookie("auth")
		if err != nil {
			fmt.Fprintf(w, "File is not public and you are not authenticated")
			return false
		}
		values := strings.Split(cookie.Value, "#|#")
		if len(values) < 2 {
			fmt.Fprintf(w, "File is not public and you are not authenticated")
			return false
		}
		valid, token := ValidateSession(values[0], values[1])
		if !valid {
			fmt.Fprintf(w, "File is not public and you are not authenticated")
			return false
		}

		if !token.GrantToken {
			fmt.Fprintf(w, "You don't have token granting privilages")
			return false
		}

		identifier := r.URL.Query().Get("identifier")
		credentials := r.URL.Query().Get("credentials")
		uploadNumber, err := strconv.ParseInt(r.URL.Query().Get("uploadNumber"), 10, 64)
		logo.RuntimeError(err)
		uploadSize, err := strconv.ParseInt(r.URL.Query().Get("uploadSize"), 10, 64)
		logo.RuntimeError(err)
		reader, err := strconv.ParseBool(r.URL.Query().Get("reader"))
		logo.RuntimeError(err)
		writer, err := strconv.ParseBool(r.URL.Query().Get("writer"))
		logo.RuntimeError(err)
		admin, err := strconv.ParseBool(r.URL.Query().Get("admin"))
		logo.RuntimeError(err)

		newTokenReadFiles := make([]string, len(token.OwnedFiles)+100)
		newTokenOwneFiles := make([]string, len(token.OwnedFiles)+100)
		if reader {
			copy(newTokenOwneFiles, token.ReadPermission)
			token.Readers = append(token.Readers, identifier)
		}
		if writer {
			copy(newTokenReadFiles, token.OwnedFiles)
			token.Equals = append(token.Equals, identifier)
		}

		newToken := MakeToken(identifier, credentials, newTokenReadFiles, uploadSize, uploadNumber, newTokenOwneFiles, admin, []string{}, []string{})

		if _, ok := (*tokenMap)[identifier]; ok {
			fmt.Fprintf(w, "Token with said name already exists")
			return false
		} else {
			fmt.Fprintf(w, "Token was added")
			(*tokenMap)[identifier] = newToken
		}
		return true
	})
}
