package config

import (
    "fmt"
    "runtime"
    "os"
    "path"
    "log"
    "github.com/BurntSushi/toml"
)

type Config struct {
    LeasePath string
}

var config Config

func LeasePath() string {
    return config.LeasePath
}

func InitConfig() error {
    path, err := findConfigPath()
    if err != nil {
        return err
    }
    log.Println("load config:", path)

    if _, err := toml.DecodeFile(path, &config); err != nil {
        return err
    }

    return nil
}

func findConfigPath() (string, error) {

    const configFileNAme = "config.tml"
    errret := fmt.Errorf("%s not found at: ", configFileNAme)

    // static path
    if runtime.GOOS == "linux" {
        p := "/etc/gobell/" + configFileNAme
        if exists(p) {
            return p, nil
        }
        errret = fmt.Errorf("%v\n\t%v", errret, p)
    }

    // runpath
    runPath, err := os.Executable()
    if err == nil {
        p := path.Dir(runPath) + "/" + configFileNAme
        if exists(p) {
            return p, nil
        }
        errret = fmt.Errorf("%v\n\t%v", errret, p)
    }

    // current dir
    pwd, err := os.Getwd()
    if err == nil {
        p := pwd + "/" + configFileNAme
        if exists(p) {
            return p, nil
        }
        errret = fmt.Errorf("%v\n\t%v", errret, p)
    }

    return "", errret
}

func exists(filename string) bool {
    _, err := os.Stat(filename)
    return err == nil
}
