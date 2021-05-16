package main

import (
	"archive/zip"
	"bufio"
	"bytes"
	"image"
	"image/jpeg"
	"io"
	"log"
	"math/rand"
	"os"
)

func generateImage(x1, y1 int) *image.RGBA {
	im := image.NewRGBA(image.Rect(0, 0, 100, 200))
	for pix := 0; pix < x1*y1; pix++ {
		pixelOffset := pix * 4
		im.Pix[0+pixelOffset] = uint8(rand.Intn(256))
		im.Pix[1+pixelOffset] = uint8(rand.Intn(256))
		im.Pix[2+pixelOffset] = uint8(rand.Intn(256))
		im.Pix[3+pixelOffset] = 255
	}
	return im
}
func encodeImage(fileName string, img *image.RGBA) *os.File {
	imageFile, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	err = jpeg.Encode(imageFile, img, nil)
	if err != nil {
		log.Fatal(err)
	}
	err = imageFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	return imageFile

}
func createCompressedFiles(fileName string) *os.File {
	zipFile, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer zipFile.Close()
	zipWriter := zip.NewWriter(zipFile)
	var fichiers = []struct {
		nom, cont string
	}{
		{"fichier1.txt", "fichier 1 "},
		{"fichier2.txt", "fichier 2 "},
		{"fichier3.txt", "fichier 3 "},
	}
	for _, fichier := range fichiers {
		fz, err := zipWriter.Create(fichier.nom)
		if err != nil {
			log.Fatal(err)
		}
		_, err = fz.Write([]byte(fichier.cont))
		if err != nil {
			log.Fatal(err)
		}

	}
	err = zipWriter.Close()
	if err != nil {
		log.Fatal(err)
	}

	return zipFile
}
func hideZipFileInImage(zipFile, imageFile *os.File) *os.File {
	f1, err := os.Open(imageFile.Name())
	if err != nil {
		log.Fatal(err)
	}
	defer f1.Close()
	f2, err := os.Open(zipFile.Name())
	if err != nil {
		log.Fatal(err)
	}
	defer f2.Close()
	f, err := os.Create("steganography.jpg")
	if err != nil {
		log.Fatal(err)
	}
	_, err = io.Copy(f, f1)
	if err != nil {
		log.Fatal(err)
	}
	_, err = io.Copy(f, f2)
	if err != nil {
		log.Fatal(err)
	}

	return f
}
func detectSteganography(imageFile *os.File) bool {
	//  Zip signature = "\x50\x4b\x03\x04"
	f, err := os.Open(imageFile.Name())
	if err != nil {
		log.Fatal(err)
	}
	bfReader := bufio.NewReader(f)
	fs, _ := f.Stat()
	for i := int64(0); i < fs.Size(); i++ {
		b, err := bfReader.ReadByte()
		if err != nil {
			log.Fatal(err)
		}
		if b == '\x50' {
			bs := make([]byte, 3)
			bs, err = bfReader.Peek(3)
			if err != nil {
				log.Fatal(err)
			}
			if bytes.Equal(bs, []byte{'\x4b', '\x03', '\x04'}) {
				return true
			}

		}
	}
	return false
}
func main() {
	imageFile := encodeImage("test.jpg", generateImage(100, 200))
	zipFile := createCompressedFiles("test.zip")
	steganographyFile := hideZipFileInImage(zipFile, imageFile)
	if detectSteganography(steganographyFile) {
		log.Println("Steganography detected ")
	} else {
		log.Println("No ZIP in the image file ")
	}
}
