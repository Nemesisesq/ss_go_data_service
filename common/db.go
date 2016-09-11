package common

import (
    "bytes"
    "log"
    "net/url"
    "os"
    "strconv"
    "strings"
    "time"
    "gopkg.in/mgo.v2"
    "fmt"
)

func GetCollection() (col string) {

    buffer := &bytes.Buffer{}
    buffer.WriteString(time.Now().Month().String())
    buffer.WriteString(strconv.Itoa(time.Now().Year()))

    col = buffer.String()
    return col
}

func GetDatabase() (db string) {
    mongo_uri := os.Getenv("MONGODB_URI")

    u, err := url.Parse(mongo_uri)

    if err != nil {
        log.Panic(err)
    }
    db = strings.Trim(u.Path, "/")

    return db

}

func GetSession() (func() *mgo.Session) {
    fmt.Println("Hello from GetSession")

    mongo_uri := os.Getenv("MONGODB_URI")

    session, err := mgo.Dial(mongo_uri)

    if err != nil {
        log.Panic(err)
    }

    session.SetMode(mgo.Monotonic, true)

    Clone := func() *mgo.Session {

        return session.Copy()
    }

    return Clone
}
