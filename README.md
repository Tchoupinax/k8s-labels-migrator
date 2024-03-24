# kubernetes-labels-migrator

## Usage

If you need help, the help flag can show you available flags (`kubernetes-labels-migrator -help`).

```bash
kubernetes-labels-migrator \
      -deployment="my-application" \
      -label="kubernetes.io/app" \
      -namespace="default" \
      -value="my-application"
```

![Screenshot of a terminal with all the execution](.github/docs/demo.png "CLI demo")

## Motivation

Deployment's labels are immutable. That means, if you have to edit, rename or add labels to a deployed application, you have to delete the deployment. On a production application, it is very annoying.

To be able to add labels without any downtime, we have to use a [blue-green strategy](https://www.redhat.com/en/topics/devops/what-is-blue-green-deployment). This projects aims to do this operation automatically and without human manual operation. The production must be safe!

## Explanation

TODO

## Zero downtime testing

Using [Vegeta](https://github.com/tsenart/vegeta), we can overload the service to detect if one or more request is thrown during the process.

⚠️ During the process, pods and service will be updated, that means if you expect to port-forward the service, the connection will be lost during the migration. We suggest you to target an `Ingress` or any others way to reach the service.

### Monotoring from your terminal

#### Requirements

The following command requires two projects to work:

- [jplot](https://github.com/rs/jplot)
- [jaggr](https://github.com/rs/jaggr)

To install these project, you can enter these commands:

```bash
# MacOS
brew install rs/tap/jplot
brew install rs/tap/jaggr
```

Then, enter this command

```
echo 'GET YOUR_URL' | \
    vegeta attack -rate 200 -duration 10m | vegeta encode | \
    jaggr @count=rps \
          hist\[100,200,300,400,500\]:code \
          p25,p50,p95:latency \
          sum:bytes_in \
          sum:bytes_out | \
    jplot rps+code.hist.100+code.hist.200+code.hist.300+code.hist.400+code.hist.500 \
          latency.p95+latency.p50+latency.p25 \
          bytes_in.sum+bytes_out.sum
```
