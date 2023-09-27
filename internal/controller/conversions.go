/*
Copyright 2023 The Nephio Authors.

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

package controller

func int32Ptr(val int) *int32 {
	var a int32
	a = int32(val)
	return &a
}

func int64Ptr(val int) *int64 {
	var a int64
	a = int64(val)
	return &a
}

func intPtr(val int) *int {
	a := val
	return &a
}

func int16Ptr(val int) *int16 {
	var a int16
	a = int16(val)
	return &a
}

func boolPtr(val bool) *bool {
	a := val
	return &a
}

func stringPtr(val string) *string {
	a := val
	return &a
}
