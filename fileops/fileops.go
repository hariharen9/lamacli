package fileops

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
	"lamacli/ui/styles"
)

// ReadFile reads the content of a file at the given path.
func ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// WriteFile writes data to a file at the given path.
func WriteFile(path string, data []byte) error {
	return os.WriteFile(path, data, 0644)
}

// DeleteFile deletes a file at the given path after user confirmation.
func DeleteFile(path string) error {
	var confirmed bool

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title(fmt.Sprintf("Are you sure you want to delete %s?", path)).
				Value(&confirmed),
		),
	).WithTheme(huh.ThemeBase16())

	err := form.Run()
	if err != nil {
		return err
	}

	if confirmed {
		return os.Remove(path)
	}

	fmt.Println(styles.SubtleStyle.Render("Deletion cancelled."))
	return nil
}
