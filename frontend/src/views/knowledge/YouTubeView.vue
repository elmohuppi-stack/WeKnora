<script setup lang="ts">
import { ref, computed, onMounted, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import { useI18n } from "vue-i18n";
import { marked } from "marked";
import markedKatex from "marked-katex-extension";
import "katex/dist/katex.min.css";
import hljs from "highlight.js";
import "highlight.js/styles/github.css";
import { MessagePlugin, Tabs } from "tdesign-vue-next";
import {
  getKnowledgeDetails,
  getKnowledgeDetailsCon,
} from "@/api/knowledge-base/index";
import { listWikiPages } from "@/api/wiki/index";
import { sanitizeHTML } from "@/utils/security";

const { t } = useI18n();
const route = useRoute();
const router = useRouter();

const kbId = computed(() => route.params.kbId as string);
const knowledgeId = computed(() => route.params.knowledgeId as string);

// --- State ---
const knowledge = ref<any>(null);
const chunks = ref<any[]>([]);
const wikiContent = ref<string>("");
const wikiTitle = ref<string>("");
const loading = ref(true);
const wikiLoading = ref(false);
const activeTab = ref<"transcript" | "wiki">("wiki");
const viewMode = ref<"merged" | "chunks">("merged");

// --- YouTube video ID extraction ---
const getYouTubeVideoId = (url: string): string | null => {
  if (!url) return null;
  const patterns = [
    /(?:youtube\.com\/watch\?v=)([a-zA-Z0-9_-]{11})/,
    /(?:youtu\.be\/)([a-zA-Z0-9_-]{11})/,
    /(?:youtube\.com\/embed\/)([a-zA-Z0-9_-]{11})/,
    /(?:youtube\.com\/shorts\/)([a-zA-Z0-9_-]{11})/,
  ];
  for (const pattern of patterns) {
    const match = url.match(pattern);
    if (match) return match[1];
  }
  return null;
};

const youtubeEmbedUrl = computed(() => {
  if (!knowledge.value?.source) return "";
  const videoId = getYouTubeVideoId(knowledge.value.source);
  if (!videoId) return "";
  return `https://www.youtube.com/embed/${videoId}?rel=0&cc_load_policy=1`;
});

// --- Format metadata ---
const formatDuration = (seconds: number): string => {
  if (!seconds) return "";
  const min = Math.floor(seconds / 60);
  const sec = seconds % 60;
  return `${min}:${sec.toString().padStart(2, "0")}`;
};

const formatDate = (dateStr: string): string => {
  if (!dateStr) return "";
  const d = new Date(dateStr);
  return d.toLocaleDateString();
};

// --- Fetch data ---
const fetchData = async () => {
  loading.value = true;
  try {
    const result: any = await getKnowledgeDetails(knowledgeId.value);
    if (result.success && result.data) {
      knowledge.value = result.data;
    }
  } catch (err) {
    console.error("Failed to load knowledge details:", err);
    MessagePlugin.error(t("knowledgeBase.loadFailed"));
  }
  loading.value = false;
};

const fetchChunks = async () => {
  try {
    const result: any = await getKnowledgeDetailsCon(knowledgeId.value, 1);
    if (result.success && result.data) {
      chunks.value = result.data;
    }
  } catch (err) {
    console.error("Failed to load chunks:", err);
  }
};

const fetchWikiPage = async () => {
  if (!kbId.value) return;
  wikiLoading.value = true;
  try {
    // Try to find a wiki page with source_ref matching this knowledge
    const result: any = await listWikiPages(kbId.value, {
      page_type: "youtube_transcript",
      page_size: 50,
    });
    if (result.success && result.data?.pages) {
      const page = result.data.pages.find((p: any) =>
        p.source_refs?.some((ref: string) => ref.startsWith(knowledgeId.value)),
      );
      if (page) {
        wikiContent.value = page.content || "";
        wikiTitle.value = page.title || "";
      }
    }
  } catch (err) {
    console.error("Failed to load wiki page:", err);
  }
  wikiLoading.value = false;
};

// --- Markdown rendering ---
const renderMarkdown = (content: string): string => {
  if (!content) return "";
  try {
    // Configure marked with KaTeX and highlight.js
    const renderer = new marked.Renderer();
    renderer.code = ({ text, lang }: { text: string; lang?: string }) => {
      const language = lang || "";
      let highlighted = text;
      if (language && hljs.getLanguage(language)) {
        try {
          highlighted = hljs.highlight(text, { language }).value;
        } catch (e) {
          highlighted = text;
        }
      }
      return `<pre><code class="hljs ${language ? `language-${language}` : ""}">${highlighted}</code></pre>`;
    };

    marked.setOptions({
      renderer,
      gfm: true,
      breaks: true,
    });

    let html = marked.parse(content) as string;
    // Process wiki links [[slug|title]] or [[slug]]
    html = html.replace(
      /\[\[([^|[\]]+)(?:\|([^[\]]+))?\]\]/g,
      (_match: string, slug: string, title?: string) => {
        const display = title || slug;
        return `<a class="wiki-content-link" data-slug="${slug}">${display}</a>`;
      },
    );
    return sanitizeHTML(html);
  } catch (e) {
    console.error("Markdown rendering error:", e);
    return content;
  }
};

const renderedWikiContent = computed(() => renderMarkdown(wikiContent.value));

// --- Combined transcript from chunks ---
const mergedContent = computed(() => {
  if (!chunks.value.length) return "";
  return chunks.value.map((chunk: any) => chunk.content || "").join("\n\n");
});

const renderedMergedContent = computed(() => {
  if (!mergedContent.value) return "";
  return renderMarkdown(mergedContent.value);
});

// --- Lifecycle ---
onMounted(() => {
  fetchData();
  fetchChunks();
  fetchWikiPage();
});

// Refresh wiki when tab changes to wiki
watch(activeTab, (tab) => {
  if (tab === "wiki" && !wikiContent.value) {
    fetchWikiPage();
  }
});

// --- Navigation ---
const goBack = () => {
  router.push(`/platform/knowledge-bases/${kbId.value}`);
};
</script>

<template>
  <div class="youtube-view">
    <!-- Top navigation bar -->
    <div class="youtube-nav-bar">
      <a href="#" class="youtube-nav-back" @click.prevent="goBack">
        <t-icon name="arrow-left" size="14px" />
        <span>{{ t("common.back") }}</span>
      </a>
    </div>

    <!-- Main split layout -->
    <div class="youtube-main" v-if="!loading && knowledge">
      <!-- Left panel: Video player + Metadata -->
      <div class="youtube-left">
        <div class="youtube-player-wrap">
          <div class="youtube-player-container">
            <iframe
              v-if="youtubeEmbedUrl"
              :src="youtubeEmbedUrl"
              class="youtube-player-iframe"
              title="YouTube video player"
              frameborder="0"
              allow="
                accelerometer;
                autoplay;
                clipboard-write;
                encrypted-media;
                gyroscope;
                picture-in-picture;
                web-share;
              "
              allowfullscreen
            ></iframe>
            <div v-else class="youtube-player-placeholder">
              <t-icon name="video" size="48px" />
              <p>{{ t("knowledgeBase.videoNotAvailable") }}</p>
            </div>
          </div>
        </div>

        <!-- Metadata -->
        <div class="youtube-metadata">
          <h1 class="youtube-title">
            {{
              knowledge.title ||
              knowledge.file_name ||
              t("knowledgeBase.untitledDocument")
            }}
          </h1>

          <div class="youtube-meta-row" v-if="knowledge.description">
            <span class="youtube-meta-label">{{
              t("knowledgeBase.documentDescription")
            }}</span>
            <p class="youtube-description">{{ knowledge.description }}</p>
          </div>

          <div class="youtube-meta-grid">
            <div
              class="youtube-meta-item"
              v-if="knowledge.metadata?.channel_name"
            >
              <t-icon name="user" size="14px" />
              <span>{{ knowledge.metadata.channel_name }}</span>
            </div>
            <div class="youtube-meta-item" v-if="knowledge.metadata?.duration">
              <t-icon name="time" size="14px" />
              <span>{{ formatDuration(knowledge.metadata.duration) }}</span>
            </div>
            <div
              class="youtube-meta-item"
              v-if="knowledge.metadata?.published_at"
            >
              <t-icon name="calendar" size="14px" />
              <span>{{ formatDate(knowledge.metadata.published_at) }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Right panel: Content (Transcript / Wiki) -->
      <div class="youtube-right">
        <div class="youtube-content-header">
          <div class="youtube-content-tabs">
            <t-button
              size="small"
              :variant="activeTab === 'wiki' ? 'base' : 'outline'"
              :theme="activeTab === 'wiki' ? 'primary' : 'default'"
              @click="activeTab = 'wiki'"
            >
              <t-icon name="article" size="14px" style="margin-right: 4px" />
              {{ t("knowledgeBase.wikiArticle") }}
            </t-button>
            <t-button
              size="small"
              :variant="activeTab === 'transcript' ? 'base' : 'outline'"
              :theme="activeTab === 'transcript' ? 'primary' : 'default'"
              @click="activeTab = 'transcript'"
            >
              <t-icon name="file" size="14px" style="margin-right: 4px" />
              {{ t("knowledgeBase.transcript") }}
            </t-button>
          </div>
        </div>

        <!-- Wiki article -->
        <div v-if="activeTab === 'wiki'" class="youtube-content-body">
          <div v-if="wikiLoading" class="youtube-content-loading">
            <t-loading size="small" />
            <span>{{ t("knowledgeBase.loading") }}</span>
          </div>
          <div v-else-if="wikiContent" class="youtube-wiki-content">
            <h2 class="youtube-wiki-title">{{ wikiTitle }}</h2>
            <div class="youtube-wiki-body" v-html="renderedWikiContent"></div>
          </div>
          <div v-else class="youtube-content-empty">
            <t-icon name="file-unknown" size="36px" />
            <p>{{ t("knowledgeBase.noWikiArticle") }}</p>
          </div>
        </div>

        <!-- Transcript -->
        <div v-if="activeTab === 'transcript'" class="youtube-content-body">
          <div class="youtube-transcript-tabs">
            <t-button
              size="small"
              :variant="viewMode === 'merged' ? 'base' : 'outline'"
              :theme="viewMode === 'merged' ? 'primary' : 'default'"
              @click="viewMode = 'merged'"
            >
              {{ t("knowledgeBase.viewMerged") }}
            </t-button>
            <t-button
              size="small"
              :variant="viewMode === 'chunks' ? 'base' : 'outline'"
              :theme="viewMode === 'chunks' ? 'primary' : 'default'"
              @click="viewMode = 'chunks'"
            >
              {{ t("knowledgeBase.viewChunks") }}
            </t-button>
          </div>

          <!-- Merged view -->
          <div v-if="viewMode === 'merged'" class="youtube-transcript-body">
            <div v-if="!mergedContent" class="youtube-content-empty">
              <p>{{ t("common.noData") }}</p>
            </div>
            <div
              v-else
              class="youtube-markdown"
              v-html="renderedMergedContent"
            ></div>
          </div>

          <!-- Chunks view -->
          <div v-else class="youtube-transcript-body">
            <div v-if="!chunks.length" class="youtube-content-empty">
              <p>{{ t("common.noData") }}</p>
            </div>
            <div v-else class="youtube-chunk-list">
              <div
                v-for="(chunk, index) in chunks"
                :key="index"
                class="youtube-chunk-item"
              >
                <div class="youtube-chunk-header">
                  <span class="youtube-chunk-index">
                    {{ t("knowledgeBase.segment") }} {{ index + 1 }}
                  </span>
                </div>
                <div
                  class="youtube-chunk-content"
                  v-html="renderMarkdown(chunk.content || '')"
                ></div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Loading state -->
    <div v-else-if="loading" class="youtube-loading">
      <t-loading />
    </div>
  </div>
</template>

<style scoped lang="less">
.youtube-view {
  height: 100vh;
  display: flex;
  flex-direction: column;
  background: var(--td-bg-color-page);
  overflow: hidden;
}

.youtube-nav-bar {
  flex-shrink: 0;
  padding: 12px 24px;
  background: var(--td-bg-color-container);
  border-bottom: 1px solid var(--td-component-border);
  display: flex;
  align-items: center;

  .youtube-nav-back {
    display: flex;
    align-items: center;
    gap: 4px;
    color: var(--td-text-color-secondary);
    text-decoration: none;
    font-size: 14px;

    &:hover {
      color: var(--td-brand-color);
    }
  }
}

.youtube-main {
  flex: 1;
  display: flex;
  overflow: hidden;
}

// Left panel
.youtube-left {
  width: 50%;
  min-width: 400px;
  overflow-y: auto;
  border-right: 1px solid var(--td-component-border);
  background: var(--td-bg-color-container);
}

.youtube-player-wrap {
  padding: 24px;
  position: sticky;
  top: 0;
  background: var(--td-bg-color-container);
  z-index: 1;
}

.youtube-player-container {
  position: relative;
  width: 100%;
  padding-bottom: 56.25%;
  height: 0;
  overflow: hidden;
  border-radius: 8px;
  background: #000;
}

.youtube-player-iframe {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  border: none;
}

.youtube-player-placeholder {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: var(--td-text-color-disabled);

  p {
    margin-top: 8px;
    font-size: 14px;
  }
}

.youtube-metadata {
  padding: 0 24px 24px;
}

.youtube-title {
  font-size: 20px;
  font-weight: 600;
  margin: 0 0 12px;
  line-height: 1.4;
  color: var(--td-text-color-primary);
}

.youtube-meta-row {
  margin-bottom: 16px;

  .youtube-meta-label {
    font-size: 13px;
    font-weight: 600;
    color: var(--td-text-color-secondary);
  }

  .youtube-description {
    margin: 6px 0 0;
    font-size: 13px;
    line-height: 1.6;
    color: var(--td-text-color-primary);
    white-space: pre-wrap;
    word-break: break-word;
  }
}

.youtube-meta-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
  padding-top: 12px;
  border-top: 1px solid var(--td-component-border);
}

.youtube-meta-item {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 13px;
  color: var(--td-text-color-secondary);
}

// Right panel
.youtube-right {
  width: 50%;
  min-width: 400px;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.youtube-content-header {
  flex-shrink: 0;
  padding: 16px 24px;
  background: var(--td-bg-color-container);
  border-bottom: 1px solid var(--td-component-border);
}

.youtube-content-tabs {
  display: flex;
  gap: 8px;
}

.youtube-content-body {
  flex: 1;
  overflow-y: auto;
  padding: 24px;
}

.youtube-content-loading {
  display: flex;
  align-items: center;
  gap: 8px;
  color: var(--td-text-color-placeholder);
  font-size: 13px;
}

.youtube-content-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 48px 0;
  color: var(--td-text-color-disabled);

  p {
    margin-top: 12px;
    font-size: 14px;
  }
}

// Wiki content
.youtube-wiki-content {
  max-width: 800px;
}

.youtube-wiki-title {
  font-size: 24px;
  font-weight: 600;
  margin: 0 0 20px;
  line-height: 1.4;
}

.youtube-wiki-body {
  line-height: 1.8;
  font-size: 15px;
  color: var(--td-text-color-primary);
  word-break: break-word;

  :deep(h2) {
    margin-top: 28px;
    margin-bottom: 12px;
    font-size: 20px;
    font-weight: 600;
    border-bottom: 1px solid var(--td-component-border);
    padding-bottom: 8px;
  }

  :deep(h3) {
    margin-top: 24px;
    margin-bottom: 8px;
    font-size: 17px;
    font-weight: 600;
  }

  :deep(p) {
    margin: 0 0 12px;
  }

  :deep(ul),
  :deep(ol) {
    margin: 0 0 12px;
    padding-left: 24px;
  }

  :deep(blockquote) {
    margin: 0 0 12px;
    padding: 8px 16px;
    border-left: 4px solid var(--td-brand-color);
    background: var(--td-bg-color-container-hover);
    border-radius: 4px;
    color: var(--td-text-color-secondary);
  }

  :deep(code) {
    padding: 2px 6px;
    border-radius: 3px;
    background: var(--td-bg-color-container-hover);
    font-size: 0.9em;
  }

  :deep(pre) {
    margin: 0 0 16px;
    padding: 16px;
    border-radius: 6px;
    background: var(--td-bg-color-container-hover);
    overflow-x: auto;

    code {
      padding: 0;
      background: none;
    }
  }

  :deep(table) {
    width: 100%;
    border-collapse: collapse;
    margin: 0 0 16px;

    th,
    td {
      padding: 8px 12px;
      border: 1px solid var(--td-component-border);
      text-align: left;
    }

    th {
      background: var(--td-bg-color-container-hover);
      font-weight: 600;
    }
  }
}

// Transcript
.youtube-transcript-tabs {
  display: flex;
  gap: 8px;
  margin-bottom: 16px;
}

.youtube-transcript-body {
  max-width: 800px;
}

.youtube-markdown {
  line-height: 1.8;
  font-size: 15px;
  color: var(--td-text-color-primary);
  word-break: break-word;
}

.youtube-chunk-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.youtube-chunk-item {
  border: 1px solid var(--td-component-border);
  border-radius: 6px;
  overflow: hidden;

  .youtube-chunk-header {
    padding: 8px 12px;
    background: var(--td-bg-color-container-hover);
    border-bottom: 1px solid var(--td-component-border);
  }

  .youtube-chunk-index {
    font-size: 12px;
    font-weight: 600;
    color: var(--td-text-color-placeholder);
  }

  .youtube-chunk-content {
    padding: 12px;
    font-size: 14px;
    line-height: 1.6;
    color: var(--td-text-color-primary);
    word-break: break-word;
  }
}

// Loading
.youtube-loading {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
}
</style>
