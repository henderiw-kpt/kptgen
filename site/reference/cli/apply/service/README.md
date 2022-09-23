---
title: "`service`"
linkTitle: "service"
type: docs
description: >
  Add a service to the package.
---

<!--mdtogo:Short
    Add a service to the package.
-->

`service` adds a service to the package.

### Synopsis

<!--mdtogo:Long-->

```
kptgen apply service TARGET_DIR [flags]
```

#### Args

```
TARGET_DIR:
  The target directory of the package
```

#### Flags

```
--fn-config:
  Path to the file containing `functionConfig` for the operation.
```

<!--mdtogo-->

### Examples

{{% hide %}}

<!-- @makeWorkplace @verifyExamples-->

```
# Set up workspace for the test.
TEST_HOME=$(mktemp -d)
cd $TEST_HOME
```

{{% /hide %}}

<!--mdtogo:Examples-->

<!-- @ @verifyStaleExamples-->

```shell
# Add a service with the characteristics in the fn-config to the package
$ kptgen apply service ./blueprint/admin --fn-config service-fn-config
```

<!--mdtogo-->
