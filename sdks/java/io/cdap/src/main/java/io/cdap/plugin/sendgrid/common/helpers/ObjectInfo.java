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
package io.cdap.plugin.sendgrid.common.helpers;

import io.cdap.plugin.sendgrid.common.APIResponseType;
import io.cdap.plugin.sendgrid.common.objects.DataSourceGroupType;
import java.util.List;
import java.util.stream.Collectors;

/** Entity meta-info holder for {@link ObjectDefinition}. */
@SuppressWarnings("rawtypes")
public class ObjectInfo {
  private String cdapObjectName;

  private String sendGridAPIUrl;
  private DataSourceGroupType dataSourceGroupType;
  private Class objectClass;
  private List<ObjectFieldInfo> fieldDefinitions;
  private APIResponseType apiResponseType;
  private List<String> requiredArguments;
  private ObjectDefinition.ObjectDefinitionType objectType;

  /**
   * Constructor for ObjectInfo object.
   *
   * @param cdapObjectName the cdap object name
   * @param fieldDefinitions the field definitions
   * @param sendGridAPIUrl the sendgrid api url
   * @param apiResponseType the api response type
   * @param objectClass the object class
   * @param dataSourceGroupType the data source group type
   * @param requiredArguments the required arguments
   * @param objectType the object type
   */
  public ObjectInfo(
      String cdapObjectName,
      List<ObjectFieldInfo> fieldDefinitions,
      String sendGridAPIUrl,
      APIResponseType apiResponseType,
      Class objectClass,
      DataSourceGroupType dataSourceGroupType,
      List<String> requiredArguments,
      ObjectDefinition.ObjectDefinitionType objectType) {
    this.cdapObjectName = cdapObjectName;
    this.fieldDefinitions = fieldDefinitions;
    this.sendGridAPIUrl = sendGridAPIUrl;
    this.apiResponseType = apiResponseType;
    this.objectClass = objectClass;
    this.dataSourceGroupType = dataSourceGroupType;
    this.requiredArguments = requiredArguments;
    this.objectType = objectType;
  }

  public String getCdapObjectName() {
    return cdapObjectName;
  }

  public String getSendGridAPIUrl() {
    return sendGridAPIUrl;
  }

  public List<ObjectFieldInfo> getFieldDefinitions() {
    return fieldDefinitions;
  }

  /**
   * Returns the list of ObjectFieldInfo.
   *
   * @param fields the list of string
   * @return list of ObjectFieldInfo
   */
  public List<ObjectFieldInfo> getFieldsDefinitions(List<String> fields) {
    return fieldDefinitions.stream()
        .filter(x -> fields.stream().anyMatch(y -> x.getName().equals(y)))
        .collect(Collectors.toList());
  }

  public Class getObjectClass() {
    return objectClass;
  }

  public DataSourceGroupType getDataSourceGroupType() {
    return dataSourceGroupType;
  }

  public APIResponseType getApiResponseType() {
    return apiResponseType;
  }

  public List<String> getRequiredArguments() {
    return requiredArguments.stream().filter(x -> !x.equals("")).collect(Collectors.toList());
  }

  public ObjectDefinition.ObjectDefinitionType getObjectType() {
    return objectType;
  }
}
