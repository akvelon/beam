/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import 'package:playground_components/playground_components.dart';

const _none = 'none';

/// [EventSnippetContext] with [unitId].
class TobEventSnippetContext extends EventSnippetContext {
  final String? unitId;

  TobEventSnippetContext({
    super.originalSnippet,
    super.sdk,
    super.snippet,
    required this.unitId,
  });

  factory TobEventSnippetContext.fromEventSnippetContext({
    required EventSnippetContext eventSnippetContext,
    required String? unitId,
  }) {
    return TobEventSnippetContext(
      originalSnippet: eventSnippetContext.originalSnippet,
      sdk: eventSnippetContext.sdk,
      snippet: eventSnippetContext.snippet,
      unitId: unitId,
    );
  }

  @override
  List<Object?> get props => [
        ...super.props,
        unitId,
      ];

  @override
  Map<String, dynamic> toJson() {
    final map = super.toJson();
    map['unitId'] = unitId ?? _none;
    return map;
  }
}
