package main

import (
	"github.com/disintegration/imaging"
	"github.com/rwcarlsen/goexif/exif"
	"image"
	"image/jpeg"
	"log"
	"os"
)

func main() {

	//https://github.com/scottleedavis/go-exif-remove

	path := "C:\\Users\\Алексей\\go\\src\\exif\\1239D0D3-D96A-4CD8-BAC7-0A40EB499BC3.jpg"
	path2 := "C:\\Users\\Алексей\\go\\src\\exif\\i2.jpg"

	var err error

	file, err := os.Open(path)
	if err != nil {
		log.Println("Не удалось открыть файл", err.Error())
		return
	}
	defer file.Close()

	fileExif, err := exif.Decode(file)
	if err != nil {
		log.Println("Ошибка получения EXIF информации", err.Error())
		return
	}

	orientTiff, err := fileExif.Get(exif.Orientation)
	if err != nil || orientTiff == nil {
		log.Println("В EXIF не найдена информация по Orientation")
		return
	}
	orient := orientTiff.String()

	if orient == "1" {
		log.Println("Обработка изображения не требуется")
		return
	}

	fileImg, err := os.Open(path)
	if err != nil {
		log.Println("Не удалось открыть файл (второй раз)", err.Error())
		return
	}
	defer fileImg.Close()

	img, err := jpeg.Decode(fileImg)
	if err != nil {
		log.Println("Ошибка получения изображения", err.Error())
		return
	}

	img = reverseOrientation(img, orient)

	err = imaging.Save(img, path2)
	if err != nil {
		log.Println("Ошибка сохранения", err.Error())
		return
	}
}

func reverseOrientation(img image.Image, o string) *image.NRGBA {

	switch o {
	case "1":
		return imaging.Clone(img)
	case "2":
		return imaging.FlipV(img)
	case "3":
		return imaging.Rotate180(img)
	case "4":
		return imaging.Rotate180(imaging.FlipV(img))
	case "5":
		return imaging.Rotate270(imaging.FlipV(img))
	case "6":
		return imaging.Rotate270(img)
	case "7":
		return imaging.Rotate90(imaging.FlipV(img))
	case "8":
		return imaging.Rotate90(img)
	}

	return imaging.Clone(img)
}
