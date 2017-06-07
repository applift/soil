---
title: Interpolation
layout: default

weight: 20
---

# Interpolation

Soil agent interpolates variables in pod constraints, `unit->source` and `blob->source` declared as `${some.variable}`. If variable is not defined Soil agent leaves it unchanged.

```hcl
pod "my-pod" {
  constraint {
    "${meta.rack}" = "rack-1"
  }

  unit "unit-1" {
    source = <<EOF
    # ${meta.rack}
    [Service]
    EnvironmentFile=/etc/test
    ...
    EOF
  }
  blob "/etc/test" {
    source = <<EOF
    RACK=${meta.rack}
    EOF
  }
}
```

`meta` variables can be declared in [Agent configuration]({{site.baseurl}}/agent/configuration).

If pod contains one or more BLOBs their hashes will be available as `${blob.<blob-id>}`. There `blob-id` is escaped path. For example blob with path `/etc/my/blob.env` hash will be available in units as `${blob.etc-my-blob.env}`.

## Agent variables

Agent variables are accessible as `${agent.*}`:

`id` Agent ID.

`drain` Agent drain state.

`pod_exec` Pod unit "Exec*".