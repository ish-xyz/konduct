**KubeTest**


When running multiple and custom Kubernetes clusters you may want to test that everything is working fine at all times.

Monitoring solutions (such as prometheus) provide a lot of exporters to do that. 

However, you may want to run some more complex tests or have a complete suite of acceptance tests before delivering a platform to your customers or end users.

**If that's the case KubeTest is exacly for you.**

KubeTest provides a CRUD+Exec interface for Kubernetes leveraging Custom Resources or simple YAML definitions.

KubeTest can:

* run in 2 modes: continuos testing (controller mode) or one off testing (interactive mode).

* export tests reports to: prometheus, pushgateway, stdout and/or text files (json/yaml).

* give you a birds eye view of your clusters.

* allow you to interact with the controller using well-defined REST APIs.

* work as performance testing framework for your Kubernetes clusters.


Here's a series of examples:

**SCENARIO**: ...
....


**TODO**:

* [...] create front docs and better examples
* [...] create validation functions for TestCase, Operations and Template types
* [...] improve operations logging
* [x] add label selector for "get"
* [x] add operations interval/startupWait
* [x] build command line
* [x] build loader for CRDs (TestCase and Template)
* [x] build exporter for prometheus (pushgateway)
* [x] improve k8s dynamic client 
* [x] implement exec operation
* [] fix stderr in dynamic client (exec)
* [x] define CRDs
* [] build exporter for prometheus (http endpoint)
* [] build API to trigger tests on demand
* [] build JSON exporter
* [] add execution time metric
* [] add stats for perf tests

* [] add parallel tests execution