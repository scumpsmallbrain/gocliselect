package gocliselect

import (
	"fmt"

	kb "atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
	col "github.com/gookit/color"
)


type Menu struct {
	Prompt  	string
	CursorPos 	int
	MenuItems 	[]*MenuItem
}

type MenuItem struct {
	Text     string
	ID       string
	SubMenu  *Menu
}

func NewMenu(prompt string) *Menu {
	return &Menu{
		Prompt: prompt,
		MenuItems: make([]*MenuItem, 0),
	}
}

// AddItem will add a new menu option to the menu list
func (m *Menu) AddItem(option string, id string) *Menu {
	menuItem := &MenuItem{
		Text: option,
		ID: id,
	}

	m.MenuItems = append(m.MenuItems, menuItem)
	return m
}

// renderMenuItems prints the menu item list.
// Setting redraw to true will re-render the options list with updated current selection.
func (m *Menu) renderMenuItems(redraw bool) {
	if redraw {
		// Move the cursor up n lines where n is the number of options, setting the new
		// location to start printing from, effectively redrawing the option list
		//
		// This is done by sending a VT100 escape code to the terminal
		// @see http://www.climagic.org/mirrors/VT100_Escape_Codes.html
		fmt.Printf("\033[%dA", len(m.MenuItems) -1)
	}

	for index, menuItem := range m.MenuItems {
		var newline = "\n"
		if index == len(m.MenuItems) - 1 {
			// Adding a new line on the last option will move the cursor position out of range
			// For out redrawing
			newline = ""
		}

		menuItemText := menuItem.Text
		cursor := "  "
		if index == m.CursorPos {
			cursor = col.Yellow.Sprint("> ")
			menuItemText = col.Yellow.Sprint(menuItemText)
		}

		fmt.Printf("\r%s %s%s", cursor, menuItemText, newline)
	}
}

func (m *Menu) renderFinalChoice() {
	fmt.Printf("\r")
	fmt.Printf("\033[%dA\033E\033[0J", len(m.MenuItems))
	col.Cyan.Printf("%s", col.Bold.Sprint(m.Prompt + ":"))
	col.Yellow.Printf(" > %s\r\n", m.MenuItems[m.CursorPos].Text)
}

// Display will display the current menu options and awaits user selection
// It returns the users selected choice
func (m *Menu) Display() string {
	defer func() {
		// Show cursor again.
		fmt.Printf("\033[?25h")
	}()

	col.Cyan.Printf("%s\n", col.Bold.Sprint(m.Prompt))

	m.renderMenuItems(false)

	// Turn the terminal cursor off
	fmt.Printf("\033[?25l")

	
	var menuItem *MenuItem;
	escaped := false
	kb.Listen(func(key keys.Key) (stop bool, err error) {
		if key.Code == keys.Escape {
			escaped = true
			return true, nil // Stop listener by returning true on Ctrl+C
		} else if key.Code == keys.Enter {
			menuItem = m.MenuItems[m.CursorPos]
			m.renderFinalChoice()
			return true, nil
		} else if key.Code == keys.Up {
			m.CursorPos = (m.CursorPos + len(m.MenuItems) - 1) % len(m.MenuItems)
			m.renderMenuItems(true)
		} else if key.Code == keys.Down {
			m.CursorPos = (m.CursorPos + 1) % len(m.MenuItems)
			m.renderMenuItems(true)
		}
		return false, nil // Return false to continue listening
	})
	if escaped { return "" }
	return menuItem.ID
}