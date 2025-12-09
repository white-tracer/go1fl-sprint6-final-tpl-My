package handlers // Или ваш пакет, например, handlers

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	service "github.com/Yandex-Practicum/go1fl-sprint6-final/internal/service"
)

func MainHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	filePath := "../index.html"

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Printf("file not found: %s", filePath)
		http.Error(w, "file not found", http.StatusNotFound)
		return
	}
	http.ServeFile(w, r, filePath)
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if r.Method != http.MethodPost {
		log.Printf("%v method not allowed", http.StatusMethodNotAllowed)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ct := r.Header.Get("Content-Type")
	log.Printf("UploadHandler: method=%s content-type=%q", r.Method, ct)

	uploadDir := "uploads"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		log.Printf("MkdirAll error: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	// 1) В Headers Content-Type Multipart Form-Data, но пример из урока не сработал. Пробую так
	if strings.HasPrefix(ct, "multipart/") {
		// безопасный лимит в памяти (10 MB)
		_ = r.ParseMultipartForm(10 << 20)

		if r.MultipartForm != nil && len(r.MultipartForm.File) > 0 {
			for _, fhs := range r.MultipartForm.File {
				if len(fhs) == 0 {
					continue
				}
				// берем первый заголовок файла
				hdr := fhs[0]
				in, err := hdr.Open()
				if err != nil {
					log.Printf("hdr.Open error: %v", err)
					continue
				}
				defer in.Close()

				fileContent, err := readContentToBuffer(in)
				if err != nil {
					log.Printf("readContentToBuffer error: %v", err)
					http.Error(w, "ошибка при чтении содержимого файла", http.StatusInternalServerError)
					return
				}
				log.Printf("Содержимое файла (%d байт) сохранено в переменной.", len(fileContent))

				// конвертируем содержимое файла в морзе или стрингу
				convertedFileContent := service.LineConverter(string(fileContent))
				log.Printf("Содержимое файла (%d байт) конвертировано и находится в переменной.", len(convertedFileContent))

				_, err = saveBufferWithUniqueName([]byte(convertedFileContent), uploadDir)
				if err != nil {
					log.Printf("saveBufferWithUniqueName error: %v", err)
					http.Error(w, "ошибка при сохранении файла с уникальным именем", http.StatusInternalServerError)
					return
				}

				fmt.Fprintln(w, convertedFileContent)
				return
			}
		} else {
			// если ничего не нашлось
			http.Error(w, "no file in multipart form", http.StatusBadRequest)
		}
		return
	}

}

// Вспомогательные функциии для хэндлеров

// readContentToBuffer - считывает все содержимое из io.Reader в []byte.
func readContentToBuffer(r io.Reader) ([]byte, error) {
	var buf bytes.Buffer
	_, err := io.Copy(&buf, r)
	if err != nil {
		return nil, fmt.Errorf("ошибка при чтении содержимого в буфер: %w", err)
	}
	return buf.Bytes(), nil
}

// saveBufferWithUniqueName - сохраняет []byte содержимое в файл с уникальным именем (time.Now().UTC().String())
// в указанную директорию.
func saveBufferWithUniqueName(content []byte, uploadDir string) (string, error) {
	uniqueFilename := fmt.Sprintf("%s.txt", time.Now().UTC().Format("20060102_150405.000"))
	dstPath := filepath.Join(uploadDir, uniqueFilename)

	// Создаём файл
	dst, err := os.Create(dstPath)
	if err != nil {
		return "", fmt.Errorf("ошибка при создании файла %s: %w", dstPath, err)
	}
	defer dst.Close()

	// Записываем содержимое из буфера в файл
	_, err = io.Copy(dst, bytes.NewReader(content))
	if err != nil {
		return "", fmt.Errorf("ошибка при записи содержимого в файл %s: %w", dstPath, err)
	}

	return uniqueFilename, nil
}
