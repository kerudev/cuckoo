package ui

import (
	rg "github.com/gen2brain/raylib-go/raygui"

	. "github.com/kerudev/cuckoo/internal/models"
)

func DrawLockButton() {
	icon := ""
	if S_IsMouseLocked.Val {
		// Icon: LOCK_CLOSE
		icon = "#137#"
		rg.SetState(rg.STATE_NORMAL)
	} else {
		// Icon: LOCK_OPEN
		icon = "#138#"
		rg.SetState(rg.STATE_DISABLED)
		defer rg.SetState(rg.STATE_NORMAL)
	}

	if rg.Button(LockButton, icon) {
		// Set to false as the button can only be pressed when the state is STATE_DISABLED
		S_IsMouseLocked.Set(false)
	}
}
