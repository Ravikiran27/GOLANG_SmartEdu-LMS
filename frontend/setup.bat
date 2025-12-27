@echo off
echo Setting up LMS Frontend...

REM Check if Node.js is installed
where node >nul 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo Node.js is not installed. Please install Node.js 18+ first.
    exit /b 1
)

echo Node.js version:
node --version

REM Navigate to frontend directory
cd /d "%~dp0"

REM Install dependencies
echo Installing dependencies...
call npm install

REM Check if .env.local exists
if not exist .env.local (
    echo .env.local not found. Creating from template...
    copy .env.example .env.local
    echo Please edit .env.local with your Firebase credentials
)

echo.
echo Frontend setup complete!
echo.
echo Next steps:
echo 1. Edit .env.local with your Firebase credentials
echo 2. Run 'npm run dev' to start the development server
echo 3. Open http://localhost:3000 in your browser
pause
