**KONDUCT**

When running multiple and complex Kubernetes-based infrastructures you may want to test that everything is working fine at all times.

Monitoring solutions (such as Prometheus) provide a lot of functionalities to do that.

However, you may want to run some more complex tests or have a complete suite of acceptance tests before delivering Kubernetes-based platform to your customers.

**If that's the case** Konduct is designed exacly for you.

Konduct provides a CRUD+Exec interface for Kubernetes leveraging Custom Resources or simple YAML definitions.

Konduct can:

- run in 2 modes: continuos testing (controller mode) or one off testing (interactive mode).

- run actual tests simulating users interactions with a Kubernetes clusters.

- export tests reports to: prometheus, pushgateway, stdout and/or text files (json/yaml).

- allow you to interact with the running controller using well-defined REST APIs.

- work as performance testing framework for your Kubernetes clusters.


Here's a series of examples:

**SCENARIO**: ...
....
