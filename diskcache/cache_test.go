/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package diskcache

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"path/filepath"
	"testing"
)

func hashBytes(b []byte) string {
	hasher := sha256.New()
	hasher.Write(b)
	return hex.EncodeToString(hasher.Sum(nil))
}

func makeRandomBytes(n int) (b []byte, err error) {
	b = make([]byte, n)
	_, err = rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// test cache.PathToKey and cache.KeyToPath
func TestPathToKeyKeyToPath(t *testing.T) {
	cache := NewCache("/some/dir", 3)
	testCases := []struct {
		Key          string
		ExpectedPath string
	}{
		{
			Key:          "key",
			ExpectedPath: filepath.Join("/some/dir", "key"),
		},
		{
			Key:          "namespaced/key",
			ExpectedPath: filepath.Join("/some/dir", "namespaced/key"),
		},
		{
			Key:          "entry/nested/a/few/times",
			ExpectedPath: filepath.Join("/some/dir", "entry/nested/a/few/times"),
		},
		{
			Key:          "foobar/cas/asdf",
			ExpectedPath: filepath.Join("/some/dir", "foobar/cas/asdf"),
		},
		{
			Key:          "foobar/ac/asdf",
			ExpectedPath: filepath.Join("/some/dir", "foobar/ac/asdf"),
		},
	}
	for _, tc := range testCases {
		path := cache.KeyToPath(tc.Key)
		if path != tc.ExpectedPath {
			t.Fatalf("expected KeyToPath(%s) to be %s", tc.Key, path)
		}
		backToKey := cache.PathToKey(path)
		if backToKey != tc.Key {
			t.Fatalf("cache.KeyToPath(cache.PathToKey(%s)) was not idempotent, got %s", tc.Key, backToKey)
		}
	}
}
