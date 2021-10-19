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

package preparators

import (
	"bufio"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// GetJavaPreparation return preparation methods that should be applied to Java code
func GetJavaPreparation(filePath string) *[]Preparator {
	preparatorArgs := make([]interface{}, 1)
	preparatorArgs[0] = filePath

	publicClassModification := Preparator{
		Prepare: removePublicClassModifier,
		Args:    preparatorArgs,
	}
	additionalPackage := Preparator{
		Prepare: removeAdditionalPackage,
		Args:    preparatorArgs,
	}
	preparation := []Preparator{publicClassModification, additionalPackage}
	return &preparation
}

// removePublicClassModifier removes public modification for class from java file
func removePublicClassModifier(args ...interface{}) error {
	filePath := args[0].(string)
	if err := replace(filePath, "public class ", "class "); err != nil {
		log.Printf("Preparation: Error during remove public modification for class, err: %s\n", err.Error())
		return err
	}
	return nil
}

// removeAdditionalPackage removes packages from java file
func removeAdditionalPackage(args ...interface{}) error {
	filePath := args[0].(string)
	if err := replace(filePath, `package ([\w]+\.)+[\w]+;`, ""); err != nil {
		log.Printf("Preparation: Error during remove additional package, err: %s\n", err.Error())
		return err
	}
	return nil
}

// createTempFile creates temporary file near with originalFile
func createTempFile(originalFilePath string) (*os.File, error) {
	// all folders which are included in filePath
	filePathSlice := strings.Split(originalFilePath, "/")
	fileName := filePathSlice[len(filePathSlice)-1]

	// find parent folder for file
	folderPath := "/"
	filePathSlice = filePathSlice[:len(filePathSlice)-1]
	for _, folder := range filePathSlice {
		if folder == "" {
			continue
		}
		folderPath = filepath.Join(folderPath, folder)
	}
	return os.Create(folderPath + "/tmp_" + fileName)
}

// replace process file by filePath and replace all patterns to new
func replace(filePath, pattern, new string) error {
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("Preparation: Error during open file: %s, err: %s\n", filePath, err.Error())
		return err
	}
	defer file.Close()

	tmp, err := createTempFile(filePath)
	if err != nil {
		log.Printf("Preparation: Error during create new temporary file, err: %s\n", err.Error())
		return err
	}
	defer tmp.Close()

	reg := regexp.MustCompile(pattern)
	scanner := bufio.NewScanner(file)
	firstLine := true
	for scanner.Scan() {
		if !firstLine {
			if _, err := io.WriteString(tmp, "\n"); err != nil {
				log.Printf("Preparation: Error during write \"\\n\" to tmp file, err: %s\n", err.Error())
				return err
			}
		}
		line := scanner.Text()
		matches := reg.FindAllString(line, -1)
		for _, str := range matches {
			line = strings.ReplaceAll(line, str, new)
		}
		if _, err = io.WriteString(tmp, line); err != nil {
			log.Printf("Preparation: Error during write \"%s\" to tmp file, err: %s\n", line, err.Error())
			return err
		}
		firstLine = false
	}
	if scanner.Err() != nil {
		return scanner.Err()
	}

	if err := os.Rename(tmp.Name(), filePath); err != nil {
		log.Printf("Preparation: Error during rename temporary file, err: %s\n", err.Error())
		return err
	}
	return nil
}
