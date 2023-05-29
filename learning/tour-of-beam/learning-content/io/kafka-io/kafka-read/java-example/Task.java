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

// beam-playground:
//   name: kafka-read
//   description: KafkaIO read example
//   multifile: false
//   default_example: false
//   context_line: 54
//   categories:
//     - IO
//   complexity: ADVANCED

import com.fasterxml.jackson.annotation.JsonSetter;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.core.type.TypeReference;
import com.fasterxml.jackson.databind.ObjectMapper;

import java.util.Collections;
import java.util.HashMap;
import java.util.Map;

import org.apache.beam.sdk.Pipeline;
import org.apache.beam.sdk.io.kafka.KafkaIO;
import org.apache.beam.sdk.io.kafka.KafkaRecord;
import org.apache.beam.sdk.options.PipelineOptions;
import org.apache.beam.sdk.options.PipelineOptionsFactory;
import org.apache.beam.sdk.transforms.Combine;
import org.apache.beam.sdk.transforms.DoFn;
import org.apache.beam.sdk.transforms.ParDo;
import org.apache.beam.sdk.transforms.Sum;
import org.apache.beam.sdk.transforms.Values;
import org.apache.beam.sdk.values.KV;
import org.apache.kafka.common.TopicPartition;
import org.apache.kafka.common.serialization.StringDeserializer;
import org.apache.kafka.common.serialization.ByteArraySerializer;
import org.apache.kafka.common.serialization.StringSerializer;

public class Task {
    public static void main(String[] args) {

        PipelineOptions options = PipelineOptionsFactory.fromArgs(args).create();
        Pipeline pipeline = Pipeline.create(options);
        Map<String, Object> consumerConfig = new HashMap<>();
        consumerConfig.put("auto.offset.reset", "earliest");

        /*
        * Read from Kafka topic: A KafkaIO Read transform is applied to the pipeline.
        * The withBootstrapServers method is used to specify the Kafka brokers to connect to.
        * The withTopicPartitions method is used to specify the Kafka topic and partition to read from.
        * The withKeyDeserializer and withValueDeserializer methods specify the deserializer for keys and values.
        * The withConsumerConfigUpdates method is used to pass the configuration Map to the Kafka consumer.
        * Extract values from KafkaRecord: The Values.create() method is used to extract the value part of the KafkaRecord.
        * */

        pipeline.apply(
                        "ReadFromKafka",
                        KafkaIO.<String, String>read()
                                .withBootstrapServers(
                                        "kafka_server:9092") // The argument is hardcoded to a predefined value. Do not
                                // change it manually. It's replaced to the correct Kafka
                                // cluster address when code starts in backend.

                                .withTopicPartitions(
                                        Collections.singletonList(
                                                new TopicPartition(
                                                        "NYCTaxi1000_simple",
                                                        0))) // The argument is hardcoded to a predefined value. Do not change
                                // it manually. It's replaced to the correct Kafka cluster address
                                // when code starts in backend.
                                .withKeyDeserializer(StringDeserializer.class)
                                .withValueDeserializer(StringDeserializer.class)
                                .withConsumerConfigUpdates(consumerConfig)
                                .withMaxNumRecords(998)
                                .withoutMetadata())
                .apply("CreateValues", Values.create())
                .apply(
                        "ExtractData",
                        ParDo.of(
                                new DoFn<String, String>() {
                                    @ProcessElement
                                    public void processElement(ProcessContext c) throws JsonProcessingException {
                                        System.out.println(c.element());
                                    }
                                }));


        pipeline.run().waitUntilFinish();
    }
}