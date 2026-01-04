import React, { /*useState*/ useEffect, /*useRef*/ } from 'react';
import * as THREE from 'three';
import { GLTFLoader } from 'three/examples/jsm/loaders/GLTFLoader.js';
import gsap from 'gsap';
import { ScrollTrigger } from 'gsap/ScrollTrigger';
import Lenis from 'lenis';
import Navbar from '../components/Navbar';
import RotatingText from '../components/RotatingText';

gsap.registerPlugin(ScrollTrigger);

const Landing: React.FC = () => {

  useEffect(() => {
    const lenis = new Lenis();
    (window as any).lenis = lenis; // Store lenis globally for access
    
    function raf(time: number) {
      lenis.raf(time);
      requestAnimationFrame(raf);
    }
    requestAnimationFrame(raf);

    const heroSection = document.querySelector('.hero-section') as HTMLElement;
    if (!heroSection) return;

    const canvas = document.createElement('canvas');
    canvas.style.position = 'fixed';
    canvas.style.top = '0';
    canvas.style.left = '0';
    canvas.style.width = '100vw';
    canvas.style.height = '100vh';
    canvas.style.pointerEvents = 'none';
    canvas.style.zIndex = '2';
    heroSection.appendChild(canvas);

    const scene = new THREE.Scene();
    const camera = new THREE.PerspectiveCamera(50, window.innerWidth / window.innerHeight, 0.1, 1000);
    camera.position.set(2, 2, 8);

    const renderer = new THREE.WebGLRenderer({ canvas, antialias: true, alpha: true });
    renderer.setSize(window.innerWidth, window.innerHeight);
    renderer.setClearColor(0x000000, 0);

    const ambientLight = new THREE.AmbientLight(0xffffff, 1);
    scene.add(ambientLight);

    const directionalLight = new THREE.DirectionalLight(0xffffff, 1);
    directionalLight.position.copy(camera.position);
    directionalLight.target.position.set(6, 0, 0);
    scene.add(directionalLight);
    scene.add(directionalLight.target);


    const textureLoader = new THREE.TextureLoader();
    const bayerTexture = textureLoader.load('/BayerDithering.png');
    bayerTexture.magFilter = THREE.NearestFilter;
    bayerTexture.minFilter = THREE.NearestFilter;
    bayerTexture.wrapS = THREE.RepeatWrapping;
    bayerTexture.wrapT = THREE.RepeatWrapping;

    let model: THREE.Group;

    const loader = new GLTFLoader();
    loader.load('/Folder.glb', (gltf: any) => {
      model = gltf.scene;
      model.scale.set(0.024, 0.024, 0.024);
      model.position.set(5.7, 0, 0);
      model.rotateY(0.15);
      model.rotateX(0.20);

      // Apply dither shader to all meshes
      model.traverse((child: any) => {
        if (child instanceof THREE.Mesh) {
          const originalMaterial = child.material as THREE.Material;
          child.material = new THREE.ShaderMaterial({
            uniforms: {
              bayerTexture: { value: bayerTexture },
              ditherScale: { value: 40 }
            },
            vertexShader: `
              varying vec2 vUv;
              varying vec3 vNormal;
              void main() {
                vUv = uv;
                vNormal = normalize(normalMatrix * normal);
                gl_Position = projectionMatrix * modelViewMatrix * vec4(position, 1.0);
              }
            `,
            fragmentShader: `
              uniform sampler2D bayerTexture;
              uniform float ditherScale;
              varying vec2 vUv;
              varying vec3 vNormal;
              
              void main() {
                vec3 lightDir = normalize(vec3(5.0, 5.0, 5.0));
                float diff = max(dot(vNormal, lightDir), 0.0);
                vec3 baseColor = vec3(${originalMaterial instanceof THREE.MeshStandardMaterial && originalMaterial.color ? 
                  `${originalMaterial.color.r}, ${originalMaterial.color.g}, ${originalMaterial.color.b}` : '0.8, 0.6, 0.2'});
                vec3 litColor = baseColor * (0.45 + 0.5 * diff);
                
                vec2 ditherUV = gl_FragCoord.xy / ditherScale;
                vec3 ditherPattern = texture2D(bayerTexture, ditherUV).rgb;
                float ditherValue = (ditherPattern.r - 1.0) * 0.2;
                
                vec3 finalColor = litColor * (1.0 + ditherValue);
                
                gl_FragColor = vec4(finalColor, 1);
              }
            `
          });
        }
      });

      scene.add(model);
      renderer.render(scene, camera);

      // GSAP ScrollTrigger Animation

      gsap.to(model.rotation, {
        x: -0.1,
        y: 3.0,
        z: 0,
        ease: 'power1.inOut',
        scrollTrigger: {
          trigger: '.image-conversions',
          start: 'top 95%',
          end: '10% 10%',
          scrub: true,
          onUpdate: () => renderer.render(scene, camera)
        }
      });

      gsap.to(model.position, {
        x: -2.5,
        y: 1,
        ease: 'power2.inOut',
        scrollTrigger: {
          trigger: '.image-conversions',
          start: 'top 95%',
          end: '10% 10%',
          scrub: true,
          onUpdate: () => renderer.render(scene, camera)
        }
      });

      gsap.to(model.position, {
        y: 2,
        ease: 'power1.inOut',
        scrollTrigger: {
          trigger: '.archive-conversions',
          start: '30% center',
          end: '70% 40%',
          scrub: true,
          onUpdate: () => renderer.render(scene, camera)
        }
      });

      // Z-axis loop animation (runs independently)
      gsap.to(model.position, {
        z: 0.15,
        duration: 1,
        ease: 'sine.inOut',
        yoyo: true,
        repeat: -1,
        onUpdate: () => renderer.render(scene, camera)
      });
    });

    // Scroll progress indicator animation
    const scrollProgress = document.querySelector('.scroll-progress') as HTMLElement;
    const scrollDots = document.querySelectorAll('.scroll-dot');
    
    ScrollTrigger.create({
      trigger: '.image-conversions',
      start: 'top center',
      onEnter: () => gsap.to(scrollProgress, { opacity: 1, duration: 0.5 })
    });

    // GIF animations for all sections
    const gifConfigs = [
      { selector: '.conversion-gif', trigger: '.image-conversions' },
      { selector: '.document-gif', trigger: '.document-conversions' },
      { selector: '.audio-gif', trigger: '.audio-conversions' },
      { selector: '.video-gif', trigger: '.video-conversions' },
      { selector: '.archive-gif', trigger: '.archive-conversions' }
    ];

    gifConfigs.forEach(({ selector, trigger }) => {
      const gifElement = document.querySelector(selector) as HTMLElement;
      if (gifElement) {
        gsap.fromTo(gifElement,
          { x: '300%', y: '-80%', opacity: 0 },
          {
            x: '400%',
            y: '0%',
            opacity: 1,
            ease: 'circ.inOut',
            scrollTrigger: {
              trigger,
              start: '20% 40%',
              end: '30% 40%',
              scrub: 1
            }
          }
        );
      }
    });

    // Animate dots based on scroll position
    const sections = ['image', 'document', 'audio', 'video', 'archive'];
    sections.forEach((section, index) => {
      ScrollTrigger.create({
        trigger: `.${section}-conversions`,
        start: 'top center',
        end: 'bottom center',
        onEnter: () => {
          (scrollDots[index] as HTMLElement).style.backgroundColor = '#ffb947';
          (scrollDots[index] as HTMLElement).style.borderColor = '#ffb947';
          (scrollDots[index] as HTMLElement).style.transform = 'scale(1.3)';
        },
        onLeave: () => {
          (scrollDots[index] as HTMLElement).style.backgroundColor = 'rgba(255, 255, 255, 0.3)';
          (scrollDots[index] as HTMLElement).style.borderColor = 'rgba(255, 255, 255, 0.5)';
          (scrollDots[index] as HTMLElement).style.transform = 'scale(1)';
        },
        onEnterBack: () => {
          (scrollDots[index] as HTMLElement).style.backgroundColor = '#ffb947';
          (scrollDots[index] as HTMLElement).style.borderColor = '#ffb947';
          (scrollDots[index] as HTMLElement).style.transform = 'scale(1.3)';
        },
        onLeaveBack: () => {
          (scrollDots[index] as HTMLElement).style.backgroundColor = 'rgba(255, 255, 255, 0.3)';
          (scrollDots[index] as HTMLElement).style.borderColor = 'rgba(255, 255, 255, 0.5)';
          (scrollDots[index] as HTMLElement).style.transform = 'scale(1)';
        }
      });
    });

    // Format buttons sequential animation
    sections.forEach(section => {
      const sectionButtons = document.querySelectorAll(`.${section}-conversions .format-btn`);
      if (sectionButtons.length > 0) {
        const tl = gsap.timeline({ repeat: -1 });
        sectionButtons.forEach((btn) => {
          tl.to(btn, {
            backgroundColor: '#ffb947',
            borderColor: '#ffb947',
            color: '#161B27',
            duration: 0.3,
            ease: 'circ.inOut'
          }).to(btn, {
            backgroundColor: 'rgba(255, 255, 255, 0.1)',
            borderColor: 'rgba(255, 255, 255, 0.2)',
            color: '#fff',
            duration: 0.3,
            delay: 0.7,
            ease: 'power2.in'
          });
        });
      }
    });

    return () => {
      lenis.destroy();
      ScrollTrigger.getAll().forEach(trigger => trigger.kill());
      heroSection.removeChild(canvas);
      renderer.dispose();
    };
  }, []);

  const handleGetStarted = () => {
    // Navigate to convert page or scroll to conversion sections
    window.location.href = '/convert';
  };

  return (
    <div className="min-h-screen relative">
      <div className="fixed top-0 left-0 right-0 z-20 text-center py-3 px-4" style={{
        background: 'rgba(255, 185, 71, 0.95)',
        color: '#161B27',
        fontSize: 'clamp(0.75rem, 2vw, 0.875rem)',
        fontWeight: '500'
      }}>
        ⚠️ Backend currently unavailable due to expired Render Database free access (Dec 14, 2024)
      </div>
      
      <div className="min-h-screen" style={{ overflow: 'auto', scrollbarWidth: 'none', msOverflowStyle: 'none', paddingTop: '4rem' }}>
      <style>{`
        html, body, * {
          scrollbar-width: none;
          -ms-overflow-style: none;
        }
        html::-webkit-scrollbar,
        body::-webkit-scrollbar,
        *::-webkit-scrollbar {
          display: none;
          width: 0;
          height: 0;
        }
      `}</style>
      
    
      <div className="scroll-progress" style={{
        position: 'fixed',
        right: '2rem',
        top: '50%',
        transform: 'translateY(-50%)',
        zIndex: 100,
        display: 'flex',
        flexDirection: 'column',
        gap: '1.5rem',
        opacity: 0
      }}>
        {['Image', 'Document', 'Audio', 'Video', 'Archive'].map((section, index) => (
          <div
            key={section}
            className={`scroll-dot scroll-dot-${index}`}
            style={{
              width: '9px',
              height: '9px',
              borderRadius: '40%',
              backgroundColor: 'rgba(255, 255, 255, 0.3)',
              border: '2px solid rgba(255, 255, 255, 0.5)',
              cursor: 'pointer',
              transition: 'all 0.3s ease',
              position: 'relative'
            }}
            onClick={() => {
              const element = document.querySelector(`.${section.toLowerCase()}-conversions`);
              if (element && (window as any).lenis) {
                (window as any).lenis.scrollTo(element, { duration: 1.5, easing: (t: number) => Math.min(1, 1.001 - Math.pow(2, -10 * t)) });
              }
            }}
          >
            <span style={{
              position: 'absolute',
              right: '20px',
              top: '50%',
              transform: 'translateY(-50%)',
              whiteSpace: 'nowrap',
              fontSize: '0.65625rem',
              color: '#fff',
              opacity: 0,
              transition: 'opacity 0.3s ease',
              pointerEvents: 'none'
            }} className="dot-label">{section}</span>
          </div>
        ))}
      </div>
      <div className="halftone-bg" style={{ 
        position: 'fixed', 
        top: 0, 
        left: 0, 
        width: '100%', 
        height: '100%', 
        zIndex: 0
      }}>
        <div className="halftone-noise" />
        <div style={{
          position: 'absolute',
          top: 0,
          left: 0,
          width: '100%',
          height: '100%',
          backgroundImage: 'url(/BayerDithering.png)',
          backgroundRepeat: 'repeat',
          backgroundSize: '20px 20px',
          opacity: 0.25,
          mixBlendMode: 'overlay'
        }} />
      </div>

      <section className="hero-section" style={{ height: '100vh', position: 'relative', overflow: 'hidden' }}>
        <Navbar/>
        <div className="hero-content" style={{ 
          position: 'absolute', 
          zIndex: 5,
          top: '50%',
          left: '5%',
          transform: 'translateY(-50%)',
          maxWidth: '500px',
          textAlign: 'left',
          padding: '0 0.5rem'
        }}>
          <div style={{
            display: 'flex',
            alignItems: 'center',
            marginBottom: '0.5rem',
            gap: '1rem'
          }}>
            <img 
              src="/Qoalation.png" 
              alt="Qoal Logo" 
              loading="eager"
              style={{
                width: '75px',
                height: '75px',
                objectFit: 'contain',
                display: 'block'
              }}
            />
            <img 
              src="/QoalText.png" 
              alt="Qoal Text" 
              loading="eager"
              style={{
                height: '30px',
                objectFit: 'contain',
                display: 'block'
              }}
            />
          </div>
          <h1 className="general-title" style={{
            color: '#fff',
            marginBottom: '0.5rem',
            textShadow: '0 2px 4px rgba(0,0,0,0.3)',
            fontSize: 'clamp(1.875rem, 3.75vw, 3rem)',
            lineHeight: '1.2',
            fontWeight: 500
          }}>Convert Your</h1>
          <div style={{ marginBottom: '1rem', fontFamily: 'Poppins', fontSize: 'clamp(1.875rem, 3.75vw, 3rem)', fontWeight: 500 }}>
            <RotatingText
              texts={['Image', 'Video', 'Archive', 'Document', 'Audio']}
              mainClassName="inline-flex px-2 sm:px-2 md:px-3 bg-[#ffb947] text-black overflow-hidden py-1 sm:py-1 md:py-1 rounded-lg"
              staggerFrom={"last"}
              initial={{ y: "100%" }}
              animate={{ y: 0 }}
              exit={{ y: "-120%" }}
              staggerDuration={0.025}
              splitLevelClassName="overflow-hidden pb-0.5 sm:pb-1 md:pb-1"
              transition={{ type: "spring", damping: 30, stiffness: 400 }}
              rotationInterval={2000}
            />
          </div>
          <p style={{
            color: '#e9e9ef',
            fontSize: 'clamp(0.75rem, 1.5vw, 0.9375rem)',
            marginBottom: '1rem',
            lineHeight: '1.6',
            textShadow: '0 1px 2px rgba(0,0,0,0.3)'
          }}>Fast, secure, and easy file conversion for all your needs</p>
          <button 
            className="get-started-btn" 
            onClick={handleGetStarted}
            style={{
              backgroundColor: '#ffb947',
              color: '#161B27',
              border: 'none',
              padding: '0.75rem 1.5rem',
              fontSize: 'clamp(0.75rem, 1.125vw, 0.825rem)',
              fontWeight: '600',
              borderRadius: '0.5rem',
              cursor: 'pointer',
              transition: 'all 0.3s ease',
              boxShadow: '0 4px 12px rgba(255, 120, 90, 0.3)'
            }}
            onMouseEnter={(e) => e.currentTarget.style.backgroundColor = '#feab25'}
            onMouseLeave={(e) => e.currentTarget.style.backgroundColor = '#ffb947'}
          >
            Get Started
          </button>
        </div>

        {/* 3D Folder Model will be added here by the useEffect */}
      </section>


      {/* Conversion Types Sections */}
      <section className="conversion-sections" style={{ position: 'relative' }}>
        <div className="section image-conversions" style={{  
          height: '90vh',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'flex-end',
          padding: '0 5%',
          position: 'relative'
        }}>
          <img 
            src="/ImgIco.gif" 
            alt="Conversion Animation" 
            className="conversion-gif"
            style={{
              position: 'absolute',
              left: '5%',
              width: '18.75%',
              height: 'auto',
              maxWidth: '150px',
              zIndex: 5,
              opacity: 0
            }}
          />
          <div style={{ textAlign: 'right', maxWidth: '500px' }}>
            <h2 style={{ fontSize: '2.25rem', marginBottom: '1rem', color: '#fff' }}>Image Conversions</h2>
            <p style={{ fontSize: '0.9375rem', marginBottom: '2rem', color: '#e9e9ef' }}>
              Transform your images between popular formats with ease. Maintain quality while optimizing file size.
            </p>
            <div style={{ display: 'flex', flexWrap: 'wrap', gap: '0.75rem', justifyContent: 'flex-end', marginBottom: '2rem' }}>
              {['JPG', 'PNG', 'WebP', 'GIF', 'BMP', 'TIFF'].map(format => (
                <div
                  key={format}
                  className="format-btn"
                  style={{
                    padding: '0.5625rem 1.125rem',
                    backgroundColor: 'rgba(255, 255, 255, 0.1)',
                    color: '#fff',
                    borderRadius: '0.5rem',
                    fontSize: '0.75rem',
                    fontWeight: '600',
                    cursor: 'pointer',
                    transition: 'all 0.3s ease',
                    border: '2px solid rgba(255, 255, 255, 0.2)'
                  }}
                  onMouseEnter={(e) => {
                    e.currentTarget.style.backgroundColor = '#ffb947';
                    e.currentTarget.style.borderColor = '#ffb947';
                    e.currentTarget.style.transform = 'translateY(-2px)';
                  }}
                  onMouseLeave={(e) => {
                    e.currentTarget.style.backgroundColor = 'rgba(255, 255, 255, 0.1)';
                    e.currentTarget.style.borderColor = 'rgba(255, 255, 255, 0.2)';
                    e.currentTarget.style.transform = 'translateY(0)';
                  }}
                >
                  {format}
                </div>
              ))}
            </div>

          </div>
        </div>

        <div className="section document-conversions" style={{ 
          height: '100vh',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'flex-end',
          padding: '0 5%',
          position: 'relative'
        }}>
          <img 
            src="/DocIco.gif" 
            alt="Document Animation" 
            className="document-gif"
            style={{
              position: 'absolute',
              left: '5%',
              width: '18.75%',
              height: 'auto',
              maxWidth: '150px',
              zIndex: 5,
              opacity: 0
            }}
          />
          <div style={{ textAlign: 'right', maxWidth: '500px' }}>
            <h2 style={{ fontSize: '2.25rem', marginBottom: '1rem', color: '#fff' }}>Document Conversions</h2>
            <p style={{ fontSize: '0.9375rem', marginBottom: '2rem', color: '#e9e9ef' }}>
              Convert documents seamlessly between formats. Perfect for office work and document management.
            </p>
            <div style={{ display: 'flex', flexWrap: 'wrap', gap: '0.75rem', justifyContent: 'flex-end', marginBottom: '2rem' }}>
              {['PDF', 'DOCX', 'TXT', 'XLSX', 'CSV'].map(format => (
                <div
                  key={format}
                  className="format-btn"
                  style={{
                    padding: '0.5625rem 1.125rem',
                    backgroundColor: 'rgba(255, 255, 255, 0.1)',
                    color: '#fff',
                    borderRadius: '0.5rem',
                    fontSize: '0.75rem',
                    fontWeight: '600',
                    cursor: 'pointer',
                    transition: 'all 0.3s ease',
                    border: '2px solid rgba(255, 255, 255, 0.2)'
                  }}
                  onMouseEnter={(e) => {
                    e.currentTarget.style.backgroundColor = '#ffb947';
                    e.currentTarget.style.borderColor = '#ffb947';
                    e.currentTarget.style.transform = 'translateY(-2px)';
                  }}
                  onMouseLeave={(e) => {
                    e.currentTarget.style.backgroundColor = 'rgba(255, 255, 255, 0.1)';
                    e.currentTarget.style.borderColor = 'rgba(255, 255, 255, 0.2)';
                    e.currentTarget.style.transform = 'translateY(0)';
                  }}
                >
                  {format}
                </div>
              ))}
            </div>
          </div>
        </div>

        <div className="section audio-conversions" style={{ 
          height: '100vh',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'flex-end',
          padding: '0 5%',
          position: 'relative'
        }}>
          <img 
            src="/AudIco.gif" 
            alt="Audio Animation" 
            className="audio-gif"
            style={{
              position: 'absolute',
              left: '5%',
              width: '18.75%',
              height: 'auto',
              maxWidth: '150px',
              zIndex: 5,
              opacity: 0
            }}
          />
          <div style={{ textAlign: 'right', maxWidth: '500px' }}>
            <h2 style={{ fontSize: '2.25rem', marginBottom: '1rem', color: '#fff' }}>Audio Conversions</h2>
            <p style={{ fontSize: '0.9375rem', marginBottom: '2rem', color: '#e9e9ef' }}>
              Convert audio files to your preferred format. Maintain sound quality across all conversions.
            </p>
            <div style={{ display: 'flex', flexWrap: 'wrap', gap: '0.75rem', justifyContent: 'flex-end', marginBottom: '2rem' }}>
              {['MP3', 'WAV', 'FLAC', 'M4A', 'OGG'].map(format => (
                <div
                  key={format}
                  className="format-btn"
                  style={{
                    padding: '0.5625rem 1.125rem',
                    backgroundColor: 'rgba(255, 255, 255, 0.1)',
                    color: '#fff',
                    borderRadius: '0.5rem',
                    fontSize: '0.75rem',
                    fontWeight: '600',
                    cursor: 'pointer',
                    transition: 'all 0.3s ease',
                    border: '2px solid rgba(255, 255, 255, 0.2)'
                  }}
                  onMouseEnter={(e) => {
                    e.currentTarget.style.backgroundColor = '#ffb947';
                    e.currentTarget.style.borderColor = '#ffb947';
                    e.currentTarget.style.transform = 'translateY(-2px)';
                  }}
                  onMouseLeave={(e) => {
                    e.currentTarget.style.backgroundColor = 'rgba(255, 255, 255, 0.1)';
                    e.currentTarget.style.borderColor = 'rgba(255, 255, 255, 0.2)';
                    e.currentTarget.style.transform = 'translateY(0)';
                  }}
                >
                  {format}
                </div>
              ))}
            </div>
          </div>
        </div>

        <div className="section video-conversions" style={{ 
          height: '100vh',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'flex-end',
          padding: '0 5%',
          position: 'relative'
        }}>
          <img 
            src="/VidIco.gif" 
            alt="Video Animation" 
            className="video-gif"
            style={{
              position: 'absolute',
              left: '5%',
              width: '18.75%',
              height: 'auto',
              maxWidth: '150px',
              zIndex: 5,
              opacity: 0
            }}
          />
          <div style={{ textAlign: 'right', maxWidth: '500px' }}>
            <h2 style={{ fontSize: '2.25rem', marginBottom: '1rem', color: '#fff' }}>Video Conversions</h2>
            <p style={{ fontSize: '0.9375rem', marginBottom: '2rem', color: '#e9e9ef' }}>
              Transform videos between formats effortlessly. Optimize for web, mobile, or desktop playback.
            </p>
            <div style={{ display: 'flex', flexWrap: 'wrap', gap: '0.75rem', justifyContent: 'flex-end', marginBottom: '2rem' }}>
              {['MP4', 'AVI', 'MOV', 'WebM', 'MKV'].map(format => (
                <div
                  key={format}
                  className="format-btn"
                  style={{
                    padding: '0.5625rem 1.125rem',
                    backgroundColor: 'rgba(255, 255, 255, 0.1)',
                    color: '#fff',
                    borderRadius: '0.5rem',
                    fontSize: '0.75rem',
                    fontWeight: '600',
                    cursor: 'pointer',
                    transition: 'all 0.3s ease',
                    border: '2px solid rgba(255, 255, 255, 0.2)'
                  }}
                  onMouseEnter={(e) => {
                    e.currentTarget.style.backgroundColor = '#ffb947';
                    e.currentTarget.style.borderColor = '#ffb947';
                    e.currentTarget.style.transform = 'translateY(-2px)';
                  }}
                  onMouseLeave={(e) => {
                    e.currentTarget.style.backgroundColor = 'rgba(255, 255, 255, 0.1)';
                    e.currentTarget.style.borderColor = 'rgba(255, 255, 255, 0.2)';
                    e.currentTarget.style.transform = 'translateY(0)';
                  }}
                >
                  {format}
                </div>
              ))}
            </div>
          </div>
        </div>

        <div className="section archive-conversions" style={{ 
          height: '100vh',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'flex-end',
          padding: '0 5%',
          position: 'relative'
        }}>
          <img 
            src="/ArcIco.gif" 
            alt="Archive Animation" 
            className="archive-gif"
            style={{
              position: 'absolute',
              left: '5%',
              width: '18.75%',
              height: 'auto',
              maxWidth: '150px',
              zIndex: 5,
              opacity: 0
            }}
          />
          <div style={{ textAlign: 'right', maxWidth: '500px' }}>
            <h2 style={{ fontSize: '2.25rem', marginBottom: '1rem', color: '#fff' }}>Archive Conversions</h2>
            <p style={{ fontSize: '0.9375rem', marginBottom: '2rem', color: '#e9e9ef' }}>
              Compress and convert archives with ease. Support for all major compression formats.
            </p>
            <div style={{ display: 'flex', flexWrap: 'wrap', gap: '0.75rem', justifyContent: 'flex-end', marginBottom: '2rem' }}>
              {['ZIP', 'RAR', 'TAR.GZ'].map(format => (
                <div
                  key={format}
                  className="format-btn"
                  style={{
                    padding: '0.5625rem 1.125rem',
                    backgroundColor: 'rgba(255, 255, 255, 0.1)',
                    color: '#fff',
                    borderRadius: '0.5rem',
                    fontSize: '0.75rem',
                    fontWeight: '600',
                    cursor: 'pointer',
                    transition: 'all 0.3s ease',
                    border: '2px solid rgba(255, 255, 255, 0.2)'
                  }}
                  onMouseEnter={(e) => {
                    e.currentTarget.style.backgroundColor = '#ffb947';
                    e.currentTarget.style.borderColor = '#ffb947';
                    e.currentTarget.style.transform = 'translateY(-2px)';
                  }}
                  onMouseLeave={(e) => {
                    e.currentTarget.style.backgroundColor = 'rgba(255, 255, 255, 0.1)';
                    e.currentTarget.style.borderColor = 'rgba(255, 255, 255, 0.2)';
                    e.currentTarget.style.transform = 'translateY(0)';
                  }}
                >
                  {format}
                </div>
              ))}
            </div>
          </div>
        </div>
      </section>

      {/* Footer */}
      <footer className="footer" style={{ 
        position: 'relative',
        minHeight: '10vh',
        padding: '3rem 5%'
      }}>
        <div className="footer-content" style={{
          position: 'relative',
          zIndex: 1,
          display: 'flex',
          flexDirection: 'column',
          alignItems: 'center',
          justifyContent: 'center',
          height: '100%',
          gap: '1rem'
        }}>
          <p style={{ color: '#cae2e2', fontSize: 'clamp(0.675rem, 2.25vw, 1.125rem)', textAlign: 'center' }}>&copy; Qoal it iz.<br></br> Made with Gsap n Threejs |0_0|</p>
          <div className="footer-links" style={{ display: 'flex', gap: '2rem' }}>
            <a href="https://www.github.com/F4tal1t/Qoal" style={{ color: '#ffb947', textDecoration: 'none', fontSize: 'clamp(0.6375rem, 1.875vw, 0.75rem)' }}>Github</a>
            <a href="https://dibby.me" style={{ color: '#ffb947', textDecoration: 'none', fontSize: 'clamp(0.6375rem, 1.875vw, 0.75rem)' }}>Creator's Portfolio</a>
          </div>
        </div>
      </footer>
    </div>
    </div>
  );
};

export default Landing;