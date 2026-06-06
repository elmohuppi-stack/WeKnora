import { ref, computed, onMounted, onUnmounted } from "vue";
import { useUIStore } from "@/stores/ui";

/**
 * Composable for mobile viewport detection.
 * Sets uiStore.isMobile based on a breakpoint (default 767px).
 * Handles resize events automatically and cleans up on unmount.
 */
export function useMobile(breakpoint = 767) {
  const uiStore = useUIStore();
  const mql = ref<MediaQueryList | null>(null);

  const onChange = (e: MediaQueryListEvent | MediaQueryList) => {
    uiStore.setMobile(e.matches);
  };

  onMounted(() => {
    mql.value = window.matchMedia(`(max-width: ${breakpoint}px)`);
    // Set initial value
    onChange(mql.value);
    // Listen for changes
    mql.value.addEventListener(
      "change",
      onChange as (e: MediaQueryListEvent) => void,
    );
  });

  onUnmounted(() => {
    if (mql.value) {
      mql.value.removeEventListener(
        "change",
        onChange as (e: MediaQueryListEvent) => void,
      );
    }
  });

  return {
    isMobile: computed(() => uiStore.isMobile),
  };
}
