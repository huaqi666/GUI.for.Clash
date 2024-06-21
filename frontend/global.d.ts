export {}

declare global {
  interface Window {
    Plugins: any
  }

  /**
   * Data format returned by wails events
   */
  type WailsEventsResponse = { name: string; sender: string; data: string; Cancelled: boolean }

  /**
   * The variable is initialized in `globalMethods.ts:20`
   */
  var AsyncFunction: FunctionConstructor
}
