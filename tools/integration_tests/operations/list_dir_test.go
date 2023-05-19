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

// Provides integration tests for list directory.
package operations_test

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/googlecloudplatform/gcsfuse/tools/integration_tests/util/setup"
)

func createDirectoryWithFile(dirPath string, filePath string, t *testing.T) {
	err := os.Mkdir(dirPath, setup.FilePermission_0600)
	if err != nil {
		t.Errorf("Mkdir at %q: %v", dirPath, err)
		return
	}

	_, err = os.Create(filePath)
	if err != nil {
		t.Errorf("Error in creating file %v:", err)
	}
}

func checkIfListedCorrectDirectory(dirPath string, obj fs.DirEntry, t *testing.T) {
	var dirName string

	if strings.Contains(dirPath, SecondSubDirectoryForListTest) {
		dirName = SecondSubDirectoryForListTest
	} else if strings.Contains(dirPath, FirstSubDirectoryForListTest) {
		dirName = FirstSubDirectoryForListTest
	} else if strings.Contains(dirPath, DirectoryForListTest) {
		dirName = DirectoryForListTest
	} else if strings.Contains(dirPath, setup.MntDir()) {
		dirName = setup.MntDir()
	}

	switch dirName {
	case setup.MntDir():
		{
			if obj.Name() != DirectoryForListTest || obj.IsDir() != true {
				t.Errorf("Listed incorrect Object.")
			}
		}
	case DirectoryForListTest:
		{
			if (obj.Name() != FileInDirectoryForListTest && obj.IsDir() == true) && (obj.Name() != FirstSubDirectoryForListTest && obj.IsDir() != true) && (obj.Name() != SecondSubDirectoryForListTest && obj.IsDir() != true) {
				t.Errorf("Listed incorrect Object")
			}
		}
	case FirstSubDirectoryForListTest:
		{
			if obj.Name() != FileInFirstSubDirectoryForListTest && obj.IsDir() == true {
				t.Errorf("Listed incorrect Object")
			}
		}
	case SecondSubDirectoryForListTest:
		{
			if obj.Name() != FileInSecondSubDirectoryForListTest && obj.IsDir() == true {
				t.Errorf("Listed incorrect Object")
			}
		}
	}
}

// List directory recursively
func listDirectory(path string, t *testing.T) {
	//Reading contents of the directory
	files, err := os.ReadDir(path)

	// error handling
	if err != nil {
		log.Fatal(err)
	}

	for _, obj := range files {
		checkIfListedCorrectDirectory(path, obj, t)
		if obj.IsDir() {
			subDirectoryPath := filepath.Join(path, obj.Name()) // path of the subdirectory
			listDirectory(subDirectoryPath, t)                  // calling listFiles() again for subdirectory
		}
	}
}

func TestListDirectoryRecursively(t *testing.T) {
	// Clean the bucket for list testing.
	os.RemoveAll(setup.MntDir())

	// Directory structure
	// testBucket/directoryForListTest                                                                     -- Dir
	// testBucket/directoryForListTest/fileInDirectoryForListTest		                                       -- File
	// testBucket/directoryForListTest/firstSubDirectoryForListTest                                        -- Dir
	// testBucket/directoryForListTest/firstSubDirectoryForListTest/fileInFirstSubDirectoryForListTest     -- File
	// testBucket/directoryForListTest/secondSubDirectoryForListTest                                       -- Dir
	// testBucket/directoryForListTest/secondSubDirectoryForListTest/fileInSecondSubDirectoryForListTest   -- File

	// testBucket/directoryForListTest
	// testBucket/directoryForListTest/fileInDirectoryForListTest
	dirPath := path.Join(setup.MntDir(), DirectoryForListTest)
	filePath := path.Join(dirPath, FileInDirectoryForListTest)
	createDirectoryWithFile(dirPath, filePath, t)

	// testBucket/directoryForListTest/firstSubDirectoryForListTest
	// testBucket/directoryForListTest/firstSubDirectoryForListTest/fileInFirstSubDirectoryForListTest
	subDirPath := path.Join(dirPath, FirstSubDirectoryForListTest)
	subDirFilePath := path.Join(subDirPath, FileInFirstSubDirectoryForListTest)
	createDirectoryWithFile(subDirPath, subDirFilePath, t)

	// testBucket/directoryForListTest/secondSubDirectoryForListTest
	// testBucket/directoryForListTest/secondSubDirectoryForListTest/fileInSecondSubDirectoryForListTest
	subDirPath = path.Join(dirPath, SecondSubDirectoryForListTest)
	subDirFilePath = path.Join(subDirPath, FileInSecondSubDirectoryForListTest)
	createDirectoryWithFile(subDirPath, subDirFilePath, t)

	// Test directory listing recursively.
	listDirectory(setup.MntDir(), t)

	// Clean the bucket after list testing.
	os.RemoveAll(setup.MntDir())
}

func TestLargeDirectoryForListing(t *testing.T) {
	// Clean the bucket for list testing.
	//os.RemoveAll(setup.MntDir())

	dirPath := path.Join(setup.MntDir(), DirectoryForListTest)
	//setup.CreateDirectoryWithNFiles(12000, dirPath, t)

	files, err := os.ReadDir(dirPath)
	if err != nil {
		t.Errorf("Error in listing directory.")
	}

	fmt.Println(len(files))
	if len(files) != 12000 {
		t.Errorf("Listed incorrect number of files from directory: %v, expected 12000", len(files))
	}

	// Clean the bucket after list testing.
	//os.RemoveAll(setup.MntDir())
}
