package tests

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
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

	cmd := exec.Command("chmod", "997", tempFile) // Заменить на функцию runChmodCmd()

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

// Проверка на то, что смена прав у файла который является символьной ссылкой, меняет права файла на который ссылается ссылка
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

	assert.Equal(t, os.FileMode(0777), info.Mode().Perm(), "Permissions don`t match 777")
}

// Протестировать возомжно назначениярава через u-x
func TestChangePermWithSymAddXForUser(t *testing.T) {
	tempFile := "changePermWithSymAddXForUser.txt"
	defer os.Remove(tempFile)

	createTempFile(t, tempFile)

	err := runChmodCmd("u+x", tempFile)

	assert.NoError(t, err, "Error change permissions with help symbolic notation")

	info, err := os.Stat(tempFile)

	assert.NoError(t, err, "Error getting file information")

	assert.Equal(t, os.FileMode(0744), info.Mode().Perm(), "File permissions don`t match expected 0744")
}

// Протестировать возомжно отобрать права через o+w
func TestChangePermWithSymAddRForOther(t *testing.T) {
	tempFile := "ChangePermWithSymAddRForOther.txt"
	defer os.Remove(tempFile)

	createTempFile(t, tempFile)

	err := runChmodCmd("o-r", tempFile)

	assert.NoError(t, err, "Error change permissions with help symbolic notation")

	info, err := os.Stat(tempFile)

	assert.NoError(t, err, "Error getting file information")
	assert.Equal(t, os.FileMode(0640), info.Mode().Perm(), "File permissions don`t match expected 744")

}

// Протестировать это
// chmod -v 755 t.txt
// mode of 't.txt' changed from 0640 (rw-r-----) to 0755 (rwxr-xr-x)
func TestChmodVPerm(t *testing.T) {
	tempFile := "ChmodVPermFile.txt"
	defer os.Remove(tempFile)

	createTempFile(t, tempFile)

	cmd := exec.Command("chmod", "-v", "u=r,g=rx,o=rw", tempFile) //456

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()

	assert.NoError(t, err, "Error exec chmod")

	expOutput := fmt.Sprintf("mode of '%s' changed from 0644 (rw-r--r--) to 0456 (r--r-xrw-)", tempFile)
	cmdOutput := strings.TrimSpace(out.String())

	assert.Equal(t, expOutput, cmdOutput, "Output of the chmod doesn`t match what is expected")

	info, err := os.Stat(tempFile)
	assert.NoError(t, err, "Error getting file information")
	assert.Equal(t, os.FileMode(0456), info.Mode().Perm(), "File permissions don`t match expected 0456")

	// вывод у cmd будет  mode of 'ChmodVPermFile.txt' changed from 0644 (rw-r--r--) to 0456 (r--r-xrw-) и я хчоу убедиться что дейсвтительно такой вывод

}

// Протестировать это chmod go-r директория — удалить права на чтение для группы и остальных пользователей для каталога
// удаление права у юзера и группы на чтенеие
func TestChmodRemoveUserGroupR(t *testing.T) {
	tempFile := "chmodRemoveUserGroupRFile.txt"
	defer os.Remove(tempFile)

	createTempFile(t, tempFile)

	err := runChmodCmd("ug-r", tempFile)

	assert.NoError(t, err, "Error change permissions")

	info, err := os.Stat(tempFile)

	assert.NoError(t, err, "Error getting file information")

	assert.Equal(t, os.FileMode(0204), info.Mode().Perm(), "File permissions don`t match expected 204")

}

// Протестировать это chmod -R u+rwx,go-rwx каталог — добавит владельцу права на чтение, запись и выполнение, а группе и остальным пользователям уберет все права для всех файлов и каталогов в указанной директории и её подкаталогах.
func TestChmodAddandRemovePermUserGroupOther(t *testing.T) {
	// t.Skip()
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
			assert.Equal(t, os.FileMode(0722), info.Mode().Perm(), "Permissions do not match 0711 for path: "+path)
		} else {
			assert.Equal(t, os.FileMode(0622), info.Mode().Perm(), "Permissions do not match 0700 for path: "+path)
		}

		return nil
	})

	assert.NoError(t, err, "Error during filepath.Walk")

}

