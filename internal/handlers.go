package bigxxby

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

var artists []Artist = GetContent()

func StaticHandler(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Path
	if r.Method != "GET" {
		ErrorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}
	// fmt.Println((r.Header))
	if r.Header.Get("Sec-Fetch-Dest") == "document" { // если пользовтель просит получить ответ в виде документа выдаем ошибку , в других случаях выдаем 200
		fmt.Println("Permission denied")
		// http.NotFound(w, r)
		ErrorHandler(w, r, 404)
		return
	}

	// Предотвращение навигации вверх по дереву каталогов
	if strings.Contains(filePath, "..") {
		fmt.Println("ERROR:" + " PATH CONTAINS ..")
		// http.NotFound(w, r)
		ErrorHandler(w, r, 404)
		return
	}

	// Открываем файл
	file, err := http.Dir(".").Open(filePath)
	if err != nil {
		fmt.Println("ERROR: ", err)
		// http.NotFound(w, r)
		ErrorHandler(w, r, 404)
		return
	}
	defer file.Close()

	// Получаем информацию о файле
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("ERROR: ", err)
		// http.NotFound(w, r)
		ErrorHandler(w, r, 404)
		return
	}

	// Проверяем, является ли файл каталогом
	if fileInfo.IsDir() {
		// Если это каталог, добавляем индексный файл
		fmt.Println("INDEX")
		filePath += "/index.html"
	}

	// Открываем файл
	file, err = http.Dir(".").Open(filePath)
	if err != nil {
		fmt.Println("ERROR: ", err)
		// http.NotFound(w, r)
		ErrorHandler(w, r, 404)
		return
	}
	defer file.Close()
	// Копируем содержимое файла в ответ
	http.ServeContent(w, r, filePath, fileInfo.ModTime(), file)
}

func MainHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		log.Printf("Invalid URL: %s", r.URL.Path)

		// http.NotFound(w, r)
		ErrorHandler(w, r, 404)
		return
	}
	if r.Method != "GET" {
		ErrorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}

	// log.Printf("Request method: %s", r.Method)
	// log.Printf("Request URL: %s", r.URL.Path)

	data, err := template.ParseFiles("./templates/index.html")
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		// http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		ErrorHandler(w, r, 500)
		return
	}

	err = data.Execute(w, artists)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		ErrorHandler(w, r, 500)
		// http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func ArtistIdHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		ErrorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}
	

	route := regexp.MustCompile(`/artists/(\d+)$`) // проверка на /artist/любая цифра
	if !route.MatchString(r.URL.Path) {
		ErrorHandler(w, r, 404)
		return
	}
	matches := route.FindStringSubmatch(r.URL.Path)
	if matches[1][0] == '0' {

		ErrorHandler(w, r, 404)
		return
	}
	html, err := template.ParseFiles("templates/artist.html")
	if err != nil {
		fmt.Println("Internal server error", err.Error())
		ErrorHandler(w, r, 500)
		return
	}
	if len(matches) < 2 {
		// http.NotFound(w, r)
		ErrorHandler(w, r, 404)
		return
	}
	id, err := strconv.Atoi(matches[1])
	if err != nil {
		// http.NotFound(w, r)
		ErrorHandler(w, r, 404)
		return
	}

	if id <= 0 || id > len(artists) {
		// http.NotFound(w, r)
		ErrorHandler(w, r, 404)
		return
	}

	artist := artists[id-1]

	err = html.Execute(w, artist)
	if err != nil {
		log.Println(err.Error())
		ErrorHandler(w, r, 500)
		return
	}
}

type Error struct {
	Status    int
	ErrString string
}

func ErrorHandler(w http.ResponseWriter, r *http.Request, status int) {
	var err Error
	switch status {
	case 404:
		err.Status = http.StatusNotFound
		err.ErrString = "Not Found"
	case 500:
		err.Status = http.StatusInternalServerError
		err.ErrString = "Internal Server Error"
	case 502:
		err.Status = http.StatusBadGateway
		err.ErrString = "Bad Gateway"
		// http.Error(w, "", 502)
	case http.StatusMethodNotAllowed:
		err.Status = http.StatusMethodNotAllowed

		err.ErrString = "Method not allowed"
		// http.Error(w, "", 405)
	}

	// http.Error(w, "", status)

	html, errParsing := template.ParseFiles("templates/error.html")
	// html.
	if errParsing != nil {
		fmt.Fprint(w, "Internal server error")
		return
	}
	// r.Header.Set()
	// fmt.Println("check")

	w.WriteHeader(status)
	execErr := html.Execute(w, err)
	if execErr != nil {
		// fmt.Println("check")
		log.Println("Error Executing template ", execErr.Error())

		fmt.Fprint(w, "Internal server error")
		return
	}
	// err = html.Execute(w, err)
	// temp, tempError := template.ParseFiles("static/error.html")
	// if tempError != nil {
	// 	log.Println("Error Parsing template")
	// 	http.Error(w, "Internal Server Error", 500)
	// 	return
	// }
	// execError := temp.Execute(w, err)
	// if execError != nil {
	// 	log.Println("Error Executing template")
	// 	http.Error(w, "Internal Server Error", 500)
	// 	return
	// }
}
