package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"

	rg "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"

	. "github.com/kerudev/cuckoo/internal/models"
	. "github.com/kerudev/cuckoo/internal/utils"
)

var S_FileScroll = NewState(int32(0))
var S_FileActive = NewState(int32(-1))
var S_FileFocused = NewState(int32(-1))

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

	// Return early if the file picker is not active
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
		} else if slices.Contains(allowed, filepath.Ext(name)) {
			data = append(data, name)
		}
	}

	files := append(dirs, data...)
	count := Clamp(int32(len(files)), 0, 5)

	filePicker := Grid
	filePicker.Height = count * ListViewItemH

	def_LISTVIEW_BORDER_WIDTH := rg.GetStyle(rg.LISTVIEW, rg.BORDER_WIDTH)

	rg.SetStyle(rg.LISTVIEW, rg.BORDER_WIDTH, rg.PropertyValue(GridBorder))
	rg.ListViewEx(filePicker.ToFloat32(), files, &S_FileFocused.Val, &S_FileScroll.Val, &S_FileActive.Val)
	rg.SetStyle(rg.LISTVIEW, rg.BORDER_WIDTH, rg.PropertyValue(def_LISTVIEW_BORDER_WIDTH))

	if S_FileScroll.HasChanged() && S_FileScroll.Val >= 0 {
		S_FileName.Set(filepath.Join(filepath.Dir(S_FileName.Val), files[S_FileScroll.Val]))
	}

	filePickerBG := Grid
	filePickerBG.Y = filePicker.Y + filePicker.Height
	filePickerBG.Height = Grid.Height - filePicker.Height

	rg.DrawRectangle(filePickerBG.ToFloat32(), 0, rl.Black, rl.Fade(rl.White, 0.8))
}
