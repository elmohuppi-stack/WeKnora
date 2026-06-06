<template>
  <teleport to="body">
    <transition name="upload-sheet">
      <div
        v-if="visible"
        class="upload-sheet-overlay"
        @click.self="emit('close')"
      >
        <div class="upload-sheet-panel">
          <div class="upload-sheet-header">
            <span class="upload-sheet-title">{{
              t("knowledgeBase.addKnowledge") || "Wissen hinzufügen"
            }}</span>
            <div class="upload-sheet-close" @click="emit('close')">
              <svg
                width="20"
                height="20"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
              >
                <line x1="18" y1="6" x2="6" y2="18" />
                <line x1="6" y1="6" x2="18" y2="18" />
              </svg>
            </div>
          </div>
          <div class="upload-sheet-items">
            <div class="upload-sheet-item" @click="emit('uploadFile')">
              <div class="upload-sheet-item-icon upload-icon-file">
                <svg
                  width="24"
                  height="24"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="2"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                >
                  <path
                    d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"
                  />
                  <polyline points="14 2 14 8 20 8" />
                  <line x1="12" y1="18" x2="12" y2="12" />
                  <line x1="9" y1="15" x2="15" y2="15" />
                </svg>
              </div>
              <div class="upload-sheet-item-text">
                <span class="upload-sheet-item-title">{{
                  t("knowledgeBase.uploadFile") || "Datei hochladen"
                }}</span>
                <span class="upload-sheet-item-desc"
                  >PDF, DOCX, TXT, MD, Bilder uvm.</span
                >
              </div>
            </div>
            <div class="upload-sheet-item" @click="emit('importURL')">
              <div class="upload-sheet-item-icon upload-icon-url">
                <svg
                  width="24"
                  height="24"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="2"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                >
                  <path
                    d="M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.71"
                  />
                  <path
                    d="M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.71-1.71"
                  />
                </svg>
              </div>
              <div class="upload-sheet-item-text">
                <span class="upload-sheet-item-title">{{
                  t("knowledgeBase.importFromURL") || "URL importieren"
                }}</span>
                <span class="upload-sheet-item-desc">{{
                  t("knowledgeBase.urlImportTip") ||
                  "Webseite als Quelle hinzufügen"
                }}</span>
              </div>
            </div>
            <div class="upload-sheet-item" @click="emit('importYouTube')">
              <div class="upload-sheet-item-icon upload-icon-youtube">
                <svg
                  width="24"
                  height="24"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="2"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                >
                  <path
                    d="M22.54 6.42a2.78 2.78 0 0 0-1.94-2C18.88 4 12 4 12 4s-6.88 0-8.6.46a2.78 2.78 0 0 0-1.94 2A29 29 0 0 0 1 12a29 29 0 0 0 .46 5.58 2.78 2.78 0 0 0 1.94 2C5.12 20 12 20 12 20s6.88 0 8.6-.46a2.78 2.78 0 0 0 1.94-2A29 29 0 0 0 23 12a29 29 0 0 0-.46-5.58z"
                  />
                  <polygon points="9.75 15.02 15.5 12 9.75 8.98 9.75 15.02" />
                </svg>
              </div>
              <div class="upload-sheet-item-text">
                <span class="upload-sheet-item-title">{{
                  t("knowledgeBase.importYouTube") || "YouTube importieren"
                }}</span>
                <span class="upload-sheet-item-desc">{{
                  t("knowledgeBase.youtubeUrlTip") ||
                  "Transkript automatisch erfassen"
                }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </transition>
  </teleport>
</template>

<script setup lang="ts">
import { useI18n } from "vue-i18n";

defineProps<{
  visible: boolean;
}>();

const emit = defineEmits<{
  (e: "close"): void;
  (e: "uploadFile"): void;
  (e: "importURL"): void;
  (e: "importYouTube"): void;
}>();

const { t } = useI18n();
</script>

<style scoped lang="less">
.upload-sheet-overlay {
  position: fixed;
  inset: 0;
  z-index: 2000;
  background: rgba(0, 0, 0, 0.4);
  display: flex;
  align-items: flex-end;
  justify-content: center;
}

.upload-sheet-panel {
  width: 100%;
  max-width: 500px;
  background: var(--td-bg-color-container, #fff);
  border-radius: 16px 16px 0 0;
  padding-bottom: env(safe-area-inset-bottom, 16px);
  animation: slide-up 0.25s ease-out;
}

.upload-sheet-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px 12px;
  border-bottom: 1px solid var(--td-border-level-1-color, #f0f0f0);
}

.upload-sheet-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--td-text-color-primary, #333);
}

.upload-sheet-close {
  cursor: pointer;
  padding: 4px;
  color: var(--td-text-color-secondary, #999);
  -webkit-tap-highlight-color: transparent;

  &:active {
    opacity: 0.6;
  }
}

.upload-sheet-items {
  padding: 8px 0;
}

.upload-sheet-item {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 14px 20px;
  cursor: pointer;
  -webkit-tap-highlight-color: transparent;
  touch-action: manipulation;
  transition: background 0.15s;

  &:active {
    background: var(--td-bg-color-container-hover, #f5f5f5);
  }
}

.upload-sheet-item-icon {
  width: 40px;
  height: 40px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;

  &.upload-icon-file {
    background: #e8f5e9;
    color: #4caf50;
  }

  &.upload-icon-url {
    background: #e3f2fd;
    color: #2196f3;
  }

  &.upload-icon-youtube {
    background: #fce4ec;
    color: #f44336;
  }
}

.upload-sheet-item-text {
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
}

.upload-sheet-item-title {
  font-size: 15px;
  font-weight: 500;
  color: var(--td-text-color-primary, #333);
}

.upload-sheet-item-desc {
  font-size: 12px;
  color: var(--td-text-color-secondary, #999);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

@keyframes slide-up {
  from {
    transform: translateY(100%);
  }
  to {
    transform: translateY(0);
  }
}

.upload-sheet-enter-active,
.upload-sheet-leave-active {
  transition: opacity 0.2s ease;

  .upload-sheet-panel {
    transition: transform 0.25s ease;
  }
}

.upload-sheet-enter-from,
.upload-sheet-leave-to {
  opacity: 0;

  .upload-sheet-panel {
    transform: translateY(100%);
  }
}
</style>
