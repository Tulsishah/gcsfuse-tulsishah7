// Copyright 2023 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Provide a helper functions.
package file_operations

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"syscall"

	"github.com/googlecloudplatform/gcsfuse/tools/integration_tests/util/setup"
)

func CopyFile(srcFile string, newFileName string) (err error) {
	if _, err = os.Stat(newFileName); err == nil {
		err = fmt.Errorf("Copied file %s already present", newFileName)
		return
	}

	// File copying with io.Copy() utility.
	source, err := os.OpenFile(srcFile, syscall.O_DIRECT, setup.FilePermission_0600)
	if err != nil {
		err = fmt.Errorf("File %s opening error: %v", srcFile, err)
		return
	}
	defer source.Close()

	destination, err := os.OpenFile(newFileName, os.O_WRONLY|os.O_CREATE|syscall.O_DIRECT, setup.FilePermission_0600)
	if err != nil {
		err = fmt.Errorf("Copied file creation error: %v", err)
		return
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	if err != nil {
		err = fmt.Errorf("Error in file copying: %v", err)
		return
	}
	return
}

func ReadFile(filePath string) (content []byte, err error) {
	file, err := os.OpenFile(filePath, os.O_RDONLY|syscall.O_DIRECT, setup.FilePermission_0600)
	if err != nil {
		err = fmt.Errorf("Error in the opening the file %v", err)
		return
	}
	defer file.Close()

	content, err = os.ReadFile(file.Name())
	if err != nil {
		err = fmt.Errorf("ReadAll: %v", err)
		return
	}
	return
}

func RenameFile(fileName string, newFileName string) (err error) {
	if _, err = os.Stat(newFileName); err == nil {
		err = fmt.Errorf("Renamed file %s already present", newFileName)
		return
	}

	if err = os.Rename(fileName, newFileName); err != nil {
		err = fmt.Errorf("Rename unsuccessful: %v", err)
		return
	}

	if _, err = os.Stat(fileName); err == nil {
		err = fmt.Errorf("Original file %s still exists", fileName)
		return
	}
	if _, err = os.Stat(newFileName); err != nil {
		err = fmt.Errorf("Renamed file %s not found", newFileName)
		return
	}
	return
}

func WriteAtEndOfFile(fileName string, content string) (err error) {
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY|syscall.O_DIRECT, setup.FilePermission_0600)
	if err != nil {
		err = fmt.Errorf("Open file for append: %v", err)
		return
	}

	_, err = f.WriteString(content)
	f.Close()

	return
}

func WriteAtStartOfFile(fileName string, content string) (err error) {
	f, err := os.OpenFile(fileName, os.O_RDWR|syscall.O_DIRECT, setup.FilePermission_0600)
	if err != nil {
		err = fmt.Errorf("Open file for write at start: %v", err)
		return
	}

	_, err = f.WriteAt([]byte(content), 0)
	f.Close()

	return
}

func MoveFile(srcDirPath string, destDirPath string) (err error) {
	cmd := exec.Command("mv", srcDirPath, destDirPath)

	err = cmd.Run()
	if err != nil {
		err = fmt.Errorf("Moving file operation is failed: %v", err)
	}
	return err
}
