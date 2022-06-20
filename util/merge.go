package util

func MergeMap(dst map[string]interface{}, src map[string]interface{}) {
	for k, v := range src {
		if dv, ok := dst[k]; ok {
			dvm, dm := dv.(map[string]interface{})
			svm, sm := v.(map[string]interface{})
			if dm && sm {
				MergeMap(dvm, svm)
				continue
			}
		}
		dst[k] = v
	}
}

func MergeMaps(dst map[string]interface{}, src ...map[string]interface{}) {
	for _, s := range src {
		MergeMap(dst, s)
	}
}
