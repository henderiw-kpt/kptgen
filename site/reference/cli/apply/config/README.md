---
title: "`config`"
linkTitle: "config"
type: docs
description: >
  Add a set of config to the kpt package.
---

<!--mdtogo:Short
    Add a set of configuration specifed in the fn-config files to the package.
-->

`config` adds a set of configuration specifed in the fn-config files to the package.

### Synopsis

<!--mdtogo:Long-->

```
kptgen apply config TARGET_DIR [flags]
```

#### Args

```
TARGET_DIR:
  The target directory of the package
```

#### Flags

```
--fn-config-dir:
  Path to the fn-config dir containing `functionConfig` for the operation.
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
$ kptgen apply config ./blueprint/admin --fn-config-dir ./blueprint/admin/fn-config
```

<!--mdtogo-->
