package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
)

func combinZipFile() {
	url_zip := "https://qunoss1.qun100.com/85749309592576001/1668986149975/%E3%80%9020%E5%8F%B7%E6%A0%B8%E9%85%B8%E6%A3%80%E6%B5%8B%E7%BB%9F%E8%AE%A1%E3%80%91%E6%96%87%E4%BB%B6%E5%AF%BC%E5%87%BA.zip?t=1668986152235"
	zipFileName := downFile_main(url_zip)
	println(zipFileName)
	unzipmain(zipFileName, "output")
}
func GetImageObj(filePath string) (img image.Image, err error) {
	f1Src, err := os.Open(filePath)

	if err != nil {
		return nil, err
	}
	defer f1Src.Close()

	buff := make([]byte, 512) // why 512 bytes ? see http://golang.org/pkg/net/http/#DetectContentType
	_, err = f1Src.Read(buff)

	if err != nil {
		return nil, err
	}

	filetype := http.DetectContentType(buff)

	fmt.Println(filetype)

	fSrc, err := os.Open(filePath)
	defer fSrc.Close()

	switch filetype {
	case "image/jpeg", "image/jpg":
		img, err = jpeg.Decode(fSrc)
		if err != nil {
			fmt.Println("jpeg error")
			return nil, err
		}

	case "image/gif":
		img, err = gif.Decode(fSrc)
		if err != nil {
			return nil, err
		}

	case "image/png":
		img, err = png.Decode(fSrc)
		if err != nil {
			return nil, err
		}
	default:
		return nil, err
	}
	return img, nil
}
func MergeImage(worksDir string, fileList []fs.FileInfo, newName string) (string, error) {

	imageList := []image.Image{}
	newWidth := 0
	newHeight := 0
	for _, fi := range fileList {
		src, err := GetImageObj(worksDir + fi.Name())
		if err != nil {
			return "", err
		}

		srcB := src.Bounds().Max
		fmt.Println(fi.Name(), newWidth, newHeight, srcB.X, srcB.Y)
		newWidth += srcB.X
		if srcB.Y > newHeight {
			newHeight = srcB.Y
		}

		imageList = append(imageList, src)
	}

	des := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight)) // 底板

	lastX := 0
	for _, src := range imageList {

		srcB := src.Bounds().Max
		r := des.Bounds().Add(image.Pt(lastX, 0))
		draw.Draw(des, r, src, src.Bounds().Min, draw.Src) //将另外一张图片信息存入jpg

		lastX += srcB.X

	}
	fSave, err := os.Create(newName)
	if err != nil {
		return "", err
	}

	defer fSave.Close()

	var opt jpeg.Options
	opt.Quality = 80

	//newImage := resize.Resize(1024, 0, des, resize.Lanczos3)
	newImage := des
	err = jpeg.Encode(fSave, newImage, &opt) // put quality to 80%
	if err != nil {
		return "", err
	}
	return newName, nil
}
func combinPIC(rootFolder string, destFolder string) {
	os.MkdirAll(destFolder, os.ModePerm)
	dir, err := ioutil.ReadDir(rootFolder)
	if err != nil {
		panic(err)
	}
	if len(dir) == 1 {
		println("only one")
		if !dir[0].IsDir() {
			panic("don't found files")
		}
		println(dir[0].Name())
	}
	PthSep := string(os.PathSeparator)
	newDirs := rootFolder + PthSep + dir[0].Name()
	dir, err = ioutil.ReadDir(newDirs)
	if err != nil {
		panic(err)
	}
	for _, fi := range dir {
		if !fi.IsDir() {
			continue
		}
		a := newDirs + PthSep + fi.Name()
		println(a)
		files, err := ioutil.ReadDir(a)
		if err != nil {
			panic(err)
		}
		MergeImage(a+PthSep, files, destFolder+PthSep+fi.Name()+"combined.jpeg")
		for _, fi2 := range files {
			b := a + PthSep + fi2.Name()
			println(b)
		}

	}
}
func downloadExcel() {
	fullURLFile_excel := "https://qunoss3.qun100.com/85774500838248448/1668992155758/%E6%95%B0%E6%8D%AE%E7%BB%9F%E8%AE%A1%E6%8A%A5%E8%A1%A8%282022.11.19-2022.11.21%29.xlsx?t=1668992156127"

	println(fullURLFile_excel)

	fileName := downFile_main(fullURLFile_excel)
	println(fileName)
}
func main() {
	//combinZipFile()

	combinPIC("output", "combined")
}
