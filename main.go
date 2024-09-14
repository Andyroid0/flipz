package main

import (
	"fmt"
	"image/png"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	images     []*ebiten.Image
	frameCount int
	currentImg int
}

var imageDirectory string

var framesPerImage = 2
var frameRate = 24

func (g *Game) Update() error {
	g.frameCount++

	if g.frameCount >= framesPerImage {
		g.frameCount = 0
		g.currentImg++
		if g.currentImg >= len(g.images) {
			g.currentImg = 0
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Clear()
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(g.images[g.currentImg], op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 800, 600
}

func LoadImages(directory string) ([]*ebiten.Image, error) {
	var images []*ebiten.Image
	err := filepath.WalkDir(directory, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(path) == ".png" {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			img, err := png.Decode(file)
			if err != nil {
				return err
			}

			ebitenImg := ebiten.NewImageFromImage(img)
			images = append(images, ebitenImg)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return images, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: flipz <image-directory> <frames-per-send> <on-2s-y-or-n>")
		return
	}

	imageDirectory = os.Args[1]
	if len(os.Args) > 2 {
		num, frErr := strconv.Atoi(os.Args[2])
		if frErr == nil {
			frameRate = num
		}
	}
	if len(os.Args) > 3 {
		if os.Args[3] == "y" || os.Args[3] == "Y" || os.Args[3] == "Yes" || os.Args[3] == "yes" {
			framesPerImage = 2
		} else {
			framesPerImage = 1
		}
	}

	images, err := LoadImages(imageDirectory)
	if err != nil {
		log.Fatal(err)
	}
	if len(images) == 0 {
		log.Fatal("No images found in the directory.")
	}

	game := &Game{
		images:     images,
		frameCount: 0,
		currentImg: 0,
	}

	ebiten.SetTPS(frameRate)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
