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
package org.apache.beam.sdk.extensions.python;

import org.apache.beam.sdk.options.Default;
import org.apache.beam.sdk.options.Description;
import org.apache.beam.sdk.options.PipelineOptions;

/** Pipeline options for {@link PythonExternalTransform}. */
public interface PythonExternalTransformOptions extends PipelineOptions {

  @Description("Use Docker Compose based Beam Transform Service to expand transforms.")
  @Default.Boolean(false)
  boolean getUseTransformService();

  void setUseTransformService(boolean useTransformService);

  @Description("Custom Beam version for bootstrap Beam venv.")
  String getCustomBeamRequirement();

  /**
   * Set custom Beam version for bootstrap Beam venv.
   *
   * <p>For example: 2.50.0rc1, "/path/to/apache-beam.whl"
   */
  void setCustomBeamRequirement(String customBeamRequirement);
}
