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

import logging
import os
import yaml

from dataclasses import dataclass
from typing import List
from config import Config
from api.v1.api_pb2 import SDK_UNSPECIFIED, STATUS_UNSPECIFIED, Sdk

BEAM_PLAYGROUND_TITLE = "Beam-playground:\n"
BEAM_PLAYGROUND = "Beam-playground"
CATEGORIES = "categories"


@dataclass
class Example:
    """
    Class which contains all information about beam example
    """
    name: str
    pipeline_id: str
    sdk: SDK_UNSPECIFIED
    filepath: str
    code: str
    output: str
    status: STATUS_UNSPECIFIED


def find_examples(work_dir: str, categories_path: str) -> List[Example]:
    """
    Find and return beam examples.

    Search throws all child files of work_dir directory files with beam tag:
    Beam-playground:
        name: NameOfExample
        description: Description of NameOfExample.
        multifile: false
        categories:
            - category-1
            - category-2
    If some example contain beam tag with incorrect format raise an error.

    Args:
        work_dir: directory where to search examples.
        categories_path: path to the file with all supported categories.

    Returns:
        List of Examples.
    """
    failed = False
    examples = []
    for root, _, files in os.walk(work_dir):
        for filename in files:
            filepath = os.path.join(root, filename)
            extension = filepath.split(os.extsep)[-1]
            if extension in Config.SUPPORTED_SDK:
                tag = get_tag(filepath)
                if tag:
                    if _validate(tag, categories_path) is False:
                        logging.error(filepath + "contains beam playground tag with incorrect format")
                        failed = True
                    else:
                        examples.append(_get_example(filepath, filename))
    if failed:
        raise ValueError("some of the beam examples contain beam playground tag with incorrect format")
    return examples


def get_statuses(examples: List[Example]):
    """
    Receive statuses for examples and update example.status and example.pipeline_id

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
    """
    Parse file by filepath and find beam tag

    Args:
        filepath: path of the file

    Returns:
        If file contains tag, returns tag as a map.
        If file doesn't contain tag, returns None
    """
    add_to_yaml = False
    yaml_string = ""

    with open(filepath) as parsed_file:
        lines = parsed_file.readlines()

    for line in lines:
        line = line.replace("// ", "").replace("# ", "")
        if add_to_yaml is False:
            if line == BEAM_PLAYGROUND_TITLE:
                add_to_yaml = True
                yaml_string += line
        else:
            yaml_with_new_string = yaml_string + line
            try:
                yaml.load(yaml_with_new_string, Loader=yaml.SafeLoader)
                yaml_string += line
            except Exception:
                break

    if add_to_yaml:
        tag_object = yaml.load(yaml_string, Loader=yaml.SafeLoader)
        return tag_object[BEAM_PLAYGROUND]

    return None


def _get_example(filepath: str, filename: str) -> Example:
    """
    Return an Example by filepath and filename.

    Args:
         filepath: path of the example's file.
         filename: name of the example's file.

    Returns:
        Parsed Example object.
    """
    name = _get_name(filename)
    sdk = _get_sdk(filename)
    with open(filepath) as parsed_file:
        content = parsed_file.read()

    return Example(name, "", sdk, filepath, content, "", STATUS_UNSPECIFIED)


def _validate(tag: dict, categories_path: str) -> bool:
    """
    Validate all tag's fields

    Validate that tag contains all required fields and all fields have required format.

    Args:
        tag: beam tag to validate
        categories_path: path to the file with all supported categories.

    Returns:
        In case tag is valid, True
        In case tag is not valid, False
    """
    if tag.get("name") is None:
        logging.error("tag doesn't contain name field: " + tag.__str__())
        return False
    if tag.get("description") is None:
        logging.error("tag doesn't contain description field: " + tag.__str__())
        return False
    if tag.get("multifile") is None:
        logging.error("tag doesn't contain multifile field: " + tag.__str__())
        return False
    multifile = tag.get("multifile")
    if str(multifile).lower() not in ["true", "false"]:
        logging.error("tag's field multifile is incorrect: " + tag.__str__())
        return False
    if tag.get("categories") is None:
        logging.error("tag doesn't contain categories field: " + tag.__str__())
        return False
    categories = tag.get("categories")
    if type(categories) is not list:
        logging.error("tag's field categories is incorrect: " + tag.__str__())
        return False
    with open(categories_path) as supported_categories:
        yaml_object = yaml.load(supported_categories.read(), Loader=yaml.SafeLoader)
        supported_categories = yaml_object[CATEGORIES]
        result = True
        for category in categories:
            if category not in supported_categories:
                logging.error("tag contains unsupported category: " + category)
                result = False
        if result is False:
            return False
    return True


def _get_name(filename: str) -> str:
    """
    Return name of the example by his filepath.

    Get name of the example by his filename.

    Args:
        filename: filename of the beam example file.

    Returns:
        example's name.
    """
    return filename.split(os.extsep)[0]


def _get_sdk(filename: str) -> Sdk:
    """
    Return SDK of example by his filename.

    Get extension of the example's file and returns associated SDK.

    Args:
        filename: filename of the beam example.

    Returns:
        Sdk according to file extension.
    """
    extension = filename.split(os.extsep)[-1]
    if extension in Config.SUPPORTED_SDK:
        return Config.SUPPORTED_SDK[extension]
    else:
        raise ValueError(extension + " is not supported")
