<template>
  <div class="mobile-layout">
    <div class="mobile-content">
      <RouterView />
    </div>
    <BottomTabBar />
    <MobileMoreSheet
      :visible="moreSheetVisible"
      @close="moreSheetVisible = false"
    />
    <GlobalCommandPalette />
    <GlobalInvitationBell />
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import BottomTabBar from "@/components/mobile/BottomTabBar.vue";
import MobileMoreSheet from "@/components/mobile/MobileMoreSheet.vue";
import GlobalCommandPalette from "@/components/GlobalCommandPalette.vue";
import GlobalInvitationBell from "@/components/GlobalInvitationBell.vue";

const moreSheetVisible = ref(false);

const openMoreSheet = () => {
  moreSheetVisible.value = true;
};

// Provide openMoreSheet so BottomTabBar can call it
import { provide } from "vue";
provide("mobile:openMoreSheet", openMoreSheet);
</script>

<style scoped lang="less">
.mobile-layout {
  display: flex;
  flex-direction: column;
  height: 100dvh;
  width: 100%;
  overflow: hidden;
  background: var(--td-bg-color-page, #f5f5f5);
}

.mobile-content {
  flex: 1;
  min-height: 0;
  overflow-y: auto;
  overflow-x: hidden;
  -webkit-overflow-scrolling: touch;
}
</style>
