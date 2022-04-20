package main

import (
	"encoding/json"
	"github.com/disintegration/imaging"
	"github.com/gabriel-vasile/mimetype"
	"github.com/gorilla/mux"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func UploadSingleFile(w http.ResponseWriter, r *http.Request) {
	u, err := getUser(r.FormValue("user"))

	_ = r.ParseMultipartForm(32 << 20)

	slug := r.FormValue("slug")
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(Response{
			IsError: true,
		})
		return
	}
	defer file.Close()

	currentTime := time.Now()
	pathToFolder := "/" + currentTime.Format("2006/01/02") + "/"
	path := "./public/" + u.Dir + pathToFolder
	filename := randomToken(16) + filepath.Ext(fileHeader.Filename)
	if slug != "" {
		filename = slug + "-" + randomToken(5) + filepath.Ext(fileHeader.Filename)
	}

	err = os.MkdirAll(path, 0777)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(Response{
			IsError: true,
		})
		return
	}

	out, err := os.OpenFile(path+filename, os.O_WRONLY|os.O_CREATE, 0644)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(Response{
			IsError: true,
		})
		return
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(Response{
			IsError: true,
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(
		Response{Path: pathToFolder + filename},
	)
}

func UploadTempFile(w http.ResponseWriter, r *http.Request) {
	_, err := getUser(r.FormValue("user"))

	_ = r.ParseMultipartForm(32 << 20)

	slug := r.FormValue("slug")
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(Response{
			IsError: true,
		})
		return
	}
	defer file.Close()

	currentTime := time.Now()
	pathToFolder := "/temp/" + currentTime.Format("2006/01/02") + "/"
	path := "./public" + pathToFolder
	filename := randomToken(16) + filepath.Ext(fileHeader.Filename)
	if slug != "" {
		filename = slug + "-" + randomToken(5) + filepath.Ext(fileHeader.Filename)
	}

	err = os.MkdirAll(path, 0777)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(Response{
			IsError: true,
		})
		return
	}

	out, err := os.OpenFile(path+filename, os.O_WRONLY|os.O_CREATE, 0644)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(Response{
			IsError: true,
		})
		return
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(Response{
			IsError: true,
		})
		return
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(Response{Path: pathToFolder + filename})
}

func UploadCommonFile(w http.ResponseWriter, r *http.Request) {
	_, err := getUser(r.FormValue("user"))

	_ = r.ParseMultipartForm(32 << 20)

	slug := r.FormValue("slug")
	file, fileHeader, err := r.FormFile("file")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(Response{
			IsError: true,
		})
		return
	}
	defer func() {
		err = file.Close()
	}()

	currentTime := time.Now()
	pathToFolder := "/common/" + currentTime.Format("2006/01/02") + "/"
	path := "./public" + pathToFolder
	filename := randomToken(16) + filepath.Ext(fileHeader.Filename)
	if slug != "" {
		filename = slug + "-" + randomToken(5) + filepath.Ext(fileHeader.Filename)
	}

	err = os.MkdirAll(path, 0777)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(Response{
			IsError: true,
		})
		return
	}

	out, err := os.OpenFile(path+filename, os.O_WRONLY|os.O_CREATE, 0644)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(Response{
			IsError: true,
		})
		return
	}
	defer func() {
		err = out.Close()
	}()
	_, err = io.Copy(out, file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(Response{
			IsError: true,
		})
		return
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(Response{Path: pathToFolder + filename})
}

