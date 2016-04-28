package utils

import (
    "fmt"
    "os"
    _"net/http"
    "encoding/json"
    "github.com/mitchellh/go-homedir"
    "path"
    "io/ioutil"
    "gitlab.com/cloudthing/go-api-client"
    "strings"
)


type Auths map[string]*api.Token

const BaseAuthsFileName = ".cloudthing-cli/auths.json"
const BaseApikeyFileName = ".cloudthing-cli/apikey-"

var AuthsFileName string
var ApikeyFileName string

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
    ApikeyFileName = path.Join(expanded, BaseApikeyFileName)
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

func DeleteAuth(server string) error {
    auths, err := loadAuths()
    if err != nil {
        return err
    }
    
    if _, ok := (*auths)[server]; ok {
        delete(*auths, server)
    } 

    buf, err := json.Marshal(auths)
    if err != nil {
        return err
    }

    if err := ioutil.WriteFile(AuthsFileName, buf, 0600); err != nil {
        return err
    }   
    return nil
}

func SaveApikey(key, secret string) (string, error) {
    file := fmt.Sprintf("%s%s.key", ApikeyFileName, key)
    apikey := []byte(fmt.Sprintf("apikey.id=%s\napikey.secret=%s", key, secret))
    if err := ioutil.WriteFile(file, apikey, 0600); err != nil {
        return "", err
    }   
    return file, nil
}

func LoadApikey(file string) (string, string, error) {
    if _, err := os.Stat(file); os.IsNotExist(err) {
        return "", "", err
    }
    buf, err := ioutil.ReadFile(file)
    if err != nil {
        return "", "", err
    }

    var key, secret string
    lines := strings.Split(string(buf),"\n")
    for _, line := range lines {
        col := strings.Split(line, "=")
        if col[0] == "apikey.id" {
            key = col[1]
        }
        if col[0] == "apikey.secret" {
            secret = col[1]
        }
    }
    if key != "" && secret != "" {
        return key, secret, nil
    }
    return "","", fmt.Errorf("API key file is malformed")
}