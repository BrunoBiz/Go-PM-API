- [Intro](#intro)
  - [Start](#start)
  - [Stop](#stop)
  - [Restart](#restart)
  - [Details](#details)


# Intro
This document presents the outcomes of LinuxGSM and GSM, in response to each of the main options provided by the API, in a way to compare both and find common semantic meaning between them.

The purpose of these tests if to streamline how the response in the API will be handled, while still being agnostic to which one of the managers is currently running the server.

The data collected here will represent if a manager is idempotent, if a request will result in an error or not, and the server state before and after it's sent.

Below is a list of all the tests that will be performed to collect the data for this document:

| **Initial state** | **Operation** | **LinuxGSM**  | **GSM**       |
|-------------------|---------------|---------------|---------------|
| on                | start         | actual output | actual output |
| off               | start         | actual output | actual output |
| on                | stop          | actual output | actual output |
| off               | stop          | actual output | actual output |
| on                | restart       | actual output | actual output |
| off               | restart       | actual output | actual output |
| on                | details       | actual output | actual output |
| off               | details       | actual output | actual output |

Whenever possible, this is the data that will be collected from each test:
- Exit code
- Stdout
- Stderr
- Server state after test
- Wheter the command succeded

## Start

### LinuxGSM <!-- omit from toc -->

### GSM <!-- omit from toc -->



## Stop

## Restart

## Details