// Протестировать Ключ —reference и его использование
func TestChmodReferenceOpt(t *testing.T) {
	tempFile := "sourceFile.txt"  //откуда
	tempFile1 := "assignPerm.txt" //куда переделать на нормальные названия

	defer os.Remove(tempFile)
	defer os.Remove(tempFile1)

	createTempFile(t, tempFile)
	createTempFile(t, tempFile1)

	err := runChmodCmd("400", tempFile) // меняем права у файла от которого будем копировать права к другому файлу

	assert.NoError(t, err, "Error change permissions")

	//Проверяем что права успешно изменились
	info, err := os.Stat(tempFile)

	assert.NoError(t, err, "Error getting file information")

	assert.Equal(t, os.FileMode(0400), info.Mode().Perm(), "File permissions don`t match expected 111")

	// key := fmt.Sprintf("--reference=%s %s", tempFile, tempFile1)

	cmd := exec.Command("chmod", "--reference", tempFile, tempFile1)

	err = cmd.Run()

	er := fmt.Sprintf("Error copying permissions from %s to %s", tempFile, tempFile1)

	assert.NoError(t, err, er)

	info, err = os.Stat(tempFile1)

	assert.NoError(t, err, "Error getting file information")

	assert.Equal(t, os.FileMode(0400), info.Mode().Perm(), "File permissions don`t match 0400")

}

// Протестировать такую штуку chmod o-r,a-w month.txt text.txt, изменение прав сразу у несколькоких файлов
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
	//дописать завтра

	info, err = os.Stat(tempFile1)

	assert.NoError(t, err, "Error getting file information")

	assert.Equal(t, os.FileMode(0745), info.Mode().Perm(), "File permissions don`t match expected 0745 for %s", tempFile1)
	//дописать завтра
}

// (3) Добавить тест на chmod a+x
// Ты проверил u+x, o-r, но нет теста на a+x, который делает файл исполняемым для всех.
// Почему полезно?
// Покрывает важный chmod-сценарий, когда исполняемый бит ставится для всех.

// (4) Проверить chmod на файлы без прав (0000)
// Сейчас ты тестируешь chmod на обычных файлах, но что если файл вообще не имеет прав (0000)?
// Пример:
// очему это важно?

// Файлы без прав (0000) могут заблокировать доступ, chmod должен уметь их менять.

//Протестировать chmod u+s l.txt

// Настройка Sticky Bit
// Последний специальный бит разрешения – это Sticky Bit . В случае, если этот бит установлен для папки, то файлы в этой папке могут быть удалены только их владельцем. Пример использования этого бита в операционной системе это системная папка /tmp . Эта папка разрешена на запись любому пользователю, но удалять файлы в ней могут только пользователи, являющиеся владельцами этих файлов.

// root@ruvds-hrc [~]#  ls -ld /tmp
// drwxrwxrwt 8 root root 4096 Mar 25 10:22 /tmp
// Символ «t» указывает, что на папку установлен Sticky Bit.
func TestChmodStickyBit(t *testing.T) {
	tempDir := "chmodStickyBitDir"

	err := os.Mkdir(tempDir, 0755)
	assert.NoError(t, err, "Error creating directory")

	// defer os.RemoveAll(tempDir)

	err = runChmodCmd("+t", tempDir)

	assert.NoError(t, err, "Error change permissions")

	info, err := os.Stat(tempDir)

	assert.NoError(t, err, "Error getting dir information")

	fmt.Printf("PERM: %o", info.Mode().Perm())

	if os.ModeSticky != 0 {
		fmt.Println(os.ModeSticky, reflect.TypeOf(os.ModeSticky))
		fmt.Println("Sticky bit is set")
	} else {
		fmt.Println("Sticky bit is not set")
	}

	//доделать, что то не рабоатет

	// assert.Equal(t, os.ModeSticky)

	// assert.Equal(t, os.FileMode(01755), info.Mode().Perm(), "Permissions don`t match expected 01755")

}
