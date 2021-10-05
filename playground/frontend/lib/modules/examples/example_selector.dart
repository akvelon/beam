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
import 'package:playground/config/theme.dart';
import 'package:playground/modules/examples/components/examples_components.dart';
import 'package:playground/modules/examples/models/selector_size_model.dart';

class ExampleSelector extends StatefulWidget {
  final Function changeSelectorVisibility;
  final bool isSelectorOpened;
  final List examples;

  const ExampleSelector({
    Key? key,
    required this.changeSelectorVisibility,
    required this.isSelectorOpened,
    required this.examples,
  }) : super(key: key);

  @override
  State<ExampleSelector> createState() => _ExampleSelectorState();
}

class _ExampleSelectorState extends State<ExampleSelector>
    with TickerProviderStateMixin {
  final GlobalKey selectorKey = LabeledGlobalKey('ExampleSelector');
  late OverlayEntry? _examplesDropdown;
  late AnimationController _animationController;
  late Animation<Offset> _offsetAnimation;

  @override
  void initState() {
    super.initState();
    _animationController = AnimationController(
      vsync: this,
      duration: const Duration(milliseconds: 80),
    );
    _offsetAnimation = Tween<Offset>(
      begin: const Offset(0.0, -0.02),
      end: const Offset(0.0, 0.0),
    ).animate(_animationController);
  }

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      key: selectorKey,
      onTap: () {
        if (widget.isSelectorOpened) {
          _animationController.reverse();
          _examplesDropdown?.remove();
        } else {
          _animationController.forward();
          _examplesDropdown = _createExamplesDropdown();
          Overlay.of(context)?.insert(_examplesDropdown!);
        }
        widget.changeSelectorVisibility();
      },
      child: Container(
        width: 105.0,
        height: 40.0,
        decoration: BoxDecoration(
          color: ThemeColors.greyColor(context),
          borderRadius: BorderRadius.circular(4.0),
        ),
        child: Row(
          mainAxisAlignment: MainAxisAlignment.center,
          children: const [
            Text(
              'Catalog',
              style: TextStyle(color: Colors.black, fontSize: 16.0),
            ),
            Icon(Icons.keyboard_arrow_down, color: Colors.black54),
          ],
        ),
      ),
    );
  }

  OverlayEntry _createExamplesDropdown() {
    SelectorPositionModel _posModel = _findSelectorPositionData();
    final TextEditingController _textController = TextEditingController();

    return OverlayEntry(
      builder: (context) {
        return Stack(children: [
          GestureDetector(
            onTap: () {
              _animationController.reverse();
              _examplesDropdown?.remove();
              widget.changeSelectorVisibility();
            },
            child: Container(
              color: Colors.transparent,
              height: double.infinity,
              width: double.infinity,
            ),
          ),
          Positioned(
            left: _posModel.xAlignment,
            top: _posModel.yAlignment + 50.0,
            child: SlideTransition(
              position: _offsetAnimation,
              child: Material(
                elevation: 2.0,
                child: Container(
                  height: 444.0,
                  width: 400.0,
                  decoration: BoxDecoration(
                    color: Colors.white,
                    borderRadius: BorderRadius.circular(6.0),
                  ),
                  child: Column(
                    children: [
                      SearchField(controller: _textController),
                      const TypeFilter(),
                      ExampleList(items: widget.examples),
                    ],
                  ),
                ),
              ),
            ),
          ),
        ]);
      },
    );
  }

  SelectorPositionModel _findSelectorPositionData() {
    RenderBox? rBox =
        selectorKey.currentContext?.findRenderObject() as RenderBox;
    SelectorPositionModel positionModel = SelectorPositionModel(
      xAlignment: rBox.localToGlobal(Offset.zero).dx,
      yAlignment: rBox.localToGlobal(Offset.zero).dy,
    );
    return positionModel;
  }
}
