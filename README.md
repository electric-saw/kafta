![kafta logo](img/kafta.png)

Kafta is a modern non-JVM command line for managing Kafka Clusters written in golang.Usability was inspired by kubectl, the interface is simple and supports managing several clusters at the same time in a simple way.

# Table of Contents
- [Overview](#overview)
- [Concepts](#concepts)
- [Installing](#installing)
- [Configuration](#configuration)
- [Commands](#commands)
  * [Topic](#topic)
  * [Consumer Group](#consumergroup)
  * [Cluster](#cluster)
  * [Broker](#broker)
- [Next features](#nextfeatures)


# Overview

Kafkta providing a simple interface to manage topics, brokers, consumer-groups and many things like that. Interfaces similar to kubectl & go tools. 
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
project is to be simple to use, it is bad to need to install java, pass the kafka cluster address in every command, mess with xml's. 
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

    go get -u github.com/electric-saw/kafta

# Configuration

Kafta will create a config file in ~/.kafta/config. This yaml is used to support kafka multi-clusters and avoid passing all addresses every time.

To set up a new cluster, create a new config via kafta, you'll need to provide some information, don't worry, it's all in terminal, 
you don't need to edit any XML \o/

Follow the example:

```
$ kafta config set-context production
Bootstrap servers: b-1.mydomain:9092,b-2.mydomain:9092,b-3.mydomain:9092
Schema registry: https://schema-registry-schema-registry.com
Use SASL: y
SASL Algorithm: sha512
User: myuser
✔ Password: ******
```

To list the contexts, run:

```
$ kafta config get-contexts

╭─────────┬────────────────┬───────────────────┬─────────────────────────────────────────────┬──────╮
│ CURRENT │ NAME           │ CLUSTER           │ SCHEMA REGISTRY                             │ KSQL │
├─────────┼────────────────┼───────────────────┼─────────────────────────────────────────────┼──────┤
│         │ production     │ b-1.mydomain:9092 │ https://schema-registry-schema-registry.com │      │
│         │                │ b-2.mydomain:9092 │                                             │      │
│         │                │ b-3.mydomain:9092 │                                             │      │
├─────────┼────────────────┼─────────────────────────────────────────────────────────────────┼──────┤
│ *       │ dev            │ b-1.mydomain:9092 │ https://schema-registry-schema-registry.com │      │
│         │                │ b-3.mydomain:9092 │                                             │      │
│         │                │ b-4.mydomain:9092 │                                             │      │
╰─────────┴────────────────┴───────────────────┴─────────────────────────────────────────────┴──────╯

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
INTERNAL   NAME       PARTITIONS
           my-topic   10


ID         ISR        LEADER   REPLICAS   OFFLINE REPLICAS
0          [1 4 3]    1        [1 4 3]    []
1          [3 1 2]    3        [3 1 2]    []
2          [2 3 4]    2        [2 3 4]    []
3          [4 2 1]    4        [4 2 1]    []
4          [1 3 2]    1        [1 3 2]    []
5          [3 2 4]    3        [3 2 4]    []
6          [2 4 1]    2        [2 4 1]    []
7          [4 1 3]    4        [4 1 3]    []
8          [1 4 3]    1        [1 4 3]    []
9          [3 1 2]    3        [3 1 2]    []
```

List topics:

```
$ kafta topic list 

╭─────────┬────────────┬────╮
│ NAME    │ PARTITIONS │ RF │
├─────────┼────────────┼────┤
│ topic1  │         20 │  3 │
├─────────┼────────────┼────┤
│ topic2  │         10 │  3 │
├─────────┼────────────┼────┤
│ topic3  │          1 │  3 │
╰─────────┴────────────┴────╯
```

## Consumer Group

List all consumers, run:

```
$ kafta consumer list
╭──────────────────┬──────────┬────────╮
│ NAME             │ TYPE     │ STATE  │
├──────────────────┼──────────┼────────┤
│ app1             │ consumer │ Stable │
├──────────────────┼──────────┼────────┤
│ app2             │ consumer │ Empty  │
├──────────────────┼──────────┼────────┤
│ schema-registry  │ sr       │ Stable │
╰──────────────────┴──────────┴────────╯
```

Lag of consumer:

```
$ kafta consumer lag app1
CONSUMER   TOPIC      TOTAL LAG
app1       my-topic   41
```

Sarama get Lag on topic *TODO*

Describe consumer, run:

```
$ kafta consumer describe app1
ID     PROTOCOL   PROTOCOL TYPE   STATE    MEMBER COUNT
app1   range      consumer        Stable   1


MEMBER ID   MEMBER HOST      TOPICS
app1        /10.111.22.203   my-topic   map[my-topic:[0 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15]]
```

To remove a consumer, it's simple. Just run:

```
$ kafta consumer delete app1
Consumer app1 deleted
```

## Cluster

Describe kafka cluster

```
$ ./bin/kafta cluster describe
╭────┬────────────────────────────────────┬────────────╮
│ ID │ ADDR                               │ CONTROLLER │
├────┼────────────────────────────────────┼────────────┤
│  1 │ b-1.mydomain.com:9092              │            │
├────┼────────────────────────────────────┼────────────┤
│  2 │ b-2.mydomain.com:9092              │            │
├────┼────────────────────────────────────┼────────────┤
│  3 │ b-4.mydomain.com:9092              │ *          │
╰────┴────────────────────────────────────┴────────────╯
```

## Broker

Describe broker, run:

```
$ ./bin/kafta broker describe 2
╭────────────────┬────────────────────────────────────────────────────────────╮
│ BROKER DETAILS │                                                            │
├────────────────┼────────────────────────────────────────────────────────────┤
│ Id             │ 2                                                          │
├────────────────┼────────────────────────────────────────────────────────────┤
│ Host           │ b-2.mydomain.com                                           │
├────────────────┼────────────────────────────────────────────────────────────┤
│ Address        │ b-2.mydomain.com:9092                                      │
├────────────────┼────────────────────────────────────────────────────────────┤
│ Connected      │ false                                                      │
├────────────────┼────────────────────────────────────────────────────────────┤
│ Rack           │ sae1-az1                                                   │
├────────────────┼────────────────────────────────────────────────────────────┤
│ IsController   │ false                                                      │
╰────────────────┴────────────────────────────────────────────────────────────╯
```

Get broker config, run:

```
[lucas.viecelli@fedora kafta]$ ./bin/kafta broker get-configs 2
╭────────────────────────────────────┬─────────────────────────────────────┬─────────╮
│ NAME                               │ VALUE                               │ DEFAULT │
├────────────────────────────────────┼─────────────────────────────────────┼─────────┤
│ log.cleaner.min.compaction.lag.ms  │ 0                                   │ true    │
├────────────────────────────────────┼─────────────────────────────────────┼─────────┤
│ offsets.topic.num.partitions       │ 50                                  │ true    │
├────────────────────────────────────┼─────────────────────────────────────┼─────────┤
│ log.flush.interval.messages        │ 9223372036854775807                 │ true    │
├────────────────────────────────────┼─────────────────────────────────────┼─────────┤
│ controller.socket.timeout.ms       │ 30000                               │ true    │
├────────────────────────────────────┼─────────────────────────────────────┼─────────┤
│ ....                                                                               │
╰────────────────────────────────────────────────────────────────────────────────────╯
```

# Next features

* Full support and management for schema-registry
* Full support and management for KSQL
* Tail data of topics
* Producer data of topics
* Change configs for topics (WIP)
* Dump/Restore data of topics
