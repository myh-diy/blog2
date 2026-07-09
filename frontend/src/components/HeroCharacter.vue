<template>
  <svg viewBox="0 0 400 400" class="w-full h-full" xmlns="http://www.w3.org/2000/svg">
    <!-- Background decoration -->
    <circle cx="200" cy="200" r="180" fill="var(--brand-50)" class="dark:opacity-20"/>
    <circle cx="320" cy="80" r="30" fill="var(--accent-100)" class="dark:opacity-30"/>
    <circle cx="60" cy="100" r="20" fill="var(--brand-100)" class="dark:opacity-30"/>
    <circle cx="340" cy="300" r="24" fill="var(--brand-100)" class="dark:opacity-30"/>

    <!-- Character body -->
    <path d="M130 380 Q200 420 270 380 L270 300 Q200 280 130 300 Z" fill="#f8fafc"/>
    <path d="M130 300 Q200 280 270 300 L270 380 Q200 420 130 380 Z" fill="var(--brand-100)" opacity="0.5"/>

    <!-- Hair back -->
    <path d="M120 180 Q90 280 110 340 L290 340 Q310 280 280 180 Q200 120 120 180" fill="var(--brand-200)"/>

    <!-- Face -->
    <ellipse cx="200" cy="200" rx="75" ry="85" fill="#fefce8"/>

    <!-- Hair front -->
    <path d="M125 160 Q140 120 200 120 Q260 120 275 160 Q285 200 275 230 Q260 180 200 170 Q140 180 125 230 Q115 200 125 160" fill="var(--brand-300)"/>
    <path d="M180 125 Q200 140 220 125" fill="none" stroke="var(--brand-400)" stroke-width="3" stroke-linecap="round"/>

    <!-- Eyes -->
    <ellipse cx="170" cy="205" rx="12" ry="16" fill="#334155"/>
    <ellipse cx="230" cy="205" rx="12" ry="16" fill="#334155"/>
    <circle cx="174" cy="200" r="4" fill="white"/>
    <circle cx="234" cy="200" r="4" fill="white"/>

    <!-- Blush -->
    <ellipse cx="150" cy="230" rx="12" ry="8" fill="var(--brand-200)" opacity="0.6"/>
    <ellipse cx="250" cy="230" rx="12" ry="8" fill="var(--brand-200)" opacity="0.6"/>

    <!-- Mouth -->
    <path d="M190 245 Q200 255 210 245" fill="none" stroke="#334155" stroke-width="2.5" stroke-linecap="round"/>

    <!-- Leaf hair accessory -->
    <path d="M260 140 Q290 120 300 150 Q280 160 260 140" fill="var(--accent-400)"/>
    <path d="M280 145 Q300 135 310 155" fill="none" stroke="var(--accent-600)" stroke-width="1.5" stroke-linecap="round"/>

    <!-- Arms -->
    <ellipse cx="115" cy="320" rx="18" ry="35" fill="#fefce8" transform="rotate(20 115 320)"/>
    <ellipse cx="285" cy="320" rx="18" ry="35" fill="#fefce8" transform="rotate(-20 285 320)"/>

    <!-- Book/notepad -->
    <rect x="155" y="330" width="90" height="70" rx="6" fill="white" stroke="#e2e8f0" stroke-width="2"/>
    <rect x="160" y="335" width="80" height="60" rx="4" fill="var(--brand-50)"/>
    <line x1="170" y1="350" x2="230" y2="350" stroke="var(--brand-300)" stroke-width="3" stroke-linecap="round"/>
    <line x1="170" y1="365" x2="220" y2="365" stroke="var(--accent-300)" stroke-width="3" stroke-linecap="round"/>
    <line x1="170" y1="380" x2="210" y2="380" stroke="var(--brand-300)" stroke-width="3" stroke-linecap="round"/>

    <!-- Floating sparkles -->
    <path d="M70 250 L75 235 L80 250 L95 255 L80 260 L75 275 L70 260 L55 255 Z" fill="var(--accent-300)"/>
    <path d="M330 180 L333 170 L336 180 L346 183 L336 186 L333 196 L330 186 L320 183 Z" fill="var(--brand-300)"/>
    <circle cx="90" cy="320" r="5" fill="var(--brand-200)"/>
    <circle cx="310" cy="260" r="4" fill="var(--accent-200)"/>
  </svg>
</template>
