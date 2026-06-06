<template>
  <teleport to="body">
    <transition name="upload-sheet">
      <div
        v-if="visible"
        class="upload-sheet-overlay"
        @click.self="goBackOrClose"
      >
        <div class="upload-sheet-panel">
          <!-- Header -->
          <div class="upload-sheet-header">
            <div class="upload-sheet-header-left">
              <div
                v-if="page !== 'menu'"
                class="upload-sheet-back"
                @click="page = 'menu'"
              >
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
                  <line x1="19" y1="12" x2="5" y2="12" />
                  <polyline points="12 19 5 12 12 5" />
                </svg>
              </div>
            </div>
            <span class="upload-sheet-title">{{ pageTitle }}</span>
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

          <!-- Menu page -->
          <div v-if="page === 'menu'" class="upload-sheet-items">
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
            <div class="upload-sheet-item" @click="page = 'url'">
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
            <div class="upload-sheet-item" @click="page = 'youtube'">
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

          <!-- URL input page -->
          <div v-if="page === 'url'" class="upload-sheet-youtube-form">
            <div class="youtube-input-label">
              {{ t("knowledgeBase.urlLabel") || "URL" }}
            </div>
            <input
              ref="urlInputRef"
              v-model="urlValue"
              class="youtube-url-input"
              :placeholder="
                t('knowledgeBase.urlPlaceholder') || 'https://example.com'
              "
              type="url"
              inputmode="url"
              autofocus
              @keydown.enter="submitUrl"
            />
            <div class="youtube-input-tip">
              {{
                t("knowledgeBase.urlTip") || "Webseite als Quelle hinzufügen"
              }}
            </div>
            <button
              class="youtube-submit-btn"
              :disabled="!urlValue.trim()"
              @click="submitUrl"
            >
              {{ t("knowledgeBase.import") || "Importieren" }}
            </button>
          </div>

          <!-- YouTube URL input page -->
          <div v-if="page === 'youtube'" class="upload-sheet-youtube-form">
            <div class="youtube-input-label">
              {{ t("knowledgeBase.youtubeUrlLabel") || "YouTube-URL" }}
            </div>
            <input
              ref="youtubeInputRef"
              v-model="youtubeUrl"
              class="youtube-url-input"
              :placeholder="
                t('knowledgeBase.youtubeUrlPlaceholder') ||
                'https://youtube.com/watch?v=...'
              "
              type="url"
              inputmode="url"
              autofocus
              @keydown.enter="submitYoutubeUrl"
            />
            <div class="youtube-input-tip">
              {{
                t("knowledgeBase.youtubeUrlTip") ||
                "Transkript automatisch erfassen"
              }}
            </div>
            <button
              class="youtube-submit-btn"
              :disabled="!youtubeUrl.trim()"
              @click="submitYoutubeUrl"
            >
              {{ t("knowledgeBase.import") || "Importieren" }}
            </button>
          </div>
        </div>
      </div>
    </transition>
  </teleport>
</template>

<script setup lang="ts">
import { ref, computed, watch, nextTick } from "vue";
import { useI18n } from "vue-i18n";

const props = defineProps<{
  visible: boolean;
}>();

const emit = defineEmits<{
  (e: "close"): void;
  (e: "uploadFile"): void;
  (e: "submitUrl", url: string): void;
  (e: "submitYoutubeUrl", url: string): void;
}>();

const { t } = useI18n();

const page = ref<"menu" | "url" | "youtube">("menu");
const youtubeUrl = ref("");
const urlValue = ref("");
const youtubeInputRef = ref<HTMLInputElement | null>(null);
const urlInputRef = ref<HTMLInputElement | null>(null);

const pageTitle = computed(() => {
  if (page.value === "youtube") {
    return t("knowledgeBase.importYouTube") || "YouTube importieren";
  }
  if (page.value === "url") {
    return t("knowledgeBase.importFromURL") || "URL importieren";
  }
  return t("knowledgeBase.addKnowledge") || "Wissen hinzufügen";
});

// Reset to menu whenever the sheet opens
watch(
  () => props.visible,
  (val) => {
    if (val) {
      page.value = "menu";
      youtubeUrl.value = "";
      urlValue.value = "";
    }
  },
);

const goBackOrClose = () => {
  if (page.value !== "menu") {
    page.value = "menu";
  } else {
    emit("close");
  }
};

const submitUrl = () => {
  const u = urlValue.value.trim();
  if (!u) return;
  emit("submitUrl", u);
};

const submitYoutubeUrl = () => {
  const url = youtubeUrl.value.trim();
  if (!url) return;
  emit("submitYoutubeUrl", url);
};
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

/* YouTube URL input page */
.upload-sheet-youtube-form {
  padding: 20px;
}

.youtube-input-label {
  font-size: 14px;
  font-weight: 500;
  color: var(--td-text-color-primary, #333);
  margin-bottom: 8px;
}

.youtube-url-input {
  width: 100%;
  height: 44px;
  padding: 0 14px;
  font-size: 15px;
  border: 1px solid var(--td-border-level-2-color, #ddd);
  border-radius: 8px;
  outline: none;
  background: var(--td-bg-color-secondarycontainer, #f5f5f5);
  color: var(--td-text-color-primary, #333);
  box-sizing: border-box;
  -webkit-appearance: none;

  &:focus {
    border-color: var(--td-brand-color, #0052d9);
    background: var(--td-bg-color-container, #fff);
  }

  &::placeholder {
    color: var(--td-text-color-placeholder, #bbb);
  }
}

.youtube-input-tip {
  font-size: 12px;
  color: var(--td-text-color-secondary, #999);
  margin-top: 6px;
  margin-bottom: 16px;
}

.youtube-submit-btn {
  width: 100%;
  height: 44px;
  border: none;
  border-radius: 8px;
  font-size: 16px;
  font-weight: 500;
  background: var(--td-brand-color, #0052d9);
  color: #fff;
  cursor: pointer;
  -webkit-tap-highlight-color: transparent;

  &:disabled {
    opacity: 0.4;
    cursor: not-allowed;
  }

  &:active:not(:disabled) {
    opacity: 0.85;
  }
}

/* Back button */
.upload-sheet-header-left {
  min-width: 28px;
}

.upload-sheet-back {
  cursor: pointer;
  padding: 4px;
  color: var(--td-text-color-secondary, #999);
  -webkit-tap-highlight-color: transparent;

  &:active {
    opacity: 0.6;
  }
}
</style>
