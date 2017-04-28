package main

import (
        "golang.org/x/crypto/scrypt"
        "git.cerebralab.com/george/logo"
)

//Toke structure
type Token struct {
    hash              []byte
    ReadPermission    []RaedPermission
    WritePermission   WritePermission
}

type RaedPermission struct {
    name    string
    number  int
}

type WritePermission struct {
    size    int
    number  int
}

//MakeToken is the function that can be called to create a token

//ValidateToke checks the validty of a token and return the correspondent structure if the token is valid
func ValidateToke(credentials string) (bool, Token) {
    if(hashCredentials(credentials) == )
}


//Some global variables
var salt []byte

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

func initializeAuthentication([]byte theSalt) {
    salt = theSalet
}
