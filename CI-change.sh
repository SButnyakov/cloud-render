#!/bin/bash

sed -i '/react-scripts build/c \ \ \ \ "build": "CI=false DISABLE_ESLINT_PLUGIN=true react-scripts build",' ./frontend/package.json
