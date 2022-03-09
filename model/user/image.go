package user

import (
	"github.com/google/uuid"
	"image"
	"image/png"
	"log"
	"os"
	"path"
	"strings"
)

//Image 画像のUUID.png
type Image string

const ImageDir = "resource/userResource/"

/*
	GetImage
	画像パスを取得
	パスがなければデフォルト画像を表示
*/
func (receiver Image) GetImage() string {
	if receiver == "" {
		return defaultImage()
	}
	open, err := os.Open(string(receiver))
	defer open.Close()
	if err != nil {
		log.Printf("File not found : %s", receiver)
		if receiver == "image.png" {
			return ""
		}
		return defaultImage()
	}
	return string(receiver)
}

/*
	defaultImage
	デフォルト画像パスを渡す
*/
func defaultImage() string {
	var s Image
	s = ImageDir + "image.png"
	return s.GetImage()
}

/*
 */
func (r Image) SaveImage(img image.Image) string {
	//var f *os.File
	f, _ := os.Create(path.Join(ImageDir, strings.Replace(uuid.NewString(), "-", "", -1)+".png"))
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Println("Cant Close File!! : ", f.Name())
			f.Close()
		}
	}(f)
	if !strings.Contains(r.GetImage(), defaultImage()) {
		os.Remove(r.GetImage())
	}
	png.Encode(f, img)
	return f.Name()
}
