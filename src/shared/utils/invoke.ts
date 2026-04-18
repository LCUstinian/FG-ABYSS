import { invoke as tauriInvoke } from '@tauri-apps/api/core'
import type { AppError } from '@/shared/types/error'

export async function invoke<T>(cmd: string, args?: Record<string, unknown>): Promise<T> {
  try {
    return await tauriInvoke<T>(cmd, args)
  } catch (e) {
    throw e as AppError
  }
}
