#
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
#

"""
This file contains metric cell classes. A metric cell is used to accumulate
in-memory changes to a metric. It represents a specific metric in a single
context.
"""

# pytype: skip-file

import logging
import threading
import time
from datetime import datetime
from typing import Iterable
from typing import Optional
from typing import Set

try:
  import cython
except ImportError:

  class fake_cython:
    compiled = False

  globals()['cython'] = fake_cython

__all__ = [
    'MetricCell', 'MetricCellFactory', 'DistributionResult', 'GaugeResult'
]

_LOGGER = logging.getLogger(__name__)


class MetricCell(object):
  """For internal use only; no backwards-compatibility guarantees.

  Accumulates in-memory changes to a metric.

  A MetricCell represents a specific metric in a single context and bundle.
  All subclasses must be thread safe, as these are used in the pipeline runners,
  and may be subject to parallel/concurrent updates. Cells should only be used
  directly within a runner.
  """
  def __init__(self):
    self._lock = threading.Lock()
    self._start_time = None

  def update(self, value):
    raise NotImplementedError

  def get_cumulative(self):
    raise NotImplementedError

  def to_runner_api_monitoring_info(self, name, transform_id):
    if not self._start_time:
      self._start_time = datetime.utcnow()
    mi = self.to_runner_api_monitoring_info_impl(name, transform_id)
    mi.start_time.FromDatetime(self._start_time)
    return mi

  def to_runner_api_monitoring_info_impl(self, name, transform_id):
    raise NotImplementedError

  def reset(self):
    # type: () -> None
    raise NotImplementedError

  def __reduce__(self):
    raise NotImplementedError


class MetricCellFactory(object):
  def __call__(self):
    # type: () -> MetricCell
    raise NotImplementedError


class CounterCell(MetricCell):
  """For internal use only; no backwards-compatibility guarantees.

  Tracks the current value and delta of a counter metric.

  Each cell tracks the state of an integer metric independently per context
  per bundle. Therefore, each metric has a different cell in each bundle,
  cells are aggregated by the runner.

  This class is thread safe.
  """
  def __init__(self, *args):
    super().__init__(*args)
    self.value = 0

  def reset(self):
    # type: () -> None
    self.value = 0

  def combine(self, other):
    # type: (CounterCell) -> CounterCell
    result = CounterCell()
    result.inc(self.value + other.value)
    return result

  def inc(self, n=1):
    self.update(n)

  def dec(self, n=1):
    self.update(-n)

  def update(self, value):
    # type: (int) -> None
    if cython.compiled:
      ivalue = value
      # Since We hold the GIL, no need for another lock.
      # And because the C threads won't preempt and interleave
      # each other.
      # Assuming there is no code trying to access the counters
      # directly by circumventing the GIL.
      self.value += ivalue
    else:
      with self._lock:
        self.value += value

  def get_cumulative(self):
    # type: () -> int
    with self._lock:
      return self.value

  def to_runner_api_monitoring_info_impl(self, name, transform_id):
    from apache_beam.metrics import monitoring_infos
    if not name.urn:
      # User counter case.
      return monitoring_infos.int64_user_counter(
          name.namespace,
          name.name,
          self.get_cumulative(),
          ptransform=transform_id)
    else:
      # Arbitrary URN case.
      return monitoring_infos.int64_counter(
          name.urn, self.get_cumulative(), labels=name.labels)


