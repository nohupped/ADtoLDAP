#!/usr/bin/python3
import sys
import os
import time
class Monitor:
    def __init__(self, path: str, warningThreshold: int, ErrorThreshold: int):
        self.path = path
        self.warningThreshold = warningThreshold
        self.errorThreshold = ErrorThreshold
        self.seekbuffer = 1

    def __check__thresholds(self):
        if self.warningThreshold > self.errorThreshold:
            print("Warning threshold must be less than error threshold")
            sys.exit(2)

    def __checklogexists(self):
        try:
            self.f = open(self.path)
            self.seekto = os.stat(self.path)[6]
        except:
            print("I/O error on file", self.path, ".Check if file exists and has right permission for the script to read")
            sys.exit(2)

    def __get__last__line(self):
        try:
            self.f.seek((self.seekto - self.seekbuffer), 0)
        except:
            print("Negative seek, empty log file. Possibly the check coincided when the log file was rotated, or the daemon is dead.")
            sys.exit(2)
        line = None
        counter = 0
        while 1:
            byte = self.f.read(1)
            if byte == "\n" and counter == 1:
                line = self.f.readlines()
                break
            elif byte == "\n":
                counter += 1
                self.seekbuffer += 1
                self.f.seek((self.seekto - self.seekbuffer), 0)
            else:
                self.seekbuffer += 1
                self.f.seek((self.seekto - self.seekbuffer), 0)
        self.lastline = line


    def __get__time(self):
        try:
            datetime = self.lastline[0].split(" ")
            self.timestamp = " ".join([datetime[6], datetime[7]])
            self.pattern = '%Y/%m/%d %H:%M:%S'
        except:
            print("Error finding time from the log line. This could be a bug. Please getback")
            sys.exit(2)

    def __get__epoch(self):
        self.LogTimeEpoch = int(time.mktime(time.strptime(self.timestamp, self.pattern)))
        self.CurrentTimeEpoch = int(time.time())

    def __generate_alert(self):
        timediff = self.CurrentTimeEpoch - self.LogTimeEpoch
        if timediff < self.warningThreshold:
            self.Exitstatus = 0
        elif timediff > self.warningThreshold and timediff < self.errorThreshold:
            self.Exitstatus = 1
        else:
            self.Exitstatus = 2

    def init(self):
        self.__check__thresholds()
        self.__checklogexists()
        self.__get__last__line()
        self.__get__time()
        self.__get__epoch()
        self.__generate_alert()


    def toString(self):
        self.init()
        attrs = vars(self)
        print(attrs.items())
