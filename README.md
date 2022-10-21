![kafta logo](img/kafta.png)

Kafta is a modern non-JVM command line for managing Kafka Clusters written in golang. Usability was inspired by kubectl, the interface is simple and supports managing several clusters at the same time in a simple way.

# Table of Contents
- [Table of Contents](#table-of-contents)
- [Overview](#overview)
- [Concepts](#concepts)
- [Installing](#installing)
- [Configuration](#configuration)
- [Commands](#commands)
  - [Topic](#topic)
  - [Consumer Group](#consumer-group)
  - [Cluster](#cluster)
  - [Broker](#broker)
  - [Producer](#producer)
  - [Consumer](#consumer)
  - [Subjects List](#subjects-list)
- [Next features](#next-features)


# Overview

Kafta provides a simple interface to manage topics, brokers, consumer-groups and many things like that. Interfaces are similar to kubectl & go tools.
It is built using [sarama](https://github.com/Shopify/sarama).

Kafta provides:
* Easy commands CLIs: `kafta topic list`, `kafta cluster describe`, etc.
* Fully manage topics, consumers, kafka clusters
* yaml configuration file with all clusters that can be used. You set the current one and can easily switch to another
* Is a binary, not need install JVM or something like that
* Auth use sasl
* Intelligent suggestions (`kafta clustr`... did you mean `kafta cluster`?)
* Help flag recognition of `-h`, `--help`

# Concepts

Kafta is built on a structure of commands, arguments & flags. Kafta will always be interacting on one cluster at a time,
the reason for this is not having to pass which cluster is in each command, as it is with most kafka CLIs.

Kafta was created by developers for developers. We feel the pain of maintaining a kafka cluster using the bash's provided by apache-kafka,
it's confusing and the experience is miserable. To facilitate the adoption of the kafka, the kafta began to be born. The focus of this
project is to be simple to use, it is bad to need to install java, pass the kafka cluster address in every command & mess with xml's.
Kafta is a golang project that is easy to install, easy to configure and the main thing is simple to use.

To see all options exists relate to same command, run:

```
$ kafta topic
Topics management

Usage:
  kafta topic [command]

Available Commands:
  create      Create topics
  delete      Delete topics
  describe    Describe a topic
  list        List topics
```

# Installing

Using Kafta is easy. First, use `go get` to install the latest version
of the library. This command will install the `kafta` executable
along with the library and its dependencies:
- go < 1.18: `go get -u github.com/electric-saw/kafta`
- go >= 1.18: `go install  github.com/electric-saw/kafta/cmd/kafta@latest`


# Configuration

Kafta will create a config file in ~/.kafta/config. This yaml is used to support kafka multi-clusters and avoid passing all addresses every time.

To set up a new cluster, create a new config via kafta, you'll need to provide some information, don't worry, it's all in terminal,
you don't need to edit any XML \o/

Follow the example:

```
$ kafta config set-context production
Bootstrap servers: b-1.mydomain:9092,b-2.mydomain:9092,b-3.mydomain:9092
Schema registry: https://schema-registry.com
Use SASL: y
SASL Algorithm: sha512
User: myuser
✔ Password: ******
```

To list the contexts, run:

```
$ kafta config get-contexts
+---------------+---------------------------+-----------------------------+------+---------+
| NAME          | CLUSTER                   | SCHEMA REGISTRY             | KSQL | CURRENT |
+---------------+---------------------------+-----------------------------+------+---------+
| dev           | b-1.mydomain:9092         | https://schema-registry.com |      | true    |
| production    | b-3.productiondomain:9092 | https://schema-registry.com |      | false   |
+---------------+---------------------------+-----------------------------+------+---------+
```

To change the current cluster, run:

```
$ kafta config use-context production
Switched to context "production".
```

List current context, run:

```
$ kafta config current-context
production
```

Delete context, run:

```
$ kafta config delete-context production
deleted context production from /home/myuser/.kafta/config

```
* it's only config files, kafta don't delete anything of kafka cluster

# Commands

After the context was configureted, it's time to use commands for manage kafka cluster

## Topic

Create new topic:

```
$ kafta topic create my-topic --rf 3 --partitions 10
Topic created
```

There are default values ​​for partition and replication factor, which is why it can only be used without specifying RF or partition.
The topic will be created with RF=3 and partition=10. Example:

```
$ kafta topic create my-topic
Topic created
```

Describe topic:

```
$ kafta topic describe my-topic
+--------------+------------+----------+
| NAME         | PARTITIONS | INTERNAL |
+--------------+------------+----------+
| my-topic     |         10 | false    |
+--------------+------------+----------+

+----+---------+--------+----------+------------------+
| ID | ISR     | LEADER | REPLICAS | OFFLINE REPLICAS |
+----+---------+--------+----------+------------------+
|  0 | [3 2 1] |      3 | [3 1 2]  | []               |
|  1 | [2 1 0] |      2 | [2 0 1]  | []               |
|  2 | [5 3 1] |      1 | [1 5 3]  | []               |
|  3 | [5 4 0] |      0 | [0 4 5]  | []               |
|  4 | [5 1 0] |      5 | [5 1 0]  | []               |
|  5 | [5 4 0] |      4 | [4 0 5]  | []               |
|  6 | [5 4 3] |      3 | [3 5 4]  | []               |
|  7 | [4 3 2] |      2 | [2 4 3]  | []               |
|  8 | [3 2 1] |      1 | [1 3 2]  | []               |
|  9 | [2 1 0] |      0 | [0 2 1]  | []               |
+----+---------+--------+----------+------------------+
```

List topics:

```
$ kafta topic list
+-------------------------+------------+--------------------+
| NAME                    | PARTITIONS | REPLICATION FACTOR |
+-------------------------+------------+--------------------+
| my-topic                |         10 |                  3 |
| topic1                  |         10 |                  3 |
| topic2                  |          6 |                  3 |
+-------------------------+------------+--------------------+
```

## Consumer Group

List all consumers, run:

```
$ kafta consumer list
+------------------------+----------+--------+
| NAME                   | TYPE     | STATE  |
+------------------------+----------+--------+
| app1                   | consumer | Stable |
| app2                   | consumer | Stable |
| app3                   | consumer | Empty  |
+------------------------+----------+--------+
```

Lag of consumer:

```
$ kafta consumer lag app1
+---------------+----------+-----------+
| CONSUMER      | TOPIC    | TOTAL LAG |
+---------------+----------+-----------+
| app1          | my-topic |        41 |
+---------------+----------+-----------+
```

Sarama get Lag on topic *TODO*

Describe consumer, run:

```
$ kafta consumer describe app1
+--------------+----------+---------------+--------+--------------+
| ID           | PROTOCOL | PROTOCOL TYPE | STATE  | MEMBER COUNT |
+--------------+----------+---------------+--------+--------------+
| app1         | range    | consumer      | Stable |            4 |
+--------------+----------+---------------+--------+--------------+

+-----------------------+---------------+-----------------+------------+
| MEMBER ID             | MEMBER HOST   | TOPIC           | PARTITIONS |
+-----------------------+---------------+-----------------+------------+
| app1-service-1        | /100.10.22.21 | my-topic        | [0 1 2]    |
| app1-service-1        | /100.10.22.22 | my-topic        | [3 4 5]    |
| app1-service-1        | /100.10.22.23 | my-topic        | [6 7]      |
| app1-service-1        | /100.10.22.24 | my-topic        | [8 9]      |
+-----------------------+---------------+-----------------+------------+
```

To remove a consumer, it's simple. Just run:

```
$ kafta consumer delete app1
Consumer app1 deleted
```

## Cluster

Describe kafka cluster

```
$ kafta cluster describe
+----+-----------------------+------+------------+
| ID | ADDRESS               | RACK | CONTROLLER |
+----+-----------------------+------+------------+
|  0 | b-1.mydomain.com:9092 | 0    | false      |
|  1 | b-2.mydomain.com:9092 | 1    | false      |
|  2 | b-3.mydomain.com:9092 | 2    | false      |
|  3 | b-4.mydomain.com:9092 | 0    | false      |
|  4 | b-5.mydomain.com:9092 | 1    | false      |
|  5 | b-6.mydomain.com:9092 | 2    | true       |
+----+-----------------------+------+------------+

```

## Broker

Get broker config, run:

```
$ kafta broker get-configs 2
+-----------------------------------------+--------------------------------+---------+
| NAME                                    | VALUE                          | DEFAULT |
+-----------------------------------------+--------------------------------+---------+
| log.cleaner.delete.retention.ms         | 86400000                       | true    |
| log.message.timestamp.type              | CreateTime                     | true    |
| auto.create.topics.enable               | false                          | false   |
| log.cleanup.policy                      | delete                         | true    |
| log.cleaner.min.compaction.lag.ms       | 0                              | true    |
| message.max.bytes                       | 2097164                        | false   |
| default.replication.factor              | 3                              | false   |
| num.partitions                          | 6                              | false   |
+-----------------------------------------+--------------------------------+---------+
```

## Producer
Produce data for topic, run:
```
kafta console producer topic.test
> message test
```
You can produce a message with key, using the ":" between the key and message:
```
kafta console producer topic.test
> key1:message test
```
## Consumer
Consume data from topic, you can optionally enter the consumer group and the flag `--verbose` for debug the consumer:
```
kafta console consumer topic.test  [group=group.test] [--verbose]

2022/05/17 19:48:47 Initializing Consumer with group [group.test]...
2022/05/17 19:48:50 Consumer running, waiting for events...
2022/05/17 19:48:57 Partition: 0 Key: 1 Message: teste event
2022/05/17 19:49:05 Partition: 0 Key: w Message: {"json":"test"}
2022/05/17 19:50:05 Partition: 0 Key: 1 Message: event 2
```

## Subjects List
List subjects in SchemaRegistry

```
$ kafta schema subjects-list
+-------------------------------------------+
| NAME                                      |
+-------------------------------------------+
| topic1-dlq-value                          |
| topic1-value                              |
| topic2-dlq-value                          |
| topic2-value                              |
+-------------------------------------------+
```

# 💻 Code Contributors

<a href="https://github.com/electric-saw/kafta/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=electric-saw/kafta" />
</a>


# Next features

* Full support and management for schema-registry
* Full support and management for KSQL
* Tail data of topics (WIP)
* Producer data of topics (WIP)
* Change configs for topics (WIP)
* Dump/Restore data of topics