class DistributionCell(MetricCell):
  """For internal use only; no backwards-compatibility guarantees.

  Tracks the current value and delta for a distribution metric.

  Each cell tracks the state of a metric independently per context per bundle.
  Therefore, each metric has a different cell in each bundle, that is later
  aggregated.

  This class is thread safe.
  """
  def __init__(self, *args):
    super().__init__(*args)
    self.data = DistributionData.identity_element()

  def reset(self):
    # type: () -> None
    self.data = DistributionData.identity_element()

  def combine(self, other):
    # type: (DistributionCell) -> DistributionCell
    result = DistributionCell()
    result.data = self.data.combine(other.data)
    return result

  def update(self, value):
    if cython.compiled:
      # We will hold the GIL throughout the entire _update.
      self._update(value)
    else:
      with self._lock:
        self._update(value)

  def _update(self, value):
    if cython.compiled:
      ivalue = value
    else:
      ivalue = int(value)
    self.data.count = self.data.count + 1
    self.data.sum = self.data.sum + ivalue
    if ivalue < self.data.min:
      self.data.min = ivalue
    if ivalue > self.data.max:
      self.data.max = ivalue

  def get_cumulative(self):
    # type: () -> DistributionData
    with self._lock:
      return self.data.get_cumulative()

  def to_runner_api_monitoring_info_impl(self, name, transform_id):
    from apache_beam.metrics import monitoring_infos
    return monitoring_infos.int64_user_distribution(
        name.namespace,
        name.name,
        self.get_cumulative(),
        ptransform=transform_id)


class AbstractMetricCell(MetricCell):
  """For internal use only; no backwards-compatibility guarantees.

  Tracks the current value and delta for a metric with a data class.

  This class is thread safe.
  """
  def __init__(self, data_class):
    super().__init__()
    self.data_class = data_class
    self.data = self.data_class.identity_element()

  def reset(self):
    self.data = self.data_class.identity_element()

  def combine(self, other: 'AbstractMetricCell') -> 'AbstractMetricCell':
    result = type(self)()  # type: ignore[call-arg]
    result.data = self.data.combine(other.data)
    return result

  def set(self, value):
    with self._lock:
      self._update_locked(value)

  def update(self, value):
    with self._lock:
      self._update_locked(value)

  def _update_locked(self, value):
    raise NotImplementedError(type(self))

  def get_cumulative(self):
    with self._lock:
      return self.data.get_cumulative()

  def to_runner_api_monitoring_info_impl(self, name, transform_id):
    raise NotImplementedError(type(self))


class GaugeCell(AbstractMetricCell):
  """For internal use only; no backwards-compatibility guarantees.

  Tracks the current value and delta for a gauge metric.

  Each cell tracks the state of a metric independently per context per bundle.
  Therefore, each metric has a different cell in each bundle, that is later
  aggregated.

  This class is thread safe.
  """
  def __init__(self):
    super().__init__(GaugeData)

  def _update_locked(self, value):
    # Set the value directly without checking timestamp, because
    # this value is naturally the latest value.
    self.data.value = int(value)
    self.data.timestamp = time.time()

  def to_runner_api_monitoring_info_impl(self, name, transform_id):
    from apache_beam.metrics import monitoring_infos
    return monitoring_infos.int64_user_gauge(
        name.namespace,
        name.name,
        self.get_cumulative(),
        ptransform=transform_id)


class StringSetCell(AbstractMetricCell):
  """For internal use only; no backwards-compatibility guarantees.

  Tracks the current value for a StringSet metric.

  Each cell tracks the state of a metric independently per context per bundle.
  Therefore, each metric has a different cell in each bundle, that is later
  aggregated.

  This class is thread safe.
  """
  def __init__(self):
    super().__init__(StringSetData)

  def add(self, value):
    self.update(value)

  def _update_locked(self, value):
    self.data.add(value)

  def to_runner_api_monitoring_info_impl(self, name, transform_id):
    from apache_beam.metrics import monitoring_infos
    return monitoring_infos.user_set_string(
        name.namespace,
        name.name,
        self.get_cumulative(),
        ptransform=transform_id)


