import { ref, onMounted, onUnmounted, type Ref } from 'vue'

/**
 * Composable for detecting when an element enters the viewport.
 * Used to trigger scroll-based entrance animations.
 */
export function useScrollReveal(threshold = 0.15): {
  elementRef: Ref<HTMLElement | null>
  isVisible: Ref<boolean>
} {
  const elementRef = ref<HTMLElement | null>(null)
  const isVisible = ref(false)
  let observer: IntersectionObserver | null = null

  onMounted(() => {
    if (!elementRef.value) return

    observer = new IntersectionObserver(
      ([entry]) => {
        if (entry && entry.isIntersecting) {
          isVisible.value = true
          // Once visible, stop observing (animate only once)
          observer?.unobserve(entry.target)
        }
      },
      { threshold },
    )

    observer.observe(elementRef.value)
  })

  onUnmounted(() => {
    observer?.disconnect()
  })

  return { elementRef, isVisible }
}
