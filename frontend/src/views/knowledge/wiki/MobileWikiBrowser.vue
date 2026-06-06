<template>
  <div class="mobile-wiki">
    <!-- Page list view -->
    <template v-if="!selectedPage">
      <div class="mobile-wiki-header">
        <div class="mobile-wiki-search">
          <t-input
            v-model="searchQuery"
            :placeholder="
              $t('knowledgeEditor.wikiBrowser.searchPlaceholder') ||
              'Wiki durchsuchen...'
            "
            size="small"
            clearable
            @enter="handleSearch"
          >
            <template #prefix-icon><t-icon name="search" /></template>
          </t-input>
        </div>
        <div class="mobile-wiki-filter-row">
          <t-radio-group
            v-model="activeTypeFilter"
            variant="default"
            size="small"
          >
            <t-radio-button value="all">{{
              $t("common.all") || "Alle"
            }}</t-radio-button>
            <t-radio-button value="summary">{{
              $t("knowledgeEditor.wikiBrowser.filterSummary") ||
              "Zusammenfassung"
            }}</t-radio-button>
            <t-radio-button value="entity">{{
              $t("knowledgeEditor.wikiBrowser.filterEntity") || "Entität"
            }}</t-radio-button>
            <t-radio-button value="concept">{{
              $t("knowledgeEditor.wikiBrowser.filterConcept") || "Konzept"
            }}</t-radio-button>
          </t-radio-group>
        </div>
      </div>

      <div class="mobile-wiki-list" ref="listRef">
        <div v-if="pagesLoading" class="mobile-wiki-loading">
          <t-loading size="small" />
        </div>
        <template v-else-if="filteredPages.length > 0">
          <div
            v-for="page in filteredPages"
            :key="page.slug"
            class="mobile-wiki-list-item"
            @click="selectPage(page)"
          >
            <div class="wiki-item-title">{{ page.title }}</div>
            <div class="wiki-item-meta">
              <span class="wiki-item-type" :class="'type-' + page.page_type">
                {{ pageTypeLabel(page.page_type) }}
              </span>
              <span v-if="page.summary" class="wiki-item-summary">{{
                page.summary
              }}</span>
            </div>
          </div>
        </template>
        <div v-else class="mobile-wiki-empty">
          <t-icon name="file" size="32px" />
          <p>
            {{
              $t("knowledgeEditor.wikiBrowser.noResults") ||
              "Keine Seiten gefunden"
            }}
          </p>
        </div>
      </div>
    </template>

    <!-- Page reader view -->
    <template v-else>
      <div class="mobile-wiki-reader">
        <div class="mobile-wiki-reader-header">
          <button class="mobile-wiki-back-btn" @click="selectedPage = null">
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
              <polyline points="15 18 9 12 15 6" />
            </svg>
          </button>
          <h1 class="mobile-wiki-reader-title">{{ selectedPage.title }}</h1>
        </div>
        <div
          class="mobile-wiki-reader-content"
          ref="contentRef"
          @click="handleContentClick"
        >
          <div
            v-if="selectedPage.in_links && selectedPage.in_links.length"
            class="mobile-wiki-backlinks"
          >
            <span class="backlinks-label"
              >{{
                $t("knowledgeEditor.wikiBrowser.backlinks") || "Rückverweise"
              }}:</span
            >
            <span
              v-for="link in selectedPage.in_links.slice(0, 5)"
              :key="link"
              class="backlink-item"
              @click="handleBacklinkClick(link)"
              >{{ link }}</span
            >
            <span
              v-if="selectedPage.in_links.length > 5"
              class="backlinks-more"
            >
              +{{ selectedPage.in_links.length - 5 }}
              {{ $t("knowledgeEditor.wikiBrowser.more") || "mehr" }}
            </span>
          </div>
          <div class="mobile-wiki-rendered" v-html="renderedContent"></div>
          <div
            v-if="selectedPage.source_refs && selectedPage.source_refs.length"
            class="mobile-wiki-sources"
          >
            <div class="sources-label">
              {{ $t("knowledgeEditor.wikiBrowser.sources") || "Quellen" }}:
            </div>
            <div
              v-for="ref in selectedPage.source_refs"
              :key="ref"
              class="source-item source-item--clickable"
              @click="emit('open-source-doc', ref)"
            >
              <t-icon name="file" size="14px" />
              <span>{{ ref }}</span>
            </div>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, nextTick } from "vue";