class DistributionResult(object):
  """The result of a Distribution metric."""
  def __init__(self, data):
    # type: (DistributionData) -> None
    self.data = data

  def __eq__(self, other):
    # type: (object) -> bool
    if isinstance(other, DistributionResult):
      return self.data == other.data
    else:
      return False

  def __hash__(self):
    # type: () -> int
    return hash(self.data)

  def __repr__(self):
    # type: () -> str
    return (
        'DistributionResult(sum={}, count={}, min={}, max={}, '
        'mean={})'.format(self.sum, self.count, self.min, self.max, self.mean))

  @property
  def max(self):
    # type: () -> Optional[int]
    return self.data.max if self.data.count else None

  @property
  def min(self):
    # type: () -> Optional[int]
    return self.data.min if self.data.count else None

  @property
  def count(self):
    # type: () -> Optional[int]
    return self.data.count

  @property
  def sum(self):
    # type: () -> Optional[int]
    return self.data.sum

  @property
  def mean(self):
    # type: () -> Optional[float]

    """Returns the float mean of the distribution.

    If the distribution contains no elements, it returns None.
    """
    if self.data.count == 0:
      return None
    return self.data.sum / self.data.count


class GaugeResult(object):
  def __init__(self, data):
    # type: (GaugeData) -> None
    self.data = data

  def __eq__(self, other):
    # type: (object) -> bool
    if isinstance(other, GaugeResult):
      return self.data == other.data
    else:
      return False

  def __hash__(self):
    # type: () -> int
    return hash(self.data)

  def __repr__(self):
    return '<GaugeResult(value={}, timestamp={})>'.format(
        self.value, self.timestamp)

  @property
  def value(self):
    # type: () -> Optional[int]
    return self.data.value

  @property
  def timestamp(self):
    # type: () -> Optional[int]
    return self.data.timestamp


class GaugeData(object):
  """For internal use only; no backwards-compatibility guarantees.

  The data structure that holds data about a gauge metric.

  Gauge metrics are restricted to integers only.

  This object is not thread safe, so it's not supposed to be modified
  by other than the GaugeCell that contains it.
  """
  def __init__(self, value, timestamp=None):
    # type: (Optional[int], Optional[int]) -> None
    self.value = value
    self.timestamp = timestamp if timestamp is not None else 0

  def __eq__(self, other):
    # type: (object) -> bool
    if isinstance(other, GaugeData):
      return self.value == other.value and self.timestamp == other.timestamp
    else:
      return False

  def __hash__(self):
    # type: () -> int
    return hash((self.value, self.timestamp))

  def __repr__(self):
    # type: () -> str
    return '<GaugeData(value={}, timestamp={})>'.format(
        self.value, self.timestamp)

  def get_cumulative(self):
    # type: () -> GaugeData
    return GaugeData(self.value, timestamp=self.timestamp)

  def get_result(self):
    # type: () -> GaugeResult
    return GaugeResult(self.get_cumulative())

  def combine(self, other):
    # type: (Optional[GaugeData]) -> GaugeData
    if other is None:
      return self

    if other.timestamp > self.timestamp:
      return other
    else:
      return self

  @staticmethod
  def singleton(value, timestamp=None):
    # type: (Optional[int], Optional[int]) -> GaugeData
    return GaugeData(value, timestamp=timestamp)

  @staticmethod
  def identity_element():
    # type: () -> GaugeData
    return GaugeData(0, timestamp=0)


