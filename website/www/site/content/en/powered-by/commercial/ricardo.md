---
type: commercial
categories:
  - Online & Traditional Retail
icon: /images/logos/powered-by/ricardo.png
---
<!--
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
-->
<div class="post-header1">
    Data streams facilitation with Apache Beam
</div>
<div class="post-block">
    <div class="post-block-image">
        <img src="/images/ricardo_logo.png"/>
    </div>
    <div class="post-block-text">
        “Without Beam, without all this data and real time information, we could not provide the services we are providing and process the volumes of data we are processing”
        <div class="post-block-text-author">
            <div class="post-block-text-author-image">
                <img src="/images/tobias_kaymak_photo.png">
            </div>
            <div class="post-block-text-author-column">
                <div class="post-block-text-author-name">
                    Tobias Kaymak
                </div>
                <div class="post-block-text-author-position">
                    Data Engineer
                </div>
            </div>
        </div>
    </div>
</div>

## Background

Ricardo is a leading online auction marketplace in Switzerland, with over 4 million registered buyers and sellers and more than 5 million articles changing hands through the platform every year. Ricardo needs to process high volumes of streaming events and manage over 5 TB of articles, assets, and analytical data.

With the scale that came from 20 years in the market, Ricardo made the decision to migrate from their on-premises data center to Google Cloud to easily grow and evolve further and reduce operational costs through managed cloud services. Data intelligence and engineering teams took the lead on this transformation and development of new AI/ML-enabled customer experiences.

## Challenge

Migrating from an on-premises data center to the cloud presented Ricardo with an opportunity to modernize their marketplace from heavy legacy reliance on transactional SQL and take advantage of the event-based streaming architecture.

Ricardo’s data intelligence team identified two key success factors: a carefully designed data model and a framework that provides unified stream and batch data pipelines execution, both on-premises and in the cloud.

Ricardo needed a data processing framework that can scale easily, enrich unbounded streams with historic data from multiple sources, provide granular control on data freshness, and provide an abstract pipeline operational infrastructure, thus helping their team focus on creating new value for customers and business.

## Journey to the Cloud

Ricardo’s data intelligence team began modernizing their stack in 2018. They selected frameworks that provide reliable and scalable data processing both on-premises and in the cloud. Apache Beam offered both Flink runner that could execute on-premises and Dataflow runner for managed cloud service for the same pipelines developed using Apache Beam Java SDK. Apache Flink is well known for its reliability and cost-efficiency and an on-premises cluster was spun up as the initial environment.

<div class="post-quote">
    <div class="post-quote-position">
        <div class="post-quote-text"> 
            “We wanted to implement a solution that would multiply our possibilities, and that’s exactly where Beam comes in. One of the major drivers in this decision was the ability to evolve without adding too much operational load.” 
        </div>
    </div>
    <div class="post-quote-author">
        <div class="post-quote-author-image">
            <img src="/images/quote.png">
        </div>
        <div class="post-quote-author-name">
            Tobias
        </div>
    </div>
</div>

Beam pipelines for core business workloads to ingest events data from Apache Kafka into BigQuery were running stable in just one month.

The flexibility to refresh data every hour, minute, or stream data real-time, depending on the specific use case and need, helped the team improve data freshness which was a significant advancement for Ricardo’s eCommerce platform analytics and reporting.

While having full control over Flink provisioning, the data intelligence team was able to optimize operating costs through cluster resource utilization by executing steaming pipelines with Beam Flink runner. Dataflow runner enabled the team to scale the infrastructure on-demand and benefit from cost advantages for batch pipelines that support most advanced use cases.

With recent changes that streamline connectivity between Flink cluster workloads in GKE, Dataflow workers, and peered connection to managed Kafka service, data intelligence team is looking forward to completing all verifications and moving additional workloads to Dataflow to take advantage of streaming inserts and real-time event response.

<div class="post-quote">
    <div class="post-quote-position">
        <div class="post-quote-text"> 
            “I knew Beam, I knew it works. When you need to move from Kafka to BigQuery and you know that Beam is exactly the right tool, you just need to choose the right executor for it.” 
        </div>
    </div>
    <div class="post-quote-author">
        <div class="post-quote-author-image">
            <img src="/images/quote.png">
        </div>
        <div class="post-quote-author-name">
            Tobias
        </div>
    </div>
