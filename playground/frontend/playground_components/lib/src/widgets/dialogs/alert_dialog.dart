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

import 'package:flutter/material.dart';

import '../../../playground_components.dart';

class BeamAlertDialog extends StatelessWidget {
  final String? body;
  final String continueLabel;
  final VoidCallback onContinue;
  final String title;

  const BeamAlertDialog({
    this.body,
    required this.continueLabel,
    required this.onContinue,
    required this.title,
  });

  @override
  Widget build(BuildContext context) {
    return OverlayBody(
      child: Container(
        width: BeamSizes.popupWidth,
        padding: const EdgeInsets.all(BeamSizes.size16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          mainAxisSize: MainAxisSize.min,
          children: [
            Text(
              title,
              style: Theme.of(context).textTheme.headlineMedium,
            ),
            if (body != null)
              Padding(
                padding: const EdgeInsets.only(top: BeamSizes.size8),
                child: Text(body!),
              ),
            const SizedBox(height: BeamSizes.size8),
            Row(
              mainAxisAlignment: MainAxisAlignment.end,
              children: [
                TextButton(
                  onPressed: () {
                    Navigator.pop(context);
                  },
                  // TODO(nausharipov): review: translate in PGC?
                  child: const Text('Cancel'),
                ),
                const SizedBox(width: BeamSizes.size8),
                TextButton(
                  onPressed: onContinue,
                  child: Text(continueLabel),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }
}
