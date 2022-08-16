export = Creator;
export as namespace Creator;

declare namespace Creator {
  interface ConstItem {
    key: string
    label: string
    name?: string
    title?: string
    value: any
  }
}