# LifeTrack API
**Neutrino _LifeTrack_** is a platform where you can track your life habits and get useful insights from them. 
This improves your productivity and keeps you on track in case you need some help for life management.

## Features
**_LifeTrack_** offers a simple, yet useful system to manage your activities.

You may group your activities by categories to get aggregated data from related activities, _for example:_

`Category = "Science" -> Activities = ["Physics", "Chemistry", "Math"] -> Occrrence = ["01/31/2020", "02/12/2020"]`

_This could also be represented with **LaTeX notation** as_:

![x](https://latex.codecogs.com/gif.latex?%5Calpha%20%3D%20%5Cbegin%7Bpmatrix%7D%20a_%7B1%7D%20%5C%5C%20a_%7B2%7D%20%5C%5C%20a_%7B3%7D%20%5C%5C%20...%20%5C%5C%20a_%7Bn%7D%20%5Cend%7Bpmatrix%7D%20%5Cotimes%20%5Cbegin%7Bpmatrix%7D%20o_%7B1%7D%20%5C%5C%20o_%7B2%7D%20%5C%5C%20o_%7B3%7D%20%5C%5C%20...%20%5C%5C%20o_%7Bn%7D%20%5Cend%7Bpmatrix%7D%3B%20o_%7Bn%7D%20%5Csubset%20a_%7Bn%7D)

**Where**
- **alpha** = _Category matrix_
- **a** = _Activity vector_
- **b** = _Occrrence vector_

**Mathematical Domain**
- o is a subset of a
- alpha is the tensor product between a and o vectorial fields

These are the main _features_:

- **Occurrence**: _Event_ that happened inside an activity.
- **Activity**: _Task you use to do_ every week.
- **Category**: _Group of activities_.
- **Insight**: _Track your quantified activities_ every week, month and year.
- **Reminder**: _Get reminded about your tasks_ every day or week.

## Technology
**Neutrino _LifeTrack_** makes use of _serverless ecosystems_ to keep itself simple to manage and orchestrate.
We make use of blazingly fast technology such as Go, AWS DynamoDB, Memcached and CassandraDB to satisfy our needs.

- **Go**: The Go Programming Language is a language created by Google which is fast and makes easy to handle high-concurrent systems.
- **Hashicorp Terraform**: Hashicorp Terraform is a IaaC solution which offers infrastructure orchestration with pure code (HCL).
- **AWS**: The Amazon Web Service platform, offers many integrated infrastructure services.
  - **Route53**: Route53 is a DNS web service which helps us to publish our services through _domain names_.
  - **API Gateway**: API Gateway is the front door of most of webservices, handles and routes incoming requests.
  - **Lambda**: AWS Lambda is a serverless function stored in the cloud.
  - **DynamoDB**: DynamoDB is a Key-Value high-available database which is managed by AWS automatically.
  - **S3** _(Standard/Infrequent Access)_: S3 is the standard file storage of AWS, it is self managed and is high-available.
  - **CloudFront**: CloudFront is a Content Distributed Network or CDN, which propagates static data through edge-points all over the world to reduce latency.
  - **ACM**: Amazon Certificate Manages is an SSL certification manager which are mainly used to secure connections with the TLS protocol.
  - **CLi** _v2_: The AWS Command-Line Tool.
  - **IAM**: Identity Access Manager, manages users and policies to give/restrict access to AWS resources.
  - **SNS**: Simple Notification Service, offers a publish/subscribe mechanism to publish messages to n-consumers.
  - **SQS**: Simple Queue Service, offers a distributed and high-available queue system that is mainly focused on ETL jobs, Cron Jobs or as a Message Broker when used with SNS.
  - **ElastiCache**: Caching service which offers self-managed computation for either Redis or Memcached systems.
  - **KeySpaces**: Document-oriented service wich offers self-managed computation for Apache Cassandra.
  - **XRay**: AWS XRay is a remote-request distributed tracing mechanism.
  - **CloudWatch**: CloudWatch is the standard AWS monitoring ecosystem, it offers many functionalities apart from just monitoring such as event triggering straight to SNS to notify specified consumers.

- **Firebase**: Google's application ecosystem.
  - **Authentication** _OAuth2 (Google, Facebook & Apple)_: OAuth2 is a authentication mechanism using the latest security standards.
