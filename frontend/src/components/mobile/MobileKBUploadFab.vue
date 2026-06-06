<template>
  <button class="mobile-fab" @click="$emit('fabClick')">
    <svg
      width="24"
      height="24"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      stroke-width="2.5"
      stroke-linecap="round"
      stroke-linejoin="round"
    >
      <line x1="12" y1="5" x2="12" y2="19" />
      <line x1="5" y1="12" x2="19" y2="12" />
    </svg>
  </button>
</template>

<script setup lang="ts">
defineEmits<{
  (e: "fabClick"): void;
}>();
</script>

<style scoped lang="less">
.mobile-fab {
  position: fixed;
  right: 20px;
  bottom: calc(72px + env(safe-area-inset-bottom, 0px));
  width: 52px;
  height: 52px;
  border-radius: 50%;
  background: var(--td-brand-color, #0052d9);
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  border: none;
  outline: none;
  padding: 0;
  margin: 0;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.25);
  cursor: pointer;
  z-index: 50;
  -webkit-tap-highlight-color: transparent;
  touch-action: manipulation;
  transition:
    transform 0.15s ease,
    box-shadow 0.15s ease;

  &:active {
    transform: scale(0.92);
    box-shadow: 0 2px 6px rgba(0, 0, 0, 0.2);
  }
}
</style>
