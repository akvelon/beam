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

import 'package:flutter/cupertino.dart';
import 'package:playground/config.g.dart';
import 'package:playground/modules/analytics/analytics_events.dart';
import 'package:playground/modules/examples/models/example_model.dart';
import 'package:playground/modules/sdk/models/sdk.dart';
import 'package:provider/provider.dart';
import 'package:usage/usage_html.dart';

class AnalyticsService {
  late final Analytics _analytics;

  AnalyticsService() {
    _analytics = AnalyticsHtml(kAnalyticsUA, 'beam', '1.0');
  }

  static get(BuildContext context) {
    return Provider.of<AnalyticsService>(context, listen: false);
  }

  void trackSelectSdk(SDK oldSdk, SDK newSdk) {
    _analytics.sendEvent(
      kCategory,
      kSelectSdkEvent,
      parameters: {'oldSdk': oldSdk.displayName, 'newSdk': newSdk.displayName},
    );
  }

  void trackSelectExample(ExampleModel oldExample) {
    _analytics.sendEvent(
      kCategory,
      kSelectExampleEvent,
      parameters: {'example': oldExample.path},
    );
  }

  void trackClickNewExample() {
    _analytics.sendEvent(kCategory, kClickNewExampleEvent);
  }

  void trackResetExample() {
    _analytics.sendEvent(kCategory, kClickResetEvent);
  }

  void trackClickToggleTheme(bool isDark) {
    _analytics.sendEvent(
      kCategory,
      kClickToggleThemeEvent,
      parameters: {'isDark': isDark.toString()},
    );
  }

  void trackOpenShortcutsModal() {
    _analytics.sendEvent(kCategory, kOpenShortcutsModalEvent);
  }

  void trackOpenLink(String link) {
    _analytics.sendEvent(
      kCategory,
      kOpenShortcutsModalEvent,
      parameters: {'link': link},
    );
  }

  void trackClickEnjoyPlayground(bool isEnjoying) {
    _analytics.sendEvent(
      kCategory,
      kClickEnjoyPlaygroundEvent,
      parameters: {'isEnjoying': isEnjoying.toString()},
    );
  }

  void trackClickReportIssue() {
    _analytics.sendEvent(kCategory, kClickReportIssueEvent);
  }

  void trackClickRunEvent(ExampleModel? example) {
    _analytics.sendEvent(
      kCategory,
      kClickRunEvent,
      parameters: {'example': example?.path ?? ''},
    );
  }
}
