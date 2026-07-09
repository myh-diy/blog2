<template>
  <svg viewBox="0 0 400 400" class="w-full h-full" xmlns="http://www.w3.org/2000/svg">
    <!-- Background soft circle -->
    <circle cx="200" cy="200" r="170" fill="var(--brand-50, #fff1f8)" opacity="0.6"/>

    <!-- Back hair -->
    <path fill="var(--brand-200, #ffc9e9)" d="M110 180 Q100 260 120 310 Q200 340 280 310 Q300 260 290 180 Q290 110 200 100 Q110 110 110 180"/>

    <!-- Body / Dress -->
    <path fill="var(--brand-200, #ffc9e9)" d="M150 280 Q200 270 250 280 L270 360 Q200 380 130 360 Z"/>
    <path fill="var(--brand-100, #ffe4f3)" d="M160 285 Q200 278 240 285 L255 350 Q200 365 145 350 Z"/>

    <!-- Neck -->
    <rect x="185" y="250" width="30" height="40" rx="8" fill="#fff1e6"/>

    <!-- Face -->
    <ellipse cx="200" cy="200" rx="65" ry="72" fill="#fff1e6"/>

    <!-- Front hair / bangs -->
    <path fill="var(--brand-200, #ffc9e9)" d="M135 160 Q150 120 200 120 Q250 120 265 160 Q270 200 255 210 Q240 170 200 165 Q160 170 145 210 Q130 200 135 160"/>
    <path fill="var(--brand-300, #ff9dd4)" d="M180 125 Q200 140 220 125 Q210 155 200 150 Q190 155 180 125" opacity="0.6"/>

    <!-- Eyes -->
    <ellipse cx="170" cy="205" rx="11" ry="15" fill="#475569"/>
    <ellipse cx="230" cy="205" rx="11" ry="15" fill="#475569"/>
    <circle cx="174" cy="200" r="5" fill="white"/>
    <circle cx="234" cy="200" r="5" fill="white"/>
    <circle cx="172" cy="212" r="3" fill="var(--brand-400, #fb7299)" opacity="0.4"/>
    <circle cx="232" cy="212" r="3" fill="var(--brand-400, #fb7299)" opacity="0.4"/>

    <!-- Blush -->
    <ellipse cx="150" cy="230" rx="14" ry="9" fill="var(--brand-300, #ff9dd4)" opacity="0.5"/>
    <ellipse cx="250" cy="230" rx="14" ry="9" fill="var(--brand-300, #ff9dd4)" opacity="0.5"/>

    <!-- Mouth -->
    <path d="M190 248 Q200 256 210 248" fill="none" stroke="#475569" stroke-width="2.5" stroke-linecap="round"/>

    <!-- Bow hair accessory -->
    <g transform="translate(250, 130)">
      <ellipse cx="0" cy="0" rx="18" ry="12" transform="rotate(-30)" fill="var(--accent-400, #f472b6)"/>
      <ellipse cx="0" cy="0" rx="18" ry="12" transform="rotate(30)" fill="var(--accent-400, #f472b6)"/>
      <circle cx="0" cy="0" r="6" fill="var(--accent-500, #ec4899)"/>
    </g>

    <!-- Hands -->
    <ellipse cx="135" cy="310" rx="14" ry="22" transform="rotate(20 135 310)" fill="#fff1e6"/>
    <ellipse cx="265" cy="310" rx="14" ry="22" transform="rotate(-20 265 310)" fill="#fff1e6"/>

    <!-- Sparkles -->
    <path fill="var(--accent-400, #f472b6)" d="M60 120 Q65 105 70 120 Q85 125 70 130 Q65 145 60 130 Q45 125 60 120" opacity="0.7"/>
    <path fill="var(--brand-300, #ff9dd4)" d="M340 100 Q344 88 348 100 Q360 104 348 108 Q344 120 340 108 Q328 104 340 100" opacity="0.6"/>
    <circle cx="80" cy="280" r="5" fill="var(--brand-200, #ffc9e9)" opacity="0.5"/>
    <circle cx="330" cy="260" r="4" fill="var(--accent-300, #f9a8d4)" opacity="0.5"/>
    <path fill="var(--accent-400, #f472b6)" d="M320 340 Q325 330 330 340 Q340 345 330 350 Q325 360 320 350 Q310 345 320 340" opacity="0.5"/>
  </svg>
</template>
