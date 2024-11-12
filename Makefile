help: # Show this help
	@egrep -h '\s#\s' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?# "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

ximg: # Convert image to base64
	go run cmd/ximg/main.go "https://tse2-mm.cn.bing.net/th/id/OIP-C.lSR62jSbA_RtjWMojuC-FgHaKx?rs=1&pid=ImgDetMain"

