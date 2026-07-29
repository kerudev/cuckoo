package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

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
	backButtonClicked := rg.Button(NewRectangleFromInt32(Offset.X, Offset.Y, BoxSize, BoxSize), "<")

	fileButton := NewRectangleFromInt32(Offset.X+BoxPad, Offset.Y, S_Screen.Val.W-Offset.X*2-BoxPad, BoxSize)
	fileButtonClicked := rg.Button(fileButton, "")

	fileButton.X += float32(BoxSize) / 2

	rg.DrawText(S_FileName.Val, fileButton, int32(rg.TEXT_ALIGN_LEFT), rg.GetStyle(rg.BUTTON, rg.TEXT_COLOR_NORMAL).AsColor())

	// Toggle file picker
	if fileButtonClicked {
		S_FilePicker.Set(!S_FilePicker.Val)
	}

	// Return early if the file picker is not active
	if !S_FilePicker.Val {
		S_FileName.Set(S_LastFile.Val)
		return
	}

	dirName := ""
	stat, _ := os.Stat(S_FileName.Val)

	if stat.IsDir() {
		dirName = S_FileName.Val
	} else {
		dirName = filepath.Dir(S_FileName.Val)
	}

	// Navigate to the previous directory
	if backButtonClicked {
		S_FileName.Set(filepath.Dir(dirName))
	}

	// Read and show all files in directory
	// TODO read directory only when its timestamp changes
	dirFiles, err := os.ReadDir(dirName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	dirs := []string{}
	data := []string{}

	for _, file := range dirFiles {
		name := file.Name()

		if file.IsDir() {
			dirs = append(dirs, fmt.Sprintf("#1# %s", name))
		} else if slices.Contains(allowed, filepath.Ext(name)) {
			data = append(data, fmt.Sprintf("#10# %s", name))
		}
	}

	files := append(dirs, data...)
	count := Clamp(int32(len(files)), 0, 6)

	filePicker := Grid.ToFloat32()
	filePicker.Height = float32(count * ListViewItemH)

	def_LISTVIEW_BORDER_WIDTH := rg.GetStyle(rg.LISTVIEW, rg.BORDER_WIDTH)

	rg.SetStyle(rg.LISTVIEW, rg.BORDER_WIDTH, rg.PropertyValue(GridBorder))
	rg.ListViewEx(filePicker, files, &S_FileFocused.Val, &S_FileScroll.Val, &S_FileActive.Val)
	rg.SetStyle(rg.LISTVIEW, rg.BORDER_WIDTH, def_LISTVIEW_BORDER_WIDTH)

	// Change file name when the picker is open
	if S_FileScroll.HasChanged() && S_FileScroll.Val >= 0 {
		// #10# file -> #10#, file
		_, newName, _ := strings.Cut(files[S_FileScroll.Val], " ")
		newPath := ""

		if stat.IsDir() {
			newPath = filepath.Join(S_FileName.Val, newName)
		} else {
			newPath = filepath.Join(filepath.Dir(S_FileName.Val), newName)
			S_LastFile.Set(newPath)
		}

		S_FileName.Set(newPath)
	}

	// Draw white background over grid
	filePickerBG := Grid.ToFloat32()
	filePickerBG.Y = filePicker.Y + filePicker.Height
	filePickerBG.Height = float32(Grid.Height) - filePicker.Height

	rg.DrawRectangle(filePickerBG, 0, rl.Black, rl.Fade(rl.White, 0.8))
}
