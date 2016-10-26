#!/usr/bin/python3
import SyncMonitor
monitor = SyncMonitor.Monitor("/var/log/ldapsync.log", 30, 40)
monitor.init()


