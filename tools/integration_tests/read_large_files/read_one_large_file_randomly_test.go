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

package read_large_files

import (
	"math/rand"
	"os"
	"path"
	"strconv"
	"testing"

	"github.com/googlecloudplatform/gcsfuse/tools/integration_tests/util/operations"
	"github.com/googlecloudplatform/gcsfuse/tools/integration_tests/util/setup"
)

func TestReadLargeFileRandomly(t *testing.T) {
	// Clean the mountedDirectory before running test.
	setup.CleanMntDir()

	// Create file of 500 MB with random data in local disk.
	fileInLocalDisk := path.Join(os.Getenv("HOME"), FiveHundredMBFile)
	setup.RunScriptForTestData("testdata/write_content_of_fix_size_in_file.sh", fileInLocalDisk, strconv.Itoa(FiveHundredMB))

	// Copy the file in mounted directory.
	file := path.Join(setup.MntDir(), FiveHundredMBFile)
	err := operations.CopyFile(fileInLocalDisk, file)
	if err != nil {
		t.Errorf("Error in copying file:%v", err)
	}

	file2 := path.Join("/tmp", "a.txt")
	err = operations.CopyFile(fileInLocalDisk, file2)
	if err != nil {
		t.Errorf("Error in copying file:%v", err)
	}

	for i := 0; i < 20; i++ {
		offset := rand.Int63n(MaxReadbleByteFromFile - MinReadbleByteFromFile)
		// Randomly read the data from file.
		//content, err := operations.ReadChunkFromFile(file, chunkSize, offset)
		//if err != nil {
		//	t.Errorf("Error in reading file: %v", err)
		//}

		//	Read actual content from file located in local disk.
		_, err := operations.ReadChunkFromFile(file2, chunkSize, offset)
		if err != nil {
			t.Errorf("Error in reading file: %v", err)
		}

		//	Compare actual content and expect content.
		//if bytes.Equal(actualContent, content) == false {
		//	t.Errorf("Error in reading file randomly.")
		//}
	}

	// Removing file after testing.
	operations.RemoveFile(fileInLocalDisk)
}
