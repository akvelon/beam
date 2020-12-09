# Dataflow flex templates - Streaming Beam KafkaToPubsub

The [KafkaToPubsub](https://link_to_pipeline) pipeline is a streaming pipeline which ingests data in JSON format
from Kafka, and outputs the resulting records to PubSub. Input topics, output topic, Bootstrap
servers are specified by the user as template parameters.

## Pipeline requirements

* Kafka Bootstrap Server(s).
* Kafka Topic(s) exists.
* The PubSub output topic exists.

## Example usage

### 1. Set the pipeline vars

```sh
export KAFKA_TOPICS="<my-topic1>,<my-topic2>"
export PUB_SUB_TOPIC="projects/<my-project>/topics/<my-topic>"
export BUCKET="<my-bucket>"
export PROJECT="$(gcloud config get-value project)"
```

### 2. (Optional) Enable to use Kaniko cache by default

You can speed up subsequent builds with
[Kaniko cache](https://cloud.google.com/cloud-build/docs/kaniko-cache)
in Cloud Build.

```sh
# (Optional) Enable to use Kaniko cache by default.
gcloud config set builds/use_kaniko True
```

### 3. Build a Docker image using a `Dockerfile` and save it into
[Container Registry](https://cloud.google.com/container-registry/).
Images starting with `gcr.io/PROJECT/` are saved into your project's
Container Registry, where the image is accessible to other Google Cloud products.

```sh
export TEMPLATE_IMAGE="gcr.io/$PROJECT/samples/dataflow/kafka-to-pubsub:latest"

# Build the image into Container Registry, this is roughly equivalent to:
#   gcloud auth configure-docker
#   docker image build -t $TEMPLATE_IMAGE .
#   docker push $TEMPLATE_IMAGE
gcloud builds submit --tag "$TEMPLATE_IMAGE" .
```

### 4. Creating a Flex Template

To run a template, you need to create a *template spec* file containing all the
necessary information to run the job, such as the SDK information and metadata.

The [`metadata.json`](metadata.json) file contains additional information for
the template such as the "name", "description", and input "parameters" field.

The template file must be created in a Cloud Storage location,
and is used to run a new Dataflow job.

```sh
export TEMPLATE_PATH="gs://$BUCKET/samples/dataflow/templates/streaming-kafka-beam-sql.json"

# Build the Flex Template.
gcloud dataflow flex-template build $TEMPLATE_PATH \
--image "$TEMPLATE_IMAGE" \
--sdk-language "PYTHON" \
--metadata-file "metadata.json"
```

The template is now available through the template file in the Cloud Storage
location that you specified.

### 5. Running a Dataflow Flex Template pipeline

```sh
# Run the Flex Template.
curl -X POST \
"https://dataflow.googleapis.com/v1b3/projects/$PROJECT/locations/us-central1/flexTemplates:launch" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $(gcloud auth print-access-token)" \
-d '{
"launch_parameter": {
"jobName": "streaming-kafka-to-pubsub-'$(date +%Y%m%d-%H%M%S)'",
"parameters": {
"bootstrap_servers":"'broker_1:9092,broker_2:9092'",
"input_topics": "'$KAFKA_TOPICS'",
"output_topic": "'$PUB_SUB_TOPIC'",
"experiments":"use_runner_v2"
},
"containerSpecGcsPath": "'$TEMPLATE_PATH'"
}
}'
```

## Limitations

* You must use a Google-provided base image to package your containers using Docker.
* You can use only Dataflow runner, use 313.0.1 (2020-10-07) Dataflow or later. Related issues:
    * Google Cloud SDK [Fixed bug with gcloud dataflow flex-template build](https://cloud.google.com/sdk/docs/release-notes#31301_2020-10-07)
    * [BEAM-10484](https://issues.apache.org/jira/browse/BEAM-10484)
    * [BEAM-6868](https://issues.apache.org/jira/browse/BEAM-6868)
    * fixed in v2.0 [BEAM-1573](https://issues.apache.org/jira/browse/BEAM-1573)
* Kafka must not contain messages with null keys due to an open issue 
    * Kafka XLang fails for ?empty? key/values [BEAM-10529](https://issues.apache.org/jira/browse/BEAM-10529)

📝 Docs: [Using Flex Templates](https://cloud.google.com/dataflow/docs/guides/templates/using-flex-templates)
