package main

import (
    "github.com/labstack/echo/v4"
    "sync"
    "net/http"
)

type (
    user struct {
        ID        int    `json:"id"`
        Firstname string `json:"first-name"`
        Lastname  string `json:"last-name"`
        email     string `json:"email"`
        password  string `json:"password"`
    }
)

var (
    users = map[int]*user{}
    seq   = 1
    lock  = sync.Mutex{}
)

func createuser(c echo.Context) error {
    lock.Lock() // lock web processing for action
    defer lock.Unlock() 

    u := &user { //Creating new user object
        ID: seq,
    }

    if err := c.Bind(u); err != nil {
        return err
    }

    users[u.ID] = u
    seq++
    return c.JSON(http.StatusCreated, u)
}

func main() {
    e := echo.New()

    e.File("/", "html/index.html") //Defining Homepage routing
    e.Static("/css", "css") //Defining assets
    e.POST("/users", createuser)

    e.Logger.Fatal(e.Start(":1323"))
}
