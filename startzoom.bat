@echo off
set MEETINGID=%~1
set PASSWORD=%~2
start zoommtg:"//zoom.us/join?action=join&confno=%MEETINGID%&pwd=%PASSWORD%"