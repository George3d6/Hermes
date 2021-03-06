package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

//Used for holding the configuration
type configurationHolder struct {
	Port      string `json:"port"`
	LogPath   string `json:"logPath"`
	FilePath  string `json:"filePath"`
	StatePath string `json:"statePath"`
}

var Configuration = configurationHolder{}

func main() {
	var configurationFile string = "config.json"
	if len(os.Args) > 1 {
		configurationFile = os.Args[1]
	}

	file, err := os.Open(configurationFile)
	if err != nil {
		log.Println(err)
	}
	decoder := json.NewDecoder(file)

	err = decoder.Decode(&Configuration)
	if err != nil {
		panic("Cannot decode configuration")
	}

	initHandlers(os.Args[2], os.Args[3])

	var server = &http.Server{
		Addr:              ":" + Configuration.Port,
		ReadTimeout:       300 * time.Second,
		WriteTimeout:      300 * time.Second,
		ReadHeaderTimeout: 300 * time.Second,
		MaxHeaderBytes:    500000000}

	//Static ressources serving
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./client/static"))))
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./client/static/logo_mini.png")
	})

	//API calls
	http.HandleFunc("/post/file/", uploadFile)
	http.HandleFunc("/post/token/", createToken)

	http.HandleFunc("/get/list/", listFiles)
	http.HandleFunc("/get/file/", getFile)
	http.HandleFunc("/get/authentication/", engageAuthSession)

	http.HandleFunc("/delete/file/", removeFile)

	//veiw rendering
	http.HandleFunc("/", serveHome)

	//Start Server
	log.Printf("Server will start running on port: %s\n", Configuration.Port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
