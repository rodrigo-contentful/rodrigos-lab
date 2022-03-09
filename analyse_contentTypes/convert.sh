#!/bin/bash

echo "Will process a JSON content model"

SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" &> /dev/null && pwd)

cd $SCRIPT_DIR;

if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        # ...
        echo "gnu";
        ./bin/contentAnalyser_linux -d task/
elif [[ "$OSTYPE" == "darwin"* ]]; then
        # Mac OSX
        echo "mac";
        ./bin/contentAnalyser_mac -d task/
elif [[ "$OSTYPE" == "cygwin" ]]; then
        # POSIX compatibility layer and Linux environment emulation for Windows
        echo "win cywin";
        ./bin/contentAnalyser_win.exe -d task/
elif [[ "$OSTYPE" == "msys" ]]; then
        # Lightweight shell and GNU utilities compiled for Windows (part of MinGW)
        echo "win other";
        ./bin/contentAnalyser_win.exe -d task/
elif [[ "$OSTYPE" == "win32" ]]; then
        # I'm not sure this can happen.
        echo "win 32";
        ./bin/contentAnalyser_win.exe -d task/
elif [[ "$OSTYPE" == "freebsd"* ]]; then
        # ...
        echo 'freebds'
        ./bin/contentAnalyser_linux -d task/
else
        # Unknown.
         echo 'unknown'
fi

sleep 5