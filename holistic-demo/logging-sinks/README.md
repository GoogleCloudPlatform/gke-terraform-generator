# Logging with Stackdriver on Kubernetes Engine

## Table of Contents

* [Introduction](#introduction)
* [Architecture](#architecture)
* [Prerequisites](#prerequisites)
  * [Enable GCP APIs](#enable-gcp-apis)
  * [Install Cloud SDK](#install-cloud-sdk)
  * [Install Terraform](#install-terraform)
* [Deployment](#deployment)
  * [How does it work?](#how-does-it-work)
  * [Running Terraform](#running-terraform)
* [Validation](#validation)
  * [Generating Logs](#generating-logs)
  * [Logs in the Stackdriver UI](#logs-in-the-stackdriver-ui)
  * [Viewing Log Exports](#viewing-log-exports)
  * [Logs in Cloud Storage](#logs-in-cloud-storage)
  * [Logs in BigQuery](#logs-in-bigquery)
* [Teardown](#teardown)
  * [Next Steps](#next-steps)
* [Troubleshooting](#troubleshooting)
* [Relevant Material](#relevant-material)

## Introduction
Stackdriver Logging can be used aggregate logs from all GCP resources as well as any custom resources (on other platforms) to allow for one centralized store for all logs and metrics.  Logs are aggregated and then viewable within the provided Stackdriver Logging UI. They can also be [exported to Sinks](https://cloud.google.com/logging/docs/export/configure_export_v2) to support more specialized of use cases.  Currently, Stackdriver Logging supports exporting to the following sinks:
* Cloud Storage
* Pub/Sub
* BigQuery

This document will describe the steps required to deploy a sample application to Kubernetes Engine that forwards log events to [Stackdriver Logging](https://cloud.google.com/logging/). It makes use of [Terraform](https://www.terraform.io/), a declarative [Infrastructure as Code](https://en.wikipedia.org/wiki/Infrastructure_as_Code) tool that enables configuration files to be used to automate the deployment and evolution of infrastructure in the cloud.  The configuration will also create a Cloud Storage bucket and a BigQuery dataset for exporting log data to.

## Architecture

The Terraform configurations are going to build a Kubernetes Engine cluster that will generate logs and metrics that can be ingested by Stackdriver.  The scripts will also build out Logging Export Sinks for Cloud Storage, BigQuery, and Cloud Pub/Sub.  The diagram of how this will look along with the data flow can be seen in the following graphic.

![Logging Architecture](docs/logging-architecture.png)

## Prerequisites

The steps described in this document require the installation of several tools and the proper configuration of authentication to allow them to access your GCP resources.

### Cloud Project

You'll need access to a Google Cloud Project with billing enabled. See **Creating and Managing Projects** (https://cloud.google.com/resource-manager/docs/creating-managing-projects) for creating a new project. To make cleanup easier it's recommended to create a new project.

### Required GCP APIs

The following APIs will be enabled in the
* Cloud Resource Manager API
* Kubernetes Engine API
* Stackdriver Logging API
* Stackdriver Monitoring API
* BigQuery API

### Tools
1. [Terraform >= 0.11.7](https://www.terraform.io/downloads.html)
2. [Google Cloud SDK version >= 204.0.0](https://cloud.google.com/sdk/docs/downloads-versioned-archives)
3. [kubectl matching the latest GKE version](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
4. bash or bash compatible shell
5. [GNU Make 3.x or later](https://www.gnu.org/software/make/)
6. A Google Cloud Platform project where you have permission to create networks

#### Install Cloud SDK
The Google Cloud SDK is used to interact with your GCP resources.
[Installation instructions](https://cloud.google.com/sdk/downloads) for multiple platforms are available online.

#### Install kubectl CLI

The kubectl CLI is used to interteract with both Kubernetes Engine and kubernetes in general.
[Installation instructions](https://cloud.google.com/kubernetes-engine/docs/quickstart)
for multiple platforms are available online.

#### Install Terraform

Terraform is used to automate the manipulation of cloud infrastructure. Its
[installation instructions](https://www.terraform.io/intro/getting-started/install.html) are also available online.

### Configure Authentication

The Terraform configuration will execute against your GCP environment and create various resources.  The script will use your personal account to build out these resources.  To setup the default account the script will use, run the following command to select the appropriate account:

`gcloud auth application-default login`


## Validation

If no errors are displayed during deployment, after a few minutes you should see your Kubernetes Engine cluster in the [GCP Console](https://console.cloud.google.com/kubernetes) with the sample application deployed.

Now that the application is deployed to Kubernetes Engine we can generate log data and use the [Stackdriver UI](https://console.cloud.google.com/logs) and other tools to view it.

### Generating Logs

The sample application that Terraform deployed serves up a simple web page.  Each time you open this application in your browser the application will publish log events to Stackdriver Logging. Refresh the page a few times to produce several log events.

To get the URL for the application page you must perform the following steps:

1. In the GCP console navigate to the **Networking -> Network services** page.
2. On the default **Load balancing** page that shows up, click on the TCP load balancer that was setup.
3. On the **Load balancer details** page there is a top section labeled **Frontend**.  Note the IP:Port value as this will be used in the upcoming steps.

Using the IP:Port value you can now access the application.  Go to a browser and enter the URL.  The browser should return a screen that looks similar to the following:

![Sample application screen](docs/application-screen.png)

### Logs in the Stackdriver UI

Stackdriver provides a UI for viewing log events. Basic search and filtering features are provided, which can be useful when debugging system issues. The Stackdriver Logging UI is best suited to exploring more recent log events. Users requiring longer-term storage of log events should consider some of the tools in following sections.

To access the Stackdriver Logging console perform the following steps:

1. In the GCP console navigate to the **Stackdriver -> Logging** page.
2. The default logging console will load.  On this page change the resource filter to be **GKE Container -> stackdriver-logging -> default** (the **stackdriver-logging** is the cluster; and the **default** is the namespace).  Your screen should look similar the screenshot below.
3. On this screen you can expand the bulleted log items to view more complete details about the log entry.

![Logging Console](docs/loggingconsole.png)

In the logging console you can perform any type of text search, or try the various filters by log type, log level, timeframe, etc.

### Viewing Log Exports

The Terraform configuration built out two Log Export Sinks.  To view the sinks perform the following steps:

1. In the GCP console navigate to the **Stackdriver -> Logging** page.
2. The default logging console will load.  On the left navigation click on the **Exports** menu option.
3. This will bring you to the **Exports** page.  You should see two Sinks in the list of log exports.
4. You can edit/view these sinks by clicking on the context menu to the right and selecting the **Edit sink** option.
5. Additionally, you could create additional custom export sinks by clicking on the **Create Export** option in the top of the navigation window.

### Logs in Cloud Storage

Log events can be stored in [Cloud Storage](https://cloud.google.com/storage/), an object storage system suitable for archiving data. Policies can be configured for Cloud Storage buckets that, for instance, allow aging data to expire and be deleted while more recent data can be stored with a variety of storage classes affecting price and availability.

The Terraform configuration created a Cloud Storage Bucket named stackdriver-gke-logging-<random-Id> to which logs will be exported for medium to long-term archival.  In this example, the Storage Class for the bucket is defined as Nearline because the logs should be infrequently accessed in a normal production environment (this will help to manage the costs of medium-term storage).  In a production scenario this bucket may also include a lifecycle policy that moves the content to Coldline storage for cheaper long-term storage of logs.

To access the Stackdriver logs in Cloud Storage perform the following steps:

**Note:** Logs from Cloud Storage Export are not populated immediately.  It may take up to 2-3 hours for logs to appear.

1. In the GCP console navigate to the **Storage -> Storage** page.
2. This loads the Cloud Storage Browser page.  On the page find the Bucket with the name stackdriver-gke-logging-<random-Id>, and click on the name (which is a hyperlink).
3. This will show the details of the bucket.  You should see a list of directories corresponding to pods running in the cluster (eg autoscaler, dnsmasq, etc.).

![Cloud Storage Bucket](docs/cloudstoragebucket.png)

On this page you can click into any of the named folders to browse specific log details like heapster, kubedns, sidecar, etc.

### Logs in BigQuery

Stackdriver log events can be configured to be published to [BigQuery](https://cloud.google.com/bigquery/), a data warehouse tool that supports fast, sophisticated, querying over large  data sets.

The Terraform configuration will create a BigQuery [DataSet](https://cloud.google.com/bigquery/docs/reference/rest/v2/datasets) named gke_logs_dataset.  This dataset will be setup to include all Kubernetes Engine related logs for the last hour (by setting a Default Table Expiration for the dataset).  A Stackdriver Export will be created that pushes Kubernetes Engine container logs to the dataset.

To access the Stackdriver logs in BigQuery perform the following steps:

**Note:** The BigQuery Export is not populated immediately.  It may take a few minutes for logs to appear.

1. In the GCP console navigate to the **Big Data -> BigQuery** page.
2. This loads a new browser tab with the BigQuery console.
3. On the left hand console you will have a display of the datasets you have access to.  You should see a dataset named **gke_logs_dataset**.  Expand this dataset to view the tables that exist (**Note:** The dataset is created immediately, but the tables are what is generated as logs are written and new tables are needed).
4. Click on one of the tables to view the table details.  Your screen should look similar to the screenshot below.
5. Review the schema of the table to note the column names and their data types.  This information can be used in the next step when we query the table to look at the data.

![BigQuery](docs/bigquery.png)

5. Click on the **Query Table** towards the top right to perform a custom query against the table.
6. This opens the query window.  You can simply add an asterisk (*) after the **Select** in the window to pull all details from the current table. **Note:** A 'Select *' query is generally very expensive and not advised.  For this tutorial the dataset is limited to only the last hour of logs so the overall dataset is relatively small.
7. Click the **Run Query** button to execute the query and return some results from the table.
8. A popup window till ask you to confirm running the query.  Click the **Run Query** button on this window as well.
9. The results window should display some rows and columns.  You can scroll through the various rows of data that are returned, or download the results to a local file.
10. Execute some custom queries that filter for specific data based on the results that were shown in the original query.


## Teardown

When you are finished with this example you will want to clean up thebigquery data set before tearing down the cluster.

```

git add teardown.sh
git commit -m ‘<Description>’
git push origin master

```

Since Terraform tracks the resources it created it is able to tear them all down.

### Next Steps

Having used Terraform to deploy an application to Kubernetes Engine, generated logs, and viewed them in Stackdriver, you might consider exploring [Stackdriver Monitoring](https://cloud.google.com/monitoring/) and [Stackdriver Tracing](https://cloud.google.com/trace/). Examples for these topics are available [here](https://github.com/GoogleCloudPlatform?q=gke-tracing-demo++OR+gke-monitoring-tutorial) and build on the work performed with this document.

## Troubleshooting

** The install script fails with a `Permission denied` when running Terraform. **
The credentials that Terraform is using do not provide the
necessary permissions to create resources in the selected projects. Ensure
that the account listed in `gcloud config list` has necessary permissions to
create resources. If it does, regenerate the application default credentials
using `gcloud auth application-default login`.

** Cloud Storage Bucket not populated **
Once the Terraform configuration is complete the Cloud Storage Bucket will be created
but it is not always populated immediately with log data from the Kubernetes Engine cluster.  The logs
details rarely populate in the bucket immediately.  Give the process some time because it can take
up to 2 to 3 hours before the first entries start appearing (https://cloud.google.com/logging/docs/export/using_exported_logs).

** No tables created in the BigQuery dataset **
Once the Terraform configuration is complete the BigQuery Dataset will be created
but it will not always have tables created in it by the time you go to review the results.  The
tables are rarely populated immediately.  Give the process some time (minimum of 5 minutes)
before determining that something is not working properly.

## Relevant Material
* [Kubernetes Engine Logging](https://cloud.google.com/kubernetes-engine/docs/how-to/logging)
* [Viewing Logs](https://cloud.google.com/logging/docs/view/overview)
* [Advanced Logs Filters](https://cloud.google.com/logging/docs/view/advanced-filters)
* [Overview of Logs Exports](https://cloud.google.com/logging/docs/export/)
* [Procesing Logs at Scale Using Cloud Dataflow](https://cloud.google.com/solutions/processing-logs-at-scale-using-dataflow)
* [Terraform Google Cloud Provider](https://www.terraform.io/docs/providers/google/index.html)

**This is not an officially supported Google product**
