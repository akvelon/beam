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
package org.apache.beam.sdk.io.sparkreceiver;

import static org.apache.beam.sdk.util.Preconditions.checkStateNotNull;

import java.util.Queue;
import java.util.concurrent.ConcurrentLinkedQueue;
import java.util.concurrent.TimeUnit;
import org.apache.beam.sdk.coders.Coder;
import org.apache.beam.sdk.io.range.OffsetRange;
import org.apache.beam.sdk.transforms.DoFn;
import org.apache.beam.sdk.transforms.DoFn.UnboundedPerElement;
import org.apache.beam.sdk.transforms.SerializableFunction;
import org.apache.beam.sdk.transforms.splittabledofn.ManualWatermarkEstimator;
import org.apache.beam.sdk.transforms.splittabledofn.OffsetRangeTracker;
import org.apache.beam.sdk.transforms.splittabledofn.RestrictionTracker;
import org.apache.beam.sdk.transforms.splittabledofn.WatermarkEstimator;
import org.apache.beam.sdk.transforms.splittabledofn.WatermarkEstimators;
import org.apache.beam.sdk.transforms.windowing.BoundedWindow;
import org.apache.spark.SparkConf;
import org.apache.spark.streaming.receiver.Receiver;
import org.checkerframework.checker.nullness.qual.Nullable;
import org.joda.time.Instant;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

/** A SplittableDoFn which reads from {@link Receiver} that implements {@link HasOffset}. */
@UnboundedPerElement
public class ReadFromSparkReceiverWithOffsetDoFn<V> extends DoFn<byte[], V> {

  private static final Logger LOG =
      LoggerFactory.getLogger(ReadFromSparkReceiverWithOffsetDoFn.class);
  private static final int START_POLL_TIMEOUT_MS = 1000;

  private final SerializableFunction<Instant, WatermarkEstimator<Instant>>
      createWatermarkEstimatorFn;
  private final SerializableFunction<V, Long> getOffsetFn;
  private final SerializableFunction<V, Instant> getWatermarkFn;
  private final ReceiverBuilder<V, ? extends Receiver<V>> sparkReceiverBuilder;

  public ReadFromSparkReceiverWithOffsetDoFn(SparkReceiverIO.Read<V> transform) {
    createWatermarkEstimatorFn = WatermarkEstimators.Manual::new;

    ReceiverBuilder<V, ? extends Receiver<V>> sparkReceiverBuilder =
        transform.getSparkReceiverBuilder();
    checkStateNotNull(sparkReceiverBuilder, "Spark Receiver Builder can't be null!");
    this.sparkReceiverBuilder = sparkReceiverBuilder;

    SerializableFunction<V, Instant> getWatermarkFn = transform.getWatermarkFn();
    checkStateNotNull(getWatermarkFn, "Watermark fn fn can't be null!");
    this.getWatermarkFn = getWatermarkFn;

    SerializableFunction<V, Long> getOffsetFn = transform.getGetOffsetFn();
    checkStateNotNull(getOffsetFn, "Get offset fn can't be null!");
    this.getOffsetFn = getOffsetFn;
  }

  @GetInitialRestriction
  public OffsetRange initialRestriction(@Element byte[] element) {
    return new OffsetRange(0, Long.MAX_VALUE);
  }

  @GetInitialWatermarkEstimatorState
  public Instant getInitialWatermarkEstimatorState(@Timestamp Instant currentElementTimestamp) {
    return currentElementTimestamp;
  }

  @NewWatermarkEstimator
  public WatermarkEstimator<Instant> newWatermarkEstimator(
      @WatermarkEstimatorState Instant watermarkEstimatorState) {
    return createWatermarkEstimatorFn.apply(ensureTimestampWithinBounds(watermarkEstimatorState));
  }

  @GetSize
  public double getSize(@Element byte[] element, @Restriction OffsetRange offsetRange) {
    return restrictionTracker(element, offsetRange).getProgress().getWorkRemaining();
    // Before processing elements, we don't have a good estimated size of records and offset gap.
  }

  @NewTracker
  public OffsetRangeTracker restrictionTracker(
      @Element byte[] element, @Restriction OffsetRange restriction) {
    return new OffsetRangeTracker(restriction);
  }

