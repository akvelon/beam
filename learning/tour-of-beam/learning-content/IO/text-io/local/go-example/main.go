// Licensed to the Apache Software Foundation (ASF) under one or more
// contributor license agreements.  See the NOTICE file distributed with
// this work for additional information regarding copyright ownership.
// The ASF licenses this file to You under the Apache License, Version 2.0
// (the "License"); you may not use this file except in compliance with
// the License.  You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// beam-playground:
//   name: text-io-local
//   description: TextIO local example.
//   multifile: false
//   context_line: 30
//   categories:
//     - Quickstart
//   complexity: ADVANCED
//   tags:
//     - hellobeam


package main

import (
	"regexp"
	"github.com/apache/beam/sdks/v2/go/pkg/beam"
	"github.com/apache/beam/sdks/v2/go/pkg/beam/io/textio"
    "context"
    "github.com/apache/beam/sdks/v2/go/pkg/beam/log"
    "github.com/apache/beam/sdks/v2/go/pkg/beam/x/beamx"
)

var wordRE = regexp.MustCompile(`[a-zA-Z]+('[a-z])?`)


func main() {
p, s := beam.NewPipelineWithRoot()
	lines := textio.Read(s, "myfile.txt")

words := beam.ParDo(s, func(line string, emit func(string)) {
     for _, word := range wordRE.FindAllString(line, -1) {
        emit(word)
      }
}, lines)

    textio.Write(s, "output.txt", words)

err := beamx.Run(context.Background(), p)
    if err != nil {
        log.Exitf(context.Background(), "Failed to execute job: %v", err)
    }
}