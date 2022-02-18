package user

import (
	"github.com/google/uuid"
	"image"
	"image/png"
	"log"
	"os"
	"strings"
	"time"
)

//Image 画像のUUID.png
type Image string

/*
	GetImage
	画像パスを取得
	パスがなければデフォルト画像を表示
*/
func (receiver Image) GetImage() string {
	if receiver == "" {
		return defaultImage()
	}
	open, err := os.Open("resource/userResource/" + string(receiver))
	defer open.Close()
	if err != nil {
		log.Printf("File not found : %s", receiver)
		if receiver == "image.png" {return ""}
		return defaultImage()
	}
	return "resource/userResource/" + string(receiver)
}

/*
	defaultImage
	デフォルト画像パスを渡す
*/
func defaultImage() string {
	var s Image
	s = "image.png"
	return s.GetImage()
}

/*
 */
func SaveImage(img image.Image) string {
	var f *os.File
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Println("Cant Close File!! : ", f.Name())
			f.Close()
		}
	}(f)
	for {
		var err error
		fname := strings.ReplaceAll(uuid.NewString(), "-", "")
		f, err = os.Create(fname + ".png")
		if err == nil {
			break
		}
		time.Sleep(time.Second * 20)
	}
	err := png.Encode(f, img)
	if err != nil {
		log.Println("Cannot Create Image!! : ", f.Name())
	}
	return f.Name()
}
