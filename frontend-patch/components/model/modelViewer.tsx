import { useEffect, useRef } from "react";
import * as THREE from "three";
import { GLTFLoader } from "three/examples/jsm/loaders/GLTFLoader";
import { OrbitControls } from "three/examples/jsm/controls/OrbitControls";
import "./modelViewer.scss";

interface IModelViewerParams {
  modelUrl: string;
  width?: number;
  height?: number;
}

export default function ModelViewer(params: IModelViewerParams) {
  const containerRef = useRef<HTMLDivElement>(null);
  const width = params.width || 300;
  const height = params.height || 300;

  useEffect(() => {
    if (!containerRef.current || !params.modelUrl) return;

    // Scene
    const scene = new THREE.Scene();
    scene.background = new THREE.Color(0x1a1a2e);

    // Camera
    const camera = new THREE.PerspectiveCamera(75, width / height, 0.1, 1000);
    camera.position.set(0, 0, 3);

    // Renderer
    const renderer = new THREE.WebGLRenderer({ antialias: true });
    renderer.setSize(width, height);
    renderer.shadowMap.enabled = true;
    containerRef.current.innerHTML = "";
    containerRef.current.appendChild(renderer.domElement);

    // Lighting
    scene.add(new THREE.AmbientLight(0xffffff, 0.6));
    const dirLight = new THREE.DirectionalLight(0xffffff, 1);
    dirLight.position.set(5, 5, 5);
    scene.add(dirLight);

    // Controls
    const controls = new OrbitControls(camera, renderer.domElement);
    controls.enableDamping = true;
    controls.autoRotate = true;
    controls.autoRotateSpeed = 2.0;

    // Load model
    const loader = new GLTFLoader();
    loader.load(
      params.modelUrl,
      (gltf) => {
        const model = gltf.scene;
        const box = new THREE.Box3().setFromObject(model);
        const center = box.getCenter(new THREE.Vector3());
        const size = box.getSize(new THREE.Vector3());
        model.position.sub(center);
        const maxDim = Math.max(size.x, size.y, size.z);
        model.scale.setScalar(2 / maxDim);
        scene.add(model);
      },
      undefined,
      (error) => {
        console.error("Model load error:", error);
      }
    );

    // Animation loop
    let animId: number;
    const animate = () => {
      animId = requestAnimationFrame(animate);
      controls.update();
      renderer.render(scene, camera);
    };
    animate();

    // Cleanup
    return () => {
      cancelAnimationFrame(animId);
      renderer.dispose();
    };
  }, [params.modelUrl, width, height]);

  return (
    <div className="modelViewer" ref={containerRef} style={{ width, height }} />
  );
}
