import { ref } from 'vue'

export function usePagination(defaultPageSize = 10) {
  const page = ref(1)
  const pageSize = ref(defaultPageSize)
  const total = ref(0)
  return { page, pageSize, total }
}
