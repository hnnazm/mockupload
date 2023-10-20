package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "DELETE, POST, GET, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
}

func handler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	image, _, err := r.FormFile("image")

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	defer image.Close()

	savedImage, err := os.Create("image.jpg")

  defer savedImage.Close()

	_, err = io.Copy(savedImage, image)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

  w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
  w.Write([]byte(`{"success":1,"file":{"url":"http://localhost:8081/image.jpg"}}`))
}

func main() {
	http.HandleFunc("/upload", handler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
