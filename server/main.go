package main

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"git.cerebralab.com/george/logo"
)

//Used for holding the configuration
type configurationHolder struct {
	Port     string `json:"port"`
	LogPath  string `json:"logPath"`
	FilePath string `json:"filePath"`
}

var configuration = configurationHolder{}

//Keep it down to one file ?
func main() {
	var configurationFile string = "config.json"
	if len(os.Args) > 1 {
		configurationFile = os.Args[1]
	};

	file, err := os.Open(configurationFile)
	logo.RuntimeFatal(err)
	decoder := json.NewDecoder(file)

	err = decoder.Decode(&configuration)
	if err != nil {
		panic("Cannot decode configuration")
	}

	logo.InitLoggers(configuration.LogPath)
	initHandlers()
	var server = &http.Server{
		Addr:         ":" + configuration.Port,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	//Static ressources serving
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./client/static"))))
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./client/static/favicon.ico")
	})

	//API calls
	http.HandleFunc("/authenticate/", engageAuthSession)
	http.HandleFunc("/upload/", uploadFile)
	http.HandleFunc("/list/", listFiles)

	//veiw rendering
	http.HandleFunc("/", serveHome)

	//Start Server
	logo.LogDebug("Server will start running on port: " + configuration.Port)
	logo.RuntimeFatal(server.ListenAndServe())
}
