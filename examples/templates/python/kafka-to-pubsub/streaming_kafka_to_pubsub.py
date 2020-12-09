# coding=UTF8
#
# Licensed to the Apache Software Foundation (ASF) under one
# or more contributor license agreements.  See the NOTICE file
# distributed with this work for additional information
# regarding copyright ownership.  The ASF licenses this file
# to you under the Apache License, Version 2.0 (the
# "License"); you may not use this file except in compliance
# with the License.  You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

"""An Apache Beam streaming pipeline.

It reads messages from Kafka writes the results to Pub/Sub.
"""

import argparse
import logging
import apache_beam as beam
from apache_beam.options.pipeline_options import PipelineOptions
from apache_beam.io.external.kafka import ReadFromKafka, WriteToKafka


def run(args, bootstrap_servers, input_topics, output_topic):
    """Build and run the pipeline."""
    options = PipelineOptions(args, save_main_session=True, streaming=True)
    topics_list = (input_topics.split(","))

    # Log pipeline parameters
    logging.info("%sbootstrap servers: %s\ninput topics: %s\noutput pubsub topic: %s",
                 "Starting Kafka-To-PubSub pipeline with parameters ",
                 bootstrap_servers,
                 input_topics,
                 output_topic)

    """Steps:
    1) Read messages from Kafka
    2) Write to PubSub
    """

    # Create the pipeline
    with beam.Pipeline(options=options) as p:
        (p
         | "Read from Kafka" >> ReadFromKafka(
                    consumer_config={"bootstrap.servers": bootstrap_servers,
                                     "auto.offset.reset": "latest"},
                    topics=topics_list)
         | "Get values" >> beam.Map(lambda kv: kv[1])
         | "Write to PubSub" >> beam.io.WriteToPubSub(output_topic)
         )


if __name__ == "__main__":
    logging.getLogger().setLevel(logging.INFO)
    parser = argparse.ArgumentParser()
    parser.add_argument(
        "--bootstrap_servers",
        help="Kafka Bootstrap Servers: "
             "\"broker_1:9092,broker_2:9092\"")
    parser.add_argument(
        "--input_topics",
        help="Kafka topic to read the input from "
             "\"topic1,topic2\"")
    parser.add_argument(
        "--output_topic",
        help="The Cloud Pub/Sub topic to publish to. "
             "The name should be in the format of "
             "\"projects/<project-id>/topics/<topic-name>\".")
    known_args, pipeline_args = parser.parse_known_args()
    run(pipeline_args,
        known_args.bootstrap_servers,
        known_args.input_topics,
        known_args.output_topic)
