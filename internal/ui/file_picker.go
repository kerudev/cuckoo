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
	if ShowHelp || S_PickerIsOn.Eq(false) {
		rg.SetState(rg.STATE_DISABLED)
	}

	isDir, err := IsDir(S_PickerPath.Val)
	if err != nil {
		ErrorText = "Selected path doesn't exist (maybe it got deleted or moved?)"
	}

	S_PathExists.Set(err == nil)

	// Navigate to the previous directory
	var dirName string
	if isDir {
		dirName = S_PickerPath.Val
	} else {
		dirName = filepath.Dir(S_PickerPath.Val)
	}

	// Icon: ARROW_LEFT_FILL
	if rg.Button(BackButton, "#118#") {
		S_PickerPath.Set(filepath.Dir(dirName))
		S_FileScroll.Set(-1)
	}

	if ShowHelp || S_PickerIsOn.Eq(false) {
		rg.SetState(rg.STATE_NORMAL)
	}

	if ShowHelp {
		rg.SetState(rg.STATE_DISABLED)
		defer rg.SetState(rg.STATE_NORMAL)
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
		S_PickerIsOn.Set(!S_PickerIsOn.Val)
	}

	rg.DrawText(
		icon,
		FileButtonIcon,
		int32(rg.TEXT_ALIGN_LEFT),
		rg.GetStyle(rg.BUTTON, rg.TEXT_COLOR_NORMAL).AsColor(),
	)

	fmt.Println("===========================")
	fmt.Println(rl.MeasureText(S_PickerPath.Val, FontSize))
	fmt.Println(int32(FileButtonText.Width))

	// TODO draw text with ... on the start when text is too long
	rg.DrawText(
		S_PickerPath.Val,
		FileButtonText,
		int32(rg.TEXT_ALIGN_LEFT),
		rg.GetStyle(rg.BUTTON, rg.TEXT_COLOR_NORMAL).AsColor(),
	)

	// Warn the user when the file has changed
	if S_PathExists.Eq(true) {
		S_FileLastUpdate.Set(GetUnix(S_FilePath.Val))

		if S_FileLastUpdate.HasChanged() && S_FileLastUpdate.Old > 0 {
			ErrorText = "The file has been updated. Select it again if you need to reload data."
		}
	}

	// Return early if the file picker is not active
	if S_PickerIsOn.Eq(false) {
		// Reset S_FileName when File Picker was open previously
		if S_PickerIsOn.HasChanged() && S_FilePath.Not("") {
			S_PickerPath.Set(S_FilePath.Val)
		}

		return
	}

	filePicker := Grid.ToFloat32()

	if DirFilesCount == 0 || S_DirLastUpdate.Eq(0) {
		// Draw warning message button when there are no files in the directory
		filePicker.Height = float32(ListViewItemH)*1.5 + 4

		rg.DrawRectangle(filePicker, 2, rl.Black, rl.RayWhite)
		rg.DrawText(EmptyPickerMessage, filePicker, int32(rg.TEXT_ALIGN_CENTER), rl.Black)
	} else {
		// Draw file picker
		filePicker.Height = float32(ListViewItemH*DirFilesCount) + 4

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
				S_PickerPath.Set(filepath.Join(S_PickerPath.Val, newName))
			} else {
				// Append name to the parent path when the selected path is a file
				S_PickerPath.Set(filepath.Join(filepath.Dir(S_PickerPath.Val), newName))
			}

			pathIsDir, _ := IsDir(S_PickerPath.Val)
			if pathIsDir {
				// Change S_FileScroll so nothing is selected by default
				S_FileScroll.Set(-1)
			} else {
				// If the path is a dir, change S_FileScroll so nothing is selected by default
				S_FilePath.Set(S_PickerPath.Val)
				S_FileLastUpdate.Set(GetUnix(S_FilePath.Val))
			}
		}
	}

	// Update dirName as the path may have changed
	isDir, _ = IsDir(S_PickerPath.Val)
	if isDir {
		dirName = S_PickerPath.Val
	} else {
		dirName = filepath.Dir(S_PickerPath.Val)
	}

	// Get last update constantly
	S_DirLastUpdate.Set(GetUnix(dirName))

	// Read the directory only when it has changes
	if S_PathExists.HasChanged() || S_DirLastUpdate.HasChanged() || S_PickerPath.HasChanged() || S_PickerIsOn.HasChanged() {
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

	// Draw white background over grid
	filePickerBG := Grid.ToFloat32()
	filePickerBG.Y = filePicker.Y + filePicker.Height
	filePickerBG.Height = float32(Grid.Height) - filePicker.Height

	rg.DrawRectangle(filePickerBG, 0, rl.Blank, rl.Fade(rl.White, 0.8))
}
