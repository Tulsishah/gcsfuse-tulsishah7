# Copyright 2023 Google Inc. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#	http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

TEST_BUCKET=$1
OBJECT_NAME=$2
OBJECT_PATH_IN_DISK=$3

gcloud alpha storage cp gs://$TEST_BUCKET/$OBJECT_NAME $OBJECT_PATH_IN_DISK 2>&1 | tee ~/output.txt
if grep "The following URLs matched no objects or files" ~/output.txt; then
  echo "Object does not exist."
fi
rm ~/output.txt