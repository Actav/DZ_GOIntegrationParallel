package creaters

import "os"

func FolderIfNotExists(folderName string) error {
	if _, err := os.Stat(folderName); os.IsNotExist(err) {
		if err := os.MkdirAll(folderName, 0755); err != nil {
			return err
		}
	}

	return nil
}
