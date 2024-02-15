#!/bin/zsh
source ~/.zsh_profile

if [ "$#" = 1 ]
then
	blender -b -P /Users/vladlavrov/Documents/Study/convert_server/utils/scripts/2gltf2-master/2gltf2.py -- "$1"
else
	echo To glTF 2.0 converter.
	echo Supported file formats: .abc .blend .dae .fbx. .obj .ply .stl .usd .wrl .x3d
	echo  
	echo 2gltf2.sh [filename]
fi