class DistributionData(object):
  """For internal use only; no backwards-compatibility guarantees.

  The data structure that holds data about a distribution metric.

  Distribution metrics are restricted to distributions of integers only.

  This object is not thread safe, so it's not supposed to be modified
  by other than the DistributionCell that contains it.
  """
  def __init__(self, sum, count, min, max):
    # type: (int, int, int, int) -> None
    if count:
      self.sum = sum
      self.count = count
      self.min = min
      self.max = max
    else:
      self.sum = self.count = 0
      self.min = 2**63 - 1
      # Avoid Wimplicitly-unsigned-literal caused by -2**63.
      self.max = -self.min - 1

  def __eq__(self, other):
    # type: (object) -> bool
    if isinstance(other, DistributionData):
      return (
          self.sum == other.sum and self.count == other.count and
          self.min == other.min and self.max == other.max)
    else:
      return False

  def __hash__(self):
    # type: () -> int
    return hash((self.sum, self.count, self.min, self.max))

  def __repr__(self):
    # type: () -> str
    return 'DistributionData(sum={}, count={}, min={}, max={})'.format(
        self.sum, self.count, self.min, self.max)

  def get_cumulative(self):
    # type: () -> DistributionData
    return DistributionData(self.sum, self.count, self.min, self.max)

  def get_result(self) -> DistributionResult:
    return DistributionResult(self.get_cumulative())

  def combine(self, other):
    # type: (Optional[DistributionData]) -> DistributionData
    if other is None:
      return self

    return DistributionData(
        self.sum + other.sum,
        self.count + other.count,
        self.min if self.min < other.min else other.min,
        self.max if self.max > other.max else other.max)

  @staticmethod
  def singleton(value):
    # type: (int) -> DistributionData
    return DistributionData(value, 1, value, value)

  @staticmethod
  def identity_element():
    # type: () -> DistributionData
    return DistributionData(0, 0, 2**63 - 1, -2**63)


class StringSetData(object):
  """For internal use only; no backwards-compatibility guarantees.

  The data structure that holds data about a StringSet metric.

  StringSet metrics are restricted to set of strings only.

  This object is not thread safe, so it's not supposed to be modified
  by other than the StringSetCell that contains it.

  The summation of all string length for a StringSetData cannot exceed 1 MB.
  Further addition of elements are dropped.
  """

  _STRING_SET_SIZE_LIMIT = 1_000_000

  def __init__(self, string_set: Optional[Set] = None, string_size: int = 0):
    self.string_set = string_set or set()
    if not string_size:
      string_size = 0
      for s in self.string_set:
        string_size += len(s)
    self.string_size = string_size

  def __eq__(self, other: object) -> bool:
    if isinstance(other, StringSetData):
      return (
          self.string_size == other.string_size and
          self.string_set == other.string_set)
    else:
      return False

  def __hash__(self) -> int:
    return hash(self.string_set)

  def __repr__(self) -> str:
    return 'StringSetData{}:{}'.format(self.string_set, self.string_size)

  def get_cumulative(self) -> "StringSetData":
    return StringSetData(set(self.string_set), self.string_size)

  def get_result(self) -> set[str]:
    return set(self.string_set)

  def add(self, *strings):
    """
    Add strings into this StringSetData and return the result StringSetData.
    Reuse the original StringSetData's set.
    """
    self.string_size = self.add_until_capacity(
        self.string_set, self.string_size, strings)
    return self

  def combine(self, other: "StringSetData") -> "StringSetData":
    """
    Combines this StringSetData with other, both original StringSetData are left
    intact.
    """
    if other is None:
      return self

    if not other.string_set:
      return self
    elif not self.string_set:
      return other

    combined = set(self.string_set)
    string_size = self.add_until_capacity(
        combined, self.string_size, other.string_set)
    return StringSetData(combined, string_size)

  @classmethod
  def add_until_capacity(
      cls, combined: set, current_size: int, others: Iterable[str]):
    """
    Add strings into set until reach capacity. Return the all string size of
    added set.
    """
    if current_size > cls._STRING_SET_SIZE_LIMIT:
      return current_size

    for string in others:
      if string not in combined:
        combined.add(string)
        current_size += len(string)
        if current_size > cls._STRING_SET_SIZE_LIMIT:
          _LOGGER.warning(
              "StringSet metrics reaches capacity. Further incoming elements "
              "won't be recorded. Current size: %d, last element size: %d.",
              current_size,
              len(string))
          break
    return current_size

  @staticmethod
  def singleton(value: str) -> "StringSetData":
    return StringSetData({value})

  @staticmethod
  def identity_element() -> "StringSetData":
    return StringSetData()
