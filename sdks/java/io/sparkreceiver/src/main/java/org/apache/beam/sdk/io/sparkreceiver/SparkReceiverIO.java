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
import static org.apache.beam.vendor.guava.v26_0_jre.com.google.common.base.Preconditions.checkArgument;

import com.google.auto.value.AutoValue;
import org.apache.beam.sdk.annotations.Experimental;
import org.apache.beam.sdk.transforms.Impulse;
import org.apache.beam.sdk.transforms.PTransform;
import org.apache.beam.sdk.transforms.ParDo;
import org.apache.beam.sdk.transforms.SerializableFunction;
import org.apache.beam.sdk.values.PBegin;
import org.apache.beam.sdk.values.PCollection;
import org.apache.spark.streaming.receiver.Receiver;
import org.checkerframework.checker.nullness.qual.Nullable;
import org.joda.time.Instant;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

/** Streaming sources for Spark {@link Receiver}. */
@SuppressWarnings("nullness")
public class SparkReceiverIO {

  private static final Logger LOG = LoggerFactory.getLogger(SparkReceiverIO.class);

  public static <V> Read<V> read() {
    return new AutoValue_SparkReceiverIO_Read.Builder<V>().build();
  }

  /** A {@link PTransform} to read from Spark {@link Receiver}. */
  @AutoValue
  @AutoValue.CopyAnnotations
  public abstract static class Read<V> extends PTransform<PBegin, PCollection<V>> {

    abstract @Nullable ReceiverBuilder<V, ? extends Receiver<V>> getSparkReceiverBuilder();

    abstract @Nullable Class<V> getValueClass();

    abstract @Nullable SerializableFunction<V, Long> getGetOffsetFn();

    abstract @Nullable SparkConsumer<V> getSparkConsumer();

    abstract @Nullable SerializableFunction<V, Instant> getWatermarkFn();

    abstract Builder<V> toBuilder();

    @Experimental(Experimental.Kind.PORTABILITY)
    @AutoValue.Builder
    abstract static class Builder<V> {

      abstract Builder<V> setValueClass(Class<V> valueClass);

      abstract Builder<V> setSparkReceiverBuilder(
          ReceiverBuilder<V, ? extends Receiver<V>> sparkReceiverBuilder);

      abstract Builder<V> setGetOffsetFn(SerializableFunction<V, Long> getOffsetFn);

      abstract Builder<V> setSparkConsumer(SparkConsumer<V> sparkConsumer);

      abstract Builder<V> setWatermarkFn(SerializableFunction<V, Instant> watermarkFn);

      abstract Read<V> build();
    }

    public Read<V> withValueClass(Class<V> valueClass) {
      checkArgument(valueClass != null, "Value class can not be null");
      return toBuilder().setValueClass(valueClass).build();
    }

    public Read<V> withSparkReceiverBuilder(
        ReceiverBuilder<V, ? extends Receiver<V>> sparkReceiverBuilder) {
      checkArgument(sparkReceiverBuilder != null, "Spark receiver builder can not be null");
      return toBuilder().setSparkReceiverBuilder(sparkReceiverBuilder).build();
    }

    public Read<V> withGetOffsetFn(SerializableFunction<V, Long> getOffsetFn) {
      checkArgument(getOffsetFn != null, "Get offset function can not be null");
      return toBuilder().setGetOffsetFn(getOffsetFn).build();
    }

    public Read<V> withSparkConsumer(SparkConsumer<V> sparkConsumer) {
      checkArgument(sparkConsumer != null, "Spark consumer can not be null");
      return toBuilder().setSparkConsumer(sparkConsumer).build();
    }

    public Read<V> withWatermarkFn(SerializableFunction<V, Instant> watermarkFn) {
      checkArgument(watermarkFn != null, "Watermark function can not be null");
      return toBuilder().setWatermarkFn(watermarkFn).build();
    }

    @Override
    public PCollection<V> expand(PBegin input) {
      validateTransform();
      return input.apply(new ReadFromSparkReceiverViaSdf<>(this));
    }

    public void validateTransform() {
      ReceiverBuilder<V, ? extends Receiver<V>> sparkReceiverBuilder = getSparkReceiverBuilder();
      checkStateNotNull(sparkReceiverBuilder, "withSparkReceiverBuilder() is required");
      if (!HasOffset.class.isAssignableFrom(sparkReceiverBuilder.getSparkReceiverClass())) {
        checkStateNotNull(getSparkConsumer(), "withSparkConsumer() is required");
      }
      checkStateNotNull(getValueClass(), "withValueClass() is required");
      checkStateNotNull(getGetOffsetFn(), "withGetOffsetFn() is required");
    }
  }

  static class ReadFromSparkReceiverViaSdf<V> extends PTransform<PBegin, PCollection<V>> {

    private final Read<V> sparkReceiverRead;

    ReadFromSparkReceiverViaSdf(Read<V> sparkReceiverRead) {
      this.sparkReceiverRead = sparkReceiverRead;
    }

    @Override
    public PCollection<V> expand(PBegin input) {
      final ReceiverBuilder<V, ? extends Receiver<V>> sparkReceiverBuilder =
          sparkReceiverRead.getSparkReceiverBuilder();
      checkStateNotNull(sparkReceiverBuilder, "withSparkReceiverBuilder() is required");
      if (!HasOffset.class.isAssignableFrom(sparkReceiverBuilder.getSparkReceiverClass())) {
        LOG.info(
            "{} started reading", ReadFromSparkReceiverWithoutOffsetDoFn.class.getSimpleName());
        return input
            .apply(Impulse.create())
            .apply(ParDo.of(new ReadFromSparkReceiverWithoutOffsetDoFn<>(sparkReceiverRead)));
      } else {
        LOG.info("{} started reading", ReadFromSparkReceiverWithOffsetDoFn.class.getSimpleName());
        return input
            .apply(Impulse.create())
            .apply(ParDo.of(new ReadFromSparkReceiverWithOffsetDoFn<>(sparkReceiverRead)));
      }
    }
  }
}
