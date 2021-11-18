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

import unittest

from api.v1.api_pb2 import SDK_UNSPECIFIED, STATUS_UNSPECIFIED
from helper import find_examples, Example


class Test(unittest.TestCase):
    def test_find_examples(self):
        result = find_examples("")

        assert len(result) == 1
        assert result[0] == Example("", "", SDK_UNSPECIFIED, "", "", "", STATUS_UNSPECIFIED)
