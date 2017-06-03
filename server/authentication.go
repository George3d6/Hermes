package main

import (
        "sync"
        "encoding/json"
        "bytes"
        "math/rand"
        "strconv"
        "time"
        "math"

        "golang.org/x/crypto/scrypt"

        "git.cerebralab.com/george/logo"
)

//Some constants
var MaxInt int64 = int64(math.Pow(2, 62))

//Some global variables
var salt        []byte           = []byte("")
var tokenMap    map[string]Token = make(map[string]Token)
var mutex       sync.RWMutex     = sync.RWMutex{}

//Toke structure
type Token struct {
    Identifier        string            `json:"identifier"`
    hash              []byte            `json:"hash"`
    ReadPermission    []string          `json:"readPermission"`
    UploadNumber      int64             `json:"uploadNumber"`
    UploadSize        int64             `json:"uploadSize"`
    OwnedFiles        []string          `json:"ownedFiles"`
    GrantToken        bool              `json"grantToken"`
    Readers           []string          `json:"readers"`
    Equals            []string          `json:equals`
    sessionIdHash     []byte
}

//Serialize gives a string (as a byte slice) represntation of a Token struct
func (token * Token) Serialize() []byte {
    serialization, err := json.Marshal(token)
	logo.RuntimeError(err)
	return serialization
}

//DeserializeFileModel takes a byte slice and create a Token
func DeserializeToken(serialization []byte) Token {
	var newToken Token
	logo.RuntimeError(json.Unmarshal(serialization, &newToken))
	return newToken
}

//MakeToken is the function that can be called to create a token
func MakeToken(identifier string, credentials string, readPermission []string, uploadSize int64, uploadNumber int64, ownedFiles []string, grantToken bool, readers []string, equals []string) Token {
    return Token{Identifier: identifier, hash: hashCredentials(credentials), ReadPermission: readPermission, UploadNumber: uploadNumber, UploadSize: uploadSize,
                OwnedFiles: ownedFiles, GrantToken: grantToken, Readers: readers, Equals: equals}
}

//AddToken adds a token to the map of tokens
func AddToken(newToken Token) bool {
    mutex.Lock()
    if _, exists := tokenMap[newToken.Identifier]; exists {
        mutex.Unlock()
        return false
    }
    tokenMap[newToken.Identifier] = newToken
    mutex.Unlock()
    return true
}

//ModifyToken modifies a token in the map
func ModifyToken(newToken Token) bool {
    mutex.Lock()
    if _, exists := tokenMap[newToken.Identifier]; !exists {
        mutex.Unlock()
        return false
    }
    tokenMap[newToken.Identifier] = newToken
    mutex.Unlock()
    return true
}

//GetPublicToken checks the public token
func GetPublicToken() Token {
    mutex.RLock()
    token := tokenMap["public"]
    mutex.RUnlock()
    return token
}

//UpdatePublicToken update the public token, multiple threads can do it at the same time so we can't do it though the modify method
func UpdatePublicToken(filename string) {
    mutex.Lock()
    newPublicToken := tokenMap["public"]
    newPublicToken.OwnedFiles = append(newPublicToken.OwnedFiles, filename)
    tokenMap["public"] = newPublicToken
    mutex.Unlock()
    return
}

//ValidateSession checks the validty of a token and return the correspondent structure if the token is valid
func ValidateSession(identifier string, sessionId string) (bool, Token) {
    mutex.RLock()
    requestedToken := tokenMap[identifier]
    mutex.RUnlock()
    if bytes.Equal(hashCredentials(sessionId), requestedToken.sessionIdHash) {
        return true, requestedToken
    }
    logo.LogDebug("Someone used identifier '" + identifier + "' in order to try accessing a token for which he didn't have credentials");
    return false, requestedToken
}

//ValidateToke validates an ongoing session
func ValidateToke(identifier string, credentials string) (bool, string) {
    mutex.RLock()
    requestedToken := tokenMap[identifier]
    mutex.RUnlock()
    if bytes.Equal(hashCredentials(credentials), requestedToken.hash) {
        random := rand.New(rand.NewSource(time.Now().Unix() - time.Now().UnixNano()))
        sessionId := strconv.Itoa(random.Int())
        requestedToken.sessionIdHash = hashCredentials(sessionId)
        mutex.Lock()
        tokenMap[identifier] = requestedToken
        mutex.Unlock()
        return true, sessionId
    }
    logo.LogDebug("Someone used identifier '" + identifier + "' in order to try accessing a token for which he didn't have credentials");
    return false, ""
}

//UploadToken uploads the token map

//hashCredentials is a function used for hashing a string
//It will be used internally for storing of all tokens and/or passwords until the need arises for better security
func hashCredentials(credentials string)[]byte {
    //Using values recommended on https://godoc.org/golang.org/x/crypto/scrypt  for N, r,p
    //Generating a 32-byte hash key (again, since that's the example)
    hash, err := scrypt.Key([]byte(credentials), salt, 16384, 8, 1, 32)
    //Hashing of credentials fails, this shouldn't happen and I don't know how to handle it, crashing app
    logo.RuntimeFatal(err)
    return hash
}

func initializeAuthentication(theSalt []byte ) {
    salt = theSalt
}

func initializeAdmin(theSalt []byte, name string, password string) {
    adminToken := MakeToken(name, password, []string{}, MaxInt, MaxInt, []string{}, true, []string{}, []string{})
    publicToken := MakeToken("public", "", []string{}, 0, 0, []string{}, false, []string{}, []string{})
    AddToken(adminToken)
    AddToken(publicToken)
}
