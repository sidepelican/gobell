package user

import (
    "database/sql"
    "fmt"
    "time"
    _ "github.com/mattn/go-sqlite3"
    "os"
)

type User struct {
    userId     string
    mac        string
    name       string
    lastAppear time.Time
}

var dbPath = "./users.db"

func FindUser(userId string) (*User, error) {
    db := getDB()
    defer db.Close()

    stmt, err := db.Prepare("select * from users where user_id=?")
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
        err = rows.Scan(&user.userId, &user.mac, &user.name, &user.lastAppear)
        if err != nil {
            fmt.Println(err)
            return nil, err
        }
        return &user, nil
    }

    return nil, fmt.Errorf("userid: %v is not found", userId)
}

func InsertUser(user User) error {
    db := getDB()
    defer db.Close()

    stmt, err := db.Prepare("insert into users(user_id, mac, name, last_appear) values(?,?,?,?)")
    if err != nil {
        return err
    }

    _, err = stmt.Exec(user.userId, user.mac, user.name, user.lastAppear)
    if err != nil {
        return err
    }

    return nil
}

func UpdateLastAppear(userId string, appear time.Time) error {
    db := getDB()
    defer db.Close()

    stmt, err := db.Prepare("update users set last_appear=? where user_id=?")
    if err != nil {
        return err
    }

    _, err = stmt.Exec(appear, userId)
    if err != nil {
        return err
    }

    return nil
}

func getDB() *sql.DB {
    needInit := !exists(dbPath)
    if needInit {
        _, err := os.Create(dbPath)
        if err != nil {
            fmt.Println(err)
            return nil
        }
        fmt.Printf("new .db file created at: %v\n", dbPath)
    }

    db, err := sql.Open("sqlite3", dbPath)
    if err != nil {
        fmt.Println(err)
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
            fmt.Println(err)
            return nil
        }
        fmt.Println("new table created")
    }

    if db == nil{
        fmt.Println("db is nil. something wrong")
        return nil
    }

    return db
}

func exists(filename string) bool {
    _, err := os.Stat(filename)
    return err == nil
}
