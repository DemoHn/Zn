package exec

const (
	LIB_TYPE_STD    = 1 // 标准库
	LIB_TYPE_VENDOR = 2 // 第三方库
	LIB_TYPE_CUSTOM = 3 // 用户自定义模块
)

type LibNameInfo struct {
	// libName original string
	originalName string
	// parsed libType
	libType uint8
	// separate full libstring to subPath e.g.: "A-B-C" -> []string{"A", "B", "C"}
	libPath []string
}

func parseLibName(libName string) LibNameInfo {
	return LibNameInfo{originalName: libName}
}
