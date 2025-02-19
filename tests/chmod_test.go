package tests

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createTempFile(t *testing.T, fileName string) {
	file, err := os.Create(fileName)
	assert.NoError(t, err, "Error creating file")
	file.Close()
}

func runChmodCmd(perm, fileName string) error {
	cmd := exec.Command("chmod", perm, fileName)
	err := cmd.Run()

	return err
}

func generateDirectoriesAndFiles(baseDir string) {
	for i := range 3 {
		tempDir := filepath.Join(baseDir, fmt.Sprintf("tempDir%d", i))
		err := os.Mkdir(tempDir, 0755)
		if err != nil {
			fmt.Println("Error creating directory:", err)
			return
		}
		// defer os.RemoveAll(tempDir)

		numSubDirs := rand.Intn(6)
		numFiles := rand.Intn(6)

		for j := range numSubDirs {

			subDir := filepath.Join(tempDir, fmt.Sprintf("subdir%d", j))
			err := os.Mkdir(subDir, 0755)
			if err != nil {
				fmt.Println("Error creating subdirectory:", err)
				return
			}

			numSubFiles := rand.Intn(6)

			for k := range numSubFiles {
				filePath := filepath.Join(subDir, fmt.Sprintf("file%d.txt", k))
				file, err := os.Create(filePath)
				if err != nil {
					return
				}
				file.Close()
			}
		}

		for k := range numFiles {
			filePath := filepath.Join(tempDir, fmt.Sprintf("file%d.txt", k))
			file, err := os.Create(filePath)
			if err != nil {
				return
			}
			file.Close()
		}
	}
}

// Проверка создания файла
func TestCreateFile(t *testing.T) {
	tempFile := "testFile1.txt"
	defer os.Remove(tempFile)

	createTempFile(t, tempFile)

	_, err := os.Stat(tempFile)

	assert.NoError(t, err, "File not found after creation")
}

// Проверка изменения прав доступа файла
func TestChangeFilePerm(t *testing.T) {
	tempFile := "testFile2.txt"
	defer os.Remove(tempFile)

	createTempFile(t, tempFile)

	err := runChmodCmd("755", tempFile)

	assert.NoError(t, err, "Error saving file permissions")

	info, err := os.Stat(tempFile)

	assert.NoError(t, err, "Error getting file information")
	assert.Equal(t, os.FileMode(0755), info.Mode().Perm(), "File permissions don`t match 0755")
}

// проверка на изменение прав директории
func TestChangeDirPerm(t *testing.T) {
	tempDir := "testDir"
	defer os.RemoveAll(tempDir)

	err := os.Mkdir(tempDir, 0755)

	assert.NoError(t, err, "Error creating directory")

	cmd := exec.Command("chmod", "700", tempDir)

	err = cmd.Run()

	assert.NoError(t, err, "Error changing directory permissions")

	info, err := os.Stat(tempDir)

	assert.NoError(t, err, "Error getting directory information")
	assert.Equal(t, os.FileMode(0700), info.Mode().Perm(), "Directory permissions do not match 0700")

}

// // проверка на неправильное использование
func TestSetWrongPerm(t *testing.T) {
	tempFile := "testFile.txt"
	defer os.Remove(tempFile)

	file, err := os.Create(tempFile)

	assert.NoError(t, err, "Error creating file")
	file.Close()

	cmd := exec.Command("chmod", "997", tempFile)

	err = cmd.Run()
	assert.Error(t, err, "Expected error when trying to set wrong permissions")
}

// Проверка на изменение прав для несуществуюшего файла
func TestChangeNonExistFilePerm(t *testing.T) {
	err := runChmodCmd("755", "nonexistFile.txt")
	assert.Error(t, err, "Error expected when changing permissions on a non-existent file")
}

// Рекурсивная проверка на то что права изменились
func TestChangePermRecurs(t *testing.T) {
	baseDir := "folder"

	err := os.Mkdir(baseDir, 0755)

	assert.NoError(t, err, "Error creating directory")

	defer os.RemoveAll(baseDir)

	generateDirectoriesAndFiles(baseDir)

	cmd := exec.Command("chmod", "-R", "700", baseDir)
	err = cmd.Run()

	assert.NoError(t, err, "Error changing permissions recursively")

	err = filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		assert.NoError(t, err, "Error walking through path")

		assert.Equal(t, os.FileMode(0700), info.Mode().Perm(), "Permissions do not match 0700 for path: "+path)
		return nil
	})

	assert.NoError(t, err, "Error during filepath.Walk")
}
