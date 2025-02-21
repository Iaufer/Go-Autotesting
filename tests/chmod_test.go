package tests

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createTempFile(t *testing.T, fileName string) {
	file, err := os.Create(fileName)
	assert.NoError(t, err, "Error creating file")
	file.Close()
}

// Функция изменения прав доступа к файлу
func runChmodCmd(perm, fileName string) error {
	cmd := exec.Command("chmod", perm, fileName)
	err := cmd.Run()

	return err
}

// Функция генерации вложенных папок и файлов
func generateDirectoriesAndFiles(baseDir string) {
	for i := range 3 {
		tempDir := filepath.Join(baseDir, fmt.Sprintf("tempDir%d", i))
		err := os.Mkdir(tempDir, 0755)
		if err != nil {
			fmt.Println("Error creating directory:", err)
			return
		}
		defer os.RemoveAll(tempDir)

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

/*
Тест-кейс 1:
Проверка создания файла
Проверяет, что файл создается

1. Создается временный файл testFile1.txt
2. Проверяется, что файл действительно создан

Ожидается, что файл будет успешно создан
*/
func TestCreateFile(t *testing.T) {
	tempFile := "testFile1.txt"
	defer os.Remove(tempFile)

	createTempFile(t, tempFile)

	_, err := os.Stat(tempFile)

	assert.NoError(t, err, "File not found after creation")
}

/*
Тест-кейс 2:
Проверка изменения прав доступа к файлу
Проверяет, что права доступа изменяется при помощи chmod

1. Создается временный файл
2. Изменяются права доступа на 0755
3. Проверяется, чтоу файла права доступа равны 0755

Ожидается, что права доступа к файлу изменятся на 0755
*/
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

/*
Тест-кейс 3
Проверка изменения прав доступа для директории
Проверяет, что права доступа к директории изменятся

1. Создается директория с именем testDir
2. Изменяются права доступа на 0700
3. Проверяются, что права доступа к директории изменились на 0700

Ожидается, что права доступа директории будут изменены на 0700
*/
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

/*
Тест-кейс 4
Проверка ошибки при установке неправильных прав доступа
Проверяется, что если установать некорректные права то будет ошибка

1. Создается временый файл testFile.txt
2. Устанавливаются некорректные права доступа 0997
3. Ожидается ошибка при выполнении команды

Ожидается, что возникнет ошибка
*/
func TestSetWrongPerm(t *testing.T) {
	tempFile := "testFile.txt"
	defer os.Remove(tempFile)

	file, err := os.Create(tempFile)

	assert.NoError(t, err, "Error creating file")
	file.Close()

	err = runChmodCmd("997", tempFile)
	assert.Error(t, err, "Error changing permissions")

	cmd := exec.Command("chmod", "997", tempFile)

	err = cmd.Run()
	assert.Error(t, err, "Expected error when trying to set wrong permissions")
}

/*
Тест-кейс 5
Проверка изменения прав для несуществующего файла

1.Попытка изменить права доступа для nonexistFile.txt на 0755
2.Ожидается ошибка

Ожидается, что возникнет ошибка
*/
func TestChangeNonExistFilePerm(t *testing.T) {
	err := runChmodCmd("755", "nonexistFile.txt")
	assert.Error(t, err, "Error expected when changing permissions on a non-existent file")
}

/*
Тест-кейс 6:
Проверка рекурсивного изменения правдоступа для директории и всех ее поддиректорий и файлов

1. Создается директория и внутри нее несколько файлов и поддиректорий
2. Выполняется рекурсивное изменение прав доступа на 0700
3. Проверяется, что права доступа были изменены для всех файлов и поддиректорий

Ожидается, что права доступа будут изменены на 0700
*/
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

/*
Тест-кейс 7:
Проверка на то, что смена прав у файла который является символьной ссылкой, меняет права файла на который ссылается ссылка

1. Создается файл ChangeSymLinkPermFile.txt
2. Создается символьная ссылка testSymLinlk.txt
3. Выполняется изменение прав доступа у ссылки
4. Проверяется, что права доступа были изменены на файла на который ссылается testSymLinlk.txt

Ожидается, что изменятся права доступа у файла ChangeSymLinkPermFile.txt
*/
func TestChangeSymLinkPerm(t *testing.T) {
	tempFile := "ChangeSymLinkPermFile.txt"
	symLinkFile := "testSymLinlk.txt"

	defer os.Remove(tempFile)
	defer os.Remove(symLinkFile)

	createTempFile(t, tempFile)
	err := os.Symlink(tempFile, symLinkFile)

	assert.NoError(t, err, "Error creating symbolic link")

	err = runChmodCmd("777", symLinkFile)

	assert.NoError(t, err, "Error change permissions with help symbolic link")

	info, err := os.Stat(tempFile)
	assert.NoError(t, err, "Error getting file information")

	assert.Equal(t, os.FileMode(0777), info.Mode().Perm(), "Permissions don`t match 0777")
}

/*
Тест-кейс 8:
Проверка возможности назначения прав через u+x

1. Создается файл
2. Выполняется изменение прав доступа для пользователя u+x
3. Проверка, что права доступа были изменены корректно

Ожидается, что права доступа у файла будут равны 0764
*/
func TestChangePermWithSymAddXForUser(t *testing.T) {
	tempFile := "changePermWithSymAddXForUser.txt"
	defer os.Remove(tempFile)

	createTempFile(t, tempFile)

	err := runChmodCmd("u+x", tempFile)

	assert.NoError(t, err, "Error change permissions with help symbolic notation")

	info, err := os.Stat(tempFile)

	assert.NoError(t, err, "Error getting file information")

	assert.Equal(t, os.FileMode(0764), info.Mode().Perm(), "File permissions don`t match expected 0764")
}

/*
Тест-кейс 9:
Проверка возможности удаления прав через o-r

1. Создается файл
2. Выполняется изменение прав доступа для other o-r
3. Проверка, что права доступа были изменены корректно

Ожидается, что права доступа у файла будут равны 0660
*/
func TestChangePermWithSymAddRForOther(t *testing.T) {
	tempFile := "ChangePermWithSymAddRForOther.txt"
	defer os.Remove(tempFile)

	createTempFile(t, tempFile)

	err := runChmodCmd("o-r", tempFile)

	assert.NoError(t, err, "Error change permissions with help symbolic notation")

	info, err := os.Stat(tempFile)

	assert.NoError(t, err, "Error getting file information")
	assert.Equal(t, os.FileMode(0660), info.Mode().Perm(), "File permissions don`t match expected 0660")

}

/*
Тест-кейс 10:
Проверка использования опции chmod -v

1. Создается файл
2. Выполняется командас клюом -v
3. Изменяются права на u=r,g=rx,o=r
4. Проверяется, что права изменились корректно

Ожидается, что вывод команды chmod с ключом -v соотсветсвует ожидаемому и права файла были изменены на 0456
*/
func TestChmodVPerm(t *testing.T) {
	tempFile := "ChmodVPermFile.txt"
	defer os.Remove(tempFile)

	createTempFile(t, tempFile)

	cmd := exec.Command("chmod", "-v", "u=r,g=rx,o=rw", tempFile) //456

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()

	assert.NoError(t, err, "Error exec chmod")

	expOutput := fmt.Sprintf("mode of '%s' changed from 0664 (rw-rw-r--) to 0456 (r--r-xrw-)", tempFile)
	cmdOutput := strings.TrimSpace(out.String())

	assert.Equal(t, expOutput, cmdOutput, "Output of the chmod doesn`t match what is expected")

	info, err := os.Stat(tempFile)
	assert.NoError(t, err, "Error getting file information")
	assert.Equal(t, os.FileMode(0456), info.Mode().Perm(), "File permissions don`t match expected 0456")
}

/*
Тест-кейс 11:
Удаление прав на чтение у группы и остальных пользователей

1. Создается файл
2. Удаляются права на чтение для группы и остальных пользователей
3. Проверяется, что права файла изменились корректно

Ожидается, что права файла будут 0224
*/
func TestChmodRemoveUserGroupR(t *testing.T) {
	tempFile := "chmodRemoveUserGroupRFile.txt"
	defer os.Remove(tempFile)

	createTempFile(t, tempFile)

	err := runChmodCmd("ug-r", tempFile)

	assert.NoError(t, err, "Error change permissions")

	info, err := os.Stat(tempFile)

	assert.NoError(t, err, "Error getting file information")

	assert.Equal(t, os.FileMode(0224), info.Mode().Perm(), "File permissions don`t match expected 204")

}

/*
Тест-кейс 12:
Рекурсивное изменение прав для владельца, группы и других

1. Создается директория
2. Создаются поддиректории и файлы
3. Выполняется командаc сhmod -R go-wrx,go+w

Ожидается, что прав для директорий будут 0722, а для файлов 0622
*/
func TestChmodAddandRemovePermUserGroupOther(t *testing.T) {
	baseDir := "folder1"
	err := os.Mkdir(baseDir, 0755)

	assert.NoError(t, err, "Error creating directory")

	defer os.RemoveAll(baseDir)

	generateDirectoriesAndFiles(baseDir)

	cmd := exec.Command("chmod", "-R", "go-wrx,go+w", baseDir)
	err = cmd.Run()

	assert.NoError(t, err, "Error changing permissions recursively")

	err = filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		assert.NoError(t, err, "Error walking through path")

		if info.IsDir() {
			assert.Equal(t, os.FileMode(0722), info.Mode().Perm(), "Permissions do not match 0722 for path: "+path)
		} else {
			assert.Equal(t, os.FileMode(0622), info.Mode().Perm(), "Permissions do not match 0622 for path: "+path)
		}

		return nil
	})

	assert.NoError(t, err, "Error during filepath.Walk")

}

