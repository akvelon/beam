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

class SearchField extends StatelessWidget {
  final TextEditingController controller;

  const SearchField({Key? key, required this.controller}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Container(
      margin: const EdgeInsets.only(
        top: 12.0,
        right: 12.0,
        left: 12.0,
      ),
      width: 376.0,
      height: 40.0,
      color: Colors.white,
      child: ClipRRect(
        borderRadius: BorderRadius.circular(6.0),
        child: TextFormField(
          controller: controller,
          decoration: InputDecoration(
            suffixIcon: const Padding(
              padding: EdgeInsetsDirectional.only(
                start: 0.0,
                end: 0.0,
              ),
              child: Icon(
                Icons.search,
                color: Colors.black38,
                size: 25.0,
              ),
            ),
            focusedBorder: OutlineInputBorder(
              borderSide: const BorderSide(
                color: Colors.black12,
              ),
              borderRadius: BorderRadius.circular(6.0),
            ),
            enabledBorder: OutlineInputBorder(
              borderSide: const BorderSide(
                color: Colors.black12,
              ),
              borderRadius: BorderRadius.circular(6.0),
            ),
            filled: false,
            isDense: true,
            hintText: 'Search',
            contentPadding: const EdgeInsets.only(left: 13.0),
          ),
          style: const TextStyle(
            fontSize: 14.0,
            height: 14 / 8,
            color: Colors.black,
          ),
          cursorColor: Colors.grey,
          cursorWidth: 1.0,
          textAlignVertical: TextAlignVertical.center,
          onFieldSubmitted: (String txt) {},
          onChanged: (String txt) {},
          maxLines: 1,
          minLines: 1,
        ),
      ),
    );
  }
}
