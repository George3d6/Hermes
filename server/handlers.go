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
	"log"

	"github.com/ulikunitz/xz"
)

var globalFileList FileList

func initHandlers(adminName string, adminPassword string) {
	fmt.Println(Configuration.StatePath + "file_list.json")
	flcontent, err := ioutil.ReadFile(Configuration.StatePath + "file_list.json")
	if err != nil {
		log.Println(err)
	}
	tmcontent, err := ioutil.ReadFile(Configuration.StatePath + "token_map.json")
	if err != nil {
		log.Println(err)
	}
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
				if err != nil {
					log.Println(err)
				}
				_, err = io.Copy(tmSerializationFile, strings.NewReader(oldTmSerialization))
				if err != nil {
					log.Println(err)
				}
				err = tmSerializationFile.Close()
				if err != nil {
					log.Println(err)
				}
			}

			if oldFlSerialization != newFlSerialization {
				oldFlSerialization = newFlSerialization

				//save file list
				fileListSerializationFile, err := os.Create(Configuration.StatePath + "file_list.json")
				if err != nil {
					log.Println(err)
				}
				_, err = io.Copy(fileListSerializationFile, strings.NewReader(oldFlSerialization))
				if err != nil {
					log.Println(err)
				}
				err = fileListSerializationFile.Close()
				if err != nil {
					log.Println(err)
				}
			}
		}
	}()
	//Loop to update token's file permissions
	go func() {
		for {
			//wait... no need to do this early, just once in a while to clean up the data model
			//Duplicate/old files will not be seen by the user
			time.Sleep(60 * time.Second)
			RunUnderAuthWMutex(func(authMap *map[string]Token) interface{} {
				//Remove duplicates and files not in the main file list
				for _, token := range *authMap {
					var alreadyPresentR map[string]bool = make(map[string]bool)
					var newRFileList []string
					for _, fileName := range token.ReadPermission {
						_, exists := alreadyPresentR[fileName]
						found, _ := globalFileList.FindFile(fileName)
						if !exists && found {
							alreadyPresentR[fileName] = true
							newRFileList = append(newRFileList, fileName)
						}
					}
					token.ReadPermission = newRFileList

					var alreadyPresentW map[string]bool = make(map[string]bool)
					var newWFileList []string
					for _, fileName := range token.OwnedFiles {
						_, exists := alreadyPresentW[fileName]
						found, _ := globalFileList.FindFile(fileName)
						if !exists && found {
							alreadyPresentW[fileName] = true
							newWFileList = append(newWFileList, fileName)
						}
					}
					token.OwnedFiles = newWFileList
					(*authMap)[token.Identifier] = token
					if(token.MarkedToDie) {
						delete(*authMap, token.Identifier)
					}
				}

				globalFileList.CleanUp()

				return true
			})
		}
	}()
}

func getAuthCookie(w http.ResponseWriter, r *http.Request) (bool, []string) {
	cookie, err := r.Cookie("auth")
	if err != nil {
		fmt.Fprintf(w, `{"status":"error","message":"Not authenticated}`)
		return false, []string{}
	}
	values := strings.Split(cookie.Value, "#|#")
	if len(values) != 2 {
		fmt.Fprintf(w, `{"status":"error","message":"Authentication cookie malformed}`)
		return false, []string{}
	}
	return true, values
}

