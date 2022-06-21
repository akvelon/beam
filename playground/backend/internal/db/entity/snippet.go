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

package entity

import (
	"beam.apache.org/playground/backend/internal/utils"
	"cloud.google.com/go/datastore"
	"fmt"
	"sort"
	"strings"
	"time"
)

type Origin int32

const (
	PG_USER      Origin = 0
	PG_EXAMPLES  Origin = 1
	TOUR_OF_BEAM Origin = 2
)

// OriginValue value maps for Origin.
var (
	OriginValue = map[string]int32{
		"PG_USER":      0,
		"PG_EXAMPLES":  1,
		"TOUR_OF_BEAM": 2,
	}
)

func (s Origin) Value() int32 {
	return int32(s)
}

type CodeEntity struct {
	Name     string `datastore:"name"`
	Code     string `datastore:"code,noindex"`
	CntxLine int32  `datastore:"cntxLine"`
	IsMain   bool   `datastore:"isMain"`
}

type SnippetEntity struct {
	OwnerId    string         `datastore:"ownerId"`
	Sdk        *datastore.Key `datastore:"sdk"`
	PipeOpts   string         `datastore:"pipeOpts"`
	Created    time.Time      `datastore:"created"`
	LVisited   time.Time      `datastore:"lVisited"`
	Origin     Origin         `datastore:"origin"`
	VisitCount int            `datastore:"visitCount"`
	SchVer     *datastore.Key `datastore:"schVer"`
}

type Snippet struct {
	*IDInfo
	Snippet *SnippetEntity
	Codes   []*CodeEntity
}

// ID generates id according to content of the entity
func (s *Snippet) ID() (string, error) {
	var codes []string
	for _, v := range s.Codes {
		codes = append(codes, strings.TrimSpace(v.Code)+strings.TrimSpace(v.Name))
	}
	sort.Strings(codes)
	var content string
	for i, v := range codes {
		content += v
		if i == len(codes)-1 {
			content += fmt.Sprintf("%v%s", s.Snippet.Sdk, strings.TrimSpace(s.Snippet.PipeOpts))
		}
	}
	id, err := utils.ID(s.Salt, content, s.IdLength)
	if err != nil {
		return "", err
	}
	return id, nil
}

// ID generates id according to content of a code and its name
func (c *CodeEntity) ID(snip *Snippet) (string, error) {
	content := fmt.Sprintf("%s%s", strings.TrimSpace(c.Code), strings.TrimSpace(c.Name))
	id, err := utils.ID(snip.Salt, content, snip.IdLength)
	if err != nil {
		return "", err
	}
	return id, nil
}
