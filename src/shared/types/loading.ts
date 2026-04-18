// Use LoadingMap instead of a single boolean. Allows loading.list, loading.create,
// loading['delete-uuid'] to be independent per CLAUDE.md.
export type LoadingMap = Record<string, boolean>

export function useLoadingMap(): LoadingMap {
  return {}
}