</div>

## Evolution of Use Cases 

Thinking of a stream as data in motion, and a table as data at rest provided a fortuitous chance to take a look at some data model decisions that were made as far back as 20 years before. Articles that are on the marketplace have assets that describe them, and for performance and cost optimizations purposes, data entities that belong together were split into separate database instances. Apache Beam enabled Ricardo to join assets and articles streams and optimize BigQuery scans to reduce costs. When designing the pipeline, the team created streams for assets and articles. Since the assets stream is the primary one, they shifted the stream 5 minutes back and created a lookup schema in BigTable. This elegant solution ensures that the assets stream is always processed first while BigTable allows for matching the latest asset to an article and Apache Beam joins them both together.

<div class="post-scheme"> 
    <img src="/images/post_scheme.png">
</div>

The successful case of joining different data streams facilitated further Apache Beam adoption by Ricardo in areas like data science and ML.

<div class="post-quote">
    <div class="post-quote-position">
        <div class="post-quote-text"> 
            “Once you start laying out the simple use cases, you will always figure out the edge case scenarios. Pipeline has been running for a year now, and Beam handles it all, from super simple use cases to something crazy.” 
        </div>
    </div>
    <div class="post-quote-author">
        <div class="post-quote-author-image">
            <img src="/images/quote.png">
        </div>
        <div class="post-quote-author-name">
            Tobias
        </div>
    </div>
</div>

Apache Beam runs real-time ML pipelines to automatically detect product type and category from uploaded product images. Image recognition helps sellers create new articles and define their attributes in a couple of clicks, while improved score rating enables more meaningful search results for buyers.

As an eCommerce retailer, Ricardo faces the increasing scale and sophistication of fraud transactions and takes a strategic approach by employing Beam pipelines for fraud detection and prevention. Beam pipelines call an external intelligent API to identify the signs of fraudulent behaviour, like device characteristics or user activity. Apache Beam stateful processing feature enables Ricardo to apply an associating, accumulating operation to the streams of data. If a payment or transaction was flagged as suspicious, Beam identifies and applies this pattern and verifies whether a similar case had ever been raised to customer care before passing the new case for investigation. For instance, once the external provider identifies an event as a fair transaction, the pattern is applied to all transactions in the stream. Thus, Apache Beam saves Ricardo’s customer care team's time and effort on investigating duplicate cases by filtering out the similarities and variations in the events, right in the streams of data.

Originally implemented by Ricardo’s data intelligence team, Apache Beam has proven to be a powerful framework that supports advanced scenarios and acts as a glue between Kafka, BigQuery, and platform and external APIs, which encouraged other teams at Ricardo, like the research team, to adopt it.

<div class="post-quote">
    <div class="post-quote-position">
        <div class="post-quote-text"> 
            “This is a framework that is so good that other teams are picking up the idea and starting to work with it after we tested it.”
        </div>
    </div>
    <div class="post-quote-author">
        <div class="post-quote-author-image">
            <img src="/images/quote.png">
        </div>
        <div class="post-quote-author-name">
            Tobias
        </div>
    </div>
</div>

## Results

Apache Beam has provided Ricardo with a scalable and reliable data processing framework that supported Ricardo’s fundamental business scenarios and enabled new use cases to respond to events in real-time.

Throughout Ricardo’s transformation, Apache Beam has been a unified framework that can run batch and stream pipelines, offers on-premises and cloud managed services execution, and programming language options like Java and Python, empowered data science and research teams to advance customer experience with new real-time scenarios fast-tracking time to value.

<div class="post-quote">
    <div class="post-quote-position">
        <div class="post-quote-text"> 
            “After this first pipeline, we are working on other use cases and planning to move them to Beam. I was always trying to spread the idea that this a framework that is reliable, it actually helps you to get the stuff done in a consistent way”
        </div>
    </div>
    <div class="post-quote-author">
        <div class="post-quote-author-image">
            <img src="/images/quote.png">
        </div>
        <div class="post-quote-author-name">
            Tobias
        </div>
    </div>
</div>

Apache Beam has been a technology that multiplied possibilities, allowing Ricardo to maximize technology benefits at all stages of their modernization and cloud journey.
