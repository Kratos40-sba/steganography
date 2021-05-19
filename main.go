package main

import (
	"archive/zip"
	"bufio"
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"log"
	"math/rand"
	"os"
)

func generateImage(x1, y1 int) *image.RGBA {
	im := image.NewRGBA(image.Rect(0, 0, x1, y1))
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
	defer imageFile.Close()
	err = jpeg.Encode(imageFile, img, nil)
	if err != nil {
		log.Fatal(err)
	}
	/*
		err = imageFile.Close()
			if err != nil {
				log.Fatal(err)
			}
	*/

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
func hideZipFileInImage(zipFile, imageFile *os.File) (*os.File, int64) {
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
	return f, getSize(f)
}
func detectSteganography(steganographyFile *os.File, steganographyFileSize int64) bool {
	//  Zip signature = "\x50\x4b\x03\x04"
	f, err := os.Open(steganographyFile.Name())
	var j int64
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
				j = i
				goto out
			}

		}
	}
out:
	out, _ := os.Create("out.txt")
	br := bufio.NewWriter(out)

	for k := j; k < steganographyFileSize; k++ {
		b, _ := bfReader.ReadByte()
		_ = br.WriteByte(b)
	}
	out.Close()
	fmt.Println(steganographyFileSize, j)
	return true

	return false
}
func getSize(f *os.File) int64 {
	fs, _ := f.Stat()
	return fs.Size()
}
func main() {
	imageFile := encodeImage("test.jpg", generateImage(100, 200))
	zipFile := createCompressedFiles("test.zip")
	steganographyFile, steganographyFileSize := hideZipFileInImage(zipFile, imageFile)
	e := detectSteganography(steganographyFile, steganographyFileSize)

	if e {
		log.Printf("Steganography detected ")
	} else {
		log.Println("No ZIP in the image file ")
	}

}
