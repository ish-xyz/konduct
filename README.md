Controller to run integration tests on Kubernetes


**TODO**:

* [x] build command line
* [x] build loader for CRDs (TestCase and Template)
* [...] create validation functions for TestCase, Operations and Template types
* [] build exporter for prometheus (pushgateway)
* [x] improve k8s dynamic client 
* [x] implement exec operation
* [] fix stderr in dynamic client (exec)
* [x] define CRDs
* [] build exporter for prometheus (http endpoint)
* [] build API to trigger tests on demand
* [] build JSON exporter
* [] create front docs and better examples
* [] add execution time metric
* [] add stats for perf tests
* [] improve operations logging
* [] add label selector for "get"
* [] add operations interval/startupWait
* [] add parallel tests execution