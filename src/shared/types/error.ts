export interface AppError {
  kind:
    | 'Database' | 'Http' | 'Crypto' | 'Connection' | 'InvalidResponse'
    | 'CircuitOpen' | 'NotFound' | 'InvalidInput' | 'Io' | 'Serialize'
    | 'Locked' | 'Plugin' | 'NeedsRedeploy' | 'MemShellExpired'
  message: string
}
