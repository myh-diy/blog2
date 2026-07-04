<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import api from '../utils/api'

interface P { id: number; text: string; x: number; y: number; c: string; s: number; vx: number; vy: number; o: number }
const particles = ref<P[]>([])
const ready = ref(false)
const C = ['#f97316','#ea580c','#fb923c','#f59e0b','#d97706','#fbbf24','#f43f5e','#e879f9','#ec4899','#d946ef','#fb7185','#a78bfa']
let raf = 0

onMounted(async () => {
  let quotes: string[] = []
  try { const r = await api.get('/quotes'); quotes = r.data.quotes } catch {}
  if (!quotes.length) quotes = ['Stay curious, keep learning.','Code is poetry.','Every bug is a lesson.','Build something awesome.','Knowledge is power.','Learn by doing.','Simplicity is key.','Think different.']
  const items: P[] = []
  for (let i=0; i<16; i++) items.push({
    id:i, text:quotes[i%quotes.length],
    x:Math.random()*90+5, y:Math.random()*80+10,
    c:C[i%C.length], s:14+Math.floor(Math.random()*20),
    vx:0.015+Math.random()*0.05, vy:0.008+Math.random()*0.035,
    o:0.3+Math.random()*0.45,
  })
  particles.value = items
  ready.value = true
  animate()
})

onUnmounted(() => cancelAnimationFrame(raf))

function animate() {
  particles.value = particles.value.map(p => {
    let x=p.x+p.vx, y=p.y+p.vy
    if(x<-10||x>110) p.vx=-p.vx
    if(y<-10||y>110) p.vy=-p.vy
    return {...p, x, y}
  })
  raf = requestAnimationFrame(animate)
}
</script>

<template>
  <div v-if="ready" class="home-hero">
    <div class="home-bg"></div>
    <div class="home-dots"></div>
    <div class="home-vignette"></div>
    <div v-for="p in particles" :key="p.id" class="home-p"
      :style="{left:p.x+'%',top:p.y+'%',fontSize:p.s+'px',color:p.c,opacity:p.o,textShadow:`0 0 20px ${p.c}`}"
    >{{ p.text }}</div>
    <div class="home-center">
      <div class="home-icon">B</div>
      <h1 class="home-h1">My <span class="home-span">Blog</span></h1>
      <p class="home-sub">Learning in public</p>
      <router-link to="/posts" class="home-link">Explore posts ▾</router-link>
    </div>
  </div>
</template>

<style>
.home-hero { position: fixed; top: 0; left: 0; width: 100%; height: 100vh; z-index: 30; overflow: hidden; display: flex; align-items: center; justify-content: center; }
.home-bg { position: absolute; inset: 0; background: linear-gradient(160deg, #fff7ed, #fed7aa 30%, #ffedd5 60%, #fff7ed); }
.dark .home-bg { background: linear-gradient(160deg, #0d0714, #1a0a1f 30%, #130818 60%, #0d0714); }
.home-dots { position: absolute; inset: 0; opacity: .05; background-image: radial-gradient(circle at 1px 1px, currentColor 1px, transparent 0); background-size: 40px 40px; }
.home-vignette { position: absolute; inset: 0; background: radial-gradient(ellipse at center, transparent 40%, rgba(0,0,0,.15) 100%); pointer-events: none; }
.home-p { position: absolute; font-weight: 500; white-space: nowrap; pointer-events: none; user-select: none; }
.home-center { position: relative; z-index: 10; text-align: center; }
.home-icon { width: 80px; height: 80px; border-radius: 20px; margin: 0 auto 1.5rem; display: flex; align-items: center; justify-content: center; color: #fff; font-weight: 900; font-size: 30px; background: linear-gradient(135deg, #f97316, #f43f5e); box-shadow: 0 0 40px rgba(249,115,22,.3); }
.dark .home-icon { background: linear-gradient(135deg, #f43f5e, #d946ef); box-shadow: 0 0 50px rgba(244,63,94,.3); }
.home-h1 { font-size: clamp(3rem, 6vw, 5rem); font-weight: 900; color: #292524; letter-spacing: -.03em; }
.dark .home-h1 { color: #f5f5f4; }
.home-span { background: linear-gradient(to right, #f97316, #ea580c); -webkit-background-clip: text; -webkit-text-fill-color: transparent; background-clip: text; }
.dark .home-span { background: linear-gradient(to right, #f43f5e, #d946ef); -webkit-background-clip: text; -webkit-text-fill-color: transparent; background-clip: text; }
.home-sub { font-size: 1.125rem; font-weight: 300; letter-spacing: .15em; text-transform: uppercase; color: rgba(120,113,108,.7); margin-top: .5rem; }
.dark .home-sub { color: rgba(245,245,244,.4); }
.home-link { display: inline-flex; align-items: center; gap: .5rem; margin-top: 2rem; padding: .75rem 1.5rem; border-radius: 9999px; border: 1px solid rgba(249,115,22,.3); color: #c2410c; font-size: .875rem; text-decoration: none; transition: all .2s; }
.dark .home-link { border-color: rgba(244,63,94,.3); color: rgba(251,113,133,.8); }
.home-link:hover { background: rgba(249,115,22,.08); }
.dark .home-link:hover { background: rgba(244,63,94,.08); }
</style>