/*
Тест-кейс 13:
Проверка копирования прав доступа с одного файла на другой с помощью --reference

1. Создается два файла
2. Изменяются права доступа первого файла на 0400
3. С помощью ключа --reference копируются права с первого файла

Ожидается, что права второго файла и первого совпадут (0400)
*/
func TestChmodReferenceOpt(t *testing.T) {
	tempFile := "sourceFile.txt"
	tempFile1 := "assignPerm.txt"

	defer os.Remove(tempFile)
	defer os.Remove(tempFile1)

	createTempFile(t, tempFile)
	createTempFile(t, tempFile1)

	err := runChmodCmd("400", tempFile)

	assert.NoError(t, err, "Error change permissions")

	info, err := os.Stat(tempFile)

	assert.NoError(t, err, "Error getting file information")

	assert.Equal(t, os.FileMode(0400), info.Mode().Perm(), "File permissions don`t match expected 111")

	cmd := exec.Command("chmod", "--reference", tempFile, tempFile1)

	err = cmd.Run()

	er := fmt.Sprintf("Error copying permissions from %s to %s", tempFile, tempFile1)

	assert.NoError(t, err, er)

	info, err = os.Stat(tempFile1)

	assert.NoError(t, err, "Error getting file information")

	assert.Equal(t, os.FileMode(0400), info.Mode().Perm(), "File permissions don`t match 0400")

}

