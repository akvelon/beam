# Licensed to the Apache Software Foundation (ASF) under one or more
# contributor license agreements.  See the NOTICE file distributed with
# this work for additional information regarding copyright ownership.
# The ASF licenses this file to You under the Apache License, Version 2.0
# (the "License"); you may not use this file except in compliance with
# the License.  You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

import os
import yaml

from dataclasses import dataclass
from typing import List
from api.v1.api_pb2 import Sdk, SDK_JAVA, SDK_UNSPECIFIED, STATUS_UNSPECIFIED

SUPPORTED_SDK = {'java': SDK_JAVA}
END_OF_LICENCE = "limitations under the License."
START_OF_IMPORT = "\nimport "
BEAM_PLAYGROUND = "Beam-playground"


@dataclass
class Example:
    """ Class which contains all information about beam example
    """
    name: str
    pipeline_id: str
    sdk: SDK_UNSPECIFIED
    filepath: str
    code: str
    output: str
    status: STATUS_UNSPECIFIED


def find_examples(work_dir: str) -> List[Example]:
    """ Find and return beam examples.

    Search throws all child files of work_dir directory files with beam tag:
    Beam-playground:
        name: NameOfExample
        description: Description of NameOfExample.
        multifile: false
        categories:
            - category-1
            - category-2

    Args:
        work_dir: directory where to search examples.

    Returns:
        List of Examples.
    """
    examples = []
    for root, _, files in os.walk(work_dir):
        for filename in files:
            filepath = os.path.join(root, filename)
            if _match_pattern(filepath):
                examples.append(_get_example(filepath, filename))
    return examples


def get_statuses(examples: [Example]):
    """ Receive statuses for examples and update example.status and example.pipeline_id

    Use client to send requests to the backend:
    1. Start code processing.
    2. Ping the backend while status is STATUS_VALIDATING/STATUS_PREPARING/STATUS_COMPILING/STATUS_EXECUTING
    Update example.pipeline_id with resulting pipelineId.
    Update example.status with resulting status.

    Args:
        examples: beam examples for processing and updating statuses.
    """
    # TODO [BEAM-13267] Implement
    pass


def get_tag(filepath):
    """Parse file by filepath and find beam tag (after the licence part, before the import part)

    If file contains tag, returns tag as a map.
    If file contains tag but tag has incorrect format, raise error
    If file doesn't contain tag, returns None

    Args:
        filepath: path of the file
    """
    with open(filepath) as parsed_file:
        content = parsed_file.read()
    index_of_licence_end = content.find(END_OF_LICENCE) + len(END_OF_LICENCE)
    index_of_import_start = content.find(START_OF_IMPORT)
    content = content[index_of_licence_end:index_of_import_start]
    index_of_tag_start = content.find(BEAM_PLAYGROUND)
    if index_of_tag_start < 0:
        return None
    content = content[index_of_tag_start:]
    yaml_tag = content.replace("//", "").replace("#", "")
    try:
        object_meta = yaml.load(yaml_tag, Loader=yaml.SafeLoader)
        return object_meta[BEAM_PLAYGROUND]
    except Exception as exp:
        print(exp)  ## todo add logErr
        raise ValueError("found tag is not correct: " + exp.__str__())


def _get_example(filepath: str, filename: str) -> Example:
    """ Return an Example by filepath and filename.

    Args:
         filepath: path of the example's file.
         filename: name of the example's file.

    Returns:
        Return an Example.
    """
    name = _get_name(filename)
    sdk = _get_sdk(filename)
    with open(filepath) as parsed_file:
        content = parsed_file.read()

    return Example(name, "", sdk, filepath, content, "", STATUS_UNSPECIFIED)


def _match_pattern(filepath: str) -> bool:
    """Check file to matching

    Check that file has the correct extension and contains the beam-playground tag.

    Args:
        filepath: path to the file.

    Returns:
        True if file matched. False if not
    """
    extension = filepath.split(os.extsep)[-1]
    if extension in SUPPORTED_SDK:
        tag = get_tag(filepath)
        if tag is None:
            return False
        _validate(tag)
        return True
    return False


def _validate(tag: dict):
    """Validate all tag's fields

    If some of the fields has incorrect format, raise error

    Args:
        tag: beam tag to validate
    """
    if tag.get("name") is None:
        raise ValueError("tag doesn't contain name field: " + tag.__str__())
    if tag.get("description") is None:
        raise ValueError("tag doesn't contain description field: " + tag.__str__())
    if tag.get("multifile") is None:
        raise ValueError("tag doesn't contain multifile field: " + tag.__str__())
    multifile = tag.get("multifile")
    if str(multifile).lower() not in ["true", "false"]:
        raise ValueError("tag's field multifile is incorrect: " + tag.__str__())
    if tag.get("categories") is None:
        raise ValueError("tag doesn't contain categories field: " + tag.__str__())


def _get_name(filename) -> str:
    """ Return name of the example by his filepath.

    Get name of the example by his filename.

    Args:
        filename: filename of the beam example file.

    Returns:
        example's name.
    """
    return filename.split(os.extsep)[0]


def _get_sdk(filename) -> Sdk:
    """ Return SDK of example by his filename.

    Get extension of the example's file and returns associated SDK.

    Args:
        filename: filename of the beam example.

    Returns:
        Sdk according to file extension.
    """
    extension = filename.split(os.extsep)[-1]
    if extension in SUPPORTED_SDK:
        return SUPPORTED_SDK[extension]
    else:
        raise ValueError(extension + " is not supported now")
