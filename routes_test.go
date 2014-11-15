package routes

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

var HandlerOk = func(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello world")
	w.WriteHeader(http.StatusOK)
}

var HandlerErr = func(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "", http.StatusBadRequest)
}

var handler = GetMuxInstance()

func TestRouteSetting(t *testing.T) {
	handler.Get(`/user/:username/:id(\d+)\.:format(\w+)`, HandlerOk)
	handler.Get(`/user/:id(\d+)`, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("get user"))
	})
	handler.Get(`/user/:id(\d+)/create`, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("get create user form"))
	})
	handler.Get(`/user/:id(\d+)/edit`, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("get edit user form"))
	})
	handler.Get(`/user/`, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("get users"))
	})
	handler.Post(`/user/`, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("create user"))
	})
	handler.Put(`/user/:id(\d+)`, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("update user"))
	})
	handler.Delete(`/user/:id(\d+)`, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("delete user"))
	})
}

func TestRouteParams(t *testing.T) {
	r, _ := http.NewRequest("GET", "/user/barbery/123.json?learn=kungfu", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	var (
		username   = r.URL.Query().Get(":username")
		id         = r.URL.Query().Get(":id")
		format     = r.URL.Query().Get(":format")
		learnParam = r.URL.Query().Get("learn")
	)

	if username != "barbery" {
		t.Errorf("url param set to [%s]; want [%s]", username, "barbery")
	}
	if id != "123" {
		t.Errorf("url param set to [%s]; want [%s]", id, "123")
	}
	if learnParam != "kungfu" {
		t.Errorf("url param set to [%s]; want [%s]", learnParam, "kungfu")
	}
	if format != "json" {
		t.Errorf("url param set to [%s]; want [%s]", format, "json")
	}

}

func TestGetMethod(t *testing.T) {
	r, _ := http.NewRequest("GET", "/user/1", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)
	body := w.Body.String()
	if body != "get user" {
		t.Errorf("url content set to [%s]; want [%s]", body, "get user")
	}

	r, _ = http.NewRequest("GET", "/user/1/create", nil)
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, r)
	body = w.Body.String()
	if body != "get create user form" {
		t.Errorf("url content set to [%s]; want [%s]", body, "get create user form")
	}

	r, _ = http.NewRequest("GET", "/user/1/edit", nil)
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, r)
	body = w.Body.String()
	if body != "get edit user form" {
		t.Errorf("url content set to [%s]; want [%s]", body, "get edit user form")
	}

	r, _ = http.NewRequest("GET", "/user", nil)
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, r)
	body = w.Body.String()
	if body != "get users" {
		t.Errorf("url content set to [%s]; want [%s]", body, "get users")
	}
}

func TestPostMethod(t *testing.T) {
	r, _ := http.NewRequest("POST", "/user/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)
	body := w.Body.String()
	if body != "create user" {
		t.Errorf("url content set to [%s]; want [%s]", body, "create user")
	}
}

func TestPutMethod(t *testing.T) {
	r, _ := http.NewRequest("PUT", "/user/789", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)
	body := w.Body.String()
	id := r.URL.Query().Get(":id")
	if id != "789" {
		t.Errorf("url param set to [%s]; want [%s]", id, "789")
	}
	if body != "update user" {
		t.Errorf("url content set to [%s]; want [%s]", body, "update user")
	}
}

func TestDeleteMethod(t *testing.T) {
	r, _ := http.NewRequest("DELETE", "/user/456", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)
	body := w.Body.String()
	id := r.URL.Query().Get(":id")
	if id != "456" {
		t.Errorf("url param set to [%s]; want [%s]", id, "456")
	}
	if body != "delete user" {
		t.Errorf("url content set to [%s]; want [%s]", body, "delete user")
	}
}
