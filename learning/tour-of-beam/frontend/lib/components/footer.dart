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

import 'dart:async';

import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:get_it/get_it.dart';
import 'package:playground_components/playground_components.dart';
import 'package:url_launcher/url_launcher.dart';

import '../constants/sizes.dart';
import '../state.dart';

class Footer extends StatelessWidget {
  const Footer();

  @override
  Widget build(BuildContext context) {
    return _Body(
      child: Wrap(
        alignment: WrapAlignment.spaceBetween,
        crossAxisAlignment: WrapCrossAlignment.center,
        children: [
          Wrap(
            spacing: BeamSizes.size16,
            crossAxisAlignment: WrapCrossAlignment.center,
            children: [
              // TODO(nausharipov) review: merge with GA4 events.
              FeedbackWidget(
                title: 'ui.feedbackTitle'.tr(),
                onRatingChanged: (rating) {
                  print([
                    'onRatingChanged',
                    rating,
                  ]);
                },
                onSubmitPressed: (rating, text) {
                  print([
                    'onSubmitPressed',
                    rating,
                    text,
                  ]);
                },
              ),
              const _ReportIssueButton(),
              const _PrivacyPolicyButton(),
              const Text('ui.copyright').tr(),
            ],
          ),
          const _BeamVersion(),
        ],
      ),
    );
  }
}

class _Body extends StatelessWidget {
  final Widget child;
  const _Body({required this.child});

  @override
  Widget build(BuildContext context) {
    final themeData = Theme.of(context);
    final ext = themeData.extension<BeamThemeExtension>()!;

    return Container(
      width: double.infinity,
      height: TobSizes.footerHeight,
      padding: const EdgeInsets.symmetric(
        vertical: BeamSizes.size4,
        horizontal: BeamSizes.size16,
      ),
      decoration: BoxDecoration(
        color: ext.secondaryBackgroundColor,
        border: Border(
          top: BorderSide(color: themeData.dividerColor),
        ),
      ),
      child: child,
    );
  }
}

class _ReportIssueButton extends StatelessWidget {
  const _ReportIssueButton();

  @override
  Widget build(BuildContext context) {
    return TextButton(
      style: _linkButtonStyle,
      onPressed: () {
        unawaited(launchUrl(Uri.parse(BeamLinks.reportIssue)));
      },
      child: const Text('ui.reportIssue').tr(),
    );
  }
}

class _PrivacyPolicyButton extends StatelessWidget {
  const _PrivacyPolicyButton();

  @override
  Widget build(BuildContext context) {
    return TextButton(
      style: _linkButtonStyle,
      onPressed: () {
        unawaited(launchUrl(Uri.parse(BeamLinks.privacyPolicy)));
      },
      child: const Text('ui.privacyPolicy').tr(),
    );
  }
}

class _BeamVersion extends StatelessWidget {
  const _BeamVersion();

  Future<String?> _getBeamSdkVersion() async {
    final sdk = GetIt.instance.get<AppNotifier>().sdk;
    if (sdk == null) {
      return null;
    }
    final runnerVersion = await GetIt.instance
        .get<BuildMetadataController>()
        .getRunnerVersion(sdk);
    return runnerVersion.beamSdkVersion;
  }

  @override
  Widget build(BuildContext context) {
    return AnimatedBuilder(
      animation: GetIt.instance.get<AppNotifier>(),
      builder: (context, child) => FutureBuilder<String?>(
        // ignore: discarded_futures
        future: _getBeamSdkVersion(),
        builder: (context, snapshot) => snapshot.hasData
            ? Text(
                'ui.builtWith'.tr(
                  namedArgs: {
                    'beamSdkVersion': snapshot.data!,
                  },
                ),
                style: const TextStyle(
                  color: BeamColors.grey3,
                ),
              )
            : Container(),
      ),
    );
  }
}

final _linkButtonStyle = TextButton.styleFrom(
  textStyle: const TextStyle(
    fontSize: 12,
    fontWeight: FontWeight.w400,
  ),
);
