# harmony-one-to-bigquery



A golang application to import Harmony ONE blockchain data into GCP BigQuery. The
overall objective of this program is to request the most recent block number submitted
to the Harmony One blockchain. Then retrieve the most recent blockchain data inserted
into GCP BigQuery. Then begin backfilling the blockchain data into BigQuery by making
RPC requests for each block missing from BigQuery.

#### Example

Most recent block header specifies block number: 0xd92e14 -> 14233108

Most recent block number found in BigQuery: 0xd92e0c -> 14233100

So BigQuery is missing 8 of the most recent blocks and their transactions. So
it will attempt to retrieve the first of those 8 missing blocks by starting at 14233100
and working it's way to 14233108.

## Pre-requisites

In order to use the backfill binary in production you will need to have access to
a GCP project that has access to BigQuery. Since this program also uses "streaming
inserts" to insert data into BigQuery your GCP project will need to have billing
enabled. Once you have a project that is up and running, has access to BigQuery and
has billing enable you should generate application credential for the program to
use. A complete guide on how to get that can be found [here](https://cloud.google.com/docs/authentication/getting-started).

In order to use everything that was utilized in this program it is recommended that
you use and understand [Docker](https://www.docker.com/). Additionally it is also
recommended to have a basic understanding of [Kubernetes](https://cloud.google.com/kubernetes-engine).

These were key to getting this program running in a production environment that
allowed for consistent runtimes without the issues of trying to load ~1.4M blocks
with transactions into GCP BigQuery on a local machine.

## Building

### Locally

You can build a local version of the binary that is specific to your OS type by
just using the `go build` command.

```
$ go build -o ./bin/hmy-bq-import ./cmd/hmy-bq-import
```

The resulting binary in the `./bin` folder will allow you to run the program locally
and begin backfilling Harmony One blockchain data into your GCP BigQuery project.

### Docker (Recommended)

You can also build a dockerized version of the application by first build a linux
specific binary and placing that in the `docker/artifacts` folder.

```
$ GOOS=linux GOARCH=amd64 go build -o ./docker/artifacts/hmy-bq-import ./cmd/hmy-bq-import
```

The next step you will need is to copy your credetials file into the artifacts folder.

When you dowloaded your Google application credentials file it will have more than
likely ended up in your Downloads folder. If you copy that over to the `docker/artifacts`
folder it will mounted to the docker image on build. Replace the source part of the
cp command with where you credentials file lives.

```
$ cp $HOME/path/to/Downloads/google-application-credentials.json ./docker/artifacts/harmonyone-gcp-bigquery.json
```

Now you should be ready to perform a docker build command to build the hmy-bq-import
docker image. Now when you created your GCP project you should have been given a project-id
to reference, we will use this as it will be useful when pushing the image to GCR (Google Cloud Registry).

```
export PROJECT_ID=your-project-id
```

Next lets build the docker image from the docker folder in the project

```
$ cd ./docker
$ docker build -t grc.io/${PROJECT_ID}/hmy-bq-import:v1 .
```

Verify that the image was built:

```
$ docker images
REPOSITORY                                     TAG       IMAGE ID       CREATED        SIZE
gcr.io/${PROJECT_ID}/hmy-bq-import             v1        acafe4ca74a5   10 hours ago   23MB
```

## Running

#### Locally

You can simply run the binary using the `backfill` command:

```
./bin/hmy-bq-import backfill --gcp-project-id $PROJECT_ID --help
```

#### Docker

You can check that the docker image build works by running the following `docker run` command.

```
$ docker run -it --rm \
  -e GOOGLE_APPLICTION_CREDENTIALS=/etc/hmy/harmonyone-gcp-bigquery.json \
  -e GCP_PROJECT_ID=${PROJECT_ID \
  gcr.io/${PROJECT_ID}/hmy-bq-import:v1
```

## Using Kubernetes

Since this application is dockerized it can be run in Kubernetes. And this was deployed
to GCP Kubernetes Engine to allow for the backfill to continuously run and keep the
public dataset as close to realtime as possible.

If you wish to run this application on Kubernetes in GCP a [quickstart guide](https://cloud.google.com/kubernetes-engine/docs/quickstart)
will be able to help you do so.

## Environment Variables

| Env Var Name        | Description                                                                | Default Value            | Required |
|---------------------|----------------------------------------------------------------------------|--------------------------|----------|
| NODE_URL            | the url of the node used to pull historical data from                      | https://api.s0.t.hmny.io | N        |
| GCP_PROJECT_ID      | the project id used in GCP to store blockchain data in BigQuery            |                          | Y        |
| GCP_DATASET_ID      | the dataset id used in GCP to store blockchain data in BigQuery            | crypto_harmony           | N        |
| GCP_BLOCKS_TABLE_ID | the blocks table id used in GCP to store blockchain data in BigQuery       | blocks                   | N        |
| GCP_TXNS_TABLE_ID   | the transactions table id used in GCP to store blockchain data in BigQuery | transactions             | N        |
| CONCURRENCY         | the number concurrent go routines pulling Harmony One blockchain data      | 1                        | N        |

## LICENSE

Mozilla Public License Version 2.0
