package main

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log"
	"math"
	"strings"
	"sync"

	"golang.org/x/crypto/scrypt"
)

//Some constants
var MaxInt int64 = int64(math.Pow(2, 62))

//Some global variables
var salt []byte = []byte("")
var tokenMap map[string]Token = make(map[string]Token)
var authMutex sync.RWMutex = sync.RWMutex{}

//Toke structure
type Token struct {
	Identifier     string   `json:"identifier"`
	Hash           []byte   `json:"hash"`
	ReadPermission []string `json:"readPermission"`
	UploadNumber   int64    `json:"uploadNumber"`
	UploadSize     int64    `json:"uploadSize"`
	OwnedFiles     []string `json:"ownedFiles"`
	GrantToken     bool     `json"grantToken"`
	Readers        []string `json:"readers"`
	Equals         []string `json:equals`
	sessionIdHash  []byte
}

//SerializeTokenMap serializez the global token map
func SerializeTokenMap() []byte {
	var serialization string
	var stringifiedTokens []string = []string{}
	for _, token := range tokenMap {
		stringifiedTokens = append(stringifiedTokens, string(token.Serialize()))
	}
	serialization = strings.Join(stringifiedTokens, "#|#")
	return []byte(serialization)
}

//DeserializeToken deserializez a byte array into a token map
func DeserializeTokenMap(serialization []byte) {
	tokenArr := strings.Split(string(serialization), "#|#")
	for _, serializedToken := range tokenArr {
		token := DeserializeToken([]byte(serializedToken))
		tokenMap[token.Identifier] = token
	}
}

//IsOwner
func (token *Token) IsOwner(filename string) bool {
	for _, file := range token.OwnedFiles {
		if file == filename {
			return true
		}
	}

	for _, equal := range token.Equals {
		for _, file := range tokenMap[equal].OwnedFiles {
			if file == filename {
				return true
			}
		}
	}

	for _, reader := range token.Readers {
		for _, file := range tokenMap[reader].OwnedFiles {
			if file == filename {
				return true
			}
		}
	}

	return false
}

//IsReader
func (token *Token) IsReader(filename string) bool {
	if token.IsOwner(filename) {
		return true
	}

	for _, equal := range token.Equals {
		for _, file := range tokenMap[equal].ReadPermission {
			if file == filename {
				return true
			}
		}
	}

	for _, reader := range token.Readers {
		for _, file := range tokenMap[reader].ReadPermission {
			if file == filename {
				return true
			}
		}
	}

	for _, file := range token.ReadPermission {
		if file == filename {
			return true
		}
	}

	//Finally check if its public
	return IsPublic(filename)
}

//Serialize gives a string (as a byte slice) represntation of a Token struct
func (token *Token) Serialize() []byte {
	serialization, err := json.Marshal(token)
	if err != nil {
		log.Printf("There was a serialization error for a token %s", err)
	}
	return serialization
}

//DeserializeFileModel takes a byte slice and create a Token
func DeserializeToken(serialization []byte) Token {
	var newToken Token
	err := json.Unmarshal(serialization, &newToken)
	if err != nil {

	}
	return newToken
}

//MakeToken is the function that can be called to create a token
func MakeToken(identifier string, credentials string, readPermission []string, uploadSize int64, uploadNumber int64, ownedFiles []string, grantToken bool, readers []string, equals []string) Token {
	return Token{Identifier: identifier, Hash: hashCredentials(credentials), ReadPermission: readPermission, UploadNumber: uploadNumber, UploadSize: uploadSize,
		OwnedFiles: ownedFiles, GrantToken: grantToken, Readers: readers, Equals: equals}
}

func RunUnderAuthWMutex(task func(*map[string]Token) interface{}) interface{} {
	authMutex.Lock()
	result := task(&tokenMap)
	authMutex.Unlock()
	return result
}

func RunUnderAuthRMutex(task func(*map[string]Token) interface{}) interface{} {
	authMutex.RLock()
	result := task(&tokenMap)
	authMutex.RUnlock()
	return result
}

//ValidateSession checks the validty of a token and return the correspondent structure if the token is valid
func ValidateSession(identifier string, sessionId string) (bool, Token) {
	requestedToken := tokenMap[identifier]
	if bytes.Equal(hashCredentials(sessionId), requestedToken.sessionIdHash) {
		return true, requestedToken
	}
	log.Println("Someone used identifier '" + identifier + "' in order to try accessing a token for which he didn't have credentials")
	return false, requestedToken
}

//ValidateToke validates an ongoing session
func ValidateToke(identifier string, credentials string, cheat bool) (bool, string) {
	requestedToken := tokenMap[identifier]
	if bytes.Equal(hashCredentials(credentials), requestedToken.Hash) || cheat {
		// Previous random seed: random := rand.New(rand.NewSource(time.Now().Unix() - time.Now().UnixNano()))
		sessionIdBytes := make([]byte, 32)
		_, err := rand.Read(sessionIdBytes)
		if err != nil {
			panic(err)
			panic("This is a standard lib failure, thsi shouldn't happen, I am confused, help!, Trying again !")
			_, err := rand.Read(sessionIdBytes)
			if err != nil {
				log.Fatal("Random session id generation not working, something is terribly wrong, crashing!")
			}
		}
		//Converting directly to string => invalid cookie value
		//Converting each byte to an int should be 'random enough'
		sessionId := hex.EncodeToString(sessionIdBytes)
		requestedToken.sessionIdHash = hashCredentials(sessionId)
		tokenMap[identifier] = requestedToken
		return true, sessionId
	}
	log.Println("Someone used identifier '" + identifier + "' in order to try accessing a token for which he didn't have credentials")
	return false, ""
}

//IsPublic tells us if a file is public once the token map has been initialized
func IsPublic(filename string) bool {
	publicToken := tokenMap["public"]
	for _, val := range publicToken.ReadPermission {
		if filename == val {
			return true
		}
	}
	return false
}

//UploadToken uploads the token map

//hashCredentials is a function used for hashing a string
//It will be used internally for storing of all tokens and/or passwords until the need arises for better security
func hashCredentials(credentials string) []byte {
	//Using values recommended on https://godoc.org/golang.org/x/crypto/scrypt  for N, r,p
	//Generating a 32-byte hash key (again, since that's the example)
	hash, err := scrypt.Key([]byte(credentials), salt, 16384, 8, 1, 32)
	//Hashing of credentials fails, this shouldn't happen and I don't know how to handle it, crashing app
	if err != nil {
		log.Fatal(err)
	}
	return hash
}

func InitializeAuthentication(theSalt []byte) {
	salt = theSalt
}

//InitializeAdmin is a function to be used for debuging, creates the user admin - admin
func InitializeAdmin(theSalt []byte, name string, password string) {
	adminToken := MakeToken(name, password, []string{}, MaxInt, MaxInt, []string{}, true, []string{}, []string{})
	publicToken := MakeToken("public", "", []string{}, 0, 0, []string{}, false, []string{}, []string{})
	if _, ok := tokenMap[adminToken.Identifier]; !ok {
		tokenMap[adminToken.Identifier] = adminToken
	}
	if _, ok := tokenMap[publicToken.Identifier]; !ok {
		tokenMap[publicToken.Identifier] = publicToken
	}
}
