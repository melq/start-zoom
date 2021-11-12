@echo off
set NAME=%1

schtasks -delete -tn startzoom_%NAME% -f