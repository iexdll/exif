package main

import (
	"github.com/disintegration/imaging"
	"github.com/rwcarlsen/goexif/exif"
	"image"
	"image/jpeg"
	"log"
	"os"
	"path/filepath"
)

func main() {

	//https://github.com/scottleedavis/go-exif-remove
	//go run main.go img/jpg/11-tests.jpg

	if len(os.Args) != 2 {
		log.Println("Не передан файл для обработки")
		os.Exit(1)
	}

	path := os.Args[1]
	pathTo := "imgexif\\" + filepath.Base(path)

	var err error

	file, err := os.Open(path)
	if err != nil {
		log.Println("Не удалось открыть файл", err.Error())
		os.Exit(1)
	}
	defer file.Close()

	fileExif, err := exif.Decode(file)
	if err != nil {
		log.Println("Ошибка получения EXIF информации", err.Error())
		os.Exit(1)
	}

	orientTiff, err := fileExif.Get(exif.Orientation)
	if err != nil || orientTiff == nil {
		log.Println("В EXIF не найдена информация по Orientation")
		os.Exit(1)
	}
	orient := orientTiff.String()

	if orient == "1" {
		log.Println("Обработка изображения не требуется")
		os.Exit(1)
	}

	fileImg, err := os.Open(path)
	if err != nil {
		log.Println("Не удалось открыть файл (второй раз)", err.Error())
		os.Exit(1)
	}
	defer fileImg.Close()

	img, err := jpeg.Decode(fileImg)
	if err != nil {
		log.Println("Ошибка получения изображения", err.Error())
		os.Exit(1)
	}

	img = reverseOrientation(img, orient)

	err = imaging.Save(img, pathTo)
	if err != nil {
		log.Println("Ошибка сохранения", err.Error())
		os.Exit(1)
	}

	os.Exit(0)
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
