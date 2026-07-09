<script setup lang="ts">
const props = defineProps<{ name: 'wave' | 'search' | 'sad' | 'happy' | 'star' | 'heart' | 'upload' }>()
</script>

<template>
  <!-- Wave: 挥手小人 -->
  <svg v-if="props.name === 'wave'" viewBox="0 0 120 120" class="w-full h-full">
    <circle cx="60" cy="60" r="56" fill="#fff1f8"/>
    <circle cx="60" cy="52" r="28" fill="#ffebf7" stroke="#fb7299" stroke-width="3"/>
    <circle cx="50" cy="48" r="4" fill="#8b5cf6"/>
    <circle cx="70" cy="48" r="4" fill="#8b5cf6"/>
    <path d="M52 62 Q60 70 68 62" stroke="#fb7299" stroke-width="3" fill="none" stroke-linecap="round"/>
    <ellipse cx="40" cy="74" rx="10" ry="8" fill="#ffc9e9"/>
    <ellipse cx="80" cy="74" rx="10" ry="8" fill="#ffc9e9"/>
    <circle cx="16" cy="38" r="8" fill="#ffc9e9"/>
    <circle cx="104" cy="30" r="6" fill="#ddd6fe"/>
    <circle cx="100" cy="86" r="5" fill="#ff9dd4"/>
  </svg>

  <!-- Search: 搜索小人 -->
  <svg v-else-if="props.name === 'search'" viewBox="0 0 120 120" class="w-full h-full">
    <circle cx="60" cy="60" r="56" fill="#f5f3ff"/>
    <circle cx="60" cy="52" r="26" fill="#ede9fe" stroke="#8b5cf6" stroke-width="3"/>
    <circle cx="52" cy="48" r="3.5" fill="#fb7299"/>
    <circle cx="68" cy="48" r="3.5" fill="#fb7299"/>
    <circle cx="60" cy="60" r="4" fill="#fb7299" opacity="0.6"/>
    <path d="M80 76 L94 90" stroke="#8b5cf6" stroke-width="5" stroke-linecap="round"/>
    <circle cx="78" cy="74" r="10" stroke="#8b5cf6" stroke-width="4" fill="none"/>
    <circle cx="18" cy="36" r="6" fill="#ff9dd4"/>
    <circle cx="102" cy="28" r="5" fill="#c4b5fd"/>
  </svg>

  <!-- Sad: 难过/空状态小人 -->
  <svg v-else-if="props.name === 'sad'" viewBox="0 0 120 120" class="w-full h-full">
    <circle cx="60" cy="60" r="56" fill="#fff1f8"/>
    <circle cx="60" cy="54" r="28" fill="#ffebf7" stroke="#fb7299" stroke-width="3"/>
    <circle cx="50" cy="50" r="4" fill="#8b5cf6"/>
    <circle cx="70" cy="50" r="4" fill="#8b5cf6"/>
    <path d="M52 68 Q60 62 68 68" stroke="#fb7299" stroke-width="3" fill="none" stroke-linecap="round"/>
    <path d="M38 38 L46 44" stroke="#ddd6fe" stroke-width="4" stroke-linecap="round"/>
    <path d="M82 38 L74 44" stroke="#ddd6fe" stroke-width="4" stroke-linecap="round"/>
    <circle cx="18" cy="86" r="5" fill="#ff9dd4"/>
    <circle cx="102" cy="86" r="5" fill="#ddd6fe"/>
  </svg>

  <!-- Happy: 开心小人 -->
  <svg v-else-if="props.name === 'happy'" viewBox="0 0 120 120" class="w-full h-full">
    <circle cx="60" cy="60" r="56" fill="#fff1f8"/>
    <circle cx="60" cy="52" r="28" fill="#ffebf7" stroke="#fb7299" stroke-width="3"/>
    <path d="M46 50 Q50 46 54 50" stroke="#8b5cf6" stroke-width="3" fill="none" stroke-linecap="round"/>
    <path d="M66 50 Q70 46 74 50" stroke="#8b5cf6" stroke-width="3" fill="none" stroke-linecap="round"/>
    <path d="M48 62 Q60 74 72 62" stroke="#fb7299" stroke-width="3" fill="none" stroke-linecap="round"/>
    <ellipse cx="38" cy="72" rx="9" ry="7" fill="#ffc9e9"/>
    <ellipse cx="82" cy="72" rx="9" ry="7" fill="#ffc9e9"/>
    <circle cx="16" cy="32" r="6" fill="#ff9dd4"/>
    <circle cx="104" cy="28" r="5" fill="#ddd6fe"/>
    <path d="M86 46 L96 38" stroke="#fb7299" stroke-width="3" stroke-linecap="round"/>
  </svg>

  <!-- Star: 星星 -->
  <svg v-else-if="props.name === 'star'" viewBox="0 0 120 120" class="w-full h-full">
    <circle cx="60" cy="60" r="56" fill="#fff1f8"/>
    <path d="M60 25 L68 48 L93 48 L73 63 L81 86 L60 72 L39 86 L47 63 L27 48 L52 48 Z" fill="url(#starGrad)"/>
    <defs>
      <linearGradient id="starGrad" x1="0" y1="0" x2="1" y2="1">
        <stop offset="0%" stop-color="#ff5fb8"/>
        <stop offset="100%" stop-color="#8b5cf6"/>
      </linearGradient>
    </defs>
    <circle cx="20" cy="40" r="5" fill="#ffc9e9"/>
    <circle cx="100" cy="84" r="4" fill="#ddd6fe"/>
  </svg>

  <!-- Heart: 爱心 -->
  <svg v-else-if="props.name === 'heart'" viewBox="0 0 120 120" class="w-full h-full">
    <circle cx="60" cy="60" r="56" fill="#fff1f8"/>
    <path d="M60 88 C60 88 30 66 30 48 C30 36 39 28 50 28 C56 28 60 32 60 32 C60 32 64 28 70 28 C81 28 90 36 90 48 C90 66 60 88 60 88 Z" fill="#fb7299"/>
    <circle cx="22" cy="38" r="5" fill="#ddd6fe"/>
    <circle cx="98" cy="44" r="4" fill="#ff9dd4"/>
  </svg>

  <!-- Upload: 上传小人 -->
  <svg v-else-if="props.name === 'upload'" viewBox="0 0 120 120" class="w-full h-full">
    <circle cx="60" cy="60" r="56" fill="#f5f3ff"/>
    <circle cx="60" cy="56" r="26" fill="#ede9fe" stroke="#8b5cf6" stroke-width="3"/>
    <circle cx="52" cy="52" r="3.5" fill="#fb7299"/>
    <circle cx="68" cy="52" r="3.5" fill="#fb7299"/>
    <path d="M54 64 Q60 70 66 64" stroke="#fb7299" stroke-width="3" fill="none" stroke-linecap="round"/>
    <path d="M60 26 L60 46 M50 36 L60 26 L70 36" stroke="#8b5cf6" stroke-width="4" fill="none" stroke-linecap="round" stroke-linejoin="round"/>
    <circle cx="20" cy="88" r="5" fill="#ff9dd4"/>
    <circle cx="100" cy="90" r="4" fill="#ddd6fe"/>
  </svg>
</template>
