package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	rg "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"

	. "github.com/kerudev/cuckoo/internal/shared"
	. "github.com/kerudev/cuckoo/internal/utils"
)

var EmptyPickerMessage = fmt.Sprintf(`Nothing here!
	Click to go back or add contents to this directory (it will update automatically)
	Allowed file extensions: %s`, strings.Join(AllowedExt, ", "),
)

func DrawFilePicker() {
	if ShowHelp || !S_FilePicker.Val {
		rg.SetState(rg.STATE_DISABLED)
	}

	// Icon: ARROW_LEFT_FILL
	backButtonClicked := rg.Button(BackButton, "#118#")

	if ShowHelp || !S_FilePicker.Val {
		rg.SetState(rg.STATE_NORMAL)
	}

	if ShowHelp {
		rg.SetState(rg.STATE_DISABLED)
		defer rg.SetState(rg.STATE_NORMAL)
	}

	isDir, err := IsDir(S_FilePath.Val)
	if err != nil {
		ErrorText = "Path doesn't exist"
		return
	}

	var icon string
	if isDir {
		// Icon: FOLDER
		icon = "#217#"
	} else {
		// Icon: FILETYPE_TEXT
		icon = "#10#"
	}

	// Toggle file picker
	if rg.Button(FileButton, "") {
		S_FilePicker.Set(!S_FilePicker.Val)
	}

	rg.DrawText(
		icon+" "+S_FilePath.Val,
		FileButtonText,
		int32(rg.TEXT_ALIGN_LEFT),
		rg.GetStyle(rg.BUTTON, rg.TEXT_COLOR_NORMAL).AsColor(),
	)

	// Warn the user when the file has changed
	S_FileLastUpdate.Set(GetUnix(S_FileName.Val))

	if S_FileLastUpdate.HasChanged() && S_FileLastUpdate.Old > 0 {
		ErrorText = "The contents of have changed"
	}

	// Return early if the file picker is not active
	if !S_FilePicker.Val {
		// Reset S_FileName when File Picker was open previously
		if S_FilePicker.HasChanged() && S_FileName.Val != "" {
			S_FilePath.Set(S_FileName.Val)
		}

		return
	}

	var dirName string
	if isDir {
		dirName = S_FilePath.Val
	} else {
		dirName = filepath.Dir(S_FilePath.Val)
	}

	S_DirLastUpdate.Set(GetUnix(dirName))

	// Navigate to the previous directory
	if backButtonClicked {
		S_FilePath.Set(filepath.Dir(dirName))
		S_FileScroll.Set(-1)
	}

	// Read the directory only when it has changes
	if S_DirLastUpdate.HasChanged() {
		dirContent, err := os.ReadDir(dirName)
		if err != nil {
			ErrorText = err.Error()
		}

		var dirs []string
		var data []string
		for _, path := range dirContent {
			name := path.Name()

			var icon string
			if path.IsDir() {
				if content, _ := os.ReadDir(name); len(content) > 0 {
					// Icon: FOLDER_FILE_OPEN
					icon = "#1#"
				} else {
					// Icon: FOLDER
					icon = "#217#"
				}

				dirs = append(dirs, icon+" "+name)
			} else if slices.Contains(AllowedExt, filepath.Ext(name)) {
				if info, _ := path.Info(); info.Size() > 0 {
					// Icon: FILETYPE_TEXT
					icon = "#10#"
				} else {
					// Icon: FILE
					icon = "#218#"
				}

				data = append(data, icon+" "+name)
			}
		}

		DirFiles = append(dirs, data...)
		DirFilesCount = Clamp(int32(len(DirFiles)), MIN_FILES, MAX_FILES)
	}

	filePicker := Grid.ToFloat32()

	if DirFilesCount == 0 {
		// Draw warning message button when there are no files in the directory
		filePicker.Height = float32(ListViewItemH) * 1.5

		rg.DrawRectangle(filePicker, 2, rl.Black, rl.RayWhite)
		rg.DrawText(EmptyPickerMessage, filePicker, int32(rg.TEXT_ALIGN_CENTER), rl.Black)
	} else {
		// Draw file picker
		filePicker.Height = float32(ListViewItemH * DirFilesCount)

		// Draw file and dir list
		rg.SetStyle(rg.LISTVIEW, rg.BORDER_WIDTH, rg.PropertyValue(GridBorder))
		rg.ListViewEx(filePicker, DirFiles, &S_FileFocused.Val, &S_FileScroll.Val, &S_FileActive.Val)
		rg.SetStyle(rg.LISTVIEW, rg.BORDER_WIDTH, Style["LISTVIEW_BORDER_WIDTH"])

		if ShowHelp {
			rg.SetState(rg.STATE_NORMAL)
		}

		// Change file name when the picker is open
		if S_FileScroll.HasChanged() && S_FileScroll.Val >= 0 {
			// Extract name: #10# file -> #10#, file
			_, newName, _ := strings.Cut(DirFiles[S_FileScroll.Val], " ")

			if isDir {
				// Append name when the selected path is a dir
				S_FilePath.Set(filepath.Join(S_FilePath.Val, newName))
			} else {
				// Append name to the parent path when the selected path is a file
				S_FilePath.Set(filepath.Join(filepath.Dir(S_FilePath.Val), newName))
			}

			if pathIsDir, _ := IsDir(S_FilePath.Val); pathIsDir {
				S_FileScroll.Set(-1)
			} else {
				S_FileName.Set(S_FilePath.Val)
				S_FileLastUpdate.Set(GetUnix(S_FileName.Val))
			}
		}
	}

	// Draw white background over grid
	filePickerBG := Grid.ToFloat32()
	filePickerBG.Y = filePicker.Y + filePicker.Height
	filePickerBG.Height = float32(Grid.Height) - filePicker.Height

	rg.DrawRectangle(filePickerBG, 0, rl.Blank, rl.Fade(rl.White, 0.8))
}
