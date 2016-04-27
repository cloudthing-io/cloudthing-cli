package utils

import (
    _"fmt"
    "os"
    _"net/http"
    "encoding/json"
    "github.com/mitchellh/go-homedir"
    "path"
    "io/ioutil"
    "gitlab.com/cloudthing/go-api-client"
)


type Auths map[string]*api.Token

const BaseAuthsFileName = ".cloudthing-cli/auths.json"

var AuthsFileName string

func init() {
	dir, err := homedir.Dir()
	if err != nil {
		panic(err)
	}
	expanded, err := homedir.Expand(dir)
	if err != nil {
		panic(err)
	}
	AuthsFileName = path.Join(expanded, BaseAuthsFileName)
}

func LoadAuth(server string) *api.Token {
	a, err := loadAuths()
    if err != nil {
        return nil
    }

    if token, ok := (*a)[server]; ok {
        return token
    }
    return nil
}


func loadAuths() (*Auths, error) {
    if _, err := os.Stat(AuthsFileName); os.IsNotExist(err) {
        return nil, err
    }
    buf, err := ioutil.ReadFile(AuthsFileName)
    if err != nil {
        return nil, err
    }
    var a Auths
    if err := json.Unmarshal(buf, &a); err != nil {
        return nil, err
    }
    return &a, nil
}

func SaveAuth(server string, token *api.Token) error {
    auths, err := loadAuths()
    if err != nil {
        return err
    }
    
    (*auths)[server] = token

    buf, err := json.Marshal(auths)
    if err != nil {
        return err
    }

    if err := ioutil.WriteFile(AuthsFileName, buf, 0600); err != nil {
        return err
    }	
    return nil
}