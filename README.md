# Reconciler

## Overview

>**CAUTION:** This repository is in a very early stage. Use it at your own risk.

The Reconciler is a central system to reconcile Kyma clusters.

## Run Reconciler locally

Follow these steps to run Reconciler locally:

1. Build the Docker image:

```
docker build -f Dockerfile -t reconciler:v1 .
```


2. Run the Docker container:

```
docker run --name reconciler -it -p 8080:8080 reconciler:v1 reconciler service start
```


## Testing

The reconciler unit tests include also expensive test suites. Expensive means that the test execution might do the following:

* take an unusual amount of time (e.g. >1 min)
* generate a big amount of network traffic to remote systems (e.g. >100MB)
* allocates during the execution many disc space (e.g. > 1GB)

By default, expensive test suites are disabled. To enable them, before you execute the test suits, apply one of the following options:

* Set the environment variable `RECONCILER_EXPENSIVE_TESTS=true`
* In the GO code, execute the function `test.EnableExpensiveTests()`

## Adding a new component reconciler

If a custom logic must be executed before, during, or after the reconciliation of a component, component reconcilers are required.

The reconciler supports component reconcilers, which handle component-specific reconciliation runs.

To add another component reconciler, execute following steps:

1. Create a component reconciler by executing the script `pkg/reconciler/instances/new-reconciler.sh`.
   Provide the name of the component as parameter, for example:
   component as parameter.
   
   
   ```bash
   pkg/reconciler/instances/new-reconciler.sh istio

    The script creates a new package including the boilerplate code required to initialize a new component reconciler instance during runtime.
   new component reconciler instance during runtime.
   
 2. Edit the files inside the package:
   
      - Edit the file `action.go` and encapsulate your custom reconciliation logic in `Action` structs.

     - Edit the `$componentName.go` file:
            - Use the `WithDependencies()` method to list the components that are required before this reconciler can run.
            - Use the `WithPreReconcileAction()`, `WithReconcileAction()`, `WithPostReconcileAction()` to inject custom `Action` instances into the reconciliation process.
               
3. Re-build the CLI to add the new component reconciler to the `reconciler start` command.
   The `reconciler start` command is a convenient way to run a component reconciler as standalone server.

    Example:

        # Build CLI
        cd $GOPATH/src/github.com/kyma-incubator/reconciler/
        make build
        
        # Start the component reconciler (for example, 'istio') as standalone service
        ./bin/reconciler-darwin reconciler start istio
        
        # To get a list of all configuration options for the component reconciler, call: 
        ./bin/reconciler-darwin reconciler start istio --help
