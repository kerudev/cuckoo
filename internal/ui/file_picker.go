package ui

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"slices"

	rg "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"

	. "github.com/kerudev/cuckoo/internal/models"
)

var fileScrollIdx = int32(0)
var fileActive = int32(-1)
var fileFocused = int32(-1)

var allowed = []string{".json"}

func DrawFilePicker() {
	// backButtonClicked := rg.Button(NewRectangleFromInt32(Offset.X, Offset.Y, BoxSize, BoxSize), "<")
	// okButtonClicked := rg.Button(NewRectangleFromInt32(S_Screen.Val.W-Offset.X*2, Offset.Y, BoxSize, BoxSize), "OK")

	rg.Button(NewRectangleFromInt32(Offset.X, Offset.Y, BoxSize, BoxSize), "<")
	fileButtonClicked := rg.Button(NewRectangleFromInt32(Offset.X+BoxSize+4, Offset.Y, S_Screen.Val.W-Offset.X*2-BoxSize*2-8, BoxSize), S_FileName.Val)
	rg.Button(NewRectangleFromInt32(S_Screen.Val.W-Offset.X*2, Offset.Y, BoxSize, BoxSize), "OK")

	if fileButtonClicked {
		S_FilePicker.Set(!S_FilePicker.Val)
	}

	if !S_FilePicker.Val {
		return
	}

	dirFiles, err := os.ReadDir(filepath.Dir(S_FileName.Val))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	dirs := []string{}
	data := []string{}

	for _, file := range dirFiles {
		name := file.Name()

		if file.IsDir() {
			dirs = append(dirs, name)
		} else if slices.Contains(allowed, path.Ext(name)) {
			data = append(data, name)
		}
	}

	files := append(dirs, data...)
	fileCount := int32(len(files))

	filePicker := Grid
	filePicker.Height = fileCount * ListViewItemH

	rg.ListViewEx(filePicker.ToFloat32(), files, &fileFocused, &fileScrollIdx, &fileActive)

	filePickerBG := Grid
	filePickerBG.Height = Grid.Height - filePicker.Height
	filePickerBG.Y = filePicker.Y + filePicker.Height

	rg.DrawRectangle(filePickerBG.ToFloat32(), 0, rl.Black, rl.Fade(rl.White, 0.8))
}