func RemoveMultipleFile(w http.ResponseWriter, r *http.Request) {
	getUser(r.FormValue("user"))

	paths := r.FormValue("path")
	var listPath []string
	err := json.Unmarshal([]byte(paths), &listPath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(Response{
			IsError: true,
		})
		return
	}
	for _, path := range listPath {
		fullPath := "./public" + string(path)
		if _, err := os.Stat(fullPath); !os.IsNotExist(err) {
			err = os.Remove(fullPath)
		}
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(Response{IsError: false})
}

func MoveMultipleFile(w http.ResponseWriter, r *http.Request) {
	getUser(r.FormValue("user"))

	slug := r.FormValue("slug")
	paths := r.FormValue("path")
	var listPath []string
	err := json.Unmarshal([]byte(paths), &listPath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(Response{
			IsError: true,
		})
		return
	}
	currentTime := time.Now()
	pathToFolder := "/" + currentTime.Format("2006/01/02") + "/"
	pathFull := "./public" + pathToFolder

	var pathResults []MediaInfo
	for _, path := range listPath {

		fullTempPath := "./public" + string(path)
		filename := randomToken(16) + filepath.Ext(string(fullTempPath))
		if slug != "" {
			filename = slug + "-" + randomToken(5) + filepath.Ext(string(fullTempPath))
		}

		err = os.MkdirAll(pathFull, 0777)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(Response{
				IsError: true,
			})
			return
		}
		fullNewPath := pathFull + filename

		if file, err := os.Stat(fullTempPath); !os.IsNotExist(err) {
			f, err := os.Open(fullTempPath)
			if err == nil {
				contentType, err := mime(f)
				if err == nil {
					err2 := os.Rename(fullTempPath, fullNewPath)
					if err2 == nil {
						mediaInfo := MediaInfo{Path: pathToFolder + filename, Size: file.Size(), Mime: contentType}
						pathResults = append(pathResults, mediaInfo)
					}
				}
			}
			err = f.Close()

		}
	}
	data, _ := json.Marshal(pathResults)
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(
		Response{IsError: false, Data: string(data)},
	)
}

func ResizeImage(w http.ResponseWriter, r *http.Request) {
	u, err := getUserFromDomain(r.Host)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(Response{
			IsError: true,
		})
		return
	}
	vars := mux.Vars(r)
	width, _ := strconv.Atoi(r.FormValue("width"))
	height, _ := strconv.Atoi(r.FormValue("height"))

	pathToFolder := "./public/" + u.Dir + "/" + vars["year"] + "/" + vars["month"] + "/" + vars["date"] + "/"
	pathToResizeFolder := "./public/" + u.Dir + "/resize/" + vars["year"] + "/" + vars["month"] + "/" + vars["date"] + "/"
	path := pathToFolder + vars["file"]
	pathResize := pathToResizeFolder + strconv.Itoa(width) + "." + strconv.Itoa(height) + "." + vars["file"]

	if _, err := os.Stat(pathResize); os.IsNotExist(err) {
		err = os.MkdirAll(pathToResizeFolder, 0777)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(Response{
				IsError: true,
			})
			return
		}

		srcImg, err := imaging.Open(path)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(Response{
				IsError: true,
			})
			return
		}

		dstImage := imaging.Resize(srcImg, width, height, imaging.Lanczos)
		err = imaging.Save(dstImage, pathResize)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(Response{
				IsError: true,
			})
			return
		}
	}

	img, _ := os.Open(pathResize)
	defer img.Close()

	mime, err := mimetype.DetectFile(pathResize)
	if err != nil {
		w.Header().Set("Content-Type", "image/jpeg")
	} else {
		w.Header().Set("Content-Type", mime.String())
	}
	w.WriteHeader(http.StatusOK)
	_, _ = io.Copy(w, img)
}

func UsageCapacity(w http.ResponseWriter, r *http.Request) {
	u, _ := getUser(r.FormValue("user"))

	path := "./public/" + u.Dir
	size := dirSize(path)
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(Response{IsError: false, Quota: Quota{Size: size + u.Source*1024*1024, Quota: u.Quota}})
}

type Response struct {
	IsError  bool
	Path     string
	Data     string
	ErrorMsg string
	FileInfo FileInfo
	Quota    Quota
}
type Path struct {
	Path string
}
type FileInfo struct {
	Name string
	Size int64
	Mime string
	Ext  string
}

type Quota struct {
	Size  int64
	Quota int64
}

type MediaInfo struct {
	Path    string
	OldPath string
	Size    int64
	Mime    string
}

func randomToken(IdLength int) string {
	ALPHABET := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}
	rtn := ""
	for i := 0; i < IdLength; i++ {
		rtn += ALPHABET[rand.Intn(len(ALPHABET)-1)]
	}
	return rtn
}

func mime(out *os.File) (string, error) {
	buffer := make([]byte, 512)
	_, err := out.Read(buffer)
	if err != nil {
		return "", err
	}
	contentType := http.DetectContentType(buffer)
	return contentType, nil
}

func inArray(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func dirSize(path string) int64 {
	var dirSize int64 = 0

	readSize := func(path string, file os.FileInfo, err error) error {
		if !file.IsDir() {
			dirSize += file.Size()
		}

		return nil
	}

	_ = filepath.Walk(path, readSize)

	return dirSize
}
