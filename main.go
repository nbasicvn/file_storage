package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type User struct {
	ApiKey  string `json:"apiKey"`
	Dir     string `json:"dir"`
	Active  int    `json:"active"`
	Quota int64  `json:"quota"`
	Source int64  `json:"source"`
}

func getUser(user string) (User, error) {
	u := User{}
	jsonFile, err := os.Open("./configs/users.json")
	if err != nil {
		return u, err
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var result map[string]User
	err = json.Unmarshal([]byte(byteValue), &result)
	if err != nil {
		return u, err
	}
	u = result[user]
	if u.Dir == "" {
		return u, fmt.Errorf("user not found")
	}

	err = os.MkdirAll("./public/"+u.Dir, 0755)
	if err != nil {
		return u, err
	}
	return u, err
}
func getUserFromDomain(domain string) (User, error) {
	u := User{}
	jsonFile, err := os.Open("./configs/domains.json")
	if err != nil {
		return u, err
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var result map[string]string
	err = json.Unmarshal([]byte(byteValue), &result)
	if err != nil {
		return u, err
	}

	if result[domain] == "" {
		return u, fmt.Errorf("domain not found")
	}

	return getUser(result[domain])
}

func main() {
	router := mux.NewRouter()
	router.Use(mux.CORSMethodMiddleware(router))
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "POST" {
				user, err := getUser(r.FormValue("user"))
				if err != nil {
					w.WriteHeader(http.StatusForbidden)
					_ = json.NewEncoder(w).Encode(Response{
						IsError:  true,
						ErrorMsg: err.Error(),
					})
					return
				}
				if user.ApiKey != r.FormValue("apiKey") {
					w.WriteHeader(http.StatusForbidden)
					_ = json.NewEncoder(w).Encode(Response{
						IsError:  true,
						ErrorMsg: "",
					})
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	})
	router.HandleFunc("/api/upload-single-file", UploadSingleFile).Methods("POST")
	router.HandleFunc("/api/upload-temp-file", UploadTempFile).Methods("POST")
	router.HandleFunc("/api/upload-common-file", UploadCommonFile).Methods("POST")
	router.HandleFunc("/api/remove-multiple-file-from-path", RemoveMultipleFile).Methods("POST")
	router.HandleFunc("/api/move-temp-multiple-file-from-path", MoveMultipleFile).Methods("POST")
	router.HandleFunc("/{year}/{month}/{date}/{file}/resize", ResizeImage).Methods("GET")
	router.HandleFunc("/api/getUsageCapacity", UsageCapacity).Methods("POST")
	router.PathPrefix("/").HandlerFunc(DirHandler).Methods("GET")
	srv := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:2508",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		ErrorLog: &log.Logger{

		},
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Println(err)
	}
}
func DirHandler(w http.ResponseWriter, r *http.Request) {
	user, err := getUserFromDomain(r.Host)
	//fmt.Println(user)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		_ = json.NewEncoder(w).Encode(err)
		return
	}
	var webrootStr = "./public/" + user.Dir
	fs := NoDirListingHandler(http.StripPrefix("/", http.FileServer(http.Dir(webrootStr))))
	fs.ServeHTTP(w, r)
}

func NoDirListingHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}
		h.ServeHTTP(w, r)
	})
}
