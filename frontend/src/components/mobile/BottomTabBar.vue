<template>
  <div class="bottom-tab-bar">
    <div
      v-for="tab in tabs"
      :key="tab.key"
      class="tab-item"
      :class="{ active: activeTab === tab.key }"
      @click="onTabClick(tab)"
    >
      <div class="tab-icon" v-html="tab.icon"></div>
      <span class="tab-label">{{ tab.label }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, inject } from "vue";
import { useRouter, useRoute } from "vue-router";
import { useUIStore } from "@/stores/ui";
import { useI18n } from "vue-i18n";

const uiStore = useUIStore();
const router = useRouter();
const route = useRoute();
const { t } = useI18n();

const openMoreSheet = inject<() => void>("mobile:openMoreSheet", () => {});

const activeTab = computed(() => uiStore.mobileActiveTab);

interface TabItem {
  key: "chat" | "knowledge" | "wiki" | "more";
  label: string;
  icon: string; // SVG inline
  route?: string;
  action?: () => void;
}

const tabs = computed<TabItem[]>(() => [
  {
    key: "chat",
    label: t("menu.chat"),
    icon: '<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/></svg>',
    route: "/platform/chat",
  },
  {
    key: "knowledge",
    label: t("menu.knowledgeBase"),
    icon: '<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M4 19.5A2.5 2.5 0 0 1 6.5 17H20"/><path d="M6.5 2H20v20H6.5A2.5 2.5 0 0 1 4 19.5v-15A2.5 2.5 0 0 1 6.5 2z"/></svg>',
    route: "/platform/knowledge-bases",
  },
  {
    key: "wiki",
    label: t("knowledgeEditor.wikiBrowser.tabWiki"),
    icon: '<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="3"/><path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06A1.65 1.65 0 0 0 4.68 15a1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06A1.65 1.65 0 0 0 9 4.68a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"/></svg>',
    route: "/platform/knowledge-bases",
  },
  {
    key: "more",
    label: t("menu.settings"),
    icon: '<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="1"/><circle cx="19" cy="12" r="1"/><circle cx="5" cy="12" r="1"/></svg>',
    action: () => openMoreSheet(),
  },
]);

const onTabClick = (tab: TabItem) => {
  uiStore.setMobileActiveTab(tab.key);
  if (tab.action) {
    tab.action();
  } else if (tab.route) {
    // For 'wiki' tab, try to navigate to the current KB's wiki tab
    if (tab.key === "wiki") {
      const kbId = route.params?.kbId || route.query?.kbId;
      if (kbId) {
        router.push(`/platform/knowledge-bases/${kbId}?tab=wiki`);
      } else {
        // If no KB is loaded, go to KB list first
        router.push("/platform/knowledge-bases");
      }
    } else if (tab.key === "chat") {
      // Navigate to the most recent chat or chat creation
      router.push("/platform/creatChat");
    } else if (tab.key === "knowledge") {
      const kbId = route.params?.kbId;
      if (kbId) {
        router.push(`/platform/knowledge-bases/${kbId}`);
      } else {
        router.push("/platform/knowledge-bases");
      }
    }
  }
};
</script>

<style scoped lang="less">
.bottom-tab-bar {
  display: flex;
  align-items: center;
  justify-content: space-around;
  height: 56px;
  padding-bottom: env(safe-area-inset-bottom, 0);
  background: var(--td-bg-color-container, #fff);
  border-top: 1px solid var(--td-border-level-1-color, #e0e0e0);
  flex-shrink: 0;
  z-index: 100;
}

.tab-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 2px;
  flex: 1;
  height: 100%;
  cursor: pointer;
  -webkit-tap-highlight-color: transparent;
  touch-action: manipulation;
  color: var(--td-text-color-secondary, #999);
  transition: color 0.15s ease;

  &.active {
    color: var(--td-brand-color, #0052d9);
  }

  &:active {
    opacity: 0.7;
  }
}

.tab-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;

  :deep(svg) {
    width: 24px;
    height: 24px;
  }
}

.tab-label {
  font-size: 10px;
  line-height: 1;
  white-space: nowrap;
}
</style>