/*
Тест-кейс 14:
Проверка изменения прав сразу у двух файлов

1. Создается два файла
2. Изменяются права доступа для обоих файлов (chmod u+x,g-wx,o+x)
3. Проверяется, что права были изменены корректно

Ожидается, что права у обоих файлов будут 0745
*/
func TestModifyPermForTwoFiles(t *testing.T) {
	tempFile := "modifyPermForTwoFiles1.txt"
	tempFile1 := "modifyPermForTwoFiles2.txt"

	defer os.RemoveAll(tempFile)
	defer os.RemoveAll(tempFile1)

	createTempFile(t, tempFile)
	createTempFile(t, tempFile1)

	cmd := exec.Command("chmod", "u+x,g-wx,o+x", tempFile, tempFile1)

	err := cmd.Run()

	assert.NoError(t, err, "Error changing permissions for two files")

	info, err := os.Stat(tempFile)

	assert.NoError(t, err, "Error getting file information")

	assert.Equal(t, os.FileMode(0745), info.Mode().Perm(), "File permissions don`t match expected 0745 for %s", tempFile)

	info, err = os.Stat(tempFile1)

	assert.NoError(t, err, "Error getting file information")

	assert.Equal(t, os.FileMode(0745), info.Mode().Perm(), "File permissions don`t match expected 0745 for %s", tempFile1)
}
