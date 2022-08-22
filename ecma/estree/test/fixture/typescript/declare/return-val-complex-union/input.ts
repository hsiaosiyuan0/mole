interface a {
  hasPressHandler: () =>
    | (((event: import("react-native").GestureResponderEvent) => void) &
        (() => void))
    | undefined;
}
