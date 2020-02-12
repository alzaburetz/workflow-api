package handlers

import ("net/http"
		"os"
		"io"
		"github.com/satori/go.uuid")

func UploadFile(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)
	file, _, err := r.FormFile("File")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		WriteAnswer(&w, "", []string{"Couldn't get file from form", err.Error()}, 400)
		return
	} 
	defer file.Close()

	// if typ := http.DetectContentType(file); typ != "image/jpg" {
	// 	WriteAnswer(&w, "", []string{"Only jpg is allowed"}, 400)
	// 	return
	// }
	fileName := uuid.NewV4()
	f, err := os.OpenFile("./static/" + fileName.String() + ".png", os.O_WRONLY | os.O_CREATE, 0666)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		WriteAnswer(&w, "", []string{"Couldn't create file", err.Error()}, 500)
		return
	}
	defer f.Close()
	io.Copy(f, file)
}