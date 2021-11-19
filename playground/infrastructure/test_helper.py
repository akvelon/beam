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

import mock
import pytest

from unittest.mock import mock_open
from api.v1.api_pb2 import SDK_UNSPECIFIED, STATUS_UNSPECIFIED, SDK_JAVA, SDK_PYTHON, SDK_SCIO, SDK_GO
from helper import find_examples, Example, _get_example, _get_name, _get_sdk, get_tag, _validate


@mock.patch('helper._get_example')
@mock.patch('helper._validate')
@mock.patch('helper.get_tag')
@mock.patch('helper.os.walk')
def test_find_examples_with_valid_tag(mock_os_walk, mock_get_tag, mock_validate, mock_get_example):
    example = Example("file", "pipeline_id", SDK_UNSPECIFIED, "root/file.extension", "code", "output",
                      STATUS_UNSPECIFIED)
    mock_os_walk.return_value = [("/root", (), ("file.java",))]
    mock_get_tag.return_value = {"name": "Name"}
    mock_validate.return_value = True
    mock_get_example.return_value = example

    result = find_examples("", "")

    assert result == [example]
    mock_os_walk.assert_called_once_with("")
    mock_get_tag.assert_called_once_with("/root/file.java")
    mock_validate.assert_called_once_with({"name": "Name"}, "")
    mock_get_example.assert_called_once_with("/root/file.java", "file.java")


@mock.patch('helper._validate')
@mock.patch('helper.get_tag')
@mock.patch('helper.os.walk')
def test_find_examples_with_invalid_tag(mock_os_walk, mock_get_tag, mock_validate):
    mock_os_walk.return_value = [("/root", (), ("file.java",))]
    mock_get_tag.return_value = {"name": "Name"}
    mock_validate.return_value = False

    with pytest.raises(ValueError, match="some of the beam examples contain beam playground tag with incorrect format"):
        find_examples("", "")

    mock_os_walk.assert_called_once_with("")
    mock_get_tag.assert_called_once_with("/root/file.java")
    mock_validate.assert_called_once_with({"name": "Name"}, "")


@mock.patch('builtins.open', mock_open(read_data="...\n# Beam-playground:\n#     name: Name\n..."))
def test_get_tag_when_tag_is_exists():
    result = get_tag("")

    assert result.get("name") == "Name"


@mock.patch('builtins.open', mock_open(read_data="...\n..."))
def test_get_tag_when_tag_does_not_exist():
    result = get_tag("")

    assert result is None


@mock.patch('builtins.open', mock_open(read_data="data"))
@mock.patch('helper._get_sdk')
@mock.patch('helper._get_name')
def test__get_example(mock_get_name, mock_get_sdk):
    mock_get_name.return_value = "filepath"
    mock_get_sdk.return_value = SDK_UNSPECIFIED

    result = _get_example("/root/filepath.extension", "filepath.extension")

    assert result == Example("filepath", "", SDK_UNSPECIFIED, "/root/filepath.extension", "data", "",
                             STATUS_UNSPECIFIED)
    mock_get_name.assert_called_once_with("filepath.extension")
    mock_get_sdk.assert_called_once_with("filepath.extension")


def test__validate_without_name_field():
    tag = {}
    assert _validate(tag, "") is False


def test__validate_without_description_field():
    tag = {"name": "Name"}
    assert _validate(tag, "") is False


def test__validate_without_multifile_field():
    tag = {"name": "Name", "description": "Description"}
    assert _validate(tag, "") is False


def test__validate_with_incorrect_multifile_field():
    tag = {"name": "Name", "description": "Description", "multifile": "Multifile"}
    assert _validate(tag, "") is False


def test__validate_without_categories_field():
    tag = {"name": "Name", "description": "Description", "multifile": "true"}
    assert _validate(tag, "") is False


def test__validate_without_incorrect_categories_field():
    tag = {"name": "Name", "description": "Description", "multifile": "true", "categories": "Categories"}
    assert _validate(tag, "") is False


@mock.patch('builtins.open', mock_open(read_data="categories:\n    - category"))
def test__validate_with_not_supported_category():
    tag = {"name": "Name", "description": "Description", "multifile": "true", "categories": ["category1"]}
    assert _validate(tag, "") is False


@mock.patch('builtins.open', mock_open(read_data="categories:\n    - category"))
def test__validate_with_all_fields():
    tag = {"name": "Name", "description": "Description", "multifile": "true", "categories": ["category"]}
    assert _validate(tag, "") is True


def test__get_name():
    result = _get_name("filepath.extension")

    assert result == "filepath"


def test__get_sdk_with_supported_extension():
    assert _get_sdk("filename.java") == SDK_JAVA
    assert _get_sdk("filename.go") == SDK_GO
    assert _get_sdk("filename.py") == SDK_PYTHON


def test__get_sdk_with_unsupported_extension():
    with pytest.raises(ValueError, match="extension is not supported"):
        _get_sdk("filename.extension")
