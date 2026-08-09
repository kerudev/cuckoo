package ui

import (
	"os"
	"path/filepath"
	"slices"
	"strings"

	rg "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"

	. "github.com/kerudev/cuckoo/internal/models"
	. "github.com/kerudev/cuckoo/internal/utils"
)

func DrawFilePicker() {
	if ShowHelp || !S_FilePicker.Val {
		rg.SetState(rg.STATE_DISABLED)
	}

	backButtonClicked := rg.Button(BackButton, "#118#")

	if ShowHelp || !S_FilePicker.Val {
		rg.SetState(rg.STATE_NORMAL)
	}

	if ShowHelp {
		rg.SetState(rg.STATE_DISABLED)
		defer rg.SetState(rg.STATE_NORMAL)
	}

	fileButtonClicked := rg.Button(FileButton, "")

	if ErrorText != "" {
		errorRec := NewRectangleFromInt32(Offset.X, Offset.Y+BoxSize, S_Screen.Val.W-Offset.X*2, BoxSize)
		rg.DrawRectangle(errorRec, 2, rl.Fade(rl.Red, 0.8), rl.Fade(rl.Red, 0.4))

		errorRec.X += float32(BoxBorder) * 4
		rg.DrawText("#113# "+ErrorText, errorRec, int32(rg.TEXT_ALIGN_LEFT), rl.Black)
	}

	isDir, err := IsDir(S_FileName.Val)
	if err != nil {
		ErrorText = "Path " + S_FileName.Val + " doesn't exist"
		return
	}

	icon := ""
	if isDir {
		icon = "#217#"
	} else {
		icon = "#10#"
	}

	rg.DrawText(
		icon+" "+S_FileName.Val,
		FileButtonText,
		int32(rg.TEXT_ALIGN_LEFT),
		rg.GetStyle(rg.BUTTON, rg.TEXT_COLOR_NORMAL).AsColor(),
	)

	// Toggle file picker
	if fileButtonClicked {
		S_FilePicker.Set(!S_FilePicker.Val)
	}

	// Return early if the file picker is not active
	if !S_FilePicker.Val {
		// Reset S_FileName when File Picker was open previously
		if S_FilePicker.HasChanged() && S_LastFile.Val != "" {
			S_FileName.Set(S_LastFile.Val)
		}

		return
	}

	dirName := ""
	if isDir {
		dirName = S_FileName.Val
	} else {
		dirName = filepath.Dir(S_FileName.Val)
	}

	// Navigate to the previous directory
	if backButtonClicked {
		S_FileName.Set(filepath.Dir(dirName))
		S_FileScroll.Set(-1)
	}

	// Read and show all files in directory
	// TODO read directory only when its timestamp changes
	dirFiles, err := os.ReadDir(dirName)
	if err != nil {
		ErrorText = err.Error()
	}

	dirs := []string{}
	data := []string{}

	for _, file := range dirFiles {
		name := file.Name()

		if file.IsDir() {
			dirs = append(dirs, "#217# "+name)
		} else if slices.Contains(AllowedExt, filepath.Ext(name)) {
			data = append(data, "#10# "+name)
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

	if ShowHelp {
		rg.SetState(rg.STATE_NORMAL)
	}

	// Change file name when the picker is open
	if S_FileScroll.HasChanged() && S_FileScroll.Val >= 0 {
		// #10# file -> #10#, file
		_, newName, _ := strings.Cut(files[S_FileScroll.Val], " ")

		if isDir {
			S_FileName.Set(filepath.Join(S_FileName.Val, newName))
		} else {
			S_FileName.Set(filepath.Join(filepath.Dir(S_FileName.Val), newName))
		}

		if pathIsDir, _ := IsDir(S_FileName.Val); pathIsDir {
			S_FileScroll.Set(-1)
		} else {
			S_LastFile.Set(S_FileName.Val)
		}
	}

	// Draw white background over grid
	filePickerBG := Grid.ToFloat32()
	filePickerBG.Y = filePicker.Y + filePicker.Height
	filePickerBG.Height = float32(Grid.Height) - filePicker.Height

	rg.DrawRectangle(filePickerBG, 0, rl.Black, rl.Fade(rl.White, 0.8))
}
