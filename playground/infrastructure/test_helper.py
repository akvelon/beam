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

from unittest.mock import mock_open
from api.v1.api_pb2 import SDK_UNSPECIFIED, STATUS_UNSPECIFIED, SDK_JAVA
from helper import find_examples, Example, _get_example, _match_pattern, _get_name, _get_sdk, get_tag, _validate


@mock.patch('helper._get_example')
@mock.patch('helper._match_pattern')
@mock.patch('helper.os.path.join')
@mock.patch('helper.os.walk')
def test_find_examples(self, mock_os_walk, mock_os_path_join, mock_match_pattern, mock_get_example):
    example = Example("file", "pipeline_id", SDK_UNSPECIFIED, "root/file.extension", "code", "output",
                      STATUS_UNSPECIFIED)
    mock_os_walk.return_value = [("/root", (), ("file.extension",))]
    mock_os_path_join.return_value = "/root/file.extension"
    mock_match_pattern.return_value = True
    mock_get_example.return_value = example

    find_examples("")

    mock_os_walk.assert_called_once_with("")
    mock_os_path_join.assert_any_call("/root", "file.extension")
    mock_match_pattern.assert_called_once_with("/root/file.extension")
    mock_get_example.assert_called_once_with("/root/file.extension", "file.extension")


@mock.patch('builtins.open',
            mock_open(
                read_data="...\n# limitations under the License.\n# Beam-playground:\n#     name: Name\n\nimport ..."))
def test_get_tag_when_tag_is_exists(self):
    result = get_tag("")

    assert result.get("name") == "Name"


@mock.patch('builtins.open',
            mock_open(
                read_data="...\n# limitations under the License.\n\nimport ..."))
def test_get_tag_when_tag_does_not_exist(self):
    result = get_tag("")

    assert result is None


@mock.patch('builtins.open',
            mock_open(
                read_data="...\n# limitations under the License.\n# Beam-playground:\n# additional string\n#     name: Name\n\nimport ..."))
def test_get_tag_when_tag_has_incorrect_format(self):
    try:
        get_tag("")
        assert False
    except Exception:
        assert True


@mock.patch('builtins.open', mock_open(read_data="data"))
@mock.patch('helper._get_sdk')
@mock.patch('helper._get_name')
def test__get_example(self, mock_get_name, mock_get_sdk):
    mock_get_name.return_value = "filepath"
    mock_get_sdk.return_value = SDK_UNSPECIFIED

    result = _get_example("/root/filepath.extension", "filepath.extension")

    assert result == Example("filepath", "", SDK_UNSPECIFIED, "/root/filepath.extension", "data", "",
                             STATUS_UNSPECIFIED)
    mock_get_name.assert_called_once_with("filepath.extension")
    mock_get_sdk.assert_called_once_with("filepath.extension")


@mock.patch('helper._validate')
@mock.patch('helper.get_tag')
def test__match_pattern_with_correct_tag(self, mock_get_tag, mock_validate):
    mock_get_tag.return_value = {}

    result = _match_pattern("/root/filepath.java")

    assert result
    mock_get_tag.assert_called_once_with("/root/filepath.java")
    mock_validate.assert_called_once_with({})


@mock.patch('helper.get_tag')
def test__match_pattern_without_tag(self, mock_get_tag):
    mock_get_tag.return_value = None

    result = _match_pattern("/root/filepath.java")

    assert result is False
    mock_get_tag.assert_called_once_with("/root/filepath.java")


def test__match_pattern_with_unsupported_extension(self):
    result = _match_pattern("/root/filepath.extension")

    assert result is False


def test__validate_without_name_field(self):
    try:
        _validate({})
        assert False
    except Exception:
        assert True


def test__validate_without_description_field(self):
    try:
        _validate({"name": "Name"})
        assert False
    except Exception:
        assert True


def test__validate_without_multifile_field(self):
    try:
        _validate({"name": "Name", "description": "Description"})
        assert False
    except Exception:
        assert True


def test__validate_with_incorrect_multifile_field(self):
    try:
        _validate({"name": "Name", "description": "Description", "multifile": "Multifile"})
        assert False
    except Exception:
        assert True


def test__validate_without_categories_field(self):
    try:
        _validate({"name": "Name", "description": "Description", "multifile": "true"})
        assert False
    except Exception:
        assert True


def test__validate_with_all_fields(self):
    _validate({"name": "Name", "description": "Description", "multifile": "true", "categories": []})
    assert True


def test__get_name(self):
    result = _get_name("filepath.extension")

    assert result == "filepath"


def test__get_sdk_with_supported_extension(self):
    assert _get_sdk("filename.java") == SDK_JAVA


def test__get_sdk_with_unsupported_extension(self):
    try:
        _get_sdk("filename.extension")
        assert False
    except Exception:
        assert True