//Server the index.html file
func serveHome(w http.ResponseWriter, r *http.Request) {
	index, err := ioutil.ReadFile("./client/index.html")
	if err != nil {
		log.Println(err)
	}

	if err != nil {
		if err != nil {
			log.Println(err)
		}
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	fmt.Fprintf(w, string(index))
}

//Authentication
func engageAuthSession(w http.ResponseWriter, r *http.Request) {
	identifier := r.URL.Query().Get("identifier")
	credentials := r.URL.Query().Get("credentials")
	redirect := r.URL.Query().Get("redirect")

	RunUnderAuthWMutex(func(arg1 *map[string]Token) interface{} {
		isValid, sessionId := ValidateToke(identifier, credentials, false)
		if isValid {
			expiration := time.Now().Add(24 * time.Hour)
			authCookie := http.Cookie{Path: "/", Name: "auth", Value: string(identifier + "#|#" + sessionId), Expires: expiration, MaxAge: 3600 * 24}
			http.SetCookie(w, &authCookie)
		} else {
			if redirect == "true" {
				http.Redirect(w, r, "/", http.StatusUnauthorized)
			} else {
				fmt.Fprintf(w, `{"status":"error","message":"Your login credentials are not authorized"}`)
			}
			return true
		}
		if redirect == "true" {
			http.Redirect(w, r, "/", http.StatusFound)
		} else {
			fmt.Fprintf(w, `{"status":"ok","message":"Authentication successful"}`)
		}
		return false
	})
	return
}

//parses multipart form and saves an image called "image" to a file
func uploadFile(w http.ResponseWriter, r *http.Request) {
	//Parsing the upload arguments into the values we shall be working with
	r.ParseMultipartForm(5000000000000000)

	file, _, err := r.FormFile("file")
	if err != nil {
		fmt.Fprintf(w, `{"status":"error","message":"Could not upload file, file can't be read"}`)
		return
	}
	//defer file.Close()
	name := r.FormValue("name")
	compression := r.FormValue("compression")
	public := r.FormValue("ispublis")

	ttlString := r.FormValue("ttl")
	ttl, err := strconv.Atoi(ttlString)
	if err != nil {
		fmt.Fprintf(w, `{"status":"error","message":"Could not upload file, can't parse time to live(ttl)"}`)
		return
	}
	ttl = ttl * 60 * 60

	extension := ""
	if compression == "gz" {
		extension = ".gz"
	}
	if compression == "xz" {
		extension = ".xz"
	}

	newFileModel := FileModel{Path: Configuration.FilePath + strings.Replace(name, "/", "wtf", -1) + extension, Name: name, TTL: int64(ttl), Birth: time.Now(),
		Compression: compression, Size: GetFileSizeInBytes(file)}
	//Doing the authentication
	succ, values := getAuthCookie(w,r);
	if(!succ) { return }

	isAuthenticatedInterface := RunUnderAuthWMutex(func(tokenMap *map[string]Token) interface{} {
		valid, token := ValidateSession(values[0], values[1])
		if !valid {
			fmt.Fprintf(w, `{"status":"error","message":"Session Id invalid"}`)
			return false
		}

		if token.UploadNumber < 1 {
			fmt.Fprint(w, `{"status":"error","message":"Reached maximum upload limit for this token/user"}`)
			return false
		}

		token.UploadNumber = token.UploadNumber - 1
		token.OwnedFiles = append(token.OwnedFiles, newFileModel.Name)

		for _, equalsName := range token.Equals {
			equal := (*tokenMap)[equalsName]
			equal.OwnedFiles = append(equal.OwnedFiles, newFileModel.Name)
			(*tokenMap)[equal.Identifier] = equal
		}

		for _, readersName := range token.Readers {
			reader := (*tokenMap)[readersName]
			reader.ReadPermission = append(reader.ReadPermission, newFileModel.Name)
			(*tokenMap)[reader.Identifier] = reader
		}

		(*tokenMap)[token.Identifier] = token

		if public == "true" || public == "on" {
			var publicToken Token = (*tokenMap)["public"]
			publicToken.ReadPermission = append(publicToken.ReadPermission, name)
			(*tokenMap)["public"] = publicToken
		}
		return true
	})

	isAuthenticated, ok := isAuthenticatedInterface.(bool)
	if !ok {
		panic("Run under auth mutex in file upload should always return bool")
	}
	if !bool(isAuthenticated) {
		return
	}

	//Uploading
	permanentFile, err := os.OpenFile(newFileModel.Path, os.O_WRONLY|os.O_CREATE, 0666)

	if err != nil {
		log.Println(err)
		fmt.Fprintf(w, `{"status":"error","message":"Could not upload file, there was an internal file system error, try again in a few seconds"}`)
		return
	}

	defer permanentFile.Close()

	switch compression := newFileModel.Compression; compression {
	case "gz":
		gzipped, err := gzip.NewWriterLevel(permanentFile, 6)
		if err != nil {
			log.Println(err)
			fmt.Fprintf(w, `{"status":"error","message":"Could not upload file, there was an internal file system error, try again in a few seconds"}`)
			return
		}
		defer gzipped.Close()
		io.Copy(gzipped, file)
	case "xz":
		xzw, err := xz.NewWriter(permanentFile)
		if err != nil {
			log.Println(err)
			fmt.Fprintf(w, `{"status":"error","message":"Could not upload file, there was an internal file system error, try again in a few seconds"}`)
			return
		}
		defer xzw.Close()
		io.Copy(xzw, file)
	case "plain":
		io.Copy(permanentFile, file)
	default:
		fmt.Fprintf(w, `{"status":"error","message":"Wrong compression format"}`)
	}
	if !globalFileList.AddFile(newFileModel) {
		fmt.Fprintf(w, `{"status":"error","message":"Could not upload file, a similarly named file already exists"}`)
		return
	}

	if r.FormValue("isAsync") == "true" {
		fmt.Fprintf(w, `{"status":"error","message":"File uploaded sucessfully"}`)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
	return
}

//List files
func listFiles(w http.ResponseWriter, r *http.Request) {
	RunUnderAuthRMutex(func(tokenMap *map[string]Token) interface{} {
		stringified := ""

		isAuthenticated := false
		var token Token


		//Doing the authentication
		cookie, err := r.Cookie("auth")
		if err == nil {
			values := strings.Split(cookie.Value, "#|#")
			if len(values) == 2 {
				valid, theToken := ValidateSession(values[0], values[1])
				if valid {
					token = theToken
					isAuthenticated = true
				}
			}
		}
		globalFileList.ReadOnFileList(func(fileList []FileModel) interface{} {
			for _, file := range fileList {
				if !isAuthenticated {
					if IsPublic(file.Name) {
						stringified += file.Name + "|#|" + file.Compression + "|#|" + strconv.Itoa(int(file.GetDeathTime().Unix())) + "|#|" + strconv.Itoa(int(file.Size)) + "|#|true#|#"
					}
				} else if isAuthenticated {
					if token.IsReader(file.Name) {
						stringified += file.Name + "|#|" + file.Compression + "|#|" + strconv.Itoa(int(file.GetDeathTime().Unix())) + "|#|" + strconv.Itoa(int(file.Size)) + "|#|false#|#"
					}
				}
			}
			return true
		})

		stringified = strings.TrimRight(stringified, "#|#")
		fmt.Fprintf(w, stringified)
		return true
	})
	return
}

//getFile
func getFile(w http.ResponseWriter, r *http.Request) {
	identifier := r.URL.Query().Get("identifier")
	credentials := r.URL.Query().Get("credentials")

	RunUnderAuthWMutex(func(tokenMap *map[string]Token) interface{} {
		fileName := r.URL.Query().Get("file")
		//There's actually not a build in find '( ...
		//also the list should probably be sroted and I should implement a basic search function
		found, file := globalFileList.FindFile(fileName)
		if !found {
			fmt.Fprintf(w, `{"status":"error","message":"File not found"}`)
			return false
		}

		sufix := ""
		if file.Compression == "gz" {
			sufix = ".gz"
		}
		if file.Compression == "xz" {
			sufix = ".xz"
		}

		if IsPublic(file.Name) {
			w.Header().Set("Content-Disposition", "attachment; filename="+file.Name+sufix)
			http.ServeFile(w, r, file.Path)
			return true
		}

		//Doing the authentication
		var token Token
		if identifier != "" || credentials != "" {
			valid, _ := ValidateToke(identifier, credentials, false)
			if !valid {
				fmt.Fprintf(w, `{"status":"error","message":"File is not public and you are not authenticated"}`)
				return false
			}
			token = (*tokenMap)[identifier]
			token.MarkedToDie = true
			(*tokenMap)[identifier] = token
		} else {
			succ, values := getAuthCookie(w, r)
			if(!succ) { return false }
			valid := false
			valid, token = ValidateSession(values[0], values[1])
			if !valid {
				fmt.Fprintf(w, `{"status":"error","message":"File is not public and you are not authenticated"}`)
				return false
			}
		}
		if token.IsReader(file.Name) {
			w.Header().Set("Content-Disposition", "attachment; filename="+file.Name+sufix)
			http.ServeFile(w, r, file.Path)
			return true
		}

		fmt.Fprintf(w, `{"status":"error","message":"File not found"}`)
		return true
	})
	return
}

//removeFile
func removeFile(w http.ResponseWriter, r *http.Request) {
	//@TODO FIND WAYS TO REMOVE FROM EVERY TOKEN.... ARGH -_-
	filename := r.URL.Query().Get("file")

	//Doing the authentication
	succ, values := getAuthCookie(w, r)
	if(!succ) { return }
	valid, token := ValidateSession(values[0], values[1])
	if !valid {
		fmt.Fprintf(w, `{"status":"error","message":"Session id is invalid, please relog"}`)
		return
	}

	if token.IsOwner(filename) {
		succ, _ := globalFileList.DeleteFile(filename)
		if !succ {
			fmt.Fprintf(w, `{"status":"error","message":"Error deleting file"}`)
			return
		}
		fmt.Fprintf(w, `{"status":"ok","message":"Successfully removed file"}`)
		return
	}
	fmt.Fprintf(w, `{"status":"error","message":"You don't have permission to delete this file"}`)
}

//creatToken
func createToken(w http.ResponseWriter, r *http.Request) {
	RunUnderAuthWMutex(func(tokenMap *map[string]Token) interface{} {
		//Doing the authentication
		succ, values := getAuthCookie(w, r)
		if(!succ) { return false }

		valid, token := ValidateSession(values[0], values[1])
		if !valid {
			fmt.Fprintf(w, `{"status":"error","message":"Session id is invalid, please relog"}`)
			return false
		}

		if !token.GrantToken {
			fmt.Fprintf(w, `{"status":"error","message":"You don't have the privilage to add a token"}`)
			return false
		}

		identifier := r.URL.Query().Get("identifier")
		credentials := r.URL.Query().Get("credentials")
		uploadNumber, err := strconv.ParseInt(r.URL.Query().Get("uploadNumber"), 10, 64)
		if err != nil {
			log.Println(err)
		}
		uploadSize, err := strconv.ParseInt(r.URL.Query().Get("uploadSize"), 10, 64)
		if err != nil {
			log.Println(err)
		}
		reader, err := strconv.ParseBool(r.URL.Query().Get("reader"))
		if err != nil {
			log.Println(err)
		}
		writer, err := strconv.ParseBool(r.URL.Query().Get("writer"))
		if err != nil {
			log.Println(err)
		}
		admin, err := strconv.ParseBool(r.URL.Query().Get("admin"))
		if err != nil {
			log.Println(err)
		}

		newTokenReadFiles := make([]string, 0, len(token.OwnedFiles)+100)
		newTokenOwneFiles := make([]string, 0, len(token.OwnedFiles)+100)
		if reader {
			for _, name := range token.OwnedFiles {
				newTokenReadFiles = append(newTokenReadFiles, name)
			}
			token.Readers = append(token.Readers, identifier)
		}
		if writer {
			for _, name := range token.OwnedFiles {
				newTokenOwneFiles = append(newTokenOwneFiles, name)
			}
			token.Equals = append(token.Equals, identifier)
		}

		newToken := MakeToken(identifier, credentials, newTokenReadFiles, uploadSize, uploadNumber, newTokenOwneFiles, admin, []string{}, []string{})

		if _, ok := (*tokenMap)[identifier]; ok {
			fmt.Fprintf(w, `{"status":"error","message":"Token with said name already exists"}`)
			return false
		} else {
			fmt.Fprintf(w, "Token was added")
			(*tokenMap)[token.Identifier] = token
			(*tokenMap)[identifier] = newToken
		}
		return true
	})
}

/*
//getTokens retunrs all the tokens you granted
func getToken(w http.ResponseWriter, r *http.Request) {
	RunUnderAuthRMutex(func(arg1 *map[string]Token) interface{} {
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
		stringified = ""
		for _, tokenName := token.Readers {
			stringified += tokenName + "|#|" + unsafeCredentials + "#|#"
		}
		for _, tokenName := token.Equals {
			stringified += tokenName + "|#|" + unsafeCredentials + "#|#"
		}

		stringified = strings.TrimRight(stringified, "#|#")
		fmt.Fprintf(w, stringified)
		return true
	})
	return
}
*/
