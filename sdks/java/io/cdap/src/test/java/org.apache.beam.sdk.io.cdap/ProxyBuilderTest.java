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
package org.apache.beam.sdk.io.cdap;

import static org.junit.Assert.assertTrue;

import io.cdap.plugin.salesforce.plugin.source.streaming.SalesforceReceiver;
import io.cdap.plugin.salesforce.plugin.source.streaming.SalesforceStreamingSourceConfig;
import java.util.concurrent.atomic.AtomicBoolean;
import org.apache.spark.streaming.receiver.ReceiverSupervisor;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.junit.runners.JUnit4;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

/** Test class for {@link ProxyBuilder} */
@RunWith(JUnit4.class)
public class ProxyBuilderTest {

  private static final Logger LOG = LoggerFactory.getLogger(ProxyBuilderTest.class);

  private static final String SALESFORCE_CONFIG_JSON_STRING =
      "{\n"
          + "\"pushTopicName\": \"topicName\",\n"
          + "\"sObjectName\": \"sObject\",\n"
          + "\"datetimeAfter\": \"datetime\",\n"
          + "\"consumerKey\": \"key\",\n"
          + "\"consumerSecret\": \"secret\",\n"
          + "\"username\": \"user\",\n"
          + "\"password\": \"password\",\n"
          + "\"loginUrl\": \"https://www.google.com\",\n"
          + "\"referenceName\": \"reference\"\n"
          + "}";
  public static final String TEST_MESSAGE = "testMessage";

  /**
   * If this test passed, then proxy object for custom {@link
   * org.apache.spark.streaming.receiver.Receiver} was created successfully, and the corresponding
   * {@link ReceiverSupervisor} was wrapped into {@link WrappedSupervisor}.
   */
  @Test
  public void testCreatingProxySparkReceiverForSalesforce() {
    try {
      SalesforceStreamingSourceConfig config =
          new ConfigWrapper<>(SalesforceStreamingSourceConfig.class)
              .fromJsonString(SALESFORCE_CONFIG_JSON_STRING)
              .build();
      ProxyBuilder<String, SalesforceReceiver> builder =
          new ProxyBuilder<>(SalesforceReceiver.class);

      AtomicBoolean customStoreConsumerWasUsed = new AtomicBoolean(false);

      SalesforceReceiver proxyReciever =
          builder
              .withConstructorArgs(config.getAuthenticatorCredentials(), config.getPushTopicName())
              .withCustomStoreConsumer(args -> customStoreConsumerWasUsed.set(true))
              .build();
      proxyReciever.onStart();

      ReceiverSupervisor supervisor = proxyReciever.supervisor();
      assertTrue(supervisor instanceof WrappedSupervisor);

      proxyReciever.store(TEST_MESSAGE);
      assertTrue(customStoreConsumerWasUsed.get());
    } catch (Exception e) {
      LOG.error("Can not get proxy", e);
    }
  }
}
