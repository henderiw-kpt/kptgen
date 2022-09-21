---
title: "`clusterrole`"
linkTitle: "clusterrole"
type: docs
description: >
  Add a clusterrole to the package.
---

<!--mdtogo:Short
    Add a pod to the package.
-->

`clusterrole` adds a clusterrole to the package.

### Synopsis

<!--mdtogo:Long-->

```
kptgen add clusterrole TARGET_DIR [flags]
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
$ kptgen add clusterrole ./blueprint/admin --fn-config cluster-role-fn-config
```

<!--mdtogo-->
