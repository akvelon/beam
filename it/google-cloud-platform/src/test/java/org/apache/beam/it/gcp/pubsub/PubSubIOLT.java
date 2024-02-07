package org.apache.beam.it.gcp.pubsub;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.google.pubsub.v1.SubscriptionName;
import com.google.pubsub.v1.TopicName;
import org.apache.beam.it.common.PipelineLauncher;
import org.apache.beam.it.common.PipelineOperator;
import org.apache.beam.it.common.TestProperties;
import org.apache.beam.it.common.utils.ResourceManagerUtils;
import org.apache.beam.it.gcp.IOLoadTestBase;
import org.apache.beam.runners.direct.DirectOptions;
import org.apache.beam.sdk.io.Read;
import org.apache.beam.sdk.io.gcp.pubsub.PubsubIO;
import org.apache.beam.sdk.io.gcp.pubsub.PubsubOptions;
import org.apache.beam.sdk.io.synthetic.SyntheticSourceOptions;
import org.apache.beam.sdk.io.synthetic.SyntheticUnboundedSource;
import org.apache.beam.sdk.options.PipelineOptionsFactory;
import org.apache.beam.sdk.testing.TestPipeline;
import org.apache.beam.sdk.testing.TestPipelineOptions;
import org.apache.beam.sdk.transforms.DoFn;
import org.apache.beam.sdk.transforms.ParDo;
import org.apache.beam.sdk.values.KV;
import org.apache.beam.vendor.guava.v32_1_2_jre.com.google.common.base.Strings;
import org.apache.beam.vendor.guava.v32_1_2_jre.com.google.common.collect.ImmutableMap;
import org.junit.*;

import java.io.IOException;
import java.nio.charset.StandardCharsets;
import java.time.Duration;
import java.util.Map;
import java.util.Objects;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertNotEquals;

public class PubSubIOLT extends IOLoadTestBase {

    private Configuration configuration;
    private static PubsubResourceManager resourceManager;
    private static TopicName topicName;
    private static SubscriptionName subscription;
    private static final String READ_ELEMENT_METRIC_NAME = "read_count";
    private static int numOfSourceBundles;

    @Rule public transient TestPipeline writePipeline = TestPipeline.create();
    @Rule public transient TestPipeline readPipeline = TestPipeline.create();

    @BeforeClass
    public static void beforeClass() throws IOException {
        resourceManager =
                PubsubResourceManager.builder("io-pubsub-lt", project, CREDENTIALS_PROVIDER).build();
        topicName = resourceManager.createTopic("topic");
        subscription = resourceManager.createSubscription(topicName, "subscription");
        PipelineOptionsFactory.register(TestPipelineOptions.class);
    }

    @Before
    public void setup() {
        // parse configuration
        String testConfig =
                TestProperties.getProperty("configuration", "local", TestProperties.Type.PROPERTY);
        configuration = TEST_CONFIGS_PRESET.get(testConfig);
        if (configuration == null) {
            try {
                configuration = PubSubIOLT.Configuration.fromJsonString(testConfig, PubSubIOLT.Configuration.class);
            } catch (IOException e) {
                throw new IllegalArgumentException(
                        String.format(
                                "Unknown test configuration: [%s]. Pass to a valid configuration json, or use"
                                        + " config presets: %s",
                                testConfig, TEST_CONFIGS_PRESET.keySet()));
            }
        }

        // Explicitly set up number of bundles in SyntheticUnboundedSource since it has a bug in implementation where
        // number of lost data in streaming pipeline equals to number of initial bundles.
        numOfSourceBundles = testConfig.equals("local") ? 10 : 20;
        configuration.forceNumInitialBundles = numOfSourceBundles;

        // tempLocation needs to be set for DataflowRunner
        if (!Strings.isNullOrEmpty(tempBucketName)) {
            String tempLocation = String.format("gs://%s/temp/", tempBucketName);
            writePipeline.getOptions().as(TestPipelineOptions.class).setTempRoot(tempLocation);
            writePipeline.getOptions().setTempLocation(tempLocation);
            readPipeline.getOptions().as(TestPipelineOptions.class).setTempRoot(tempLocation);
            readPipeline.getOptions().setTempLocation(tempLocation);
        }
        writePipeline.getOptions().as(PubsubOptions.class).setProject(topicName.getProject());
        readPipeline.getOptions().as(PubsubOptions.class).setProject(topicName.getProject());
        writePipeline.getOptions().as(DirectOptions.class).setBlockOnRun(false);
        readPipeline.getOptions().as(DirectOptions.class).setBlockOnRun(false);
    }

    @AfterClass
    public static void tearDownClass() {
        ResourceManagerUtils.cleanResources(resourceManager);
    }

    private static final Map<String, Configuration> TEST_CONFIGS_PRESET;

