![language](https://img.shields.io/badge/language-go-5adaff)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

# ecsmgmt

Simple command line tool for AWS ECS (EC2 launch type) visibility.

## Get binaries

Get the latest binary for your OS from:
https://github.com/paweldudzinski/ecsmgmt/releases/tag/v1.0.1

If your OS is not listed - install as explained below.

## Install

You can build and install the official repository with [Go](https://golang.org/dl/):

`go get github.com/paweldudzinski/ecsmgmt/cmd`

This will checkout this repository into `$GOPATH/src/github.com/paweldudzinski/ecsmgmt/`, build, and install it.

It should then be available in $GOPATH/bin/ecsmgmt

## Usage

AWS credentials needs to be provided as:
* env variables (`AWS_SECRET_ACCESS_KEY`, `AWS_ACCESS_KEY_ID` and `AWS_DEFAULT_REGION`)
* command parameters (type `ecsmgmt --help` for details)

User/role identified by those credentials need to have proper ECS roles set up.

Use `ecsmgmt --help` to display available commands.

This tool is read only tool - to help you with visibility of your cluster/service/tasks/instances. You cannot change your setup (yet).

#### List clusters

Lists all clusters your user have access to.

`$ ecsmgmt list clusters`
```buildoutcfg
+---+-----------------------+--------+----------------+---------------+---------------+---------------------+
|   |      CLUSTER NAME     | STATUS | SERVICES COUNT | RUNNING TASKS | PENDING TASKS | CONTAINER INSTANCES |
+---+-----------------------+--------+----------------+---------------+---------------+---------------------+
| 1 | default               | ACTIVE | 0              | 0             | 0             | 0                   |
| 2 | my-custer             | ACTIVE | 2              | 62            | 0             | 21                  |
+---+-----------------------+--------+----------------+---------------+---------------+---------------------+
```

`CLUSER NAME` - the name of the cluster<br />
`STATUS` - one of ACTIVE, PROVISIONING, DEPROVISIONING, FAILED, INACTIVE<br />
`SERVICES COUNT` - number of services within a cluster<br />
`RUNNING TASKS` - sum of all running tasks within a cluster<br />
`PENDING TASKS` - sum of all pending tasks within a cluster<br />
`CONTAINER INSTANCES` - number of EC2 instances acting as a container instances in a cluster<br />


#### List services

Lists all services within a given cluster/

`$ ecsmgmt list services --cluster my-cluster`

```buildoutcfg
+---+----------------+--------+--------------------------------------+---------+---------+---------+
|   | SERVICE NAME   | STATUS |           TASK DEFINITION            | DESIRED | RUNNING | PENDING |
+---+----------------+--------+--------------------------------------+---------+---------+---------+
| 1 | svc-api        | ACTIVE | svc-api-task-definition:4            | 2       | 2       | 0       |
| 2 | svc-monitoring | ACTIVE | svc-monitoring-task-definition:63    | 60      | 60      | 0       |
+---+----------------+--------+--------------------------------------+---------+---------+---------+
```

`SERVICE NAME` - name of the service<br />
`STATUS` - one of ACTIVE, DRAINING, INACTIVE<br />
`TASK DEFINITION` - current active task definition name with revision number
`DESIRED` - desired tasks number<br />
`RUNNING` - running tasks number<br />
`PENDING` - desired tasks number<br />

#### List events

Lists service events (10 most recent) in a given cluster.

`$ ecsmgmt list events --cluster my-cluster --service svc-api`

```buildoutcfg
+---------------------+-----------------------------------------------------------------------+
|   DATE/TIME (UTC)   |                              EVENT                                    |
+---------------------+-----------------------------------------------------------------------+
| 2020-06-07 17:14:13 | (svc-api) has reached a steady state.                                 |
| 2020-06-07 11:13:54 | (svc-api) has reached a steady state.                                 |
| 2020-06-05 23:11:10 | (svc-api) has stopped 20 running tasks: (task a8905b25-6379-4949-9... |
| 2020-06-05 23:10:59 | (svc-api) has stopped 20 running tasks: (task d876177e-3b0a-4310-a... |
+---------------------+-----------------------------------------------------------------------+
```

`DATE/TIME (UTC)` - event occurrence date and time<br />
`EVENT` - the event<br />

#### List instances

Lists details of container instances (EC2s) that are assigned to a given cluster.

`$ ecsmgmt list instances --cluster my-cluster`

```buildoutcfg
+----+---------------------+--------+----------------+-----------------------+---------------+-------------+-------------+------------+
|    |     INSTANCE ID     | STATUS | EC2 PUBLIC IP  | WHEN REGISTERED (UTC) | RUNNING TASKS |     CPU     |   MEMORY    | CPU UTIL % |
+----+---------------------+--------+----------------+-----------------------+---------------+-------------+-------------+------------+
| 1  | i-f31a83ca839078c14 | ACTIVE | 184.169.197.23 | 2020-06-05 22:52:57   | 3             | 2048 / 2048 | 7975 / 4903 | 30.7       |
| 2  | i-37f378fcc451f192a | ACTIVE | 13.52.251.22   | 2020-04-09 00:14:47   | 3             | 2048 / 2048 | 7975 / 4903 | 1.523      |
| 3  | i-07b22ae731a41613b | ACTIVE | 18.144.62.1    | 2020-06-05 22:52:56   | 2             | 2048 / 2048 | 7975 / 5927 | 11.42      |
| 4  | i-0f8d50ad8880f36a1 | ACTIVE | 54.193.23.13   | 2020-01-27 14:08:54   | 3             | 2048 / 2048 | 7975 / 4903 | 11.21      |
| 5  | i-02e3550c8f76ce5c1 | ACTIVE | 54.215.251.2   | 2020-04-09 00:14:49   | 3             | 2048 / 2048 | 7975 / 5799 | 2.014      |
| 6  | i-02b8d80b1017bb6ad | ACTIVE | 13.57.179.33   | 2019-12-13 07:08:31   | 3             | 2048 / 2048 | 7975 / 4903 | 28.56      |
| 7  | i-0d4773eaf57e20317 | ACTIVE | 18.144.4.122   | 2020-06-05 22:52:55   | 3             | 2048 / 2048 | 7975 / 4903 | 1.025      |
| 8  | i-0c29ec126d7d37f94 | ACTIVE | 54.177.86.11   | 2020-04-09 00:14:49   | 3             | 2048 / 2048 | 7975 / 4903 | 35.94      |
+----+---------------------+--------+----------------+-----------------------+---------------+-------------+-------------+------------+
```

`INSTANCE ID` - AWS instance identifier<br />
`STATUS` -  REGISTERING, REGISTRATION_FAILED, ACTIVE, INACTIVE, DEREGISTERING, DRAINING<br />
`EC2 PUBLIC IP` - public instance IP<br />
`WHEN REGISTERED (UTC)` - when instance what registered<br />
`RUNNING TASKS` - how many tasks are running on a given instance<br />
`CPU` - how much CPU units are registered vs how much units are available<br />
`MEMORY` - how much memory units are registered vs how much units are available<br />
`CPU UTIL %` - average percent of CPU power being use at the moment of command execution<br />

#### List tasks

Lists registered tasks.

`$ ecsmgmt list tasks`

```buildoutcfg
+--------------------------------------+-----------------------------+--------+-----------------------------+
|              TASK NAME               |           FAMILY            | STATUS | CONTAINER NAME (CPU/MEMORY) |
+--------------------------------------+-----------------------------+--------+-----------------------------+
| svc-monitoring-task-definition:2     | svc-monitoring-task-family  | ACTIVE | monitoring (1/128)          |
| svc-monitoring-task-definition:4     | svc-monitoring-task-family  | ACTIVE | monitoring (auto/128)       |
| svc-api-task-definition:1            | svc-api-task-family         | ACTIVE | api (2/512)                 |
+--------------------------------------+-----------------------------+--------+-----------------------------+
```

`TASK NAME` - name of the task with revision<br />
`FAMILY` - task family<br />
`STATUS` - one of RUNNING, STOPPED<br />
`CONTAINER NAME (CPU/MEMORY)` - container name with declared allocaton of CPU/memory units
