# Disk Partition Monitoring

## Overview

Beszel can now automatically monitor all disk partitions on your system, not just the root partition and manually specified extra filesystems. This provides comprehensive visibility into disk usage across your entire system.

## Agent Configuration

### MONITOR_ALL_DISKS Environment Variable

Set `MONITOR_ALL_DISKS=true` to automatically discover and monitor all physical disk partitions.

**Example (Docker Compose):**
```yaml
services:
  beszel-agent:
    image: henrygd/beszel-agent
    environment:
      - MONITOR_ALL_DISKS=true
```

**Example (Systemd):**
```bash
# Add to /etc/beszel-agent/env
MONITOR_ALL_DISKS=true
```

### What Gets Monitored

When `MONITOR_ALL_DISKS` is enabled, the agent will automatically discover and track:

- All physical disk partitions (ext4, xfs, btrfs, ntfs, etc.)
- Root filesystem and any additional mount points
- Network-attached storage mounted as local filesystems

### What Gets Excluded

To avoid noise and unnecessary resource usage, the following are automatically excluded:

- Virtual filesystems (tmpfs, devtmpfs, procfs, sysfs)
- Container overlays (overlay, squashfs)
- Pseudo filesystems (cgroup, debugfs, configfs)
- Temporary mount points (/dev, /proc, /sys, /run, /tmp)

### Using EXTRA_FILESYSTEMS (Manual Configuration)

If you don't want to monitor all partitions, you can continue using `EXTRA_FILESYSTEMS` to explicitly list the partitions to monitor:

```yaml
environment:
  - EXTRA_FILESYSTEMS=/dev/sdb1,/dev/sdc1,/mnt/storage
```

You can also specify custom names for partitions using the `__` separator:

```yaml
environment:
  - EXTRA_FILESYSTEMS=/dev/sdb1__backup-drive,/dev/sdc1__media-storage
```

## Alert Configuration

### Excluding Partitions from Disk Alerts

When you have disk alerts enabled, you may want to exclude certain partitions from triggering alerts. For example, you might want to exclude a backup drive or a temporary storage partition.

**To exclude partitions:**

1. Navigate to the system's alerts configuration
2. Enable the "Disk" alert
3. Under "Excluded partitions", check the partitions you want to exclude
4. The disk alert will now only monitor the non-excluded partitions

**Example Use Cases:**

- Exclude `/dev/sdb1__backup` - Don't alert on backup drives that are expected to fill up
- Exclude `tmpfs` mounts - Temporary storage that fills and empties frequently
- Exclude read-only volumes - Mounted ISOs or read-only media

### Alert Behavior

- By default, disk alerts check **all** monitored partitions (root + extra filesystems)
- If **any** non-excluded partition exceeds the threshold, the alert triggers
- The alert message will indicate which partition triggered it
- You can configure different alert thresholds for different systems

## Web UI

### Systems Table

The systems table shows:
- Root disk usage percentage
- Colored indicators showing the highest usage state among extra partitions (green/yellow/red)
- Hover tooltip with detailed breakdown of all partitions

### System Detail Page

Each monitored partition gets its own section showing:
- **Usage Chart**: Disk usage over time
- **I/O Chart**: Read/write throughput over time
- Partition name and current usage percentage

Partitions are displayed in the order they are discovered, with the root partition first.

## Best Practices

1. **Start with MONITOR_ALL_DISKS** - Let the agent discover all partitions automatically
2. **Review discovered partitions** - Check the web UI to see what was discovered
3. **Configure alert exclusions** - Exclude partitions that shouldn't trigger alerts
4. **Use custom names** - For manual configuration with EXTRA_FILESYSTEMS, use the `__` separator to provide friendly names

## Troubleshooting

### Partition not showing up

- Ensure the partition is mounted
- Check that it's not in the exclusion list (tmpfs, overlay, etc.)
- On Linux, verify the device is under `/dev`
- Check agent logs for "Auto-discovered partition" messages

### Too many partitions monitored

- If `MONITOR_ALL_DISKS` discovers too many partitions, switch to manual configuration with `EXTRA_FILESYSTEMS`
- List only the partitions you want to monitor

### Partition name unclear in UI

- Use the custom naming feature in `EXTRA_FILESYSTEMS`:
  ```
  EXTRA_FILESYSTEMS=/dev/sdb1__my-backup,/dev/sdc1__media-storage
  ```

## Migration from Previous Versions

If you were previously using `EXTRA_FILESYSTEMS`:

1. Your configuration will continue to work as before
2. Optionally, enable `MONITOR_ALL_DISKS=true` to automatically discover all partitions
3. Remove `EXTRA_FILESYSTEMS` if you want automatic discovery for all partitions
4. Or keep both - `EXTRA_FILESYSTEMS` entries will be included along with auto-discovered partitions