    static {
        try {
            TEST_CONFIGS_PRESET =
                    ImmutableMap.of(
                            "local",
                            PubSubIOLT.Configuration.fromJsonString(
                                    "{\"numRecords\":500,\"valueSizeBytes\":1000,\"pipelineTimeout\":5,\"runner\":\"DirectRunner\"}",
                                    PubSubIOLT.Configuration.class), // 0.5 MB
                            "medium",
                            PubSubIOLT.Configuration.fromJsonString(
                                    "{\"numRecords\":100000,\"valueSizeBytes\":1000,\"pipelineTimeout\":15,\"runner\":\"DataflowRunner\"}",
                                    PubSubIOLT.Configuration.class), // 10 GB
                            "large",
                            PubSubIOLT.Configuration.fromJsonString(
                                    "{\"numRecords\":100000000,\"valueSizeBytes\":1000,\"pipelineTimeout\":80,\"runner\":\"DataflowRunner\"}",
                                    PubSubIOLT.Configuration.class) // 100 GB
                    );
        } catch (IOException e) {
            throw new RuntimeException(e);
        }
    }

    @Test
    public void testWriteAndRead() throws IOException {

        writePipeline
                .apply("Read from source", Read.from(new SyntheticUnboundedSource(configuration)))
                .apply("Map records", ParDo.of(new MapKVtoString()))
                .apply("Write to PubSub",
                        PubsubIO
                                .writeStrings()
                                .to(topicName.toString()));

        PipelineLauncher.LaunchConfig writeOptions =
                PipelineLauncher.LaunchConfig.builder("write-pubsub")
                        .setSdk(PipelineLauncher.Sdk.JAVA)
                        .setPipeline(writePipeline)
                        .addParameter("runner", configuration.runner)
                        .addParameter("streaming", "true")
                        .build();

        PipelineLauncher.LaunchInfo writeLaunchInfo = pipelineLauncher.launch(project, region, writeOptions);


        readPipeline
                .apply("Read from PubSub", PubsubIO.readStrings().fromSubscription(subscription.toString()))
                .apply("Counting element", ParDo.of(new CountingFn<>(READ_ELEMENT_METRIC_NAME)));

        PipelineLauncher.LaunchConfig readOptions =
                PipelineLauncher.LaunchConfig.builder("read-pubsub")
                        .setSdk(PipelineLauncher.Sdk.JAVA)
                        .setPipeline(readPipeline)
                        .addParameter("runner", configuration.runner)
                        .addParameter("streaming", "true")
                        .build();

        PipelineLauncher.LaunchInfo readLaunchInfo = pipelineLauncher.launch(project, region, readOptions);

        try {
            PipelineOperator.Result result =
                    pipelineOperator.waitUntilDone(
                            createConfig(readLaunchInfo, Duration.ofMinutes(configuration.pipelineTimeout)));

            // Check the initial launch didn't fail
            assertNotEquals(PipelineOperator.Result.LAUNCH_FAILED, result);
            // streaming read pipeline does not end itself
            // Fail the test if read pipeline (streaming) not in running state.
            assertEquals(
                    PipelineLauncher.JobState.RUNNING,
                    pipelineLauncher.getJobStatus(project, region, readLaunchInfo.jobId()));

            // check metrics
            double numRecords =
                    pipelineLauncher.getMetric(
                            project,
                            region,
                            readLaunchInfo.jobId(),
                            getBeamMetricsName(PipelineMetricsType.COUNTER, READ_ELEMENT_METRIC_NAME));
            assertEquals(configuration.numRecords, numRecords, numOfSourceBundles);
        } finally {
            if (pipelineLauncher.getJobStatus(project, region, writeLaunchInfo.jobId())
                    == PipelineLauncher.JobState.RUNNING) {
                pipelineLauncher.cancelJob(project, region, writeLaunchInfo.jobId());
            }
            if (pipelineLauncher.getJobStatus(project, region, readLaunchInfo.jobId())
                    == PipelineLauncher.JobState.RUNNING) {
                pipelineLauncher.cancelJob(project, region, readLaunchInfo.jobId());
            }
        }
    }

//    private enum ReadAndWriteType {
//        MESSAGE,
//        STRING,
//        AVRO,
//        PROTO
//    }

    private static class MapKVtoString extends DoFn<KV<byte[], byte[]>, String> {
        @ProcessElement
        public void process(ProcessContext context) {
            context.output(new String(Objects.requireNonNull(context.element()).getValue(), StandardCharsets.UTF_8));
        }
    }

    static class Configuration extends SyntheticSourceOptions {
        /** Pipeline timeout in minutes. Must be a positive value. */
        @JsonProperty public int pipelineTimeout = 20;

        /** Runner specified to run the pipeline. */
        @JsonProperty public String runner = "DirectRunner";
    }
}
