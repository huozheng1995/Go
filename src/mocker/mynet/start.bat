@echo off
echo Current directory: %CD%
cd /d %~dp0
echo Starting mynet.exe...
start mynet.exe
:: pause