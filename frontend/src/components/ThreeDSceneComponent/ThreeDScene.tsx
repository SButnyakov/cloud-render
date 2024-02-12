import { useEffect } from 'react';

import * as THREE from 'three';
// import { FBXLoader } from 'three/examples/jsm/loaders/FBXLoader';
// import { OBJLoader } from 'three/examples/jsm/loaders/OBJLoader';
// import { VOXLoader } from 'three/examples/jsm/loaders/VOXLoader';
import { GLTFLoader } from 'three/examples/jsm/loaders/GLTFLoader.js';

import SceneInit from './lib/SceneInit';

import styles from './ThreeDScene.module.css'

const ThreeDScene = ({glbFile, handleFileUploaded}: any) => {
  useEffect(() => {
    const test = new SceneInit('myThreeJsCanvas');
    test.initialize();
    test.animate();

    // const boxGeometry = new THREE.BoxGeometry(8, 8, 8);
    // const boxMaterial = new THREE.MeshNormalMaterial();
    // const boxMesh = new THREE.Mesh(boxGeometry, boxMaterial);
    // test.scene.add(boxMesh);
  
    let loadedModel: any;
    const glftLoader = new GLTFLoader();

    glftLoader.load(`${process.env.REACT_APP_API_CONVERT_BLEND}/files/${glbFile}`, (gltfScene) => {
      loadedModel = gltfScene;
      // console.log(loadedModel);

      gltfScene.scene.rotation.y = Math.PI / 8;
      gltfScene.scene.position.y = 3;
      gltfScene.scene.scale.set(10, 10, 10);
      test.scene?.add(gltfScene.scene);

      handleFileUploaded();
    });
  }, []);

  return (
    <canvas className={styles.scene} id="myThreeJsCanvas" />
  );
}

export default ThreeDScene;