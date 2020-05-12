package server

import (
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/nfnt/resize"
)

// FH 文件服务
type FH struct {
	Size int
}

const defaultSize = 400

// ServerHTTP http server
func (*FH) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fileName := r.URL.Path
	size := r.URL.Query().Get("w")
	sizeInt, err := strconv.Atoi(size)
	if size != "" && err != nil || sizeInt < 0 {
		log.Println("无效参数w", sizeInt)
		writeResponse(w, "无效参数", 500)
		return
	}
	sp := C.StoragePath
	if sp[len(sp)-1] != '/' {
		sp = sp + "/"
	}
	width, err := getImageWidth(sp + fileName)
	if err != nil {
		log.Println("没有找到有效的图片", err)
		writeResponse(w, "没有找到有效的图片", 404)
		return
	}

	if size != "" && width > sizeInt && sizeInt > 0 {
		cacheFile := sp + size + "/" + fileName
		_, err = os.Open(cacheFile)
		if err != nil {
			// should resize
			_, err := os.OpenFile(sp+size+"/", os.O_RDONLY, os.ModeDir)
			if err != nil {
				os.Mkdir(sp+size+"/", 0755)
			}
			ff, err := os.OpenFile(cacheFile, os.O_CREATE|os.O_RDWR, 0666)
			if err != nil {
				log.Println("无法创建缓存文件", err)
				writeResponse(w, "无法创建缓存文件", 500)
				return
			}
			f, err := os.Open(sp + fileName)
			if err != nil {
				log.Println("没有找到文件", err)
				writeResponse(w, "没有找到文件", 404)
				return
			}
			defer f.Close()
			defer ff.Close()
			fargs := strings.Split(fileName, ".")
			ft := "." + fargs[len(fargs)-1]
			switch strings.ToLower(ft) {
			case ".jpeg", ".jpg":
				img, err := jpeg.Decode(f)
				if err != nil {
					log.Println("图片格式无效", err)
					writeResponse(w, "图片格式无效", 500)
					return
				}
				rs := resize.Resize(uint(sizeInt), 0, img, resize.Lanczos2)
				jpeg.Encode(ff, rs, nil)
			case ".png":
				img, err := png.Decode(f)
				if err != nil {
					log.Println("图片格式无效", err)
					writeResponse(w, "图片格式无效", 500)
					return
				}
				rs := resize.Resize(uint(sizeInt), 0, img, resize.Lanczos2)
				png.Encode(ff, rs)
			case ".gif":
				img, err := gif.Decode(f)
				if err != nil {
					log.Println("图片格式无效", err)
					writeResponse(w, "图片格式无效", 500)
					return
				}
				rs := resize.Resize(uint(sizeInt), 0, img, resize.Lanczos2)
				gif.Encode(ff, rs, nil)
			default:
				log.Println("图片格式无效", ft, err)
				writeResponse(w, "图片格式无效", 500)
				return
			}
		}
		writeFile(w, cacheFile)
		return
	}
	writeFile(w, sp+fileName)
}

func fileAutoHandler() http.Handler {
	return &FH{defaultSize}
}

func writeResponse(w http.ResponseWriter, msg string, statusCode int) {
	w.WriteHeader(401)
	w.Write([]byte(msg))
	return
}

func writeFile(w http.ResponseWriter, filePath string) {
	ff, er := os.Open(filePath)
	if er != nil {
		log.Println("没有找到文件", er)
		writeResponse(w, "没有找到文件", 401)
		return
	}
	defer ff.Close()
	stat, _ := ff.Stat()
	w.Header().Set("Content-Length", strconv.FormatInt(stat.Size(), 10))
	io.Copy(w, ff)
}

// Run run http server
func Run() {
	svr := http.NewServeMux()
	svr.Handle("/files/", http.StripPrefix("/files/", fileAutoHandler()))
	svr.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		f, h, err := r.FormFile("file")
		if err != nil {
			log.Println("没有找到上传的文件", err)
			writeResponse(w, "没有找到上传的文件", 500)
		}
		u := uuid.New()

		hs := strings.Split(h.Filename, ".")
		if len(hs) < 2 {
			log.Println("没有找到上传的文件", err)
			writeResponse(w, "无法识别文件类型", 500)
		}
		saveFileName := u.String() + "." + hs[len(hs)-1]
		file, _ := os.Create("storage/" + saveFileName)

		io.Copy(file, f)
		w.Write([]byte(saveFileName))
		defer file.Close()
		defer f.Close()
	})
	http.ListenAndServe(":"+C.Port, svr)
}
