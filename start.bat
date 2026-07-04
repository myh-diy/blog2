@echo off
taskkill /F /IM blog-server.exe >nul 2>&1
timeout /t 1 /nobreak >nul
cd /d "c:\Users\myh\Documents\blog2\backend"
start "My Blog" blog-server.exe
echo Server started! Open http://localhost:8080
pause
