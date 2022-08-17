interface Options {
  /**
  The RGB output format.

  Note that when using the `css` format, the value of the alpha channel is rounded to two decimal places.

  @default 'object'
  */
  readonly format?: "object" | "array" | "css";
}
