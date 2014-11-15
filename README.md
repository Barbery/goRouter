goRouter
========
goRouter is a very flexible and lightweight router with hight performance. It is very suitable for small project or some little service.


## usage

talk is less, show you the code:

```go

package main

import (
    "fmt"
    // import router
    "github.com/Barbery/goRouter"
    "net/http"
)

func main() {
    // get the instance of router
    mux := goRouter.GetMuxInstance()

    // add routes
    mux.Get(`/user/:id(\d+)`, getUser)
    mux.Get(`/user/profile/:id(\d+)\.:format(\w+)`, getUserProfile)
    mux.Post(`/user`, postUser)
    mux.Delete(`/user/:id(\d+)`, deleteUser)
    mux.Put(`/user/:id(\d+)`, putUser)

    // run the serve
    http.ListenAndServe(":8888", mux)
}

func getUser(w http.ResponseWriter, r *http.Request) {
    id := r.URL.Query().Get(":id")
    w.Write([]byte(fmt.Sprintf(`GET user by id: %s`, id)))
}

func getUserProfile(w http.ResponseWriter, r *http.Request) {
    params := r.URL.Query()
    w.Write([]byte(fmt.Sprintf(`GET user profile by id: %s, format: %s`, params[":id"][0], params.Get(":format"))))
}

func postUser(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    w.Write([]byte(fmt.Sprintf(`POST a new user, form data: %s`, fmt.Sprintln(r.PostForm))))
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte(fmt.Sprintf(`DELETE a user by id: %s`, r.URL.Query().Get(":id"))))
}

func putUser(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    w.Write([]byte(fmt.Sprintf(`UPDATE a user by id: %s, form data: %s`, r.URL.Query().Get(":id"), fmt.Sprintln(r.PostForm))))
}


```