import { useI18n } from "vue-i18n";
import {
  listWikiPages,
  getWikiPage,
  searchWikiPages,
  type WikiPage,
} from "@/api/wiki";
import { MessagePlugin } from "tdesign-vue-next";

const props = defineProps<{
  knowledgeBaseId: string;
}>();

const emit = defineEmits<{
  (e: "open-source-doc", knowledgeId: string): void;
}>();

const { t } = useI18n();

// State
const pages = ref<WikiPage[]>([]);
const selectedPage = ref<WikiPage | null>(null);
const pagesLoading = ref(false);
const searchQuery = ref("");
const activeTypeFilter = ref("all");
const listRef = ref<HTMLElement | null>(null);
const contentRef = ref<HTMLElement | null>(null);
const currentPage = ref(1);
const hasMore = ref(true);
const loadingMore = ref(false);
const PAGE_SIZE = 50;

// Computed
const filteredPages = computed(() => {
  let result = pages.value;
  if (activeTypeFilter.value !== "all") {
    result = result.filter((p) => p.page_type === activeTypeFilter.value);
  }
  return result;
});

// Methods
const pageTypeLabel = (type: string): string => {
  const labels: Record<string, string> = {
    summary: t("knowledgeEditor.wikiBrowser.filterSummary") || "Summary",
    entity: t("knowledgeEditor.wikiBrowser.filterEntity") || "Entity",
    concept: t("knowledgeEditor.wikiBrowser.filterConcept") || "Concept",
    synthesis: t("knowledgeEditor.wikiBrowser.filterSynthesis") || "Synthesis",
    comparison:
      t("knowledgeEditor.wikiBrowser.filterComparison") || "Comparison",
  };
  return labels[type] || type;
};

const selectPage = async (page: WikiPage) => {
  try {
    const res: any = await getWikiPage(props.knowledgeBaseId, page.slug);
    // Axios-Interceptor gibt bereits data zurück → res ist das Page-Objekt
    selectedPage.value = res?.data || res || page;
  } catch {
    selectedPage.value = page;
  }
};

const renderedContent = computed(() => {
  if (!selectedPage.value?.content) return "";
  return renderMarkdown(selectedPage.value.content);
});

