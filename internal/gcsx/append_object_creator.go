// Copyright 2015 Google Inc. All Rights Reserved.
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

package gcsx

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/Tulsishah/gcsfuse-tulsishah7/v2/internal/storage/gcs"
	"golang.org/x/net/context"
)

// Create an objectCreator that accepts a source object and the contents that
// should be "appended" to it, storing temporary objects using the supplied
// prefix.
//
// Note that the Create method will attempt to remove any temporary junk left
// behind, but it may fail to do so. Users should arrange for garbage collection.
//
// Create guarantees to return *gcs.PreconditionError when the source object
// has been clobbered.
func newAppendObjectCreator(
	prefix string,
	bucket gcs.Bucket) (oc objectCreator) {
	oc = &appendObjectCreator{
		prefix: prefix,
		bucket: bucket,
	}

	return
}

////////////////////////////////////////////////////////////////////////
// Implementation
////////////////////////////////////////////////////////////////////////

type appendObjectCreator struct {
	prefix string
	bucket gcs.Bucket
}

func (oc *appendObjectCreator) chooseName() (name string, err error) {
	// Generate a good 64-bit random number.
	var buf [8]byte
	_, err = io.ReadFull(rand.Reader, buf[:])
	if err != nil {
		err = fmt.Errorf("ReadFull: %w", err)
		return
	}

	x := uint64(buf[0])<<0 |
		uint64(buf[1])<<8 |
		uint64(buf[2])<<16 |
		uint64(buf[3])<<24 |
		uint64(buf[4])<<32 |
		uint64(buf[5])<<40 |
		uint64(buf[6])<<48 |
		uint64(buf[7])<<56

	// Turn it into a name.
	name = fmt.Sprintf("%s%016x", oc.prefix, x)

	return
}

// ObjectName param is present here for consistency between fullObjectCreator
// and appendObjectCreator. ObjectName is not used in append flow since
// srcObject.Name gives the objectName.
func (oc *appendObjectCreator) Create(
	ctx context.Context,
	objectName string,
	srcObject *gcs.Object,
	mtime *time.Time,
	r io.Reader) (o *gcs.Object, err error) {
	// Choose a name for a temporary object.
	tmpName, err := oc.chooseName()
	if err != nil {
		err = fmt.Errorf("chooseName: %w", err)
		return
	}

	// Create a temporary object containing the additional contents.
	var zero int64
	tmp, err := oc.bucket.CreateObject(
		ctx,
		&gcs.CreateObjectRequest{
			Name:                   tmpName,
			GenerationPrecondition: &zero,
			Contents:               r,
		})
	if err != nil {
		err = fmt.Errorf("CreateObject: %w", err)
		return
	}

	// Attempt to delete the temporary object when we're done.
	defer func() {
		deleteErr := oc.bucket.DeleteObject(
			ctx,
			&gcs.DeleteObjectRequest{
				Name:       tmp.Name,
				Generation: 0, // Delete the latest generation of temporary object.
			})

		if err == nil && deleteErr != nil {
			err = fmt.Errorf("DeleteObject: %w", deleteErr)
		}
	}()

	MetadataMap := make(map[string]string)

	/* Copy Metadata fields from src object to new object generated by compose. */
	for key, value := range srcObject.Metadata {
		MetadataMap[key] = value
	}

	if mtime != nil {
		MetadataMap[MtimeMetadataKey] = mtime.UTC().Format(time.RFC3339Nano)
	}

	// Compose the old contents plus the new over the old.
	o, err = oc.bucket.ComposeObjects(
		ctx,
		&gcs.ComposeObjectsRequest{
			DstName:                       srcObject.Name,
			DstGenerationPrecondition:     &srcObject.Generation,
			DstMetaGenerationPrecondition: &srcObject.MetaGeneration,
			Sources: []gcs.ComposeSource{
				gcs.ComposeSource{
					Name:       srcObject.Name,
					Generation: srcObject.Generation,
				},

				gcs.ComposeSource{
					Name:       tmp.Name,
					Generation: tmp.Generation,
				},
			},
			Metadata:           MetadataMap,
			CacheControl:       srcObject.CacheControl,
			ContentDisposition: srcObject.ContentDisposition,
			ContentEncoding:    srcObject.ContentEncoding,
			ContentType:        srcObject.ContentType,
			CustomTime:         srcObject.CustomTime,
			EventBasedHold:     srcObject.EventBasedHold,
			StorageClass:       srcObject.StorageClass,
		})
	if err != nil {
		// A not found error means that either the source object was clobbered or the
		// temporary object was. The latter is unlikely, so we signal a precondition
		// error.
		var notFoundErr *gcs.NotFoundError
		if errors.As(err, &notFoundErr) {
			err = &gcs.PreconditionError{
				Err: err,
			}
		}

		err = fmt.Errorf("ComposeObjects: %w", err)
		return
	}

	return
}
