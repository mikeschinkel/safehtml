// Copyright (c) 2017 The Go Authors. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file or at
// https://developers.google.com/open-source/licenses/bsd

package safehtml

import (
	"bytes"
	"encoding/json"
)

// An JSON is an immutable string-like type that is safe to use in JSON
// contexts in DOM APIs and JSON documents.
//
// JSON guarantees that its value as a string will not cause untrusted script
// execution when evaluated as JSON in a browser.
//
// Values of this type are guaranteed to be safe to use in JSON contexts,
// such as assignment to the innerJSON DOM property, or interpolation into an
// JSON template in JSON PC_DATA context, in the sense that the use will not
// result in a Cross-site Scripting (XSS) vulnerability.
type JSON struct {
	// We declare an JSON not as a string but as a struct wrapping a string
	// to prevent construction of JSON values through string conversion.
	str string
}

func JSONFromConstant(str stringConstant) JSON {
	return JSON{str: string(str)}
}

func JSONFromValue(input string) (out JSON, err error) {
	var x any
	var j []byte
	err = json.Unmarshal([]byte(input), &x)
	if err != nil {
		goto end
	}
	j, err = json.Marshal(x)
	if err != nil {
		goto end
	}
	out = JSON{str: string(j)}
end:
	return out, err
}

// JSONer is implemented by any value that has an JSON method, which defines the
// safe JSON format for that value.
type JSONer interface {
	JSON() JSON
}

// JSONEscaped returns a JSON whose value is text, with the characters [&<>"'] escaped.
//
// text is coerced to interchange valid, so the resulting JSON contains only
// valid UTF-8 characters which are legal in JSON and XML.
func JSONEscaped(text string) JSON {
	return JSON{coerceToUTF8InterchangeValid(text)}
	//return JSON{escapeAndCoerceToInterchangeValid(text)} <== copied from HTMLEscaped()
}

// EmptyObjectJSON returns an empty JSON object '{}' as safe JSON
func EmptyObjectJSON() JSON {
	return JSON{"{}"}
}

// EmptyArrayJSON returns an empty JSON array '[]' as safe JSON
func EmptyArrayJSON() JSON {
	return JSON{"[]"}
}

// JSONConcat returns an JSON which contains, in order, the string representations
// of the given jsons.
func JSONConcat(jsons ...JSON) JSON {
	var b bytes.Buffer
	for _, json := range jsons {
		b.WriteString(json.String())
	}
	return JSON{b.String()}
}

// String returns the string form of the JSON.
func (h JSON) String() string {
	return h.str
}