const renderMarkdown = (content: string): string => {
  // Convert wiki links [[slug|title]] or [[slug]] to clickable links
  let html = content.replace(
    /\[\[([^\]|]+)(?:\|([^\]]+))?\]\]/g,
    (_match, slug, title) => {
      const display = title || slug;
      return `<a href="javascript:;" class="wiki-link" data-slug="${slug}">${display}</a>`;
    },
  );
  // Simple markdown rendering (basic)
  html = html
    .replace(/^### (.+)$/gm, "<h3>$1</h3>")
    .replace(/^## (.+)$/gm, "<h2>$1</h2>")
    .replace(/^# (.+)$/gm, "<h1>$1</h1>")
    .replace(/\*\*(.+?)\*\*/g, "<strong>$1</strong>")
    .replace(/\*(.+?)\*/g, "<em>$1</em>")
    .replace(/`(.+?)`/g, "<code>$1</code>")
    .replace(/^- (.+)$/gm, "<li>$1</li>")
    .replace(/(<li>.*<\/li>\n?)+/g, "<ul>$&</ul>")
    .replace(/\n\n/g, "</p><p>")
    .replace(/^(?!<[hul])/gm, "<p>")
    .replace(/(.+)$/gm, (match) => {
      if (match.startsWith("<")) return match;
      return match;
    });
  return `<p>${html}</p>`;
};

const handleContentClick = async (e: MouseEvent) => {
  const target = e.target as HTMLElement;
  if (target.classList.contains("wiki-link")) {
    const slug = target.getAttribute("data-slug");
    if (slug) {
      try {
        const res: any = await getWikiPage(props.knowledgeBaseId, slug);
        const page = res?.data || res;
        if (page) {
          selectedPage.value = page;
          nextTick(() => {
            contentRef.value?.scrollTo(0, 0);
          });
        } else {
          MessagePlugin.warning(
            t("knowledgeEditor.wikiBrowser.pageNotFound") ||
              "Seite nicht gefunden",
          );
        }
      } catch {
        MessagePlugin.warning(
          t("knowledgeEditor.wikiBrowser.pageNotFound") ||
            "Seite nicht gefunden",
        );
      }
    }
  }
};

const handleBacklinkClick = async (slug: string) => {
  try {
    const res: any = await getWikiPage(props.knowledgeBaseId, slug);
    const page = res?.data || res;
    if (page) {
      selectedPage.value = page;
      nextTick(() => {
        contentRef.value?.scrollTo(0, 0);
      });
    }
  } catch {
    // Silently fail for backlinks
  }
};

const handleSearch = async () => {
  if (!searchQuery.value.trim()) {
    loadPages();
    return;
  }
  pagesLoading.value = true;
  try {
    const res: any = await searchWikiPages(
      props.knowledgeBaseId,
      searchQuery.value.trim(),
      PAGE_SIZE,
    );
    const hits: WikiPage[] = res?.data?.pages || res?.pages || [];
    pages.value = hits;
    hasMore.value = hits.length >= PAGE_SIZE;
  } catch {
    pages.value = [];
  } finally {
    pagesLoading.value = false;
  }
};

const loadPages = async (loadMore = false) => {
  if (pagesLoading.value || loadingMore.value) return;
  if (loadMore) {
    if (!hasMore.value) return;
    loadingMore.value = true;
  } else {
    pagesLoading.value = true;
    currentPage.value = 1;
  }

  const params: Record<string, any> = {
    page: currentPage.value,
    page_size: PAGE_SIZE,
    sort_by: "updated_at",
    sort_order: "desc",
  };

  // API-Seitig filtern, wenn ein bestimmter Typ gewählt ist
  if (activeTypeFilter.value !== "all") {
    params.page_type = activeTypeFilter.value;
  }

  try {
    const res: any = await listWikiPages(props.knowledgeBaseId, params);
    const newPages: WikiPage[] = res?.data?.pages || res?.pages || [];
    if (loadMore) {
      pages.value = [...pages.value, ...newPages];
    } else {
      pages.value = newPages;
    }
    hasMore.value = newPages.length === PAGE_SIZE;
  } catch {
    if (!loadMore) pages.value = [];
  } finally {
    pagesLoading.value = false;
    loadingMore.value = false;
  }
};

// Infinite scroll
const handleScroll = () => {
  if (!listRef.value) return;
  const { scrollTop, scrollHeight, clientHeight } = listRef.value;
  if (
    scrollHeight - scrollTop - clientHeight < 100 &&
    hasMore.value &&
    !loadingMore.value
  ) {
    currentPage.value++;
    loadPages(true);
  }
};

watch(activeTypeFilter, () => {
  // Bei Typ-Wechsel neu von der API laden (mit page_type-Filter)
  loadPages();
});

onMounted(() => {
  loadPages();
  if (listRef.value) {
    listRef.value.addEventListener("scroll", handleScroll);
  }
});
</script>

<style scoped lang="less">
.mobile-wiki {
  display: flex;
  flex-direction: column;
  height: 100%;
  overflow: hidden;
}

.mobile-wiki-header {
  padding: 12px 12px 8px;
  display: flex;
  flex-direction: column;
  gap: 8px;
  flex-shrink: 0;
  background: var(--td-bg-color-page, #f5f5f5);
}

.mobile-wiki-search {
  width: 100%;
}

.mobile-wiki-filter-row {
  overflow-x: auto;
  -webkit-overflow-scrolling: touch;
  white-space: nowrap;
  padding-bottom: 4px;

  &::-webkit-scrollbar {
    display: none;
  }
}

.mobile-wiki-list {
  flex: 1;
  overflow-y: auto;
  padding: 0 12px 12px;
  -webkit-overflow-scrolling: touch;
}

.mobile-wiki-loading {
  display: flex;
  justify-content: center;
  padding: 32px 0;
}

.mobile-wiki-list-item {
  padding: 12px 14px;
  margin-bottom: 8px;
  background: var(--td-bg-color-container, #fff);
  border-radius: 10px;
  cursor: pointer;
  -webkit-tap-highlight-color: transparent;
  touch-action: manipulation;
  transition: background 0.15s;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);

  &:active {
    background: var(--td-bg-color-container-hover, #f5f5f5);
  }
}

.wiki-item-title {
  font-size: 15px;
  font-weight: 500;
  color: var(--td-text-color-primary, #333);
  margin-bottom: 4px;
  line-height: 1.4;
}

.wiki-item-meta {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.wiki-item-type {
  font-size: 11px;
  padding: 1px 6px;
  border-radius: 4px;
  font-weight: 500;
  flex-shrink: 0;

  &.type-summary {
    background: #e8eaff;
    color: #0052d9;
  }
  &.type-entity {
    background: #e8f5e9;
    color: #2ba471;
  }
  &.type-concept {
    background: #fff3e0;
    color: #e37318;
  }
  &.type-synthesis {
    background: #e3f2fd;
    color: #0594fa;
  }
  &.type-comparison {
    background: #fce4ec;
    color: #d54941;
  }
}

.wiki-item-summary {
  font-size: 12px;
  color: var(--td-text-color-secondary, #999);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  flex: 1;
  min-width: 0;
}

.mobile-wiki-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 48px 0;
  color: var(--td-text-color-secondary, #999);
  gap: 12px;

  p {
    font-size: 14px;
    margin: 0;
  }
}

/* Reader view */
.mobile-wiki-reader {
  display: flex;
  flex-direction: column;
  height: 100%;
  overflow: hidden;
}

.mobile-wiki-reader-header {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  border-bottom: 1px solid var(--td-border-level-1-color, #f0f0f0);
  background: var(--td-bg-color-container, #fff);
  flex-shrink: 0;
}

.mobile-wiki-back-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border: none;
  background: transparent;
  cursor: pointer;
  color: var(--td-text-color-primary, #333);
  -webkit-tap-highlight-color: transparent;
  border-radius: 50%;

  &:active {
    background: var(--td-bg-color-container-hover, #f5f5f5);
  }
}

.mobile-wiki-reader-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--td-text-color-primary, #333);
  margin: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  flex: 1;
}

.mobile-wiki-reader-content {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
  -webkit-overflow-scrolling: touch;
  line-height: 1.7;
  font-size: 15px;
  color: var(--td-text-color-primary, #333);

  h1,
  h2,
  h3 {
    margin: 16px 0 8px;
    font-weight: 600;
    line-height: 1.3;
  }

  h1 {
    font-size: 20px;
  }
  h2 {
    font-size: 17px;
  }
  h3 {
    font-size: 15px;
  }

  p {
    margin: 0 0 12px;
  }

  ul {
    padding-left: 20px;
    margin: 8px 0;

    li {
      margin-bottom: 4px;
    }
  }

  code {
    font-size: 13px;
    background: var(--td-bg-color-component, #f0f0f0);
    padding: 1px 4px;
    border-radius: 3px;
  }

  .wiki-link {
    color: var(--td-brand-color, #0052d9);
    text-decoration: none;
    font-weight: 500;

    &:hover {
      text-decoration: underline;
    }
  }
}

.mobile-wiki-backlinks {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
  padding: 8px 12px;
  margin-bottom: 12px;
  background: var(--td-brand-color-light, #f2f3ff);
  border-radius: 8px;
  font-size: 12px;
}

.backlinks-label {
  color: var(--td-text-color-secondary, #999);
  flex-shrink: 0;
}

.backlink-item {
  color: var(--td-brand-color, #0052d9);
  cursor: pointer;
  -webkit-tap-highlight-color: transparent;

  &:active {
    opacity: 0.7;
  }
}

.backlinks-more {
  color: var(--td-text-color-secondary, #999);
}

.mobile-wiki-sources {
  margin-top: 20px;
  padding-top: 12px;
  border-top: 1px solid var(--td-border-level-1-color, #f0f0f0);
}

.sources-label {
  font-size: 12px;
  font-weight: 600;
  color: var(--td-text-color-secondary, #999);
  margin-bottom: 6px;
}

.source-item {
  font-size: 12px;
  color: var(--td-text-color-secondary, #999);
  padding: 2px 0;
}

.source-item--clickable {
  display: flex;
  align-items: center;
  gap: 4px;
  cursor: pointer;
  color: var(--td-brand-color, #0052d9);
  -webkit-tap-highlight-color: transparent;
  touch-action: manipulation;

  &:active {
    opacity: 0.7;
  }
}
</style>
