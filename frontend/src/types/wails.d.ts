declare module '@wailsapp/runtime' {
  export interface Events {
    Emit(event: string, data?: any): void;
    On(event: string, callback: (data: any) => void): void;
  }
  
  export const Events: Events;
}
