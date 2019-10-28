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

	var err error

	fileLog, err := os.OpenFile("log.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer fileLog.Close()

	log.SetOutput(fileLog)

	if len(os.Args) != 2 {
		log.Println("Не передан файл для обработки")
		os.Exit(1)
	}

	path := os.Args[1]
	pathTo := "imgexif\\" + filepath.Base(path)

	file, err := os.Open(path)
	if err != nil {
		log.Println("Не удалось открыть файл", path, err.Error())
		os.Exit(1)
	}
	defer file.Close()

	fileExif, err := exif.Decode(file)
	if err != nil {
		log.Println("Ошибка получения EXIF информации", path, err.Error())
		os.Exit(1)
	}

	orientTiff, err := fileExif.Get(exif.Orientation)
	if err != nil || orientTiff == nil {
		log.Println("В EXIF не найдена информация по Orientation", path)
		os.Exit(1)
	}
	orient := orientTiff.String()

	if orient == "1" {
		log.Println("Обработка изображения не требуется", path)
		os.Exit(1)
	}

	fileImg, err := os.Open(path)
	if err != nil {
		log.Println("Не удалось открыть файл (второй раз)", path, err.Error())
		os.Exit(1)
	}
	defer fileImg.Close()

	img, err := jpeg.Decode(fileImg)
	if err != nil {
		log.Println("Ошибка получения изображения", path, err.Error())
		os.Exit(1)
	}

	img = reverseOrientation(img, orient)

	err = imaging.Save(img, pathTo)
	if err != nil {
		log.Println("Ошибка сохранения", path, err.Error())
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
