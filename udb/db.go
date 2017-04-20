package udb

import (
    "database/sql"
    "os"
    "log"
    "fmt"
    "time"
    _ "github.com/mattn/go-sqlite3"
    "path/filepath"
    "strings"
)

func getDBPath() string {
    return getRunPath() + "users.db"
}

type DBContext struct {
    db *sql.DB
}

func GetContext() *DBContext {
    return &DBContext{getDB()}
}

func (ctx *DBContext) Close() {
    ctx.db.Close()
}

func (ctx *DBContext) AllUserId() ([]string, error) {

    rows, err := ctx.db.Query("select user_id from users")
    defer rows.Close()
    if err != nil {
        return nil, err
    }

    ret := []string{}
    for rows.Next() {
        var id string
        err = rows.Scan(&id)
        if err != nil {
            return nil, err
        }
        ret = append(ret, id)
    }

    return ret, nil
}

func (ctx *DBContext) FindUser(userId string) (*User, error) {

    stmt, err := ctx.db.Prepare("select * from users where user_id=?")
    if err != nil {
        return nil, err
    }

    rows, err := stmt.Query(userId)
    defer rows.Close()
    if err != nil {
        return nil, err
    }

    for rows.Next() {
        var user User
        err = rows.Scan(&user.UserId, &user.Mac, &user.Name, &user.LastAppear)
        if err != nil {
            return nil, err
        }
        return &user, nil
    }

    return nil, fmt.Errorf("userid: %v is not found", userId)
}

func (ctx *DBContext) FindMac(mac string) (*User, error) {

    stmt, err := ctx.db.Prepare("select * from users where mac like ?")
    if err != nil {
        return nil, err
    }

    rows, err := stmt.Query(mac)
    defer rows.Close()
    if err != nil {
        return nil, err
    }

    for rows.Next() {
        var user User
        err = rows.Scan(&user.UserId, &user.Mac, &user.Name, &user.LastAppear)
        if err != nil {
            return nil, err
        }
        return &user, nil
    }

    return nil, fmt.Errorf("mac: %v is not found", mac)
}

func (ctx *DBContext) InsertUser(user User) error {

    stmt, err := ctx.db.Prepare("insert into users(user_id, mac, name, last_appear) values(?,?,?,?)")
    if err != nil {
        return err
    }

    _, err = stmt.Exec(user.UserId, user.Mac, user.Name, user.LastAppear)
    if err != nil {
        return err
    }

    return nil
}

func (ctx *DBContext) UpdateLastAppear(userId string, appear time.Time) error {

    stmt, err := ctx.db.Prepare("update users set last_appear=? where user_id=?")
    if err != nil {
        return err
    }

    _, err = stmt.Exec(appear, userId)
    if err != nil {
        return err
    }

    return nil
}

func (ctx *DBContext) EraseUser(userId string) error {

    stmt, err := ctx.db.Prepare("delete from users where user_id=?")
    if err != nil {
        return err
    }

    _, err = stmt.Exec(userId)
    if err != nil {
        return err
    }

    return nil
}

func getDB() *sql.DB {
    needInit := !exists(getDBPath())
    if needInit {
        file, err := os.Create(getDBPath())
        if err != nil {
            log.Println(err)
            return nil
        }
        log.Printf("new .db file created at: %v\n", file.Name())
    }

    db, err := sql.Open("sqlite3", getDBPath())
    if err != nil {
        log.Println(err)
        return nil
    }

    if needInit {
        q := `CREATE TABLE users (
         user_id     VARCHAR(255) PRIMARY KEY,
         mac         VARCHAR(255) NOT NULL,
         name        VARCHAR(255) NOT NULL,
         last_appear TIMESTAMP DEFAULT (DATETIME('now','localtime'))
        )`

        _, err := db.Exec(q)
        if err != nil {
            log.Println(err)
            return nil
        }
        log.Println("new table created")
    }

    if db == nil {
        log.Println("db is nil. something wrong")
        return nil
    }

    return db
}

func exists(filename string) bool {
    _, err := os.Stat(filename)
    return err == nil
}

func getRunPath() string {
    dir, err := filepath.Abs(filepath.Dir(os.Args[0])) // Get the absolute path at Executing file. Reference：http://stackoverflow.com/questions/18537257/golang-how-to-get-the-directory-of-the-currently-running-file
    if err != nil {
        log.Println(err)
        return ""
    }

    // for `$go run ~~` support
    if strings.HasPrefix(dir, "/var") {
        return ""
    }

    return dir + "/"
}