  @GetRestrictionCoder
  public Coder<OffsetRange> restrictionCoder() {
    return new OffsetRange.Coder();
  }

  @Setup
  public void setup() throws Exception {}

  @Teardown
  public void teardown() throws Exception {}

  @SuppressWarnings("unchecked")
  private static class SparkConsumerWithOffset<V> implements SparkConsumer<V> {
    private final Queue<V> recordsQueue;
    private @Nullable Receiver<V> sparkReceiver;
    private final Long startOffset;

    public SparkConsumerWithOffset(Long startOffset) {
      this.startOffset = startOffset;
      this.recordsQueue = new ConcurrentLinkedQueue<>();
    }

    @Override
    public boolean hasRecords() {
      return !recordsQueue.isEmpty();
    }

    @Override
    public @Nullable V poll() {
      return recordsQueue.poll();
    }

    @Override
    public void start(Receiver<V> sparkReceiver) {
      this.sparkReceiver = sparkReceiver;
      try {
        new WrappedSupervisor(
            sparkReceiver,
            new SparkConf(),
            objects -> {
              V record = (V) objects[0];
              recordsQueue.offer(record);
              return null;
            });
      } catch (Exception e) {
        LOG.error("Can not init Spark Receiver!", e);
      }
      ((HasOffset) sparkReceiver).setStartOffset(startOffset);
      sparkReceiver.supervisor().startReceiver();
      try {
        TimeUnit.MILLISECONDS.sleep(START_POLL_TIMEOUT_MS);
      } catch (InterruptedException e) {
        LOG.error("Interrupted", e);
      }
    }

    @Override
    public void stop() {
      if (sparkReceiver != null) {
        sparkReceiver.stop("Stopped");
      }
      recordsQueue.clear();
    }
  }

  @ProcessElement
  public ProcessContinuation processElement(
      @Element byte[] element,
      RestrictionTracker<OffsetRange, Long> tracker,
      WatermarkEstimator<Instant> watermarkEstimator,
      OutputReceiver<V> receiver) {

    SparkConsumer<V> sparkConsumer;
    Receiver<V> sparkReceiver;
    try {
      sparkReceiver = sparkReceiverBuilder.build();
    } catch (Exception e) {
      LOG.error("Can not build Spark Receiver", e);
      throw new IllegalStateException("Spark Receiver was not built!");
    }
    sparkConsumer = new SparkConsumerWithOffset<>(tracker.currentRestriction().getFrom());
    sparkConsumer.start(sparkReceiver);

    while (sparkConsumer.hasRecords()) {
      V record = sparkConsumer.poll();
      if (record != null) {
        Long offset = getOffsetFn.apply(record);
        if (!tracker.tryClaim(offset)) {
          sparkConsumer.stop();
          LOG.debug("Stop for restriction: {}", tracker.currentRestriction().toString());
          return ProcessContinuation.stop();
        }
        Instant currentTimeStamp = getWatermarkFn.apply(record);
        ((ManualWatermarkEstimator<Instant>) watermarkEstimator).setWatermark(currentTimeStamp);
        receiver.outputWithTimestamp(record, currentTimeStamp);
      }
    }
    sparkConsumer.stop();
    LOG.debug("Resume for restriction: {}", tracker.currentRestriction().toString());
    return ProcessContinuation.resume();
  }

  private static Instant ensureTimestampWithinBounds(Instant timestamp) {
    if (timestamp.isBefore(BoundedWindow.TIMESTAMP_MIN_VALUE)) {
      timestamp = BoundedWindow.TIMESTAMP_MIN_VALUE;
      LOG.debug("Timestamp was before MIN_VALUE({})", BoundedWindow.TIMESTAMP_MIN_VALUE);
    } else if (timestamp.isAfter(BoundedWindow.TIMESTAMP_MAX_VALUE)) {
      timestamp = BoundedWindow.TIMESTAMP_MAX_VALUE;
      LOG.debug("Timestamp was after MAX_VALUE({})", BoundedWindow.TIMESTAMP_MAX_VALUE);
    }
    return timestamp;
  }
}
