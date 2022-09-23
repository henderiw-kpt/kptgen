---
title: "`webhook`"
linkTitle: "webhook"
type: docs
description: >
  Add a webhook to the package.
---

<!--mdtogo:Short
    Add a webhook to the package.
-->

`webhook` adds a webhook to the package.

### Synopsis

<!--mdtogo:Long-->

```
kptgen apply webhook TARGET_DIR [flags]
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
# Add a pod with the characteristics in the fn-config to the package
$ kptgen apply webhook ./blueprint/admin --fn-config webhook-fn-config
```

<!--mdtogo-->
