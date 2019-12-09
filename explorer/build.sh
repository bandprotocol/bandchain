#!/bin/bash

echo "Builing for production..."
meteor build ../output/ --architecture os.linux.x86_64 --server-only

