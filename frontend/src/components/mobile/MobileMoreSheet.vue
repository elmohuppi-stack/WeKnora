<template>
  <teleport to="body">
    <transition name="more-sheet">
      <div
        v-if="visible"
        class="more-sheet-overlay"
        @click.self="emit('close')"
      >
        <div class="more-sheet-panel">
          <div class="more-sheet-header">
            <span class="more-sheet-title">{{ t("menu.more") || "Mehr" }}</span>
            <div class="more-sheet-close" @click="emit('close')">
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
          <div class="more-sheet-items">
            <div
              class="more-sheet-item"
              @click="navigateTo('/platform/agents')"
            >
              <div class="more-sheet-item-icon">
                <svg
                  width="22"
                  height="22"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="2"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                >
                  <path
                    d="M12 2a4 4 0 0 0-4 4v1a4 4 0 0 0 8 0V6a4 4 0 0 0-4-4z"
                  />
                  <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2" />
                </svg>
              </div>
              <span>{{ t("menu.agents") || "Agents" }}</span>
            </div>
            <div
              class="more-sheet-item"
              @click="navigateTo('/platform/organizations')"
            >
              <div class="more-sheet-item-icon">
                <svg
                  width="22"
                  height="22"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="2"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                >
                  <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2" />
                  <circle cx="9" cy="7" r="4" />
                  <path d="M23 21v-2a4 4 0 0 0-3-3.87" />
                  <path d="M16 3.13a4 4 0 0 1 0 7.75" />
                </svg>
              </div>
              <span>{{ t("menu.organizations") || "Organisationen" }}</span>
            </div>
            <div class="more-sheet-item" @click="navigateToSettings">
              <div class="more-sheet-item-icon">
                <svg
                  width="22"
                  height="22"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="2"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                >
                  <circle cx="12" cy="12" r="3" />
                  <path
                    d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06A1.65 1.65 0 0 0 4.68 15a1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06A1.65 1.65 0 0 0 9 4.68a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"
                  />
                </svg>
              </div>
              <span>{{ t("menu.settings") || "Einstellungen" }}</span>
            </div>
            <div class="more-sheet-divider" />
            <div
              class="more-sheet-item more-sheet-item--danger"
              @click="handleLogout"
            >
              <div class="more-sheet-item-icon">
                <svg
                  width="22"
                  height="22"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="2"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                >
                  <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4" />
                  <polyline points="16 17 21 12 16 7" />
                  <line x1="21" y1="12" x2="9" y2="12" />
                </svg>
              </div>
              <span>{{ t("menu.logout") || "Abmelden" }}</span>
            </div>
          </div>
        </div>
      </div>
    </transition>
  </teleport>
</template>

<script setup lang="ts">
import { useRouter } from "vue-router";
import { useUIStore } from "@/stores/ui";
import { useI18n } from "vue-i18n";

const props = defineProps<{
  visible: boolean;
}>();

const emit = defineEmits<{
  (e: "close"): void;
}>();

const router = useRouter();
const uiStore = useUIStore();
const { t } = useI18n();

const navigateTo = (path: string) => {
  emit("close");
  router.push(path);
};

const navigateToSettings = () => {
  emit("close");
  uiStore.openSettings();
};

const handleLogout = () => {
  emit("close");
  // Will be handled by the auth store/logic
  // The menu.vue handles logout via its own logic
  router.push("/login");
};
</script>

<style scoped lang="less">
.more-sheet-overlay {
  position: fixed;
  inset: 0;
  z-index: 2000;
  background: rgba(0, 0, 0, 0.4);
  display: flex;
  align-items: flex-end;
  justify-content: center;
}

.more-sheet-panel {
  width: 100%;
  max-width: 500px;
  background: var(--td-bg-color-container, #fff);
  border-radius: 16px 16px 0 0;
  padding-bottom: env(safe-area-inset-bottom, 16px);
  animation: slide-up 0.25s ease-out;
}

.more-sheet-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px 12px;
  border-bottom: 1px solid var(--td-border-level-1-color, #f0f0f0);
}

.more-sheet-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--td-text-color-primary, #333);
}

.more-sheet-close {
  cursor: pointer;
  padding: 4px;
  color: var(--td-text-color-secondary, #999);
  -webkit-tap-highlight-color: transparent;

  &:active {
    opacity: 0.6;
  }
}

.more-sheet-items {
  padding: 8px 0;
}

.more-sheet-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 20px;
  font-size: 15px;
  color: var(--td-text-color-primary, #333);
  cursor: pointer;
  -webkit-tap-highlight-color: transparent;
  touch-action: manipulation;
  transition: background 0.15s;

  &:active {
    background: var(--td-bg-color-container-hover, #f5f5f5);
  }

  &--danger {
    color: var(--td-error-color, #e34d59);
  }
}

.more-sheet-item-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 22px;
  height: 22px;
  flex-shrink: 0;
}

.more-sheet-divider {
  height: 1px;
  margin: 4px 20px;
  background: var(--td-border-level-1-color, #f0f0f0);
}

@keyframes slide-up {
  from {
    transform: translateY(100%);
  }
  to {
    transform: translateY(0);
  }
}

.more-sheet-enter-active,
.more-sheet-leave-active {
  transition: opacity 0.2s ease;

  .more-sheet-panel {
    transition: transform 0.25s ease;
  }
}

.more-sheet-enter-from,
.more-sheet-leave-to {
  opacity: 0;

  .more-sheet-panel {
    transform: translateY(100%);
  }
}
</style>
