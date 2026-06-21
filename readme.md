# Yoink - Golang Native Search Engine

<p align="center">
  <strong>A distributed web search engine built from scratch in Go</strong><br>
  Discovers, crawls, validates, indexes, and ranks pages across the web using cloud-native infrastructure.
</p>

<p align="center">
  <a href="https://go.dev/">
    <img alt="Go Version" src="https://img.shields.io/github/go-mod/go-version/DarrylMathias/yoink">
  </a>
  <a href="./LICENSE">
    <img alt="License: Apache 2.0" src="https://img.shields.io/badge/License-Apache_2.0-yellow.svg">
  </a>
  <a href="https://aws.amazon.com/">
    <img alt="AWS" src="https://img.shields.io/badge/cloud-AWS-orange.svg">
  </a>
</p>

---

## Contents

* [Overview](#overview)
* [Architecture](#architecture)
* [Infrastructure](#infrastructure)
* [Installation and Setup](#installation-and-setup)
* [Configuration](#configuration)
* [How it Works](#how-it-works)
* [Project Structure](#project-structure)
* [Roadmap](#roadmap)
* [Contributing](#contributing)
* [License](#license)

---

## Overview

Yoink is a distributed search engine natively written entirely in Go.

The project aims to explore how modern search engines work internally by building every component from scratch:

* URL discovery
* Distributed crawling
* Page validation
* Content extraction
* Deduplication
* Indexing
* Ranking
* Search retrieval

The crawler continuously discovers new pages starting from seed URLs, validates them, stores crawl metadata, and prepares content for indexing. The system is designed around cloud-native services and asynchronous message processing, allowing large-scale crawling without tightly coupled workers.

Current infrastructure processes millions of discovered URLs while maintaining a filtered set of crawlable pages for indexing.

---

## Architecture

### High-Level Flow

```text
Seed URLs
     │
     ▼
URL Discovery
     │
     ▼
Amazon SQS
     │
     ▼
Crawler Workers
     │
     ▼
Page Validation
     │
     ▼
 Redis Cache
     │
     ▼
PostgreSQL Storage
     │
     ▼
 S3 Storage
     │
     ▼
  Indexer
     │
     ▼
Search Engine
```

### Core Components

#### Discovery System

Responsible for finding new URLs from seed pages and continuously expanding the crawl graph.

Features:

* Breadth-first crawling
* Duplicate URL filtering
* Crawl queue generation
* Distributed task scheduling

#### Crawl Workers

Stateless Go workers that consume URLs from SQS and process pages independently.

Features:

* Concurrent crawling
* Retry handling
* Failure recovery
* Scalable horizontal deployment

#### Validation Pipeline

Before a page is accepted into the index, Yoink performs validation checks such as:

* URL normalization
* Duplicate detection
* Crawlability checks
* Response verification
* Content filtering

#### Storage Layer

* **PostgreSQL (RDS)** for persistent crawl metadata and indexed documents
* **Redis** for caching and high-speed lookups
* **Amazon S3** for storing large crawl artifacts and snapshots

---

## Infrastructure

Yoink currently runs on AWS using managed services.

### AWS Services

| Service                    | Purpose                 |
| -------------------------- | ----------------------- |
| Amazon SQS                 | Distributed crawl queue |
| Amazon RDS (PostgreSQL)    | Primary database        |
| Redis                      | Cache and fast lookups  |
| Amazon S3                  | Object storage          |
| EC2                        | Worker execution        |

### Design Goals

* Horizontally scalable
* Fault tolerant
* Cloud-native
* Queue-driven architecture
* Low operational complexity

---

## Installation and Setup

### Prerequisites

* Go 1.24+
* PostgreSQL
* Redis
* AWS Account
* AWS Credentials configured locally

### Clone Repository

```bash
git clone https://github.com/DarrylMathias/yoink.git
cd yoink
```

### Install Dependencies

```bash
go mod download
```

### Environment Variables

Create a `.env.local` file:

```env
HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=
DB_NAME=yoink_local
# DB_SSL_ROOT_CERT=

AWS_ACCESS_KEY_ID=
AWS_SECRET_ACCESS_KEY=
SQS_NAME=
S3_BUCKET_NAME=

REDIS_ADDRESS=
REDIS_USERNAME=
REDIS_PASSWORD=
REDIS_DATABASE=0

UPSTASH_REDIS_REST_URL=
UPSTASH_REDIS_REST_TOKEN=

# For e-mailing alerts and hearbeats
RESEND_API_KEY=

PORT=
APPLICATION=dev
```

---

## Configuration

Most configuration is handled through global variables in [`app/app.go`](./app/app.go).

These values can be adjusted depending on available infrastructure and crawl requirements.

For dev environments, `set err := env.NewEnv(".env.local")` in `app/app.go`

---

## How it Works

### 1. Discovery

The system begins with a set of seed URLs.

Each page is fetched and analyzed for outgoing links.

Discovered URLs are normalized and pushed into Amazon SQS for processing.

---

### 2. Queue Processing

Worker nodes continuously consume URLs from the queue.

Each worker:

* Downloads page content
* Validates the response
* Extracts metadata
* Identifies new links

---

### 3. Validation

The validation layer removes:

* Duplicate URLs
* Invalid responses
* Unsupported content
* Non-crawlable resources

Only valid pages move forward.

---

### 4. Storage

Validated pages are stored in PostgreSQL while frequently accessed crawl data is cached in Redis.

Large assets and crawl artifacts can be persisted to S3.

---

### 5. Indexing

The indexing layer transforms raw page content into searchable structures.

Planned indexing features include:

* Inverted indexes
* Tokenization
* Relevance scoring
* Document ranking

---

### 6. Search

The final stage serves ranked search results to users.

Future ranking systems may incorporate:

* Keyword relevance
* Link analysis
* Authority scoring
* Freshness signals

---

## Project Structure

```text
Directory structure:
└── darrylmathias-yoink/
    ├── go.mod
    ├── go.sum
    ├── LICENSE
    ├── main.go
    ├── app/
    │   └── app.go
    ├── crawler/
    │   ├── crawler.go
    │   ├── extract/
    │   │   ├── extract.go
    │   │   ├── dedup/
    │   │   │   └── dedup.go
    │   │   ├── download/
    │   │   │   └── download.go
    │   │   └── metadata/
    │   │       └── metadata.go
    │   ├── store/
    │   │   └── store.go
    │   └── validate/
    │       ├── validate.go
    │       └── hashtable/
    │           └── hash.go
    ├── models/
    │   ├── myURL.go
    │   └── page.go
    ├── seed/
    │   └── sqs_seed.go
    └── utils/
        ├── url.go
        ├── database/
        │   └── database.go
        ├── error/
        │   └── error.go
        ├── myaws/
        │   ├── config.go
        │   ├── s3/
        │   │   └── s3.go
        │   └── sqs/
        │       └── sqs.go
        ├── redis/
        │   └── redis.go
        ├── resend/
        │   └── resend.go
        └── upstash/
            └── upstash.go

```

---

## Current Status

### Completed

* AWS infrastructure setup
* PostgreSQL integration
* Redis integration
* SQS queue architecture
* URL discovery pipeline
* Validation system
* Distributed worker communication
* S3 integration
* Large-scale crawling

### In Progress

* Index generation
* Search ranking
* Query engine
* Distributed indexing
* PageRank-style scoring
* Full-text search
* Search API
* Search frontend

---

## Why Build a Search Engine?

Search engines combine nearly every major systems engineering discipline:

* Distributed systems
* Networking
* Databases
* Information retrieval
* Large-scale data processing
* Cloud infrastructure

Yoink exists primarily as a learning project to explore these concepts through a real-world implementation instead of isolated examples.

---

## Contributing

Contributions, issues, and suggestions are welcome.

Areas of particular interest:

* Distributed crawling
* Search ranking algorithms
* Information retrieval
* Index optimization
* Cloud infrastructure

---

## License

Apache 2.0 License
