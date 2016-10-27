#!/usr/bin/python3
import sys
import SyncMonitor
monitor = SyncMonitor.Monitor("/var/log/ldapsync.log", 300, 600)
monitor.init()
if monitor.Exitstatus == 0:
    print("Status OK, timediff is ",  (monitor.CurrentTimeEpoch - monitor.LogTimeEpoch), "seconds from Warning threshold.")
    sys.exit(monitor.Exitstatus)
elif monitor.Exitstatus == 1:
    print("Warning: Last timestamp updated in log file is at", monitor.timestamp, (monitor.CurrentTimeEpoch - monitor.LogTimeEpoch - monitor.warningThreshold), "seconds above threshold. Check if syncer is running.")
    sys.exit(monitor.Exitstatus)
else:
    print("Critical: Last timestamp updated in log file is at", monitor.timestamp, (monitor.CurrentTimeEpoch - monitor.LogTimeEpoch - monitor.errorThreshold), "seconds above threshold. Check if the syncer is running")
    sys.exit(monitor.Exitstatus)


