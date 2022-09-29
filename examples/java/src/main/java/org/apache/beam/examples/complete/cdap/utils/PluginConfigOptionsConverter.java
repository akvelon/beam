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
package org.apache.beam.examples.complete.cdap.utils;

import io.cdap.plugin.common.Constants;
import io.cdap.plugin.hubspot.common.BaseHubspotConfig;
import io.cdap.plugin.servicenow.source.util.ServiceNowConstants;
import java.util.Map;
import org.apache.beam.examples.complete.cdap.options.CdapHubspotOptions;
import org.apache.beam.examples.complete.cdap.options.CdapServiceNowOptions;
import org.apache.beam.vendor.guava.v26_0_jre.com.google.common.collect.ImmutableMap;

/**
 * Class for converting CDAP {@link org.apache.beam.sdk.options.PipelineOptions} to map for {@link
 * org.apache.beam.sdk.io.cdap.ConfigWrapper}.
 */
public class PluginConfigOptionsConverter {

  public static Map<String, Object> hubspotOptionsToParamsMap(CdapHubspotOptions options) {
    String apiServerUrl = options.getApiServerUrl();
    return ImmutableMap.<String, Object>builder()
        .put(
            BaseHubspotConfig.API_SERVER_URL,
            apiServerUrl != null ? apiServerUrl : BaseHubspotConfig.DEFAULT_API_SERVER_URL)
        .put(BaseHubspotConfig.API_KEY, options.getApiKey())
        .put(BaseHubspotConfig.OBJECT_TYPE, options.getObjectType())
        .put(Constants.Reference.REFERENCE_NAME, options.getReferenceName())
        .build();
  }

  public static Map<String, Object> zendeskOptionsToParamsMap(CdapServiceNowOptions options) {
    return ImmutableMap.<String, Object>builder()
        .put(ServiceNowConstants.PROPERTY_CLIENT_ID, options.getClientId())
        .put(ServiceNowConstants.PROPERTY_CLIENT_SECRET, options.getClientSecret())
        .put(ServiceNowConstants.PROPERTY_USER, options.getUser())
        .put(ServiceNowConstants.PROPERTY_PASSWORD, options.getPassword())
        .put(ServiceNowConstants.PROPERTY_API_ENDPOINT, options.getRestApiEndpoint())
        .put(ServiceNowConstants.PROPERTY_QUERY_MODE, options.getQueryMode())
        .put(ServiceNowConstants.PROPERTY_TABLE_NAME, options.getTableName())
        .put(ServiceNowConstants.PROPERTY_VALUE_TYPE, options.getValueType())
        .put(Constants.Reference.REFERENCE_NAME, options.getReferenceName())
        .build();
  }
}
