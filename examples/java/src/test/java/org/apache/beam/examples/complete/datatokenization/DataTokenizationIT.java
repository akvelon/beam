package org.apache.beam.examples.complete.datatokenization;

import com.google.cloud.Timestamp;
import java.nio.charset.StandardCharsets;
import java.util.HashSet;
import java.util.Set;
import java.util.UUID;
import java.util.function.Function;
import org.apache.beam.examples.complete.datatokenization.options.DataTokenizationOptions;
import org.apache.beam.examples.complete.datatokenization.transforms.DataProtectors.RowToTokenizedRow;
import org.apache.beam.examples.complete.datatokenization.transforms.io.FileSystemIO;
import org.apache.beam.examples.complete.datatokenization.transforms.io.FileSystemIO.FORMAT;
import org.apache.beam.examples.complete.datatokenization.utils.FailsafeElement;
import org.apache.beam.examples.complete.datatokenization.utils.SchemasUtils;
import org.apache.beam.sdk.PipelineResult;
import org.apache.beam.sdk.coders.KvCoder;
import org.apache.beam.sdk.coders.RowCoder;
import org.apache.beam.sdk.coders.VarIntCoder;
import org.apache.beam.sdk.testing.TestPipeline;
import org.apache.beam.sdk.testutils.NamedTestResult;
import org.apache.beam.sdk.testutils.metrics.IOITMetrics;
import org.apache.beam.sdk.testutils.metrics.MetricsReader;
import org.apache.beam.sdk.testutils.metrics.TimeMonitor;
import org.apache.beam.sdk.transforms.DoFn;
import org.apache.beam.sdk.transforms.JsonToRow;
import org.apache.beam.sdk.transforms.MapElements;
import org.apache.beam.sdk.transforms.ParDo;
import org.apache.beam.sdk.values.KV;
import org.apache.beam.sdk.values.PCollection;
import org.apache.beam.sdk.values.PCollectionTuple;
import org.apache.beam.sdk.values.Row;
import org.apache.beam.sdk.values.TupleTag;
import org.apache.beam.sdk.values.TypeDescriptors;
import org.junit.Before;
import org.junit.Rule;
import org.junit.Test;

/**
 * Performance test for DataTokenization example.
 */
public class DataTokenizationIT {
  private static final String FILEIOIT_NAMESPACE = DataTokenizationIT.class.getName();
  /** The tag for the main output for the UDF. */
  private static final TupleTag<Row> TOKENIZATION_OUT = new TupleTag<Row>() {};

  /** The tag for the dead-letter output of the udf. */
  private static final TupleTag<FailsafeElement<Row, Row>> TOKENIZATION_DEADLETTER_OUT =
      new TupleTag<FailsafeElement<Row, Row>>() {};

  @Rule public TestPipeline pipeline = TestPipeline.create();
  public static DataTokenizationOptions options = TestPipeline.testingPipelineOptions().as(
      DataTokenizationOptions.class);

  private static final String bigQueryDataset = "load_test";
  private static final String bigQueryTable = "test_results_table";

  @Before
  public void setUp() throws Exception {
    options.setDataSchemaPath("gs://dt-test-data/generated/schema.json");
    options.setInputFilePattern("gs://dt-test-data/generated/small/small/test-data.json");
    options.setInputFileFormat(FORMAT.JSON);
    options.setRpcUri("https://dsg-mock-zh5a53lqvq-uc.a.run.app/tokenize_all");
  }

  static class LogIt<T> extends DoFn<T, T> {
    @ProcessElement
    public void processElement(ProcessContext context) {
      System.out.println("element is" + context.element());
      context.output(context.element());
    }
  }


  @Test
  public void testSimpleMetrics() throws Exception {
    SchemasUtils schema = new SchemasUtils(options.getDataSchemaPath(), StandardCharsets.UTF_8);
    PCollection<String> jsons = new FileSystemIO(options).read(pipeline, schema.getJsonBeamSchema())
//        .apply(ParDo.of(new LogIt<>()))
        ;
    JsonToRow.ParseResult rows =
        jsons.apply(
            "JsonToRow",
            JsonToRow.withExceptionReporting(schema.getBeamSchema()).withExtendedErrorInfo());
    rows.getFailedToParseLines()
//        .apply(ParDo.of(new LogIt<>()))
        ;

    PCollectionTuple tokenizedRows =
        rows.getResults()
            .setRowSchema(schema.getBeamSchema())
//            .apply(ParDo.of(new LogIt<>()))
            .apply(
                MapElements.into(
                    TypeDescriptors.kvs(TypeDescriptors.integers(), TypeDescriptors.rows()))
                    .via((Row row) -> KV.of(0, row)))
            .setCoder(KvCoder.of(VarIntCoder.of(), RowCoder.of(schema.getBeamSchema())))
            .apply(
                "Collect write start time",
                ParDo.of(new TimeMonitor<>(FILEIOIT_NAMESPACE, "startTime")))
            .apply(
                "DsgTokenization",
                RowToTokenizedRow.newBuilder()
                    .setBatchSize(options.getBatchSize())
                    .setRpcURI(options.getRpcUri())
                    .setSchema(schema.getBeamSchema())
                    .setSuccessTag(TOKENIZATION_OUT)
                    .setFailureTag(TOKENIZATION_DEADLETTER_OUT)
                    .build())
        ;

    tokenizedRows.get(TOKENIZATION_DEADLETTER_OUT)
        .apply(ParDo.of(new LogIt<>()))
        ;
    PCollection<Row> actual = tokenizedRows.get(TOKENIZATION_OUT)
        .apply(
            "Collect write end time",
            ParDo.of(new TimeMonitor<>(FILEIOIT_NAMESPACE, "endTime")))
//        .apply(ParDo.of(new LogIt<>()))
        ;

//    PAssert.that(actual).empty();
    System.out.println(options.getRunner());
    PipelineResult pipelineResult = pipeline.run(options);
    pipelineResult.waitUntilFinish();
    collectAndPublishMetrics(pipelineResult);
//    PDone write = new FileSystemIO(options)
//        .write(tokenizedRows.get(TOKENIZATION_OUT), schema.getBeamSchema());
  }
  private void collectAndPublishMetrics(PipelineResult result) {
    String uuid = UUID.randomUUID().toString();
    Timestamp timestamp = Timestamp.now();

    Set<Function<MetricsReader, NamedTestResult>> metricSuppliers =
        fillMetricSuppliers(uuid, timestamp.toString());

    final IOITMetrics metrics =
        new IOITMetrics(metricSuppliers, result, FILEIOIT_NAMESPACE, uuid, timestamp.toString());
    metrics.publish(bigQueryDataset, bigQueryTable);
  }

  private Set<Function<MetricsReader, NamedTestResult>> fillMetricSuppliers(
      String uuid, String timestamp) {
    Set<Function<MetricsReader, NamedTestResult>> metricSuppliers = new HashSet<>();

    metricSuppliers.add(
        (reader) -> {
          long writeStartTime = reader.getStartTimeMetric("startTime");
          long readEndTime = reader.getEndTimeMetric("endTime");
          double runTime = (readEndTime - writeStartTime) / 1e3;
          return NamedTestResult.create(uuid, timestamp, "run_time", runTime);
        });
    return metricSuppliers;
  }
}
