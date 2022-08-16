package pack

var nodeBuiltin = map[string]bool{
	"assert":         true,
	"buffer":         true,
	"child_process":  true,
	"cluster":        true,
	"crypto":         true,
	"dgram":          true,
	"dns":            true,
	"domain":         true,
	"events":         true,
	"fs":             true,
	"http":           true,
	"https":          true,
	"net":            true,
	"os":             true,
	"path":           true,
	"querystring":    true,
	"readline":       true,
	"stream":         true,
	"string_decoder": true,
	"timers":         true,
	"tls":            true,
	"tty":            true,
	"url":            true,
	"util":           true,
	"v8":             true,
	"vm":             true,
	"zlib":           true,
}

// refer: https://github.com/facebook/metro/blob/main/packages/metro-react-native-babel-preset/src/configs/lazy-imports.js
var rnBuiltin = map[string]bool{
	"AccessibilityInfo":        true,
	"ActivityIndicator":        true,
	"Button":                   true,
	"DatePickerIOS":            true,
	"DrawerLayoutAndroid":      true,
	"FlatList":                 true,
	"Image":                    true,
	"ImageBackground":          true,
	"InputAccessoryView":       true,
	"KeyboardAvoidingView":     true,
	"MaskedViewIOS":            true,
	"Modal":                    true,
	"Pressable":                true,
	"ProgressBarAndroid":       true,
	"ProgressViewIOS":          true,
	"SafeAreaView":             true,
	"ScrollView":               true,
	"SectionList":              true,
	"Slider":                   true,
	"Switch":                   true,
	"RefreshControl":           true,
	"StatusBar":                true,
	"Text":                     true,
	"TextInput":                true,
	"Touchable":                true,
	"TouchableHighlight":       true,
	"TouchableNativeFeedback":  true,
	"TouchableOpacity":         true,
	"TouchableWithoutFeedback": true,
	"View":                     true,
	"VirtualizedList":          true,
	"VirtualizedSectionList":   true,

	"ReactNativeART":      true,
	"warnOnce":            true,
	"TurboModuleRegistry": true,
	"TimePickerAndroid":   true,

	// APIs
	"ActionSheetIOS":      true,
	"Alert":               true,
	"Animated":            true,
	"Appearance":          true,
	"AppRegistry":         true,
	"AppState":            true,
	"AsyncStorage":        true,
	"BackHandler":         true,
	"Clipboard":           true,
	"DeviceInfo":          true,
	"Dimensions":          true,
	"Easing":              true,
	"ReactNative":         true,
	"I18nManager":         true,
	"InteractionManager":  true,
	"Keyboard":            true,
	"LayoutAnimation":     true,
	"Linking":             true,
	"LogBox":              true,
	"NativeEventEmitter":  true,
	"PanResponder":        true,
	"PermissionsAndroid":  true,
	"PixelRatio":          true,
	"PushNotificationIOS": true,
	"Settings":            true,
	"Share":               true,
	"StyleSheet":          true,
	"Systrace":            true,
	"ToastAndroid":        true,
	"TVEventHandler":      true,
	"UIManager":           true,
	"UTFSequence":         true,
	"Vibration":           true,

	// Plugins
	"RCTDeviceEventEmitter":    true,
	"RCTNativeAppEventEmitter": true,
	"NativeModules":            true,
	"Platform":                 true,
	"processColor":             true,
	"requireNativeComponent":   true,

	// deprecated
	"YellowBox":                    true,
	"DeprecatedViewPropTypes":      true,
	"DeprecatedEdgeInsetsPropType": true,
	"DeprecatedPointPropType":      true,
	"DeprecatedColorPropType":      true,
	"StatusBarIOS":                 true,
	"ViewPagerAndroid":             true,
	"CheckBox":                     true,
	"ImageStore":                   true,
	"PickerIOS":                    true,
	"ImageEditor":                  true,
	"Picker":                       true,
	"SegmentedControlIOS":          true,
	"ToolbarAndroid":               true,
	"CameraRoll":                   true,
	"DatePickerAndroid":            true,
	"ImagePickerIOS":               true,

	"dismissKeyboard": true,
}
