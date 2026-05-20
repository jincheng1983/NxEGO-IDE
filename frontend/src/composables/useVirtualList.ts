import { ref, computed, onMounted, onUnmounted, type Ref } from 'vue'

export interface VirtualListOptions {
  itemHeight: number | ((index: number) => number)
  overscan?: number
  scrollContainer?: Ref<HTMLElement | null> | HTMLElement | null
}

export function useVirtualList<T>(
  items: Ref<T[]>,
  options: VirtualListOptions
) {
  const containerRef = ref<HTMLElement | null>(null)
  const scrollTop = ref(0)
  const containerHeight = ref(0)

  const overscan = options.overscan ?? 5

  const getItemHeight = (index: number): number => {
    if (typeof options.itemHeight === 'function') {
      return options.itemHeight(index)
    }
    return options.itemHeight
  }

  const totalHeight = computed(() => {
    let total = 0
    for (let i = 0; i < items.value.length; i++) {
      total += getItemHeight(i)
    }
    return total
  })

  const visibleRange = computed(() => {
    const height = containerHeight.value
    const top = scrollTop.value

    let offset = 0
    let start = 0
    for (let i = 0; i < items.value.length; i++) {
      const h = getItemHeight(i)
      if (offset + h > top) {
        start = i
        break
      }
      offset += h
      start = i + 1
    }

    let end = start
    let visibleHeight = 0
    for (let i = start; i < items.value.length; i++) {
      const h = getItemHeight(i)
      visibleHeight += h
      end = i + 1
      if (visibleHeight > height) break
    }

    start = Math.max(0, start - overscan)
    end = Math.min(items.value.length, end + overscan)

    return { start, end, offsetTop: offset }
  })

  const visibleItems = computed(() => {
    const { start, end } = visibleRange.value
    return items.value.slice(start, end).map((item, i) => ({
      item,
      index: start + i,
    }))
  })

  const containerStyle = computed(() => ({
    height: totalHeight.value + 'px',
    position: 'relative' as const,
  }))

  const getItemStyle = (index: number) => {
    let top = 0
    for (let i = 0; i < index; i++) {
      top += getItemHeight(i)
    }
    return {
      position: 'absolute' as const,
      top: top + 'px',
      left: 0,
      right: 0,
      height: getItemHeight(index) + 'px',
    }
  }

  const handleScroll = () => {
    const el = options.scrollContainer
      ? (options.scrollContainer instanceof HTMLElement ? options.scrollContainer : options.scrollContainer.value)
      : containerRef.value
    if (el) {
      scrollTop.value = el.scrollTop
    }
  }

  const updateContainerHeight = () => {
    const el = options.scrollContainer
      ? (options.scrollContainer instanceof HTMLElement ? options.scrollContainer : options.scrollContainer.value)
      : containerRef.value?.parentElement
    if (el) {
      containerHeight.value = el.clientHeight
    }
  }

  const scrollToIndex = (index: number, behavior: ScrollBehavior = 'smooth') => {
    let top = 0
    for (let i = 0; i < index; i++) {
      top += getItemHeight(i)
    }
    const el = options.scrollContainer
      ? (options.scrollContainer instanceof HTMLElement ? options.scrollContainer : options.scrollContainer.value)
      : containerRef.value?.parentElement
    if (el) {
      el.scrollTo({ top, behavior })
    }
  }

  onMounted(() => {
    const scrollEl = options.scrollContainer
      ? (options.scrollContainer instanceof HTMLElement ? options.scrollContainer : options.scrollContainer.value)
      : containerRef.value?.parentElement

    if (scrollEl) {
      scrollEl.addEventListener('scroll', handleScroll, { passive: true })
      updateContainerHeight()
    }

    window.addEventListener('resize', updateContainerHeight)
  })

  onUnmounted(() => {
    const scrollEl = options.scrollContainer
      ? (options.scrollContainer instanceof HTMLElement ? options.scrollContainer : options.scrollContainer.value)
      : containerRef.value?.parentElement

    if (scrollEl) {
      scrollEl.removeEventListener('scroll', handleScroll)
    }

    window.removeEventListener('resize', updateContainerHeight)
  })

  return {
    containerRef,
    visibleItems,
    visibleRange,
    totalHeight,
    containerStyle,
    getItemStyle,
    scrollToIndex,
    updateContainerHeight,
  }
}