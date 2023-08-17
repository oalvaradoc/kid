# Release v0.0.1 - 2021/02/01

# Features and Enhancements

## Init. commit

---

# Release v0.0.2 - 2021/02/07

# Features and Enhancements

## Logic optimization
Modify the version of dependency libraries

---

# Release v0.0.3 - 2021/02/10

# Features and Enhancements

## Logic optimization
1. executor logic optimization
2. Support local call (macro service function)
3. Modify and configure all configuration files to toml format (remove beego config)
4. Modify the name of the log method from the original *w* to *s* (s is the abbreviation of simple, because most of the components of MU do not print trace information when printing the log, it can be used for the log that does not need to print the trace. *s* method)
5. Increase client request retry&backoff capabilities (retry and backoff logic can be customized)
6. Unified error handling, failure return can customize whether to add debug stack information
7. Error return template can be defined 
8.Start callback logic optimization

---

# Release v0.0.4 - 2021/02/21

## Logic optimization
Modify the version of dependency libraries

---

# Release v0.0.5 - 2021/03/01

## Logic optimization
1. Compatible with new and old versions
   2.Sample starts to add ping to check whether the connection is normal
3. Change requestProps/responseProps to requestHeader/responseHeader
4. sed server can specify config path
5. Supports http methods such as get, delete, head, etc. except for registering post
6. Add response meta parameters for synchronous calls, including response headers and response original messages

---

# Release v0.0.6 - 2021/03/07

# Features and Enhancements
1.Increase response interceptor to encapsulate the return value
2.Add call interceptor (json/xml, the default is json) to automatically parse the response content of downstream services (note that the configuration file adds service.responseAutoParseKeyMapping)
3.Support custom responseTemplate, call interceptor, etc. at startup
4.Add the active delete transaction information mark in the downstream service (note that the application configuration file calls the downstream configuration modification, and the specific modification content can refer to the sample program)


---

# Release v0.0.7 - 2021/03/15

# Features and Enhancements
1. Merge xml/json client_receive interceptor
2. Set auto codec as the default group and decoder
3. Print server/client endpoint address, interceptors, codecs list when starting
4. Modify the auth module 
5.log adds IsEnable to determine whether the log level is enabled
6.Message adds the set and get methods of app properties and sets the read-write lock
7.Increase log printing when requesting downstream exceptions to facilitate troubleshooting

---

# Release v0.0.8 - 2021/03/17

## Logic optimization
Modify the version of dependency libraries

---

# Release v0.0.9 - 2021/03/22

## Logic optimization
1. Modify the GLS processing module and extract the public package
2. Add support for the old topic type (DTS)

---

# Release v0.0.10 - 2021/03/29

## Logic optimization
1. Optimize log processing and increase log printing

2. Regenerate the remote call mock file

3. Add the sample writing method of unit test for remote call

---

# Release v0.0.11 - 2021/04/06

## Logic optimization
1. New and old compatibility

2. Add response parsing rule processing for each request downstream

3. Increase audit log printing logic and corresponding interceptor

---

# Release v0.0.12 - 2021/04/12

## Logic optimization
1. Add all handlers can be set to enable validation by default. 10

2. Increase the skip response unpacking mark

3. Modify the downstream return logic, remove the additional topic information added during error, and print it in log mode

4. Modify the processing logic of the root transaction for the transaction communicator in the macro service mode

---

# Release v0.0.13 - 2021/04/19

## Logic optimization
1. Add kit executor log printing

2. When adding enable validation, you can customize the validation function

3. Delay loading client

4. Add config get method

---

# Release v0.0.14 - 2021/04/26

## Logic optimization
Added dxc, gls, remote mock files

---

# Release v0.0.15 - 2021/05/06

## Logic optimization
Set the status code to the response header when the service fails

---

# Release v0.0.16 - 2021/05/10

## Logic optimization
Modify the version of dependency libraries

---

# Release v0.0.17 - 2021/05/24

## Logic optimization
1. Add the function of separately setting the responseTemplate of the handler

2. Add the function of setting the parameter to the request header when opening the http endpoint

3. Add support for DXC Server/DXC SDK http mode

4. Increase the public method of struct field mask

5. Increase the judgment method of whether common/errors is a timeout error

---

# Release v0.0.18 - 2021/05/31

## Logic optimization
Fix the problem that the responseTemplate cannot be set separately for each handler

---

# Release v0.0.19 - 2021/06/24

## Logic optimization
1. DXC SDK adds storage and callback settings of attributes and app properties

2. Remove \u{0000} bytes that may be contained in app properties

3. Set separate error codes for different errors

---

# Release v0.0.20 - 2021/07/05

## Logic optimization
1. Add comments

2. Modify the default backoff processing logic, the backoff time of the first execution is 0

3. Add log printing of complete business parameters when dxc deserialization business parameters fail

4.dxc supports forced cancel global transaction

---

# Release v0.0.21 - 2021/07/12

## Logic optimization
1. Modify the auto codec to automatically determine the string type when the group is unpacked

2. Add the function of turning off handler interceptor according to the interceptor name

---

# Release v0.0.22 - 2021/08/02

## Logic optimization
1. Support macro service
   1.1 Increase the acquisition of DB connection pool
   1.2 Configure dynamic updates and automatically create a database connection pool based on the updated configuration
   1.3 Add APM interceptor
   1.4 Add GLS addressing
   1.5 Support the judgement of the communicator list
   1.6 Support Service config\transaction\log\heartbeat configuration hot update

2. Increase unit test

3. NewMockRemoteCallInc in the remote call mock file was renamed NewMockCallInc (modified the naming of the remote call interface according to the naming convention)

4. Support the setting of timeout messages into the dead queue

---

# Release v0.0.23 - 2021/08/11

## Logic optimization
1.Fix the problem of abnormal service config value during concurrency
2.Optimize the value processing logic for obtaining su

# Major Bug Fixes

---

# Release v0.0.24 - 2021/09/14

## Logic optimization
1.Support tcc->non tcc->tcc transaction propagation
2.Modify the version of dependency libraries

---
